package chunkfactories

import (
	"okieoth/schemaguesser/pkg/create"

	"github.com/tmc/langchaingo/schema"
)

func WordChunkFactoryFunc(chunkSize, overlap int) create.ChunkFactory {
	return func(txt string) ([]schema.Document, error) {
		// TODO
		return []schema.Document{}, nil // TODO
	}
}
