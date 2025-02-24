# Go File Organizer 🗂️

Automatically sorts files into folders (Images, Docs, Videos, etc.) using Go concurrency.  
Perfect for organizing messy directories like Downloads or Desktop.



## Features ✨
- **Concurrent file processing** (goroutines + WaitGroups)
- **Dry-run mode** to preview changes
- **Edge case handling**: invalid filenames, permissions, duplicates
- **Version flag** (`-version`)

## Installation 📦
```bash
# Install globally
go install https://github.com/bettjesse/go-file-organizer.git@latest

# Or clone and run
git clone https://github.com/bettjesse/go-file-organizer.git
cd go-file-organizer
go run main.go -dir=~/Downloads -dry-run

## Usage
# Organize a directory (dry-run first!)
go-file-organizer -dir=~/Downloads -dry-run

# Actually move files
go-file-organizer -dir=~/Downloads

# Show version
go-file-organizer -version
