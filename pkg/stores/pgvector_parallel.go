package stores

import (
	"context"
	"fmt"
	"okieoth/schemaguesser/pkg/create"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

// type ParallelStoreImpl interface {
// 	Init(embedder *embeddings.EmbedderImpl, in <-chan schema.Document) (chan StoreResult, error)
// 	Run()
// }

type PgVectorParallelStore struct {
	bucketSize int
	pgurl      string
	in         <-chan schema.Document
	out        chan create.StoreResult
	store      pgvector.Store
}

func NewPgVectorParallelStore(pgurl string, bucketSize int) *PgVectorParallelStore {
	return &PgVectorParallelStore{
		bucketSize: bucketSize,
		pgurl:      pgurl,
	}
}

func (s *PgVectorParallelStore) Init(
	embedder *embeddings.EmbedderImpl,
	in <-chan schema.Document,
) (chan create.StoreResult, error) {
	s.in = in
	s.out = make(chan create.StoreResult)

	ctx := context.Background()
	store, err := pgvector.New(
		ctx,
		pgvector.WithConnectionURL(s.pgurl),
		pgvector.WithEmbedder(embedder),
	)
	if err != nil {
		return s.out, fmt.Errorf("error while creating pgvector store: %v\n", err)
	}
	s.store = store

	return s.out, nil
}

func (s *PgVectorParallelStore) Run() {
	defer close(s.out)
	docBucket := make([]schema.Document, 0, s.bucketSize)
	docCount := 0
	for doc := range s.in {
		docCount++
		docBucket = append(docBucket, doc)
		if len(docBucket) == s.bucketSize {
			err := stroreImpl(&s.store, docBucket)
			docBucket = docBucket[:0]
			if err != nil {
				s.out <- create.StoreResult{
					StoreCount: docCount,
					Err:        &err,
				}
				return
			}
		}
	}
	if len(docBucket) > 0 {
		stroreImpl(&s.store, docBucket)
	}
	s.out <- create.StoreResult{
		StoreCount: docCount,
		Err:        nil,
	}
}
