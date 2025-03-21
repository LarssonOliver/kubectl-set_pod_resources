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

	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/cmd/set"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	generateversioned "k8s.io/kubectl/pkg/generate/versioned"
	"k8s.io/kubectl/pkg/polymorphichelpers"
	"k8s.io/kubectl/pkg/scheme"
	"maps"
)

// Set by the build process
var Version string

var (
	setPodResourcesExample = `
  # Set the requested resources of a pod named 'foo' to 1Gi of memory and 200m of CPU
  %[1]s set-pod-resources foo --requests=memory=1Gi,cpu=200m

  # Set the limits of a pod named 'foo' to 2Gi of memory and 500m of CPU
  %[1]s set-pod-resources foo --limits=memory=2Gi,cpu=500m

  # Set the requests of a container named 'bar' in a pod named 'foo' to 1Gi of memory and 200m of CPU
  %[1]s set-pod-resources foo bar --requests=memory=1Gi,cpu=200m`
)

// Command options
type SetPodResourcesOptions struct {
	configFlags *genericclioptions.ConfigFlags

	infos []*resource.Info

	namespace         string
	selector          string
	containerSelector string

	limits               string
	requests             string
	resourceRequirements v1.ResourceRequirements

	rawConfig api.Config
	args      []string

	genericiooptions.IOStreams
}

func NewSetPodResourcesOptions(streams genericiooptions.IOStreams) *SetPodResourcesOptions {
	return &SetPodResourcesOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		containerSelector: "*",

		IOStreams: streams,
	}
}

func NewCmdSetPodResources(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewSetPodResourcesOptions(streams)
	f := cmdutil.NewFactory(o.configFlags)

	cmd := &cobra.Command{
		Use:               "set-pod-resources [pod-name] [container] [flags]",
		Short:             "Resize the resources of a pod",
		Example:           fmt.Sprintf(setPodResourcesExample, "kubectl"),
		SilenceUsage:      true,
		Version:           Version,
		ValidArgsFunction: argsPodNameAndContainerCompletionFunc(f),
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl set-pod-resources",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	cmdutil.AddLabelSelectorFlagVar(cmd, &o.selector)

	cmd.Flags().StringVarP(&o.containerSelector, "containers", "c", o.containerSelector, "The names of containers in the selected pod templates to change, all containers are selected by default - may use wildcards")
	cmd.Flags().StringVar(&o.limits, "limits", o.limits, "The resource requirement requests for this container.  For example, 'cpu=100m,memory=256Mi'.")
	cmd.Flags().StringVar(&o.requests, "requests", o.requests, "The resource requirement requests for this container.  For example, 'cpu=100m,memory=256Mi'.")

	registerCompletionFuncForGlobalFlags(cmd, f)
	registerCompletionFuncForFlags(cmd, f)

	return cmd
}

// Completes the ResizePodOptions struct with the provided command and arguments
func (o *SetPodResourcesOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.args = args

	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	o.namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	builder := f.NewBuilder().
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		ContinueOnError().
		NamespaceParam(o.namespace).DefaultNamespace().
		Flatten()

	if len(o.args) > 0 {
		builder = builder.ResourceNames("pods", o.args[0])
	}

	if len(o.args) > 1 {
		o.containerSelector = o.args[1]
	}

	if len(o.selector) > 0 {
		builder = builder.LabelSelectorParam(o.selector)
	}

	builder = builder.Latest()

	o.infos, err = builder.Do().Infos()
	return err
}

// Validates the ResizePodOptions struct and returns an error if validation fails
func (o *SetPodResourcesOptions) Validate() error {
	var err error

	if len(o.args) == 0 && o.selector == "" {
		return fmt.Errorf("you must specify a pod or a pod selector")
	}

	if len(o.args) > 2 {
		return fmt.Errorf("too many arguments: you may only specify a pod name and container name")
	}

	if len(o.limits) == 0 && len(o.requests) == 0 {
		return fmt.Errorf("you must specify at least one of --limits or --requests")
	}

	o.resourceRequirements, err = generateversioned.HandleResourceRequirementsV1(map[string]string{"limits": o.limits, "requests": o.requests})
	if err != nil {
		return err
	}

	return nil
}

// Runs the resize pod command
func (o *SetPodResourcesOptions) Run() error {
	allErrs := []error{}
	patches := set.CalculatePatches(o.infos, scheme.DefaultJSONEncoder(), func(obj runtime.Object) ([]byte, error) {
		transformed := false
		_, err := polymorphichelpers.UpdatePodSpecForObjectFn(obj, func(spec *v1.PodSpec) error {
			initContainers, _ := selectContainers(spec.InitContainers, o.containerSelector)
			containers, _ := selectContainers(spec.Containers, o.containerSelector)
			containers = append(containers, initContainers...)
			if len(containers) != 0 {
				for i := range containers {
					if len(o.limits) != 0 && len(containers[i].Resources.Limits) == 0 {
						containers[i].Resources.Limits = make(v1.ResourceList)
					}
					maps.Copy(containers[i].Resources.Limits, o.resourceRequirements.Limits)

					if len(o.requests) != 0 && len(containers[i].Resources.Requests) == 0 {
						containers[i].Resources.Requests = make(v1.ResourceList)
					}
					maps.Copy(containers[i].Resources.Requests, o.resourceRequirements.Requests)
					transformed = true
				}
			} else {
				allErrs = append(allErrs, fmt.Errorf("error: unable to find container named %s", o.containerSelector))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if !transformed {
			return nil, nil
		}
		return runtime.Encode(scheme.DefaultJSONEncoder(), obj)
	})

	for _, patch := range patches {
		info := patch.Info
		name := info.ObjectName()
		if patch.Err != nil {
			allErrs = append(allErrs, fmt.Errorf("error: %s %v\n", name, patch.Err))
			continue
		}

		//no changes
		if string(patch.Patch) == "{}" || len(patch.Patch) == 0 {
			continue
		}

		_, err := resource.
			NewHelper(info.Client, info.Mapping).
			WithSubresource("resize").
			Patch(info.Namespace, info.Name, types.StrategicMergePatchType, patch.Patch, nil)
		if err != nil {
			allErrs = append(allErrs, fmt.Errorf("failed to patch resources update to pod template %v", err))
			continue
		}
	}

	return errors.NewAggregate(allErrs)
}
