package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	RiftDir    = ".rift"
	ObjectsDir = ".rift/objects"
	RefsDir    = ".rift/refs"
	HeadFile   = ".rift/HEAD"
	IndexFile  = ".rift/index"
)

type Commit struct {
	Hash      string
	Message   string
	Timestamp time.Time
	Files     []string
}

type Repository struct {
	Path string
}

func NewRepository(path string) *Repository {
	return &Repository{Path: path}
}

func (r *Repository) Init() error {
	dirs := []string{
		filepath.Join(r.Path, RiftDir),
		filepath.Join(r.Path, ObjectsDir),
		filepath.Join(r.Path, RefsDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	headFile := filepath.Join(r.Path, HeadFile)
	if err := ioutil.WriteFile(headFile, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}

	return nil
}

func (r *Repository) AddFile(filename string) error {
	// Check if file should be ignored
	ignoreChecker, err := NewIgnoreChecker(r.Path)
	if err != nil {
		return fmt.Errorf("failed to create ignore checker: %v", err)
	}
	
	if ignoreChecker.ShouldIgnore(filename) {
		return fmt.Errorf("file is ignored by .riftignore: %s", filename)
	}
	
	fullPath := filepath.Join(r.Path, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	hash, err := r.hashFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to hash file: %v", err)
	}

	if err := r.storeObject(hash, fullPath); err != nil {
		return fmt.Errorf("failed to store object: %v", err)
	}

	return r.updateIndex(filename, hash)
}

func (r *Repository) Commit(message string) error {
	index, err := r.readIndex()
	if err != nil {
		return fmt.Errorf("failed to read index: %v", err)
	}

	if len(index) == 0 {
		return fmt.Errorf("nothing to commit")
	}

	commit := Commit{
		Message:   message,
		Timestamp: time.Now(),
		Files:     make([]string, 0, len(index)),
	}

	for filename := range index {
		commit.Files = append(commit.Files, filename)
	}

	commitHash, err := r.createCommitObject(commit)
	if err != nil {
		return fmt.Errorf("failed to create commit object: %v", err)
	}

	commit.Hash = commitHash

	if err := r.updateHead(commitHash); err != nil {
		return fmt.Errorf("failed to update HEAD: %v", err)
	}

	if err := r.clearIndex(); err != nil {
		return fmt.Errorf("failed to clear index: %v", err)
	}

	fmt.Printf("Committed successfully with hash: %s\n", commitHash[:8])
	return nil
}

func (r *Repository) Status() error {
	index, err := r.readIndex()
	if err != nil {
		return fmt.Errorf("failed to read index: %v", err)
	}

	if len(index) == 0 {
		fmt.Println("Nothing staged for commit")
	} else {
		fmt.Println("Changes to be committed:")
		for filename := range index {
			fmt.Printf("  modified: %s\n", filename)
		}
	}

	return nil
}

func (r *Repository) hashFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (r *Repository) storeObject(hash, filename string) error {
	objectPath := filepath.Join(r.Path, ObjectsDir, hash)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(objectPath, content, 0644)
}

func (r *Repository) updateIndex(filename, hash string) error {
	indexPath := filepath.Join(r.Path, IndexFile)
	
	index, err := r.readIndex()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if index == nil {
		index = make(map[string]string)
	}

	index[filename] = hash

	indexContent := ""
	for file, fileHash := range index {
		indexContent += fmt.Sprintf("%s %s\n", file, fileHash)
	}

	return ioutil.WriteFile(indexPath, []byte(indexContent), 0644)
}

func (r *Repository) readIndex() (map[string]string, error) {
	indexPath := filepath.Join(r.Path, IndexFile)
	
	content, err := ioutil.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	index := make(map[string]string)
	contentStr := string(content)
	
	if contentStr == "" {
		return index, nil
	}

	return index, nil
}

func (r *Repository) createCommitObject(commit Commit) (string, error) {
	commitContent := fmt.Sprintf("message: %s\ntimestamp: %s\nfiles:\n", 
		commit.Message, commit.Timestamp.Format(time.RFC3339))
	
	for _, file := range commit.Files {
		commitContent += fmt.Sprintf("  %s\n", file)
	}

	hasher := sha256.New()
	hasher.Write([]byte(commitContent))
	hash := hex.EncodeToString(hasher.Sum(nil))

	objectPath := filepath.Join(r.Path, ObjectsDir, hash)
	return hash, ioutil.WriteFile(objectPath, []byte(commitContent), 0644)
}

func (r *Repository) updateHead(commitHash string) error {
	headPath := filepath.Join(r.Path, HeadFile)
	return ioutil.WriteFile(headPath, []byte(fmt.Sprintf("commit: %s\n", commitHash)), 0644)
}

func (r *Repository) clearIndex() error {
	indexPath := filepath.Join(r.Path, IndexFile)
	return os.Remove(indexPath)
}