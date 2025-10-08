package create

import (
	"fmt"
	"log/slog"
	"okieoth/goembeddings/internal/pkg/util"

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
	logMsg := func(msg string) string {
		const LOG_PREFIX = "EmbeddingsFactory: "
		return LOG_PREFIX + msg
	}
	text, err := e.textProvider()
	if err != nil {
		slog.Error(logMsg("error while getting text"), "error", err)
		return 0, fmt.Errorf("error while retrieving input for the embeddings: %v", err)
	}
	slog.Info(logMsg("retrieved text"), "len", len(text))
	chunks, err := e.chunkFactory(text)
	if err != nil {
		slog.Error("EmbeddingsFactory: error while creating chunk factory", "error", err)
		return 0, fmt.Errorf("error while creating chungs from intput: %v", err)
	}
	slog.Info(logMsg("retrieved chunks"), "len", len(chunks))

	serverUrl := util.GetStrVar("OLLAMA_URL", OLLAMA_URL_DEFAULT)

	llm, err := ollama.New(
		ollama.WithServerURL(serverUrl),
		ollama.WithModel(e.embeddingModel),
	)

	if err != nil {
		slog.Error(logMsg("error while creating ollama client"), "error", err)
		return 0, fmt.Errorf("error while creating ollama client: %v", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		slog.Error(logMsg("error while creating embedder"), "error", err)
		return 0, fmt.Errorf("error while creating new embedder: %v", err)
	}

	slog.Info(logMsg("start storing ..."))
	storeCount, err := e.storeImpl(chunks, embedder)
	if err != nil {
		slog.Error("EmbeddingsFactory: error while storing embeddings", "error", err)
		return storeCount, fmt.Errorf("error while storing embeddings: %v", err)
	}
	slog.Info(logMsg("storing done"), "count", storeCount)
	return storeCount, err
}
