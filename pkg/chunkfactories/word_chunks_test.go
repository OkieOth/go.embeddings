package chunkfactories_test

import (
	"okieoth/schemaguesser/pkg/chunkfactories"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

func TestWordChunkFactoryFunc(t *testing.T) {
	tests := []struct {
		input     string
		chunkSize int
		overlap   int
		chunks    []schema.Document
		hasError  bool
	}{
		{
			input:     "xxx",
			chunkSize: -1,
			overlap:   10,
			hasError:  true,
		},
		{
			input:     "xxx",
			chunkSize: 0,
			overlap:   10,
			hasError:  true,
		},
		{
			input:     "xxx",
			chunkSize: 1,
			overlap:   10,
			hasError:  true,
		},
		{
			input:     "xxx",
			chunkSize: 1,
			overlap:   0,
			hasError:  false,
			chunks: []schema.Document{schema.Document{
				PageContent: "xxx",
			}},
		},
		{
			input:     "  xxx yyyy     z",
			chunkSize: 1,
			overlap:   0,
			hasError:  false,
			chunks: []schema.Document{schema.Document{
				PageContent: "xxx",
			},
				schema.Document{
					PageContent: "yyyy",
				},
				schema.Document{
					PageContent: "z",
				},
			},
		},
	}
	for i, test := range tests {
		factoryFunc := chunkfactories.WordChunkFactoryFunc(test.chunkSize, test.overlap)
		chunks, err := factoryFunc(test.input)
		if test.hasError {
			require.NotNil(t, err, "didn't retrieve error for test #", i)
		} else {
			require.Nil(t, err, "retrieve unexpected error for test #", i)
			require.Equal(t, test.chunks, chunks, "retrieved chunks are not equal in test #", i)
		}
	}
}
