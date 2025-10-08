package sub

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetCommandChain(cmd *cobra.Command) []*cobra.Command {
	chain := make([]*cobra.Command, 0)
	for c := cmd; c != nil; c = c.Parent() {
		chain = append(chain, c)
	}
	return chain
}

func IsFlagRequired(f *pflag.Flag) bool {
	if f.Annotations != nil {
		if v, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok && len(v) > 0 && v[0] == "true" {
			return true
		}
	}
	return false
}

func SetFlagsForCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		// Iterate over all flags of the command
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			// Ask interactively for input
			if f.Name == "help" {
				return
			}
			defValue := "empty"
			if f.DefValue != "" {
				defValue = f.DefValue
			}
			skipTxt := "imput required!"
			flagRequired := IsFlagRequired(f)
			if !flagRequired {
				skipTxt = "press ⏎ to skip"
			}
			collectFlagInput(cmd, f, flagRequired, defValue, skipTxt)
		})
	}
}

func collectFlagInput(cmd *cobra.Command, f *pflag.Flag, flagRequired bool, defValue, skipTxt string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n%s (default: %s), %s: ", f.Usage, defValue, skipTxt)
		input, _ := reader.ReadString('\n') // read entire line
		input = strings.TrimSpace(input)    // remove newline and spaces

		if input != "" {
			// User provided a value -> set it
			if err := cmd.Flags().Set(f.Name, input); err != nil {
				fmt.Printf("⚠️  Could not set flag %s: %v\n", f.Name, err)
			} else {
				break
			}
		} else if flagRequired {
			fmt.Printf("⚠️  Flag %s is required, so input is needed!\n", f.Name)
		} else {
			break
		}
	}
}

func GetOptionsFromCommands(cmds ...*cobra.Command) []string {
	ret := make([]string, 0)
	for _, cmd := range cmds {
		if cmd.Use == "completion" {
			continue
		}
		ret = append(ret, cmd.Use)
	}
	return ret
}

func RunInteractive(cmd *cobra.Command) (*cobra.Command, *cobra.Command, error) {
	subCommands := cmd.Commands()
	if len(subCommands) == 0 {
		return nil, nil, fmt.Errorf("command has no sub commads")
	}
	currentCmd := cmd
	for {
		options := GetOptionsFromCommands(subCommands...)
		idx, err := fuzzyfinder.Find(
			options,
			func(i int) string { return options[i] },
		)
		if err != nil {
			return nil, nil, fmt.Errorf("error in interactive run: %v", err)
		}
		selected := options[idx]
		if strings.HasPrefix(selected, "help ") {
			helpCmd, _, err := currentCmd.Find([]string{"help"})
			return helpCmd, helpCmd, err
		}
		nextCmd, _, err := currentCmd.Find([]string{selected})
		if err != nil {
			return nil, nil, fmt.Errorf("error finding seleted sub command: %v", err)
		}
		subCommands = nextCmd.Commands()
		if len(subCommands) == 0 {
			// reached end of the chain ..
			cmdChain := GetCommandChain(nextCmd)
			possibleCompletionIndex := len(cmdChain) - 2
			if cmdChain[possibleCompletionIndex].Name() != "completion" {
				SetFlagsForCommands(cmdChain...)
				cmd = nextCmd
			} else {
				// completion command was selected
				nextCmd = cmdChain[possibleCompletionIndex]
			}
			return nextCmd, cmd, nil
		}
		currentCmd = nextCmd
	}
}
