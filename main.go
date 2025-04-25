package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/yash0000001/pie/internal/concurrency"
)

func main() {
	var root string
	var rootArr []string
	scanner := bufio.NewScanner(os.Stdin)
	color.Blue("🔗 Enter a GitHub URL's or local folder path's: ")

	for {
		scanner.Scan()
		root = strings.TrimSpace(scanner.Text())
		if root == "" {
			break
		}
		rootArr = append(rootArr, root)
	}

	grandTotal := concurrency.ProcessPaths(rootArr)

	color.Cyan("\n🧮 Final Total Lines: %d", grandTotal)

}
