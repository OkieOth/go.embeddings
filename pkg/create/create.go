package create

import (
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

type TextProviderFunc func() (string, error)
type ChunkFactory func(text string) ([]schema.Document, error)
type StoreImpl func(docs []schema.Document, embedder *embeddings.EmbedderImpl) (int, error)

type EmbeddingsFactory struct {
	textProvider   TextProviderFunc
	chunkFactory   ChunkFactory
	embeddingModel string
	storeImpl      StoreImpl
}

func NewEmbeddingsFactory(
	textProvider TextProviderFunc,
	chunkFactory ChunkFactory,
	embeddingModel string,
	storeImpl StoreImpl,
) *EmbeddingsFactory {
	return &EmbeddingsFactory{
		textProvider:   textProvider,
		chunkFactory:   chunkFactory,
		embeddingModel: embeddingModel,
		storeImpl:      storeImpl,
	}
}

func (e *EmbeddingsFactory) Run() (int, error) {
	text, err := e.textProvider()
	if err != nil {
		return 0, fmt.Errorf("error while retrieving input for the embeddings: %v", err)
	}
	chunks, err := e.chunkFactory(text)
	if err != nil {
		return 0, fmt.Errorf("error while creating chungs from intput: %v", err)
	}

	llm, err := ollama.New(ollama.WithModel(e.embeddingModel))
	if err != nil {
		return 0, fmt.Errorf("error while creating ollama client: %v", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return 0, fmt.Errorf("error while creating new embedder: %v", err)
	}

	return e.storeImpl(chunks, embedder)
}
