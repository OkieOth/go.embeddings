package chunkfactories

import (
	"errors"
	"strings"

	"github.com/tmc/langchaingo/schema"
)

type ParallelWordChunkFactory struct {
	chunkSize int
	overlap   int
	in        <-chan string
	out       chan schema.Document
}

func NewParallelWordChunkFactory(chunkSize, overlap int) (*ParallelWordChunkFactory, error) {

	var err error
	if chunkSize < 1 || overlap < 0 || overlap >= chunkSize {
		err = errors.New("not allowed values for 'chunkSize' or 'overlap'")
	}

	return &ParallelWordChunkFactory{
		chunkSize: chunkSize,
		overlap:   overlap,
	}, err
}

func (f *ParallelWordChunkFactory) Init(in <-chan string) (<-chan schema.Document, error) {
	f.in = in
	f.out = make(chan schema.Document)
	return f.out, nil
}

func (f *ParallelWordChunkFactory) Run() {
	defer close(f.out)
	var chunkTxtBuilder strings.Builder
	wordCount := 0

	handlePreparedChunk := func(chunkData *strings.Builder) {
		f.out <- schema.Document{
			PageContent: chunkTxtBuilder.String(),
		}
	}

	for txt := range f.in {
		doWordChunking(txt, &chunkTxtBuilder, &wordCount, f.chunkSize, f.overlap, handlePreparedChunk)
	}
	if chunkTxtBuilder.Len() > 0 {
		handlePreparedChunk(&chunkTxtBuilder)
	}
}
