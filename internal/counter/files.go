package counter

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
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
