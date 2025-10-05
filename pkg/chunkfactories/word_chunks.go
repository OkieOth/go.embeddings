package chunkfactories

import (
	"errors"
	"okieoth/schemaguesser/pkg/create"
	"strings"

	"github.com/tmc/langchaingo/schema"
)

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
		words := strings.Split(txt, " ")
		wordCount := 0
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
			wordCount++
			if wordCount == chunkSize {
				chunksToReturn = append(chunksToReturn, schema.Document{
					PageContent: chunkTxtBuilder.String(),
				})
				chunkTxtBuilder.Reset()
				wordCount = 0
				currentIndex = currentIndex - overlap
			}
		}
		if chunkTxtBuilder.Len() > 0 {
			chunksToReturn = append(chunksToReturn, schema.Document{
				PageContent: chunkTxtBuilder.String(),
			})
		}
		return chunksToReturn, nil // TODO
	}
}
