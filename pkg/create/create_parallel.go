package create

import (
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

type DocTransport struct {
	Docs []schema.Document
	Err  *error
}

type StoreResult struct {
	StoreCount int
	Err        *error
}

type TextProviderIter func(yield func(string) bool)
type ParallelChunkFactory interface {
	Init(in chan<- string) (<-chan DocTransport, error)
	Run()
}

type ParallelStoreImpl interface {
	Init(embedder *embeddings.EmbedderImpl, in <-chan []schema.Document) (chan StoreResult, error)
	Run()
}

type ParallelEmbeddingsFactory struct {
	textProviderIter TextProviderIter
	chunkFactory     ParallelChunkFactory
	embeddingModel   string
	storeImpl        ParallelStoreImpl
}

func NewParallelEmbeddingsFactory(
	textProvider TextProviderIter,
	chunkFactory ParallelChunkFactory,
	embeddingModel string,
	storeImpl ParallelStoreImpl,
) *ParallelEmbeddingsFactory {
	return &ParallelEmbeddingsFactory{
		textProviderIter: textProvider,
		chunkFactory:     chunkFactory,
		embeddingModel:   embeddingModel,
		storeImpl:        storeImpl,
	}
}

func (e *ParallelEmbeddingsFactory) run() (int, error) {
	sendTextChan := make(chan string)
	sendDocsChan := make(chan []schema.Document)
	rcvChunksChan, err := e.chunkFactory.Init(sendTextChan)
	if err != nil {
		return 0, fmt.Errorf("can't init chunkFactory: %v", err)
	}

	llm, err := ollama.New(ollama.WithModel(e.embeddingModel))
	if err != nil {
		return 0, fmt.Errorf("error while creating ollama client: %v", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return 0, fmt.Errorf("error while creating new embedder: %v", err)
	}

	doneChan, err := e.storeImpl.Init(embedder, sendDocsChan)
	if err != nil {
		return 0, fmt.Errorf("can't init store implementation: %v", err)
	}
	go e.chunkFactory.Run()
	go e.storeImpl.Run()

	go func() {
		for c := range rcvChunksChan {
			if c.Err != nil {
				// TODO handle error
			} else {
				sendDocsChan <- c.Docs
			}
		}
	}()
	for t := range e.textProviderIter {
		sendTextChan <- t
	}
	close(sendTextChan)
	storeResult := <-doneChan
	return storeResult.StoreCount, nil // TODO error handling
}
