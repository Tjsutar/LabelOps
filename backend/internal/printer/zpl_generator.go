package printer

import (
	"fmt"
	"labelops-backend/models"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// safeString removes problematic ZPL chars
func safeString(s string) string {
	return strings.ReplaceAll(s, "^", "")
}

// safeStringPtr safely handles pointer-to-string with nil check
func safeStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return safeString(*s)
}

// GenerateLabelZPL builds ZPL from models.Label
func GenerateLabelZPL(label models.Label) string {
	template := `^XA
^MMT
^PW812
^LL609
^LS0
^FO161,16^GB65,556,3^FS
^FT91,504^A0B,30,30^FH\^CI28^FDIN^FS^CI27
^FT130,525^A0B,30,30^FH\^CI28^FDINDIA^FS^CI27
^FT206,343^A0B,32,32^FH\^CI28^FDCHANNEL^FS^CI27
^FT53,528^A0B,30,30^FH\^CI28^FDMADE^FS^CI27
^FO16,410^GB130,160,3^FS
^FT590,218^A0B,34,33^FH\^CI28^FDV	MILL^FS^CI27
^FT483,570^A0B,25,25^FH\^CI28^FDID^FS^CI27
^FT516,570^A0B,31,30^FH\^CI28^FDVID^FS^CI27
^FT440,570^A0B,25,25^FH\^CI28^FDVGRADE^FS^CI27
^FT643,182^A0B,25,25^FH\^CI28^FDVLENGTH^FS^CI27
^FT643,310^A0B,25,25^FH\^CI28^FDLENGTH^FS^CI27
^FT718,182^A0B,25,25^FH\^CI28^FDVTIME^FS^CI27
^FT683,182^A0B,25,25^FH\^CI28^FDVDATE^FS^CI27
^FT718,310^A0B,25,25^FH\^CI28^FDTIME^FS^CI27
^FT683,310^A0B,25,25^FH\^CI28^FDDATE^FS^CI27
^FT365,570^A0B,25,25^FH\^CI25^FDVSECTION^FS^CI27
^FT411,570^A0B,25,25^FH\^CI28^FDVGRADE^FS^CI27
^FT339,570^A0B,25,25^FH\^CI28^FDSECTION^FS^CI27
^FT295,569^A0B,34,33^FH\^CI28^FDVHEAT^FS^CI27
^FT260,343^A0B,14,15^FH\^CI28^FDVTOP^FS^CI27
^FT345,340^A0B,14,15^FH\^CI28^FDVBOTTOM^FS^CI27
^FT262,570^A0B,34,33^FH\^CI28^FDHEAT NO.^FS^CI27
^FO536,1^GB0,570,3^FS
^FT718,199^A0B,25,25^FH\^CI28^FD:^FS^CI27
^FT683,199^A0B,25,25^FH\^CI28^FD:^FS^CI27
^FT643,199^A0B,25,25^FH\^CI28^FD:^FS^CI27
^FO536,1^GB0,570,3^FS
^FH\^FDVDATA^FS
^FT570,580^BQN,2,4
^FH\^FDMA,VQR^FS
^FT245,275^BQN,2,5
^FO266,261^GFA,237,664,8,:Z64:eJzF0rENxCAMBVCjFJSMwChZDekGuJW4TRiB61xE4r4NRlFOSOni5lUYgz/RpbaGykS7eBA1LXLdtHWz75bAcijWUERfY9YmHFODjuNr6EXiPajHNJpp4RvX8A1X5xOe0sXYIp5SRDylnvWjjy/uCTGVGDDlWUyVxE2WAAlLUfVb/0Wf9r3h6rzYzDqNamH91vYZYq9j37Z/y4Plw/Ji+bE8zXxZ3i71A0iyOGc=:101F
^FO18,300^GFA,289,424,8,:Z64:eJxlULsNwjAQvZMjRTREdBRIWcEFBV0oGCRjpEsyAStlBDYgRQZISWHlceezEQg3Tyff+9wj+ns+YZcwGPBm6GBYYqBCsMIUscbMutpgjdjjFfeAoDPjuSk63GF0ETgJfSuiVR32xA+hv2oiwX5tBjoKfe5VnzFhltlhFP2z0J3ot0IvRX8RegUNKHSofp7jf7STfabED6aneVSfr+bHrfnrHPPcLB93lncky89k95TB7tt5u/fg0/1L6uNi/bC3vvi7v9xn7vfTd+7/570B8NB1AQ==:5907
^FO18,33^GFA,997,2112,8,:Z64:eJyNlb1OG0EUhe/Y0a4jIbNWmlnZeF6BkgIJHiOp7CppXbrLFpZFkSKlkZBoojxC6hVCeQZ3pkIUFK6iLSyT+3NmvYZIzkjwscv83HvuubNEu5GW/MsTHT8IPYW18Fz+JAqPlHUazP8QKX+BPyhUwmejLhJmIFGv0+TAeDG3+V+N7nqu09v5TJn4inQHPwWppq4L25r6Pmzw/3FNavK2snMi8b5Vk8/m923s377l/Que97KhBOc66ERjcGrnu7WdnxaiQ0GBQNEz518VKOPTnfF6ZnHkFk8rX9dx1OxYHEL3YqQl8o860I7/NTqgxw/HT9DRVZZHQJ7DS2MCtkvMn9o82+MU+50hGI+6fIcvkHek3/ktlVfsB/Wb+Kvpt/wR/op++4n9ot+Wr/y2y2sgHJWIc62UursHqbvXOrbYn1KvFhcp5WkJJz0USt6VUUYKah3I/LCnf01MDOfg51LZ/1KA0OERnFEP8SozkA7XNYWU1q+cr+TJPtYQ2E+m1wLnFHgu0cdb9HXMY4I6nJKTSAM2r7nTtUl307I+W22pKxxtVEf6SOYjZqr9tlVdXTDS4Gm/PzP0d4p+f39Dfa3vN9T5qtbX1p0a6fKNLgf7Ad4forfC1PLprU0XkvwvttQrJG7TJ94H0Q+Jt03S7ZHq1l516VjXz6H/3Hyc39t5/gm6TfBMVgfxmZ7fsvWjZ92PRha8W5j+zpd2PyAOoSOLSyj3gnGjlH1PFlT7UvRXrsy/UqcMddojr5G4XX2vRn+cQXfon8JXx9FfS9xzFfUlj+7Y8pRLA3nuUXTe89/Y1vtLCDy1eWIeuVLYBy76oCT1dZCp3L9dXTK2+0TGBIxePTQ6Df4jzr3vjPQRNfmA7xnzSp7v7DlPQNxPuN9f9w8NvM3j/YXiN6HUVfZLat+9U6abI1vHPsuazGc45x7nxPPWb76Htu63fb+YmdTrhO/Fpc13kgf77UNhcQZ5znbeb46/PgoGsg==:4315
^FO74,20^GFA,877,2640,10,:Z64:eJzF1b2O00AQAOBZuQjFKWkpToTHSAEyjxJ0Ba2piAQ5O1xBBy1dXgTB5iydG8t5AYpFV1yDkOmCZLLszO6Mk1ycCxScC+dTsj+zs7sTgP/1JH+t1BoApffIWr0h6rGhvA6jqEJUjUSnYXj19ZGozyOXPVbVy4KWkfvuFap461r3UDn2G2YhKjiHILViYTMvbAbfKGZgvRGFtSWyypFIItUnoj5+PGYpbLfqc1RrVEx5wb5udBezU29NeXGxDdecofQ9yTU573OkE1Eqin+wBmuOWVkWpIY1wHkXGUdFoki9QGR2lLFmOsheGNYH/rWYS4/PnKHFSlSzcs2zFTLvJUj+cB1adsbwD7Hoac2ayFFMDormPUJHPErvEWW3UxV9mH+TyoKuIj6JhZz7YiA5XW9mvDtSVELjkcaG72Wn2nbbOvbu72ZD867QPQqxPJSRW41uy9Z75AZOjAknNuSgQ37eDc0kPneeQ5V6J1VqzhVp8Um0Em2d57AfMlsO3E4f1pD1PQ5hqeY1R9Ukl6GGleMbSGJNspCkXo0Tzlae1e60k164ddSo5Y1bhxljVE5TylCeZxBr3q1Ww8zAfTxdtzGSE7uloVdPNZQr1+6CchU5WZ8rA+o35oqq3i/MFVVCyhVVR1U+X0FKFbM8M6GKVi5DvrJirvx3mKsUL2JbWf3/h2UB7QIJy5Q/45EIWkkPuXnHaHcU1BXnSpUd8rMtRYW+raX0qBTf1XKPqoijL2wm4lEWDUvDXfUAxtAhje+7NRCd4ptGfiIaWa5/L1H0XM85vusHkxBofpKwPj5z2+buVjTDghl/0RBlNcmdJlrs9KcTVQisL63UgqujWwbXSStVb3owG/f+/AHLeT5B:2249
^PQ1,0,1,Y
^XZ`

	zpl := template
	zpl = strings.ReplaceAll(zpl, "VHEAT", safeString(label.HeatNo))
	zpl = strings.ReplaceAll(zpl, "VSECTION", safeString(label.Section))
	zpl = strings.ReplaceAll(zpl, "VGRADE", safeString(label.Grade))
	zpl = strings.ReplaceAll(zpl, "VID", safeString(label.BundleNo))
	zpl = strings.ReplaceAll(zpl, "VMILL", safeString(label.Mill))
	zpl = strings.ReplaceAll(zpl, "VLENGTH", fmt.Sprintf("%d", label.Length))
	zpl = strings.ReplaceAll(zpl, "VTIME", safeString(label.Time))
	zpl = strings.ReplaceAll(zpl, "VDATE", safeString(label.Date))
	zpl = strings.ReplaceAll(zpl, "VTOP", safeString(label.IsiTop))
	zpl = strings.ReplaceAll(zpl, "VBOTTOM", safeString(label.IsiBottom))
	zpl = strings.ReplaceAll(zpl, "VQR", generateQRData(label))
	zpl = strings.ReplaceAll(zpl, "VDATA", generateLowerQRData(label))

	return zpl
}

// generateQRData builds the QR string from Label
func generateLowerQRData(label models.Label) string {
	return fmt.Sprintf(
		"UNIT:%s;MILL:%s;HEAT:%s;SECTION:%s;GRADE:%s;ID:%s;LENGTH:%d;WEIGHT:%s;LOCATION:%s;PQD:%s;DATE:%s;TIME:%s;",
		"SAIL-BSP",                      // Static unit name
		safeString(label.Mill),          // Mill from label
		safeString(label.HeatNo),        // Heat number
		safeString(label.Section),       // Section
		safeString(label.Grade),         // Grade
		safeString(label.BundleNo),      // ID
		label.Length,                    // Length (numeric)
		safeStringPtr(label.Weight),   // Weight
		safeStringPtr(label.Location),      // Location
		safeString(label.PQD),           // PQD
		safeString(label.Date),          // Date
		safeString(label.Time),          // Time
	)
}


// generateQRData builds the QR URL from Label
func generateQRData(label models.Label) string {
	return fmt.Sprintf(
		"https://madeinindia.qcin.org/product-details/%s/%s_%s_%s",
		safeString(label.ID.String()), // Static UUID (change if dynamic later)
		safeString(label.Mill),             // Mill
		safeString(label.HeatNo),           // Heat number
		safeString(label.PQD),              // PQD
	)
}



// GenerateAndSaveZPL saves the generated ZPL to a file
func GenerateAndSaveZPL(label models.Label) (string, error) {
	zpl := GenerateLabelZPL(label)
	filename := fmt.Sprintf("label_%s_%d.zpl", label.LabelID, time.Now().Unix())
	path := filepath.Join("printers", "zpl", filename)

	if err := os.WriteFile(path, []byte(zpl), 0644); err != nil {
		return "", fmt.Errorf("failed to save ZPL file: %w", err)
	}
	return path, nil
}
