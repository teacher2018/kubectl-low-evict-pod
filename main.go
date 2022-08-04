package main

import (
	"evict/cmd"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	evict := cmd.NewCmdEvict(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := evict.Execute(); err != nil {
		os.Exit(1)
	}
}
