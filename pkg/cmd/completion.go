/*
Copyright 2025 Oliver Larsson

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/completion"
)

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
