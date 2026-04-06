package sh

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// AskYesNo prompts the user with the given question and returns true for
// "yes" and false for "no". It accepts "y", "yes", "n", "no" (any case).
// If the reader hits EOF, it returns the provided defaultVal.
func AskYesNo(r io.Reader, w io.Writer, question string, defaultVal bool) bool {
	reader := bufio.NewReader(r)

	for {
		defaultPrompt := "y/N"
		if defaultVal {
			defaultPrompt = "Y/n"
		}

		fmt.Fprintf(w, "%s [%s]: ", question, defaultPrompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			// on EOF or other read error, return default
			return defaultVal
		}

		response := strings.TrimSpace(strings.ToLower(line))
		if response == "" {
			return defaultVal
		}

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Fprintln(w, "Please answer yes or no (y/n).")
		}
	}
}
