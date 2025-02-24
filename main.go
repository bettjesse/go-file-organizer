package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Enhanced File struct with additional metadata.
type File struct {
	Name      string    // e.g., "report.pdf"
	Path      string    // e.g., "/Users/me/Downloads/report.pdf"
	Size      int64     // in bytes
	ModTime   time.Time // Last modified time
	IsDir     bool      // true if it's a directory
	Category  string    // e.g., "Docs", "Images"
	Extension string    // e.g., ".pdf"
}

// Categories maps file types to their valid extensions.
var Categories = map[string][]string{
	"Images": {".jpg", ".jpeg", ".png", ".gif"},
	"Docs":   {".pdf", ".docx", ".txt", ".md"},
	"Videos": {".mp4", ".mov", ".avi", ".mkv"},
	"Audio":  {".mp3", ".wav", ".ogg"},
	// Add more categories as needed.
}

// Categorize assigns a category to the File based on its extension.
func (f *File) Categorize() {
	if f.IsDir {
		f.Category = "Folder"
		return
	}
	ext := f.Extension
	for category, exts := range Categories {
		for _, e := range exts {
			if ext == e {
				f.Category = category
				return
			}
		}
	}
	f.Category = "Other" // Default category if no match is found.
}

// isFileValid checks if the File has a valid name and positive size.
func isFileValid(file File) error {
	if file.IsDir {
		return nil // Directories don't need size/name validation here
	}
	if strings.TrimSpace(file.Name) == "" {
		return errors.New("filename cannot be empty")
	}
	if file.Size <= 0 {
		return errors.New("file size must be positive")
	}
	return nil
}

// scanDir scans the directory at dirPath and returns a slice of File structs.
func scanDir(dirPath string) ([]File, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var files []File
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			// For example: permission denied.
			fmt.Printf("⚠️ Skipping %s: %v\n", entry.Name(), err)
			continue
		}

		file := File{
			Name:      entry.Name(),
			Path:      filepath.Join(dirPath, entry.Name()),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			IsDir:     entry.IsDir(),
			Extension: strings.ToLower(filepath.Ext(entry.Name())),
		}

		// Categorize the file based on its extension.
		file.Categorize()
		files = append(files, file)
	}
	return files, nil
}

// processFile processes a single file: validates it and, if in dry-run mode, prints the intended action.
func processFile(file File, dryRun bool) error {
	start := time.Now()
	defer func() {
		fmt.Printf("Processed %q in %v\n", file.Name, time.Since(start))
	}()
	if file.IsDir {
		return nil // Skip directories
	}

	if err := isFileValid(file); err != nil {
		return err
	}

	if dryRun {
		fmt.Printf("Would move %q to %s\n", file.Name, file.Category)
	} else {
		// TODO: Implement actual file moving logic (e.g., using os.Rename)
		destDir := filepath.Join(filepath.Dir(file.Path), file.Category)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		destPath := filepath.Join(destDir, file.Name)
		if err := os.Rename(file.Path, destPath); err != nil {
			return fmt.Errorf("failed to move file: %v", err)
		}
	}
	return nil
}

func main() {
	// Define command-line flags.
	version := flag.Bool("version", false, "Show version")
	dirPath := flag.String("dir", ".", "Directory to organize")
	dryRun := flag.Bool("dry-run", false, "Preview changes without moving files")
	flag.Parse()

	if *version {
		fmt.Println("v1.0.0")
		os.Exit(0)
	}
	// Scan the directory for files.
	files, err := scanDir(*dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a WaitGroup and an error channel for concurrent processing.
	var wg sync.WaitGroup
	errorChan := make(chan error)

	// Process files concurrently.
	for _, file := range files {
		wg.Add(1)
		go func(f File) {
			defer wg.Done()
			if err := processFile(f, *dryRun); err != nil {
				errorChan <- fmt.Errorf("file %q: %v", f.Name, err)
			}
		}(file)
	}

	// Close the error channel after all goroutines complete.
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Print errors received from the goroutines.
	for err := range errorChan {
		fmt.Printf("❌ Error processing file: %v\n", err)
	}

	fmt.Println("Processing complete!")
}
