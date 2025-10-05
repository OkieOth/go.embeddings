package providers_test

import (
	"okieoth/schemaguesser/pkg/providers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPdfToText(t *testing.T) {
	path := "../../resources/examples/pdf/gr07-txt.pdf"
	expectedTxtLen := 1025737
	providerFunc := providers.PdfToText(path)
	txt, err := providerFunc()
	require.Nil(t, err, "error while converting pdf to txt", err)
	require.Equal(t, expectedTxtLen, len(txt), "wrong text len size")
}
