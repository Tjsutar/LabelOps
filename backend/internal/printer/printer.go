package printer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func PrintZPLFile(zplPath string) error {
	absPath, err := filepath.Abs(zplPath)
	if err != nil {
		return err
	}
	absPath = filepath.ToSlash(absPath)

	batPath := filepath.Join("printers", "bat", "print_label.bat")
	if _, err := os.Stat(batPath); os.IsNotExist(err) {
		content := `@echo off
copy /B "%1" "\\%COMPUTERNAME%\\Your_Printer_Share_Name"
echo Label sent to printer.
pause`
		os.WriteFile(batPath, []byte(content), 0755)
	}

	cmd := exec.Command("cmd", "/C", batPath, absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to print label: %w", err)
	}
	return nil
}

// PrintZPLBatch handles printing multiple ZPL files.
func PrintZPLBatch(zplPaths []string) error {
	for _, path := range zplPaths {
		fmt.Printf("Printing: %s\n", path)
		if err := PrintZPLFile(path); err != nil {
			return fmt.Errorf("failed printing %s: %w", path, err)
		}
	}
	return nil
}
