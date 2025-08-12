package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		initRepo()
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rift add <file>")
			os.Exit(1)
		}
		addFile(os.Args[2])
	case "commit":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rift commit <message>")
			os.Exit(1)
		}
		commit(os.Args[2])
	case "status":
		status()
	case "log":
		log()
	case "clone":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rift clone <url>")
			os.Exit(1)
		}
		clone(os.Args[2])
	case "push":
		push()
	case "pull":
		pull()
	case "version", "--version", "-v":
		printVersion()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Rift CLI - Git Alternative")
	fmt.Println("Usage: rift <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init          Initialize a new repository")
	fmt.Println("  add <file>    Add file to staging area")
	fmt.Println("  commit <msg>  Commit changes with message")
	fmt.Println("  status        Show repository status")
	fmt.Println("  log           Show commit history")
	fmt.Println("  clone <url>   Clone a repository")
	fmt.Println("  push          Push changes to remote")
	fmt.Println("  pull          Pull changes from remote")
	fmt.Println("  version       Show version information")
}

func initRepo() {
	fmt.Println("Initializing new Rift repository...")
	repo := NewRepository(".")
	if err := repo.Init(); err != nil {
		fmt.Printf("Error initializing repository: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Repository initialized successfully!")
}

func addFile(filename string) {
	repo := NewRepository(".")
	
	// Handle adding all files with "."
	if filename == "." {
		fmt.Println("Adding all files...")
		if err := repo.AddAllFiles(); err != nil {
			fmt.Printf("Error adding files: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All files added to staging area")
		return
	}
	
	fmt.Printf("Adding file: %s\n", filename)
	if err := repo.AddFile(filename); err != nil {
		fmt.Printf("Error adding file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("File %s added to staging area\n", filename)
}

func commit(message string) {
	fmt.Printf("Committing with message: %s\n", message)
	repo := NewRepository(".")
	if err := repo.Commit(message); err != nil {
		fmt.Printf("Error committing: %v\n", err)
		os.Exit(1)
	}
}

func status() {
	fmt.Println("Repository status:")
	repo := NewRepository(".")
	if err := repo.Status(); err != nil {
		fmt.Printf("Error getting status: %v\n", err)
		os.Exit(1)
	}
}

func log() {
	fmt.Println("Commit history:")
	// TODO: Implement commit log
}

func clone(url string) {
	fmt.Printf("Cloning repository from: %s\n", url)
	// TODO: Implement repository cloning
}

func push() {
	fmt.Println("Pushing changes...")
	// TODO: Implement push functionality
}

func pull() {
	fmt.Println("Pulling changes...")
	// TODO: Implement pull functionality
}

func printVersion() {
	fmt.Println("Rift CLI v1.0.6")
	fmt.Println("A Git alternative written in Go")
	fmt.Println("https://github.com/jacksmethurst/rift-cli")
}