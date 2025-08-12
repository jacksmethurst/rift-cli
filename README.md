# Rift CLI

A Git alternative written in Go. Rift provides version control functionality with a simple, intuitive command-line interface.

## Features

- Repository initialization and management
- File staging and committing with SHA256 hashing
- `.riftignore` support (similar to `.gitignore`)
- Bulk file operations (`rift add .`)
- Status checking and commit history
- NextJS project support out of the box

## Installation

Install Node.js 18+, then run:

```bash
npm install -g rift-cli
```

### Manual Installation:

Download the binary for your platform from [Releases](https://github.com/jacksmethurst/rift-cli/releases) and add to your PATH.

## Usage

### Initialize a repository
```bash
rift init
```

### Add files to staging area
```bash
# Add a specific file
rift add myfile.txt

# Add all files (respects .riftignore)
rift add .
```

### Commit changes
```bash
rift commit "Your commit message"
```

### Check repository status
```bash
rift status
```

### View commit history
```bash
rift log
```

## .riftignore Support

Create a `.riftignore` file in your repository root to exclude files and directories:

```
# Dependencies
node_modules/
.env

# Build output
dist/
.next/

# OS files
.DS_Store
Thumbs.db
```

## Repository Structure

Rift stores repository data in a `.rift` directory:

```
.rift/
├── objects/     # File content snapshots (SHA256 hashed)
├── refs/        # Branch references
├── HEAD         # Current branch pointer
└── index        # Staging area
```

## Commands

- `rift init` - Initialize a new repository
- `rift add <file>` - Add file to staging area
- `rift add .` - Add all files (respecting .riftignore)
- `rift commit <message>` - Commit staged changes
- `rift status` - Show repository status
- `rift log` - Show commit history
- `rift clone <url>` - Clone a repository (coming soon)
- `rift push` - Push changes to remote (coming soon)
- `rift pull` - Pull changes from remote (coming soon)

## Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
