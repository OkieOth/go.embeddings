package create

import (
	"fmt"
	"log/slog"
	"okieoth/goembeddings/internal/pkg/util"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

const OLLAMA_URL_DEFAULT = "http://localhost:11434"

type StoreResult struct {
	StoreCount int
	Err        *error
}

type TextProviderIter func(yield func(string) bool)
type ParallelChunkFactory interface {
	Init(in <-chan string) (<-chan schema.Document, error)
	Run()
}

type ParallelStoreImpl interface {
	Init(embedder *embeddings.EmbedderImpl, in <-chan schema.Document) (chan StoreResult, error)
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

func (e *ParallelEmbeddingsFactory) Run() (int, error) {
	logMsg := func(msg string) string {
		const LOG_PREFIX = "ParallelEmbeddingsFactory: "
		return LOG_PREFIX + msg
	}

	sendTextChan := make(chan string)
	rcvChunksChan, err := e.chunkFactory.Init(sendTextChan)
	if err != nil {
		slog.Error(logMsg("can't init chunkFactory"), "error", err)
		return 0, fmt.Errorf("can't init chunkFactory: %v", err)
	}

	serverUrl := util.GetStrVar("OLLAMA_URL", OLLAMA_URL_DEFAULT)

	llm, err := ollama.New(
		ollama.WithServerURL(serverUrl),
		ollama.WithModel(e.embeddingModel),
	)
	if err != nil {
		slog.Error(logMsg("can't create ollama client"), "error", err)
		return 0, fmt.Errorf("error while creating ollama client: %v", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		slog.Error(logMsg("can't create embedder"), "error", err)
		return 0, fmt.Errorf("error while creating new embedder: %v", err)
	}

	doneChan, err := e.storeImpl.Init(embedder, rcvChunksChan)
	if err != nil {
		slog.Error(logMsg("can't init store"), "error", err)
		return 0, fmt.Errorf("can't init store implementation: %v", err)
	}
	go e.chunkFactory.Run()
	go e.storeImpl.Run()

	slog.Info(logMsg("start retrieving text ..."))
	for t := range e.textProviderIter {
		sendTextChan <- t
	}
	slog.Info(logMsg("text done."))
	close(sendTextChan)
	storeResult := <-doneChan
	slog.Info(logMsg("documents stored"), "count", storeResult)
	return storeResult.StoreCount, nil // TODO error handling
}
