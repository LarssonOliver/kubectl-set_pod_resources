package main

import (
	"os"

	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericiooptions"
    "github.com/larssonoliver/kubectl-set_pod_resources/pkg/cmd"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-set_pod_resources", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdSetPodResources(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
