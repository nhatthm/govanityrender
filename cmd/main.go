package main

import (
	"os"

	"go.nhat.io/vanityrender/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
