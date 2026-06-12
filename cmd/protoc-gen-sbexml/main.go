package main

import (
	"fmt"
	"os"

	"github.com/prnvbn/protoc-gen-sbexml/internal/protocgen"
)

func main() {
	if err := protocgen.Run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
