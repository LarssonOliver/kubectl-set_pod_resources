package main

import (
	"os"

	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericiooptions"
    "github.com/larssonoliver/kubectl-resize-pod/pkg/cmd"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-resize-pod", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdResizePod(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
