package utils

import (
	"fmt"
	"os"
)

func DeleteTempFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		// If the file doesn't exist, we consider it already deleted
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}
	return nil
}
