# pie 🥧

A concurrent line counter for code repositories.

## Overview

`pie` is a command-line tool that calculates the total number of lines of code across local directories or GitHub repositories. It processes multiple paths concurrently, making it efficient for analyzing multiple codebases simultaneously.

## Features

- **Multi-path support**: Process multiple directories or repositories at once
- **Concurrent execution**: Leverages Go's concurrency for optimal performance
- **GitHub integration**: Clone and analyze GitHub repositories directly
- **Code file recognition**: Intelligently identifies code files by extension
- **Beautiful output**: Color-coded results with helpful statistics

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/pie.git

# Navigate to project directory
cd pie

# Build the binary
go build

# Or run directly
go run main.go
```

## Usage

Run the tool and input the paths you want to analyze:

```bash
$ ./pie
🔗 Enter a GitHub URL's or local folder path's: 
/path/to/your/project
/path/to/another/project
https://github.com/username/repo
```

After entering all paths, press Enter with no input to start processing.

### Input Examples

- Local directories: `/home/user/projects/myapp`
- GitHub repositories: `https://github.com/username/repo`

## Sample Output

```
total Lines on /path/to/project/file1.js: 120 
total Lines on /path/to/project/file2.go: 85 
🎉 Well done! You wrote 205 lines of code. Go get a chai! ☕
2025-04-25 19:08:45.0388347 +0530 IST m=+5.779096701
✅ Done!

🧮 Final Total Lines: 205
```

## Supported File Types

`pie` recognizes a wide variety of code file types, including:

- Programming languages: Go, Python, JavaScript, TypeScript, Java, C/C++, Ruby, and many more
- Web and markup: HTML, CSS, SCSS, XML, YAML, Markdown
- Configuration: ENV, Dockerfile, INI, TOML
- And many others!

## How It Works

1. For each path provided, `pie` launches a separate goroutine
2. GitHub URLs are cloned to a temporary directory
3. Each directory is recursively scanned for code files
4. Line counts are calculated in parallel
5. Results are aggregated and displayed with the total count

## License

MIT License
