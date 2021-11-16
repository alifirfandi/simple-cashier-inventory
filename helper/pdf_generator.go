package helper

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

const InvoiceTemplate = `
<html lang="en">

<head>
  <meta charset="UTF-8" />
  <title>Example Pdf</title>
  <style>
    .table-list {
      margin-top: 25px;
      width: 100%;
    }

    .table-list>thead {
      color: #fff;
      background-color: #111;
    }

    .table-list>thead>th {
      text-align: center;
    }

    .table-list>tbody td {
      border-bottom: 1px solid #dddddd;
      text-align: right;
      padding: 8px;
    }

    .table-list>tbody td:first-child {
      text-align: left;
    }

    .grand-total {
      text-align: right;
    }

    .grand-total-label {
      font-size: 18px;
      font-weight: 600;
    }
  </style>
</head>

<body>
  <table>
    <tbody>
      <tr>
        <td>Invoice</td>
        <td>:</td>
        <td>{{.Invoice}}</td>
      </tr>
      <tr>
        <td>Name</td>
        <td>:</td>
        <td>{{.Name}}</td>
      </tr>
      <tr>
        <td>Date</td>
        <td>:</td>
        <td>{{.Date}}</td>
      </tr>
    </tbody>
  </table>
  <table class="table-list">
    <thead>
      <tr>
        <th>PRODUCT INFO</th>
        <th>QTY</th>
        <th>UNIT PRICE</th>
        <th>TOTAL PRICE</th>
      </tr>
    </thead>
    <tbody>
	  {{range .Items}}
      <tr>
        <td>{{.Name}}</td>
        <td>Rp. {{.Price}}</td>
        <td>{{.Qty}}</td>
        <td>Rp. {{.SubTotal}}</td>
      </tr>
	  {{end}}
    </tbody>
  </table>
  <p class="grand-total grand-total-label">Grand Total</p>
  <h2 class="grand-total grand-total-value">Rp. {{.GrandTotal}}</h2>
</body>

</html>
`

type InvoiceItemData struct {
	Name     string
	Price    int
	Qty      int
	SubTotal int
}

type InvoiceData struct {
	Invoice    string
	Name       string
	Date       string
	GrandTotal int
	Items      []InvoiceItemData
}

// InvoiceExporter
// parsing template function
// Ref: Mindinventory / Golang-HTML-TO-PDF-Converter
// https://github.com/Mindinventory/Golang-HTML-TO-PDF-Converter
func InvoiceExporter(templateStr, outputPath, fileName string, data InvoiceData) error {
	t, err := template.New("invoice").Parse(templateStr)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	return GeneratePDF(buf.String(), fmt.Sprintf("%s/%s", outputPath, fileName))
}

// GeneratePDF
// generate pdf function
// Ref: Mindinventory / Golang-HTML-TO-PDF-Converter
// https://github.com/Mindinventory/Golang-HTML-TO-PDF-Converter
func GeneratePDF(text, pdfPath string) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(text)))
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return err
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		return err
	}

	return nil
}
