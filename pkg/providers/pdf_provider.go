package providers

import (
	"bytes"
	"okieoth/schemaguesser/pkg/create"

	"github.com/ledongthuc/pdf"
)

func PdfToText(path string) create.TextProviderFunc {
	return func() (string, error) {
		f, r, err := pdf.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		var buf bytes.Buffer
		b, err := r.GetPlainText()
		if err != nil {
			panic(err)
		}
		buf.ReadFrom(b)
		content := buf.String()

		return content, nil
	}
}
