# kubectl-set_pod_resources

<div align="center">

  ![GitHub License](https://img.shields.io/github/license/larssonoliver/kubectl-set_pod_resources)
  ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/larssonoliver/kubectl-set_pod_resources)

</div>

## Overview
kubectl-set_pod_resources is a plugin designed to simplify 
InPlacePodVerticalScaling using the `kubectl` CLI tool.

### Background
Since Kubernetes 1.32, resizing of already running pods must be done using the
pod `/resize` API sub resource. This means the `kubectl set resources` command
or `kubectl patch` (including `kubectl edit`) no longer works to edit resources
of pods directly as it did pre-1.32. See this PR for additional 
info: kubernetes/kubernetes#128266. 

Sure, but **what does this mean?**

Well, resizing the CPU of a running pod could previously be done using the
following command:
```bash
kubectl set resources pod nginx --requests=cpu=2 --limits=cpu=4
```
Not too bad! 👍

Since 1.32 however, the same command will be rejected by the API server.
Instead, we have to type out this monstrosity:
```bash
kubectl patch pod nginx --subresource resize --patch '{"spec":{"containers":[{"name":"nginx","resources":{"requests":{"cpu":"2"},"limits":{"cpu":"4"}}}]}}'
```
For scripting I guess this is fine, but I regularly type this when experimenting. 🤮

### The Solution
This plugin enables the use of the `set resources` syntax, but invokes the
`/resize` subresource in the background and thus works for pods. The above
resource patch becomes:
```bash
kubectl set-pod-resources nginx --requests=cpu=2 --limits=cpu=4
```
Not too bad if I say so myself! 😎


## Building 
This project requires `go 1.24.0`.
```bash
make
```

This produces the binary `bin/kubectl-set_pod_resources`.

## Installation
Simply place the `kubectl-set_pod_resources` binary somewhere in your `PATH`.
For shell completion, also include the `bin/kubectl_complete-set_pod_resources`
somewhere in your `PATH`.

Example:
```bash
sudo cp bin/* /usr/local/bin
```

### Uninstall
Remove the previously installed files from `PATH`. That is: 
- `kubectl-set_pod_resources`
- `kubectl_complete-set_pod_resources`

## Usage
After installation, you can use the plugin with the following commands:

```bash
$ kubectl set-pod-resources --help
```
```
Resize the resources of a pod

Usage:
  kubectl set-pod-resources [pod-name] [container] [flags]

Examples:

  # Set the requested resources of a pod named 'foo' to 1Gi of memory and 200m of CPU
  kubectl set-pod-resources foo --requests=memory=1Gi,cpu=200m

  # Set the limits of a pod named 'foo' to 2Gi of memory and 500m of CPU
  kubectl set-pod-resources foo --limits=memory=2Gi,cpu=500m

  # Set the requests of a container named 'bar' in a pod named 'foo' to 1Gi of memory and 200m of CPU
  kubectl set-pod-resources foo bar --requests=memory=1Gi,cpu=200m

Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/home/olars/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
  -c, --containers string              The names of containers in the selected pod templates to change, all containers are selected by default - may use wildcards (default "*")
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
  -h, --help                           help for kubectl set-pod-resources
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --limits string                  The resource requirement requests for this container.  For example, 'cpu=100m,memory=256Mi'.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
      --requests string                The resource requirement requests for this container.  For example, 'cpu=100m,memory=256Mi'.
  -l, --selector string                Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints.
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
  -v, --version                        version for kubectl set-pod-resources
```

## Disclosure
This project is heavily based off of the Kubernetes 
[sample-cli-plugin](https://github.com/kubernetes/sample-cli-plugin) repository 
and the source for the `kubectl set resources` command 
([set_resources.go](https://github.com/kubernetes/kubectl/blob/7577f36fbc78c41770457b39947be836ab8df949/pkg/cmd/set/set_resources.go)).

## License
This project is licensed under the Apache 2.0 license - see the LICENSE file for details.

Any source taken from the [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes)
repository is similarly licensed under the Apache 2.0 license - Copyright The Kubernetes Authors.

