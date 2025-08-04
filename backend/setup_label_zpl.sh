#!/bin/bash

# Define base paths
BASE_DIR=$(pwd)
PRINTER_DIR="$BASE_DIR/internal/printer"
ZPL_DIR="$BASE_DIR/printers/zpl"
BAT_DIR="$BASE_DIR/printers/bat"
MODELS_DIR="$BASE_DIR/models"

# Create necessary folders
mkdir -p "$PRINTER_DIR" "$ZPL_DIR" "$BAT_DIR" "$MODELS_DIR"

# --- models/label.go ---
cat << 'EOF' > "$MODELS_DIR/label.go"
package models

import (
	"time"

	"github.com/google/uuid"
)

type Label struct {
	ID             uuid.UUID `json:"id" db:"id"`
	LabelID        string    `json:"label_id" db:"label_id" binding:"required"`
	Location       *string   `json:"location" db:"location"`
	BundleNo       string    `json:"bundle_no" db:"bundle_no"`
	BundleType     string    `json:"bundle_type" db:"bundle_type"`
	PQD            string    `json:"pqd" db:"pqd"`
	Unit           string    `json:"unit" db:"unit"`
	Time           string    `json:"time" db:"time"`
	Length         int       `json:"length" db:"length"`
	HeatNo         string    `json:"heat_no" db:"heat_no"`
	ProductHeading string    `json:"product_heading" db:"product_heading"`
	IsiBottom      string    `json:"isi_bottom" db:"isi_bottom"`
	IsiTop         string    `json:"isi_top" db:"isi_top"`
	ChargeDtm      string    `json:"charge_dtm" db:"charge_dtm"`
	Mill           string    `json:"mill" db:"mill"`
	Grade          string    `json:"grade" db:"grade"`
	UrlApikey      string    `json:"url_apikey" db:"url_apikey"`
	Weight         *string   `json:"weight" db:"weight"`
	Section        string    `json:"section" db:"section"`
	Date           string    `json:"date" db:"date"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Status         string    `json:"status" db:"status"`
	IsDuplicate    bool      `json:"is_duplicate" db:"is_duplicate"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
EOF

# --- internal/printer/zpl_generator.go ---
cat << 'EOF' > "$PRINTER_DIR/zpl_generator.go"
package printer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"your-module/models"
)

func safeString(s string) string {
	return strings.ReplaceAll(s, "^", "")
}

func GenerateLabelZPL(label models.Label) string {
	var zpl strings.Builder

	zpl.WriteString("^XA\n")
	zpl.WriteString("^PW812\n")
	zpl.WriteString("^LL406\n")
	zpl.WriteString("^LS0\n")

	zpl.WriteString("^FO50,50^A0N,50,50^FD" + safeString(label.ProductHeading) + "^FS\n")
	zpl.WriteString("^FO50,100^A0N,30,30^FD" + safeString(label.Grade) + "^FS\n")
	zpl.WriteString("^FO50,130^A0N,25,25^FD" + safeString(label.Section) + "^FS\n")
	zpl.WriteString("^FO50,160^A0N,25,25^FDHEAT: " + safeString(label.HeatNo) + "^FS\n")
	zpl.WriteString("^FO50,185^A0N,25,25^FDBUNDLE: " + safeString(label.BundleNo) + "^FS\n")
	zpl.WriteString("^FO50,210^A0N,25,25^FDPQD: " + safeString(label.PQD) + "^FS\n")
	zpl.WriteString("^FO50,235^A0N,25,25^FD" + safeString(label.Unit) + " | " + safeString(label.Mill) + "^FS\n")
	zpl.WriteString("^FO50,260^A0N,20,20^FDDATE: " + safeString(label.Date) + " | TIME: " + safeString(label.Time) + "^FS\n")
	zpl.WriteString("^FO50,280^A0N,20,20^FDISI: " + safeString(label.IsiTop) + " | " + safeString(label.IsiBottom) + "^FS\n")
	zpl.WriteString("^FO50,300^A0N,20,20^FDLENGTH: " + fmt.Sprintf("%d", label.Length))
	if label.Weight != nil {
		zpl.WriteString(" | WEIGHT: " + safeString(*label.Weight))
	}
	zpl.WriteString("^FS\n")
	zpl.WriteString("^FO50,320^A0N,20,20^FDCHARGE: " + safeString(label.ChargeDtm) + "^FS\n")

	qrData := generateQRData(label)
	zpl.WriteString("^FO600,50^BQN,2,5^FD" + qrData + "^FS\n")
	zpl.WriteString("^FO600,150^BY3^BCN,100,Y,N,N^FD" + safeString(label.PQD) + "^FS\n")
	zpl.WriteString("^FO50,350^A0N,15,15^FDPrinted: " + time.Now().Format("2006-01-02 15:04:05") + "^FS\n")
	zpl.WriteString("^XZ\n")

	return zpl.String()
}

func generateQRData(label models.Label) string {
	return fmt.Sprintf("TMT_BAR|%s|%s|%s|%s|%s|%s|%s",
		safeString(label.PQD),
		safeString(label.HeatNo),
		safeString(label.Grade),
		safeString(label.Section),
		safeString(label.BundleNo),
		safeString(label.Date),
		safeString(label.Time),
	)
}

func GenerateAndSaveZPL(label models.Label) (string, error) {
	zpl := GenerateLabelZPL(label)
	filename := fmt.Sprintf("label_%s_%d.zpl", label.PQD, time.Now().Unix())
	path := filepath.Join("printers", "zpl", filename)

	if err := os.WriteFile(path, []byte(zpl), 0644); err != nil {
		return "", fmt.Errorf("failed to save ZPL file: %w", err)
	}
	return path, nil
}
EOF

# --- internal/printer/printer.go ---
cat << 'EOF' > "$PRINTER_DIR/printer.go"
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
EOF

echo "✅ Label ZPL printing system created successfully."
