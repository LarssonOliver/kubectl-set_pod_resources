package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/completion"
)

// registerCompletionFuncForGlobalFlags copied from
// https://github.com/kubernetes/kubectl/blob/7577f36fbc78c41770457b39947be836ab8df949/pkg/cmd/cmd.go#L564-L585
func registerCompletionFuncForGlobalFlags(cmd *cobra.Command, f cmdutil.Factory) {
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"namespace",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completion.CompGetResource(f, "namespace", toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"context",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completion.ListContextsInConfig(toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"cluster",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completion.ListClustersInConfig(toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"user",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completion.ListUsersInConfig(toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
}

func registerCompletionFuncForFlags(cmd *cobra.Command, f cmdutil.Factory) {
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"containers",
		completion.ContainerCompletionFunc(f)))
}

func argsPodNameAndContainerCompletionFunc(f cmdutil.Factory) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	completionFunc := completion.PodResourceNameCompletionFunc(f)

	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
            containers := completion.CompGetContainers(f, args[0], toComplete)
            return containers, cobra.ShellCompDirectiveNoFileComp
		}

		pods, compDirective := completionFunc(cmd, args, fmt.Sprintf("pods/%s", toComplete))

		for i, v := range pods {
			if strings.HasPrefix(v, "pods/") {
				pods[i] = v[len("pods/"):]
			}
		}

		return pods, compDirective
	}
}
