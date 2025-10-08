package providers

import (
	"bytes"
	"io"
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

func PdfToTextIter(path string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		f, r, err := pdf.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		reader, err := r.GetPlainText()
		if err != nil {
			panic(err)
		}

		buf := make([]byte, 4096) // 4 KB buffer (tweak as needed)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				// Pass the chunk to yield; stop if yield returns false
				if !yield(string(buf[:n])) {
					return
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
		}
	}
}
