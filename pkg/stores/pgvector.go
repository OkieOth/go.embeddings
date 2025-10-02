package stores

import (
	"context"
	"fmt"
	"okieoth/schemaguesser/pkg/create"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func stroreImpl(store *pgvector.Store, chunks []schema.Document) error {
	_, errAd := store.AddDocuments(context.Background(), chunks)
	if errAd != nil {
		return errAd
	}
	// TODO information channel
	return nil
}

func StoreInPgVectorFunc(pgURL string, chunkSize int) create.StoreImpl {
	return func(docs []schema.Document, embedder *embeddings.EmbedderImpl) (int, error) {
		ctx := context.Background()
		store, err := pgvector.New(
			ctx,
			pgvector.WithConnectionURL(pgURL),
			pgvector.WithEmbedder(embedder),
		)
		if err != nil {
			return 0, fmt.Errorf("error while creating pgvector store: %v\n", err)
		}

		chunks := make([]schema.Document, 0)
		for i, d := range docs {
			chunks = append(chunks, d)
			if len(chunks) == chunkSize {
				err := stroreImpl(&store, chunks)
				if err != nil {
					return i, fmt.Errorf("error in AddDocument (already added: %d): %v\n", i, err)
				}
				chunks = chunks[:0]
			}
		}
		docCount := len(docs)
		if len(chunks) > 0 {
			err := stroreImpl(&store, chunks)
			if err != nil {
				return docCount, fmt.Errorf("error in AddDocument (already added: %d): %v\n", docCount, err)
			}
		}
		return docCount, nil
	}

}
