package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	allowedChars     = regexp.MustCompile(`[^0-9a-zA-Z._+\-]`)
	plusChar         = regexp.MustCompile(`[+]`)
	validNamePattern = regexp.MustCompile(`^[0-9a-zA-Z._\-]+$`)
)

func main() {
	dryRun := flag.Bool("dry", true, "Perform a dry run without renaming")
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  mtbls-rename [--dry=true] <folder_path>")
		fmt.Println("  mtbls-rename [--dry=false] <folder_path>")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("❌ Error: No folder path provided.")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}
	fmt.Printf("🔧 Dry run mode: %v\n", *dryRun)

	root := strings.TrimSuffix(flag.Arg(0), string(os.PathSeparator))
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		fmt.Println("Invalid folder:", root)
		return
	}

	var pathsToRename []struct {
		oldPath string
		newPath string
	}

	// Step 1: Collect all rename candidates (bottom-up)
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error reading path:", err)
			return nil
		}

		if path == root {
			return nil // skip root
		}

		oldName := filepath.Base(path)

		// Skip if already valid and no "+"
		if validNamePattern.MatchString(oldName) && !strings.Contains(oldName, "+") {
			return nil
		}

		sanitized := sanitizeName(oldName)
		if oldName == sanitized {
			return nil
		}

		dir := filepath.Dir(path)
		base := sanitized
		ext := filepath.Ext(base)
		nameOnly := strings.TrimSuffix(base, ext)

		candidate := base
		newPath := filepath.Join(dir, candidate)
		suffix := 1

		for fileExists(newPath) && newPath != path {
			candidate = fmt.Sprintf("%s_%d%s", nameOnly, suffix, ext)
			newPath = filepath.Join(dir, candidate)
			suffix++
		}

		pathsToRename = append(pathsToRename, struct {
			oldPath string
			newPath string
		}{path, newPath})

		return nil
	})

	if err != nil {
		fmt.Println("Error during scan:", err)
		return
	}

	// Reverse order for bottom-up renaming
	for i := len(pathsToRename) - 1; i >= 0; i-- {
		pair := pathsToRename[i]
		fmt.Printf("%s → %s\n", pair.oldPath, pair.newPath)

		if !*dryRun {
			err := os.Rename(pair.oldPath, pair.newPath)
			if err != nil {
				fmt.Println("  ⚠️ Rename failed:", err)
			}
		}
	}

	if *dryRun {
		if len(pathsToRename) > 0 {
			fmt.Printf("🔎 Dry run complete. %d item(s) would be renamed.\n", len(pathsToRename))
			fmt.Printf("Now run again with --dry=false argument to rename all: mtbls-rename --dry=false <folder_path>\n")
		} else {
			fmt.Printf("🔎 Dry run complete. There is no item to be renamed.\n")
		}
	} else {
		if len(pathsToRename) > 0 {
			fmt.Printf("✅ Done. Renamed %d item(s).\n", len(pathsToRename))
		} else {
			fmt.Printf("✅ Done. There is no item to be renamed.\n")
		}
	}
}

func sanitizeName(name string) string {
	name = plusChar.ReplaceAllString(name, "_PLUS_")
	name = allowedChars.ReplaceAllString(name, "__")
	return name
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
