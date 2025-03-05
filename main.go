package main

import "searchcli/cmd"

type SearchResult struct {
	C  string `json:"c"`
	D  string `json:"d"`
	R  int    `json:"r"`
	S  string `json:"s"`
	Sc string `json:"sc"`
	T  string `json:"t"`
	U  string `json:"u"`
}

func main() {
	cmd.Execute()
}
