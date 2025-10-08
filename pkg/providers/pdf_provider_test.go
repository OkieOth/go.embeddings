package providers_test

import (
	"okieoth/schemaguesser/pkg/providers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPdfToText(t *testing.T) {
	tests := []struct {
		input       string
		expectedLen int
	}{
		{
			input:       "../../resources/examples/pdf/gr07-txt.pdf",
			expectedLen: 1025737,
		},
		{
			input:       "../../resources/examples/pdf/schachnovelle.pdf",
			expectedLen: 158918,
		},
	}
	for _, test := range tests {
		providerFunc := providers.PdfToText(test.input)
		txt, err := providerFunc()
		require.Nil(t, err, "error while converting pdf to txt", err)
		require.Equal(t, test.expectedLen, len(txt), "wrong text len size")

	}
}
