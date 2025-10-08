package chunkfactories

import (
	"errors"
	"okieoth/schemaguesser/pkg/create"
	"strings"

	"github.com/tmc/langchaingo/schema"
)

type handlePreparedChunkFunc func(chunkData *strings.Builder)

// this function is separated to implement the business logic of
// the chunking only once
func doWordChunking(
	txt string, // text to be chunked
	chunkTxtBuilder *strings.Builder, // temp storage of chunk content
	wordCount *int, // current words in the chunk
	chunkSize, // considered words per chunk
	overlap int, // word overlap to create the chunks
	handlePreparedChunk handlePreparedChunkFunc, // callback to handle the filled chunk
) {
	words := strings.Split(txt, " ")
	currentIndex := 0
	for {
		if currentIndex == len(words) {
			break
		}
		w := words[currentIndex]
		trimmedWord := strings.TrimSpace(w)
		currentIndex++
		if len(trimmedWord) == 0 {
			continue
		}
		if chunkTxtBuilder.Len() > 0 {
			chunkTxtBuilder.Write([]byte(" "))
		}
		chunkTxtBuilder.Write([]byte(trimmedWord))
		*wordCount++
		if *wordCount == chunkSize {
			handlePreparedChunk(chunkTxtBuilder)
			chunkTxtBuilder.Reset()
			*wordCount = 0
			currentIndex = currentIndex - overlap
		}
	}

}

// Creates a function that creates chunks based on a set of words.
// In case you want to have meta data, write a similar implementation and
// include the meta data in the schema.Document
func WordChunkFactoryFunc(chunkSize, overlap int) create.ChunkFactory {
	return func(txt string) ([]schema.Document, error) {
		chunksToReturn := make([]schema.Document, 0)
		if chunkSize < 1 || overlap < 0 || overlap >= chunkSize {
			return chunksToReturn, errors.New("not allowed values for 'chunkSize' or 'overlap'")
		}
		var chunkTxtBuilder strings.Builder
		wordCount := 0
		handlePreparedChunk := func(chunkData *strings.Builder) {
			chunksToReturn = append(chunksToReturn, schema.Document{
				PageContent: chunkTxtBuilder.String(),
			})
		}
		doWordChunking(txt, &chunkTxtBuilder, &wordCount, chunkSize, overlap, handlePreparedChunk)
		if chunkTxtBuilder.Len() > 0 {
			handlePreparedChunk(&chunkTxtBuilder)
		}
		return chunksToReturn, nil // TODO
	}
}
