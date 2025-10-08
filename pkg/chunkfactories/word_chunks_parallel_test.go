package chunkfactories_test

import (
	"okieoth/schemaguesser/pkg/chunkfactories"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

func TestParallelWordChunkFactory(t *testing.T) {
	tests := []struct {
		input     []string
		chunkSize int
		overlap   int
		chunks    []schema.Document
		hasError  bool
	}{
		{
			input:     []string{"xxx"},
			chunkSize: -1,
			overlap:   10,
			hasError:  true,
		},
		{
			input:     []string{"xxx"},
			chunkSize: 0,
			overlap:   10,
			hasError:  true,
		},
		{
			input:     []string{"xxx"},
			chunkSize: 1,
			overlap:   10,
			hasError:  true,
		},
		{
			input:     []string{"xxx"},
			chunkSize: 1,
			overlap:   0,
			hasError:  false,
			chunks: []schema.Document{schema.Document{
				PageContent: "xxx",
			}},
		},
		{
			input:     []string{"  xxx yyyy", "     z"},
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
		factory, err := chunkfactories.NewParallelWordChunkFactory(test.chunkSize, test.overlap)
		if test.hasError {
			require.NotNil(t, err, "although error was expected ParallelWordChunkFactory was created, test #", i)
		} else {
			require.Nil(t, err, "retrieve unexpected error for test #", i)
			txtInputChan := make(chan string)
			outputChan, err := factory.Init(txtInputChan)
			require.Nil(t, err)
			go factory.Run()
			go func() {
				for _, text := range test.input {
					txtInputChan <- text
				}
				close(txtInputChan)
			}()
			retrieved := make([]schema.Document, 0)
			for output := range outputChan {
				retrieved = append(retrieved, output)
			}
			require.Equal(t, test.chunks, retrieved, "retrieved chunks are not equal in test #", i)
		}
	}
}
