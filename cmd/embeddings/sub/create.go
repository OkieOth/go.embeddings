package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var input string
var model string
var ollamaUrl string
var chunkSize int
var parallel bool

var wordCount int

var dbUrl string

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates embeddings",
	Long:  `Create embeddings and stores them in a database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command needs additional sub commands to do something useful")
	},
}

var WordsCmd = &cobra.Command{
	Use:   "word_embeddings",
	Short: "Creates word embeddings",
	Long:  `Creates word based embeddings and stores them in a database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command needs additional sub commands to do something useful")
	},
}

var SentencesCmd = &cobra.Command{
	Use:   "sentence_embeddings",
	Short: "Creates embeddings from sentences",
	Long:  `Creates sentence based embeddings and stores them in a database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command needs additional sub commands to do something useful")
	},
}

var FromPdfCmd = &cobra.Command{
	Use:   "from_pdf",
	Short: "Reads a PDF input",
	Long:  `Parses a PDF input for the texts of the embeddings`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command needs additional sub commands to do something useful")
	},
}

var FromTxtCmd = &cobra.Command{
	Use:   "from_txt",
	Short: "Reads a text input",
	Long:  `Parses a text input for the texts of the embeddings`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command needs additional sub commands to do something useful")
	},
}

var PgVectorCmd = &cobra.Command{
	Use:   "and_store_in_pgvector",
	Short: "Stores the embeddings in pgvector db",
	Long:  `Stores the created embeddings in a pgvector db`,
	Run: func(cmd *cobra.Command, args []string) {
		createEmbeddingsInPgvector(cmd)
	},
}

var ChromaCmd = &cobra.Command{
	Use:   "and_store_in_chromadb",
	Short: "Stores the embeddings in chromadb",
	Long:  `Stores the created embeddings in chromadb`,
	Run: func(cmd *cobra.Command, args []string) {
		createEmbeddingsInChromadb()
	},
}

func init() {
	CreateCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "Input source to parse.")
	CreateCmd.PersistentFlags().StringVarP(&model, "model", "m", "", "Model to use to create the embeddings.")
	CreateCmd.PersistentFlags().StringVarP(&ollamaUrl, "ollama", "o", "", "Ollama connection string to use.")
	CreateCmd.PersistentFlags().IntVarP(&chunkSize, "chunk_size", "c", 10, "Number of chunks inserted together in the db. Default is 10.")
	CreateCmd.PersistentFlags().BoolVarP(&parallel, "parallel", "p", false, "When this flag is set, the embeddings are created in a prarallel approach.")
	CreateCmd.PersistentFlags().StringVarP(&dbUrl, "db_url", "d", "", "Database connection string to use")
	CreateCmd.MarkPersistentFlagRequired("input")
	CreateCmd.MarkPersistentFlagRequired("model")
	CreateCmd.MarkPersistentFlagRequired("ollama")
	CreateCmd.MarkPersistentFlagRequired("db_url")

	WordsCmd.PersistentFlags().IntVar(&wordCount, "words", 20, "Number of word put together to one embedding. Default is 20.")

	FromPdfCmd.AddCommand(PgVectorCmd)
	FromPdfCmd.AddCommand(ChromaCmd)

	FromTxtCmd.AddCommand(PgVectorCmd)
	FromTxtCmd.AddCommand(ChromaCmd)

	WordsCmd.AddCommand(FromPdfCmd)
	WordsCmd.AddCommand(FromTxtCmd)

	SentencesCmd.AddCommand(FromPdfCmd)
	SentencesCmd.AddCommand(FromTxtCmd)

	CreateCmd.AddCommand(WordsCmd)
	CreateCmd.AddCommand(SentencesCmd)
}

func createEmbeddingsInPgvector(cmd *cobra.Command) {
	fmt.Println("TODO pgvector :)")
}

func createEmbeddingsInChromadb() {
	fmt.Println("TODO chromadb")
}
