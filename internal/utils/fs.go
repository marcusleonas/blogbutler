package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func copyFile(src, dst string, perm os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", dst, err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy data from %s to %s: %w", src, dst, err)
	}

	err = os.Chmod(dst, perm)
	if err != nil {
		fmt.Printf("Warning: failed to set permissions on %s: %v\n", dst, err)
	}

	return nil
}

func CopyDirectory(srcDir, dstDir string) error {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source directory %s does not exist", srcDir)
		}
		return fmt.Errorf("failed to stat source directory %s: %w", srcDir, err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source %s is not a directory", srcDir)
	}

	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dstDir, err)
	}

	walkErr := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Warning: accessing path %q: %v\n", path, err)
			if os.IsPermission(err) {
				if d != nil && d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			return err
		}

		if path == srcDir {
			return nil
		}

		if !d.IsDir() {
			fileName := filepath.Base(path)
			dstPath := filepath.Join(dstDir, fileName)

			info, infoErr := d.Info()
			if infoErr != nil {
				fmt.Printf("Warning: failed to get info for %s: %v. Skipping file.\n", path, infoErr)
				return nil
			}

			copyErr := copyFile(path, dstPath, info.Mode())
			if copyErr != nil {
				fmt.Printf("Error copying file %s: %v\n", path, copyErr)
			}
		}
		return nil
	})

	if walkErr != nil {
		return fmt.Errorf("error during directory walk: %w", walkErr)
	}
	return nil
}
