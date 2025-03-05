package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "search",
	Short: "searchcli is a tool that makes searching like a search engine in the cli extremely easy",
	Long:  "search cli is a basic tool that uses duck duck go bangs to create a seamless search engine experience without all the overhead",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a search query.")
			return
		}

		query := strings.Join(args, " ")

		url := "https://www.google.com/search?q=" + query

		var err error
		switch runtime.GOOS {
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		default:
			log.Fatalf("Unsupported platform %s", runtime.GOOS)
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error occurred! '%s'\n", err)
		os.Exit(1)
	}
}
