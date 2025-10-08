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

const PGURL_DEFAULT = "postgres://admin:secretpassword@localhost:5432/testdb?sslmode=disable"

func TestEmbeddingsFactory(t *testing.T) {
	pgURL := util.GetStrVar("PGURL", PGURL_DEFAULT)
	inputPdf := "../../resources/examples/pdf/wilhelm_busch.pdf"
	textProvider := providers.PdfToText(inputPdf)
	chunkFactory := chunkfactories.WordChunkFactoryFunc(10, 2)
	embeddingModel := "all-minilm:22m"
	storeImpl := stores.StoreInPgVectorFunc(pgURL, 5)
	embeddingsFactory := create.NewEmbeddingsFactory(textProvider, chunkFactory, embeddingModel, storeImpl)
	embeddingsCount, err := embeddingsFactory.Run()
	require.Nil(t, err, "error in embeddingsFactory.Run call:", err)
	require.Equal(t, 25, embeddingsCount, "wrong number of created embeddings")
}
