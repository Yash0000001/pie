package counter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var SYMBOLS = struct {
	T    string
	Bar  string
	L    string
	Pipe string
}{
	T:    "├",
	Bar:  "─",
	L:    "└",
	Pipe: "│",
}

func printTree(stats []FileStats) (totalLines int) {
	type Node struct {
		Children  map[string]*Node
		IsFile    bool
		LineCount int
	}

	root := &Node{Children: make(map[string]*Node)}

	// Step 1: Build the tree with line counts
	for _, stat := range stats {
		normalized := strings.ReplaceAll(stat.Path, `\`, `/`)
		parts := strings.Split(normalized, "/")
		curr := root
		for i, part := range parts {
			if curr.Children[part] == nil {
				curr.Children[part] = &Node{Children: make(map[string]*Node)}
			}
			curr = curr.Children[part]
			if i == len(parts)-1 {
				curr.IsFile = true
				curr.LineCount = stat.LineCount
			}
		}
	}

	// Step 2: Recursive print
	var walk func(node *Node, prefix string, isLast bool)
	walk = func(node *Node, prefix string, isLast bool) {
		keys := make([]string, 0, len(node.Children))
		for k := range node.Children {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, k := range keys {
			child := node.Children[k]
			isLastChild := i == len(keys)-1
			connector := SYMBOLS.T
			nextPrefix := prefix + SYMBOLS.Pipe + "   "

			if isLastChild {
				connector = SYMBOLS.L
				nextPrefix = prefix + "    "
			}

			if child.IsFile {
				color.Cyan("%s%s%s %s (%d lines)", prefix, connector, SYMBOLS.Bar+SYMBOLS.Bar, k, child.LineCount)
				totalLines += child.LineCount
			} else {
				color.Yellow("%s%s%s %s", prefix, connector, SYMBOLS.Bar+SYMBOLS.Bar, k)
			}

			walk(child, nextPrefix, isLastChild)
		}

	}

	color.Green(".")
	walk(root, "", true)
	return totalLines
}

func countLinesInRoot(root string) ([]FileStats, error) {
	var stats []FileStats
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// non code files
			if !isCodeFile(path) {
				return nil
			}
			lines, err := countLinesInfile(path)

			if err != nil {
				color.Red("Error reading file %s: %v\n", path, err)
				return nil
			}
			stats = append(stats, FileStats{LineCount: lines, Path: path})

		}
		return nil
	})
	return stats, err
}

func gitClone(url string, dest string) error {
	// Use appropriate null device based on OS
	nullDevice := os.DevNull
	if runtime.GOOS == "windows" {
		nullDevice = "NUL" // For Windows
	}

	// Open the null device for discarding output
	nullFile, err := os.OpenFile(nullDevice, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer nullFile.Close()

	// Run the git clone command
	cmd := exec.Command("git", "clone", url, dest)
	cmd.Stdout = nullFile // Redirect standard output to null
	cmd.Stderr = nullFile // Redirect error output to null

	return cmd.Run()
}

func ProcessPath(root string) (int, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr // now it will write without affecting color
	s.Suffix = " Counting lines..."
	s.Start()

	var flag bool
	if strings.Contains(root, "github.com") {
		flag = true
		root = root + ".git"
		color.Red("Yeah so it is a github url")
	}

	tempDir := "./temp-dir"
	defer os.RemoveAll("./temp-dir")
	if flag {
		err := gitClone(root, tempDir)
		if err != nil {
			color.Red("Error in checking your repo")
			return 0, err
		}
	}

	var stats []FileStats
	var err error
	if flag {
		stats, err = countLinesInRoot(tempDir)
	} else {
		stats, err = countLinesInRoot(root)
	}

	if err != nil {
		color.Red("Error:", err)
		return 0, err
	}

	s.Stop()

	total := printTree(stats)

	if total > 1000 {
		fmt.Printf("\n🎉 Well done! You wrote %d lines of code. Go get a chai! ☕\n", total)
	} else {
		fmt.Printf("\n🚀 Keep grinding! You wrote %d lines. Great start! 💻\n", total)
	}

	fmt.Println(time.Now())
	time.Sleep(10 * time.Second)
	fmt.Println("✅ Done!")

	return total, nil
}
