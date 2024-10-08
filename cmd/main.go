// Package main is the entry point of the application.
package main

import (
	"fmt"
	"os"

	"go.nhat.io/vanityrender/internal/cli"
	"go.nhat.io/vanityrender/internal/version"
)

const message = `
   :::     ::: :::::::::
  :+:     :+: :+:    :+:
 +:+     +:+ +:+    +:+
+#+     +:+ +#++:++#:
+#+   +#+  +#+    +#+
#+#+#+#   #+#    #+#
 ###     ###    ###  %s (rev: %s)

`

func main() {
	info := version.Info()

	_, _ = fmt.Fprintf(os.Stdout, message, info.Version, info.Revision) //nolint: errcheck

	os.Exit(cli.Execute())
}
