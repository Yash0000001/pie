package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath" // <--- this is the correct standard library package
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type FileStats struct {
	Path      string
	LineCount int
}

func isCodeFile(filename string) bool {
	ext := filepath.Ext(filename)
	codeExtensions := map[string]bool{
		// Programming languages
		".go":    true,
		".py":    true,
		".js":    true,
		".ts":    true,
		".jsx":   true,
		".tsx":   true,
		".java":  true,
		".c":     true,
		".cpp":   true,
		".cc":    true,
		".cs":    true,
		".rb":    true,
		".rs":    true,
		".php":   true,
		".kt":    true,
		".swift": true,
		".m":     true, // Objective-C
		".mm":    true, // Objective-C++
		".scala": true,
		".dart":  true,
		".lua":   true,
		".sh":    true,
		".bash":  true,
		".zsh":   true,
		".r":     true,
		".jl":    true, // Julia

		// Web and markup
		".html": true,
		".htm":  true,
		".css":  true,
		".scss": true,
		".sass": true,
		".less": true,
		".xml":  true,
		".yaml": true,
		".yml":  true,
		".toml": true,
		".md":   true,

		// Config & CI/CD
		".env":        true,
		".dockerfile": true,
		".ini":        true,
		".cfg":        true,
		".conf":       true,

		// SQL & DB
		".sql": true,
		".db":  false,

		// Data & Templates
		".csv":        true,
		".tsv":        true,
		".handlebars": true,
		".hbs":        true,
		".ejs":        true,
		".njk":        true,

		// Misc
		".makefile": true,
		".mk":       true,
		".gradle":   true,
		".bat":      true,
		".ps1":      true, // PowerShell
	}

	return codeExtensions[strings.ToLower(ext)]
}

func countLinesInfile(filePath string) (int, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0
	for scanner.Scan() {
		input := scanner.Text()

		if input != "" {
			count++
		}

	}
	return count, scanner.Err()
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

func main() {
	var root string
	color.Blue("🔗 Enter a GitHub URL or local folder path: ")

	fmt.Scan(&root) //taking input here

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr // now it will write without affecting color
	s.Suffix = " Counting lines..."
	s.Start()

	var flag bool
	if strings.Contains(root, "github.com") {
		flag = true
		color.Red("Yeah so it is a github url")
	}

	tempDir := "./temp-dir"
	defer os.RemoveAll("./temp-dir")
	if flag {
		err := gitClone(root, tempDir)
		if err != nil {
			color.Red("Error in checking your repo")
			return
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
		return
	}

	s.Stop()

	total := 0
	for _, file := range stats {
		resultPath := strings.Replace(file.Path, "temp-dir", "", -1)
		color.HiYellow("total Lines on %s: %d \n", resultPath, file.LineCount)
		total += file.LineCount
	}

	if total > 1000 {
		fmt.Printf("\n🎉 Well done! You wrote %d lines of code. Go get a chai! ☕\n", total)
	} else {
		fmt.Printf("\n🚀 Keep grinding! You wrote %d lines. Great start! 💻\n", total)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("✅ Done!")
}
