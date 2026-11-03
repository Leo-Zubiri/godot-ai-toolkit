// Command godot-ai-toolkit is the toolkit CLI: install, update, status, validate.
//
// It is distributed as a single native binary per OS/arch; end users need no
// Go, Python, Node or uv runtime to use it.
package main

import (
	"os"

	"github.com/godot-ai-toolkit/godot-ai-toolkit/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
