package utils

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"labelops-backend/models"
)

// safeString safely converts a string, handling nil/empty values
func safeString(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

// GenerateLabelZPL generates ZPL content for labels
func GenerateLabelZPL(label models.Label) string {
	var zpl strings.Builder

	// Start ZPL
	zpl.WriteString("^XA\n") // Start of label

	// Set label dimensions and orientation
	zpl.WriteString("^PW812\n") // Print width: 4 inches (812 dots)
	zpl.WriteString("^LL406\n") // Label length: 2 inches (406 dots)
	zpl.WriteString("^LS0\n")   // Left margin: 0

	// Header section
	zpl.WriteString("^FO50,50^A0N,50,50^FD")
	zpl.WriteString(safeString(label.ProductHeading))
	zpl.WriteString("^FS\n")

	// Grade and Section
	zpl.WriteString("^FO50,100^A0N,30,30^FD")
	zpl.WriteString(safeString(label.Grade))
	zpl.WriteString("^FS\n")

	zpl.WriteString("^FO50,130^A0N,25,25^FD")
	zpl.WriteString(safeString(label.Section))
	zpl.WriteString("^FS\n")

	// Heat Number
	zpl.WriteString("^FO50,160^A0N,25,25^FD")
	zpl.WriteString("HEAT: ")
	zpl.WriteString(safeString(label.HeatNo))
	zpl.WriteString("^FS\n")

	// Bundle Number
	zpl.WriteString("^FO50,185^A0N,25,25^FD")
	zpl.WriteString("BUNDLE: ")
	zpl.WriteString(safeString(label.BundleNo))
	zpl.WriteString("^FS\n")

	// PQD (Primary Quality Data)
	zpl.WriteString("^FO50,210^A0N,25,25^FD")
	zpl.WriteString("PQD: ")
	zpl.WriteString(safeString(label.PQD))
	zpl.WriteString("^FS\n")

	// Unit and Mill
	zpl.WriteString("^FO50,235^A0N,25,25^FD")
	zpl.WriteString(safeString(label.Unit))
	zpl.WriteString(" | ")
	zpl.WriteString(safeString(label.Mill))
	zpl.WriteString("^FS\n")

	// Date and Time
	zpl.WriteString("^FO50,260^A0N,20,20^FD")
	zpl.WriteString("DATE: ")
	zpl.WriteString(safeString(label.Date))
	zpl.WriteString(" | TIME: ")
	zpl.WriteString(safeString(label.Time))
	zpl.WriteString("^FS\n")

	// ISI Standards
	zpl.WriteString("^FO50,280^A0N,20,20^FD")
	zpl.WriteString("ISI: ")
	zpl.WriteString(safeString(label.IsiTop))
	zpl.WriteString(" | ")
	zpl.WriteString(safeString(label.IsiBottom))
	zpl.WriteString("^FS\n")

	// Length and Weight
	zpl.WriteString("^FO50,300^A0N,20,20^FD")
	zpl.WriteString("LENGTH: ")
	zpl.WriteString(fmt.Sprintf("%d", label.Length))
	if label.Weight != nil {
		zpl.WriteString(" | WEIGHT: ")
		zpl.WriteString(safeString(*label.Weight))
	}
	zpl.WriteString("^FS\n")

	// Charge DTM
	zpl.WriteString("^FO50,320^A0N,20,20^FD")
	zpl.WriteString("CHARGE: ")
	zpl.WriteString(safeString(label.ChargeDtm))
	zpl.WriteString("^FS\n")

	// QR Code with label data
	qrData := generateQRData(label)
	zpl.WriteString("^FO600,50^BQN,2,5^FD")
	zpl.WriteString(qrData)
	zpl.WriteString("^FS\n")

	// Barcode for scanning
	zpl.WriteString("^FO600,150^BY3^BCN,100,Y,N,N^FD")
	zpl.WriteString(safeString(label.PQD))
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
func generateQRData(label models.Label) string {
	var data strings.Builder
	data.WriteString("TMT_BAR|")
	data.WriteString(safeString(label.PQD))
	data.WriteString("|")
	data.WriteString(safeString(label.HeatNo))
	data.WriteString("|")
	data.WriteString(safeString(label.Grade))
	data.WriteString("|")
	data.WriteString(safeString(label.Section))
	data.WriteString("|")
	data.WriteString(safeString(label.BundleNo))
	data.WriteString("|")
	data.WriteString(safeString(label.Date))
	data.WriteString("|")
	data.WriteString(safeString(label.Time))
	return data.String()
}

// GenerateLabelsCSV generates CSV data for labels
// func GenerateLabelsCSV(rows interface{}) string {
// 	// This would be implemented to convert database rows to CSV format
// 	// For now, returning a placeholder
// 	return "Label ID,Location,Bundle Nos,PQD,Unit,Time,Length,Heat No,Product Heading,ISI Bottom,ISI Top,Charge DTM,Mill,Grade,Weight,Section,Date,Printed At,Status,Is Duplicate,Created At\n"
// }

// GenerateLabelsCSV generates CSV data from *sql.Rows
func GenerateLabelsCSV(rows *sql.Rows) string {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// Define CSV header
	headers := []string{
		"Label ID", "Location", "Bundle Nos", "PQD", "Unit", "Time", "Length",
		"Heat No", "Product Heading", "ISI Bottom", "ISI Top", "Charge DTM",
		"Mill", "Grade", "URL API Key", "Weight", "Section", "Date",
		"Printed At", "Status", "Is Duplicate", "Created At",
	}
	writer.Write(headers)

	for rows.Next() {
		var (
			labelID, location, bundleNos, pqd, unit, time1                    sql.NullString
			length                                                            sql.NullInt64
			heatNo, productHeading, isiBottom, isiTop, chargeDtm, mill, grade sql.NullString
			urlAPIKey, weight, section, date, status                          sql.NullString
			isDuplicate                                                       sql.NullBool
			createdAt                                                         sql.NullTime
		)

		err := rows.Scan(
			&labelID, &location, &bundleNos, &pqd, &unit, &time1, &length,
			&heatNo, &productHeading, &isiBottom, &isiTop, &chargeDtm, &mill, &grade,
			&urlAPIKey, &weight, &section, &date, &status, &isDuplicate, &createdAt,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue // skip the row on error
		}

		record := []string{
			nullToStr(labelID), nullToStr(location), nullToStr(bundleNos), nullToStr(pqd),
			nullToStr(unit), nullToStr(time1), nullInt64ToStr(length), nullToStr(heatNo),
			nullToStr(productHeading), nullToStr(isiBottom), nullToStr(isiTop),
			nullToStr(chargeDtm), nullToStr(mill), nullToStr(grade),
			nullToStr(urlAPIKey), nullToStr(weight), nullToStr(section), nullToStr(date),
			nullToStr(status), nullBoolToStr(isDuplicate), nullTimeToStr(createdAt),
		}

		writer.Write(record)
	}

	writer.Flush()
	return buffer.String()
}

// Helper functions
func nullToStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullBoolToStr(nb sql.NullBool) string {
	if nb.Valid {
		return strconv.FormatBool(nb.Bool)
	}
	return ""
}

func nullTimeToStr(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}

// GeneratePrintJobsCSV generates CSV data from print jobs *sql.Rows
func GeneratePrintJobsCSV(rows *sql.Rows) string {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// Define CSV header for print jobs
	headers := []string{
		"Job ID", "Label ID", "User ID", "Status", "Max Retries", "Retry Count",
		"Created At", "Updated At",
	}
	writer.Write(headers)

	for rows.Next() {
		var (
			jobID, labelID, userID, status sql.NullString
			maxRetries, retries            sql.NullInt32
			createdAt, updatedAt           sql.NullTime
		)

		err := rows.Scan(
			&jobID, &labelID, &userID, &status, &maxRetries, &retries,
			&createdAt, &updatedAt,
		)
		if err != nil {
			log.Println("Error scanning print job row:", err)
			continue // skip the row on error
		}

		record := []string{
			nullToStr(jobID), nullToStr(labelID), nullToStr(userID), nullToStr(status),
			nullInt32ToStr(maxRetries), nullInt32ToStr(retries),
			nullTimeToStr(createdAt), nullTimeToStr(updatedAt),
		}

		writer.Write(record)
	}

	writer.Flush()
	return buffer.String()
}

// Helper functions for nullable integers
func nullInt64ToStr(ni sql.NullInt64) string {
	if ni.Valid {
		return fmt.Sprintf("%d", ni.Int64)
	}
	return ""
}

func nullInt32ToStr(ni sql.NullInt32) string {
	if ni.Valid {
		return fmt.Sprintf("%d", ni.Int32)
	}
	return ""
}
