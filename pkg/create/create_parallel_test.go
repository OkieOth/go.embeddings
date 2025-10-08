package create_test

import (
	"okieoth/schemaguesser/internal/pkg/util"
	"okieoth/schemaguesser/pkg/chunkfactories"
	"okieoth/schemaguesser/pkg/create"
	"okieoth/schemaguesser/pkg/providers"
	"okieoth/schemaguesser/pkg/stores"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParallelEmbeddingsFactory(t *testing.T) {
	pgURL := util.GetStrVar("PGURL", PGURL_DEFAULT)
	inputPdf := "../../resources/examples/pdf/wilhelm_busch.pdf"
	textProvider := providers.PdfToTextIter(inputPdf)
	chunkFactory, err := chunkfactories.NewParallelWordChunkFactory(10, 2)
	require.Nil(t, err, "error while creating parallel chunk factory")
	embeddingModel := "all-minilm:22m"
	storeImpl := stores.NewPgVectorParallelStore(pgURL, 5)
	embeddingsFactory := create.NewParallelEmbeddingsFactory(textProvider, chunkFactory, embeddingModel, storeImpl)
	embeddingsCount, err := embeddingsFactory.Run()
	require.Nil(t, err, "error in embeddingsFactory.Run call:", err)
	require.Equal(t, 25, embeddingsCount, "wrong number of created embeddings")
}
