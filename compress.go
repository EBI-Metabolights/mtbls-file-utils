package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	includePattern string
	verbose        bool
)

func init() {
	flag.StringVar(&includePattern, "include", "*", "Glob pattern to include files (e.g. *.d, *.raw, *)")
	flag.BoolVar(&verbose, "verbose", false, "Print each added file")
}

type Task struct {
	index      int
	folderName string
	subfolder  string
	zipPath    string
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: mtbls-compress <folder_path> [--include=<folder pattern. e.g., *.d, *.raw, *] [--verbose]")
		return
	}

	root := strings.TrimSuffix(flag.Arg(0), string(os.PathSeparator))
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		fmt.Println("Invalid folder:", root)
		return
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Failed to read directory:", err)
		return
	}

	originalDir := root + "_original"
	err = os.MkdirAll(originalDir, 0755)
	if err != nil {
		fmt.Println("Failed to create _original folder:", err)
		return
	}

	// Prepare tasks
	var tasks []Task
	for _, entry := range entries {
		if entry.IsDir() && !isHidden(entry.Name()) {
			subfolder := filepath.Join(root, entry.Name())
			if !isEmptyFolder(subfolder) {
				zipPath := filepath.Join(root, entry.Name()+".zip")
				tasks = append(tasks, Task{
					index:      len(tasks) + 1,
					folderName: entry.Name(),
					subfolder:  subfolder,
					zipPath:    zipPath,
				})
			}
		}
	}

	total := len(tasks)
	if total == 0 {
		fmt.Println("No non-empty, visible folders found.")
		return
	}

	fmt.Printf("Found %d folders to process.\n", total)

	for i, task := range tasks {
		fmt.Printf("[%d/%d] Checking: %s\n", i+1, total, task.folderName)

		if fileExists(task.zipPath) {
			if isValidZip(task.zipPath) {
				fmt.Println("  ✓ Valid zip exists — skipping")
				continue
			} else {
				fmt.Println("  ⚠️ Invalid zip — re-compressing")
				os.Remove(task.zipPath)
			}
		}

		err := zipFolder(task.subfolder, task.zipPath)
		if err != nil {
			fmt.Println("  Error zipping:", err)
			continue
		}

		// Move original folder to _original directory
		newPath := filepath.Join(originalDir, task.folderName)
		err = os.Rename(task.subfolder, newPath)
		if err != nil {
			fmt.Println("  ⚠️ Failed to move folder to _original:", err)
		}
	}

	fmt.Println("✅ All done.")
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isEmptyFolder(path string) bool {
	hasFiles := false
	_ = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			match, _ := filepath.Match(strings.ToLower(includePattern), strings.ToLower(filepath.Base(p)))
			if match {
				hasFiles = true
				return io.EOF
			}
		}
		return nil
	})
	return !hasFiles
}

func isValidZip(path string) bool {
	r, err := zip.OpenReader(path)
	if err != nil {
		return false
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err == nil {
			rc.Close()
			return true
		}
	}
	return false
}

func zipFolder(folderPath, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		match, err := filepath.Match(strings.ToLower(includePattern), strings.ToLower(filepath.Base(path)))
		if err != nil || !match {
			return nil
		}

		relPath, err := filepath.Rel(filepath.Dir(folderPath), path)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Println("  +", relPath)
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}
