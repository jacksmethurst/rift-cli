package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type IgnoreChecker struct {
	patterns []string
	regexps  []*regexp.Regexp
}

func NewIgnoreChecker(repoPath string) (*IgnoreChecker, error) {
	ic := &IgnoreChecker{
		patterns: []string{},
		regexps:  []*regexp.Regexp{},
	}

	// Default ignore patterns
	defaultIgnores := []string{
		".rift",
		".rift/**",
		".DS_Store",
		"Thumbs.db",
	}

	for _, pattern := range defaultIgnores {
		ic.addPattern(pattern)
	}

	// Load .riftignore file if it exists
	riftignorePath := filepath.Join(repoPath, ".riftignore")
	if _, err := os.Stat(riftignorePath); err == nil {
		if err := ic.loadIgnoreFile(riftignorePath); err != nil {
			return nil, err
		}
	}

	return ic, nil
}

func (ic *IgnoreChecker) loadIgnoreFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ic.addPattern(line)
	}

	return scanner.Err()
}

func (ic *IgnoreChecker) addPattern(pattern string) {
	ic.patterns = append(ic.patterns, pattern)
	
	// Convert gitignore-style patterns to regex
	regexPattern := ic.convertToRegex(pattern)
	if regex, err := regexp.Compile(regexPattern); err == nil {
		ic.regexps = append(ic.regexps, regex)
	}
}

func (ic *IgnoreChecker) convertToRegex(pattern string) string {
	// Handle basic gitignore patterns
	pattern = strings.ReplaceAll(pattern, ".", `\.`)
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = strings.ReplaceAll(pattern, "?", ".")
	
	// Handle directory patterns
	if strings.HasSuffix(pattern, "/") {
		pattern = pattern + ".*"
	}
	
	// Handle patterns starting with /
	if strings.HasPrefix(pattern, "/") {
		pattern = "^" + pattern[1:]
	} else {
		pattern = "(^|.*/)" + pattern
	}
	
	pattern = pattern + "(/.*)?$"
	return pattern
}

func (ic *IgnoreChecker) ShouldIgnore(path string) bool {
	// Normalize path separators
	normalizedPath := filepath.ToSlash(path)
	
	// Remove leading ./
	if strings.HasPrefix(normalizedPath, "./") {
		normalizedPath = normalizedPath[2:]
	}
	
	for _, regex := range ic.regexps {
		if regex.MatchString(normalizedPath) {
			return true
		}
	}
	
	return false
}

func (r *Repository) GetAllFiles() ([]string, error) {
	var files []string
	
	ignoreChecker, err := NewIgnoreChecker(r.Path)
	if err != nil {
		return nil, err
	}
	
	err = filepath.Walk(r.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Get relative path
		relPath, err := filepath.Rel(r.Path, path)
		if err != nil {
			return err
		}
		
		// Skip directories and ignored files
		if info.IsDir() || ignoreChecker.ShouldIgnore(relPath) {
			if info.IsDir() && ignoreChecker.ShouldIgnore(relPath) {
				return filepath.SkipDir
			}
			return nil
		}
		
		files = append(files, relPath)
		return nil
	})
	
	return files, err
}

func (r *Repository) AddAllFiles() error {
	files, err := r.GetAllFiles()
	if err != nil {
		return err
	}
	
	addedCount := 0
	for _, file := range files {
		if err := r.AddFile(file); err != nil {
			return err
		}
		addedCount++
	}
	
	return nil
}