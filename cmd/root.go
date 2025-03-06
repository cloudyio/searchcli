package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Bang struct {
	C  string `json:"c"`
	D  string `json:"d"`
	R  int    `json:"r"`
	S  string `json:"s"`
	Sc string `json:"sc"`
	T  string `json:"t"`
	U  string `json:"u"`
}

type Config struct {
	DefaultBang string `json:"default_bang"`
}

var configFilePath = filepath.Join(os.Getenv("HOME"), ".searchcli", "config.json")
var defaultBang string

func loadConfig() Config {
	var config Config
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}
	}
	return config
}

func saveConfig(config Config) {
	dir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}
}

func loadBangs(filename string) ([]Bang, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var bangs []Bang
	if err := json.Unmarshal(bytes, &bangs); err != nil {
		return nil, err
	}

	return bangs, nil
}

func findBang(shortcut string, bangs []Bang) *Bang {
	for _, bang := range bangs {
		if bang.T == shortcut {
			return &bang
		}
	}
	return nil
}

func openURL(url string) {
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
}

var rootCmd = &cobra.Command{
	Use:                "search",
	Short:              "searchcli is a tool that makes searching like a search engine in the CLI extremely easy",
	Long:               "search CLI is a basic tool that uses DuckDuckGo bangs to create a seamless search engine experience.",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a search query.")
			return
		}

		config := loadConfig()
		bangs, err := loadBangs("bangs.json")
		if err != nil {
			log.Fatalf("Failed to load bangs: %v", err)
		}

		query := strings.Join(args, " ")
		parts := strings.SplitN(query, " ", 2)

		var shortcut string
		var searchTerm string

		if len(parts) > 1 && strings.HasPrefix(parts[0], "!") {
			shortcut = strings.TrimPrefix(parts[0], "!")
			searchTerm = parts[1]
		} else if config.DefaultBang != "" {
			shortcut = config.DefaultBang
			searchTerm = query
		} else {
			url := "https://www.google.com/search?q=" + query
			openURL(url)
			return
		}

		bang := findBang(shortcut, bangs)
		if bang != nil {
			searchURL := strings.Replace(bang.U, "{{{s", searchTerm, 1)
			openURL(searchURL)
		} else {
			fmt.Printf("Unknown bang: !%s, default bang can be reset with search set-default g\n", shortcut)
		}
	},
}

var setDefaultCmd = &cobra.Command{
	Use:   "set-default [bang]",
	Short: "Set a default bang for searches",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()
		config.DefaultBang = args[0]
		saveConfig(config)
		fmt.Printf("Default bang set to !%s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(setDefaultCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Oops. An error occurred! '%s'\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
