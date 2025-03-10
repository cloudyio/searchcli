package main

import (
	"searchcli/cmd"
)

func main() {
	//cmd.SetEmbeddedBangs(nil) // You can remove this line if embedding works fine
	cmd.Execute()
}
