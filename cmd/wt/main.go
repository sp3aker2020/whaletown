// wt is the Whale Town CLI for managing multi-agent workspaces.
package main

import (
	"os"

	"github.com/speaker20/whaletown/internal/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
