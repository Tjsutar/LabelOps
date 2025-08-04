package printer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"labelops-backend/models"
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
	filename := fmt.Sprintf("label_%s_%d.zpl", label.LabelID, time.Now().Unix())
	path := filepath.Join("printers", "zpl", filename)

	if err := os.WriteFile(path, []byte(zpl), 0644); err != nil {
		return "", fmt.Errorf("failed to save ZPL file: %w", err)
	}
	return path, nil
}
