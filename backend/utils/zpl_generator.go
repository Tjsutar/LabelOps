package utils

import (
	"fmt"
	"strings"
	"time"

	"labelops-backend/models"
)

// GenerateTMTBarZPL generates ZPL content for TMT Bar labels
func GenerateTMTBarZPL(label models.TMTBarLabel) string {
	var zpl strings.Builder

	// Start ZPL
	zpl.WriteString("^XA\n") // Start of label

	// Set label dimensions and orientation
	zpl.WriteString("^PW812\n") // Print width: 4 inches (812 dots)
	zpl.WriteString("^LL406\n") // Label length: 2 inches (406 dots)
	zpl.WriteString("^LS0\n")   // Left margin: 0

	// Header section
	zpl.WriteString("^FO50,50^A0N,50,50^FD")
	zpl.WriteString(label.ProductHeading)
	zpl.WriteString("^FS\n")

	// Grade and Section
	zpl.WriteString("^FO50,100^A0N,30,30^FD")
	zpl.WriteString(label.Grade)
	zpl.WriteString("^FS\n")

	zpl.WriteString("^FO50,130^A0N,25,25^FD")
	zpl.WriteString(label.Section)
	zpl.WriteString("^FS\n")

	// Heat Number
	zpl.WriteString("^FO50,160^A0N,25,25^FD")
	zpl.WriteString("HEAT: ")
	zpl.WriteString(label.HeatNo)
	zpl.WriteString("^FS\n")

	// Bundle Number
	zpl.WriteString("^FO50,185^A0N,25,25^FD")
	zpl.WriteString("BUNDLE: ")
	zpl.WriteString(fmt.Sprintf("%d", label.BundleNos))
	zpl.WriteString("^FS\n")

	// PQD (Primary Quality Data)
	zpl.WriteString("^FO50,210^A0N,25,25^FD")
	zpl.WriteString("PQD: ")
	zpl.WriteString(label.PQD)
	zpl.WriteString("^FS\n")

	// Unit and Mill
	zpl.WriteString("^FO50,235^A0N,25,25^FD")
	zpl.WriteString(label.Unit)
	zpl.WriteString(" | ")
	zpl.WriteString(label.Mill)
	zpl.WriteString("^FS\n")

	// Date and Time
	zpl.WriteString("^FO50,260^A0N,20,20^FD")
	zpl.WriteString("DATE: ")
	zpl.WriteString(label.Date1)
	zpl.WriteString(" | TIME: ")
	zpl.WriteString(label.Time1)
	zpl.WriteString("^FS\n")

	// ISI Standards
	zpl.WriteString("^FO50,280^A0N,20,20^FD")
	zpl.WriteString("ISI: ")
	zpl.WriteString(label.IsiTop)
	zpl.WriteString(" | ")
	zpl.WriteString(label.IsiBottom)
	zpl.WriteString("^FS\n")

	// Length and Weight
	zpl.WriteString("^FO50,300^A0N,20,20^FD")
	zpl.WriteString("LENGTH: ")
	zpl.WriteString(label.Length)
	if label.Weight != nil {
		zpl.WriteString(" | WEIGHT: ")
		zpl.WriteString(*label.Weight)
	}
	zpl.WriteString("^FS\n")

	// Charge DTM
	zpl.WriteString("^FO50,320^A0N,20,20^FD")
	zpl.WriteString("CHARGE: ")
	zpl.WriteString(label.ChargeDtm)
	zpl.WriteString("^FS\n")

	// QR Code with label data
	qrData := generateQRData(label)
	zpl.WriteString("^FO600,50^BQN,2,5^FD")
	zpl.WriteString(qrData)
	zpl.WriteString("^FS\n")

	// Barcode for scanning
	zpl.WriteString("^FO600,150^BY3^BCN,100,Y,N,N^FD")
	zpl.WriteString(label.PQD)
	zpl.WriteString("^FS\n")

	// Footer with timestamp
	zpl.WriteString("^FO50,350^A0N,15,15^FD")
	zpl.WriteString("Printed: ")
	zpl.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	zpl.WriteString("^FS\n")

	// End ZPL
	zpl.WriteString("^XZ\n")

	return zpl.String()
}

// generateQRData creates QR code data for the label
func generateQRData(label models.TMTBarLabel) string {
	var data strings.Builder
	data.WriteString("TMT_BAR|")
	data.WriteString(label.PQD)
	data.WriteString("|")
	data.WriteString(label.HeatNo)
	data.WriteString("|")
	data.WriteString(label.Grade)
	data.WriteString("|")
	data.WriteString(label.Section)
	data.WriteString("|")
	data.WriteString(fmt.Sprintf("%d", label.BundleNos))
	data.WriteString("|")
	data.WriteString(label.Date1)
	data.WriteString("|")
	data.WriteString(label.Time1)
	return data.String()
}

// GenerateTMTBarLabelsCSV generates CSV data for TMT Bar labels
func GenerateTMTBarLabelsCSV(rows interface{}) string {
	// This would be implemented to convert database rows to CSV format
	// For now, returning a placeholder
	return "Label ID,Location,Bundle Nos,PQD,Unit,Time,Length,Heat No,Product Heading,ISI Bottom,ISI Top,Charge DTM,Mill,Grade,Weight,Section,Date,Printed At,Status,Is Duplicate,Created At\n"
} 