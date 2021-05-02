# Contributing

As a starting point, please refer to the [Gardener contributor guide](https://github.com/gardener/documentation/blob/master/CONTRIBUTING.md).

## Equinix Metal Extension Provider

The rest of this document describes how to contribute to, and test, this Equinix Metal (formerly Packet) extension provider for Gardener.

The guide demonstrates how to make changes, test them and publish them using the
[Packet provider extension](https://github.com/gardener/gardener-extension-provider-packet).

Gardener uses
[an extensible architecture](https://github.com/gardener/gardener/blob/master/docs/proposals/01-extensibility.md)
which abstracts away things like cloud provider
implementations, OS-specific configuration for nodes and DNS provider logic.

Extensions are k8s controllers. They are *packaged* as Helm charts and are *registered* with the
Gardener API server by applying a `ControllerRegistration` custom k8s resource to the Gardener API
server.
In addition to being packaged as Helm charts, extensions often **deploy** Helm charts as well. For
example, the Packet provider extension deploys components such as the Packet CCM and MetalLB in
order to provide necessary services to Packet seed and shoot clusters.

## Requirements

- a running "garden" cluster with a kubeconfig for access oti

## Background

**Important!** Gardener's base control cluster, or "garden" cluster, really has
two **distinct** kinds of API servers, both of which are accessed via `kubectl` using
a `kubeconfig` file

* The "base cluster", the pre-existing normal Kubernetes cluster on which Gardener is installed, turning it into the "garden cluster"
* The "Gardener API server", a special k8s API server which doesn't have any nodes or pods and
which deals with Gardener resources only, and runs on top of the "base cluster"

The `kubeconfig` for the "base cluster" is wherever you set it when creating the initial Kubernetes
cluster, for example during [garden-setup](https://github.com/gardener/garden-setup). If you used
the standard `garden-setup` flow, then the `kubeconfig` is likely at `$GOPATH/src/github.com/gardener/sow/landscape/kubeconfig`.

The `kubeconfig` for the "Gardener API server", when using [garden-setup](https://github.com/gardener/garden-setup)
is at `$GOPATH/src/github.com/gardener/sow/landscape/export/kube-apiserver/kubeconfig`.

## Development workflow

Your development workflow in general is as follows. Relevant steps will be described in detail below.

### Setup

The setup steps are necessary just so that your extension can have what to work with against a garden cluster.
In general, you will do these once, although if you are changing how the extension works with some components,
e.g. the `CloudProfile` or `Shoot`, you may modify and redeploy them multiple times.

1. Deploy your base cluster
1. Convert your base cluster into a garden cluster
1. Deploy a seed to your "gardener API server"
1. Deploy a project to your "gardener API server"
1. Get an Equinix Metal API key and project ID, and save them to the secret file, which also contains the secret binding
1. Deploy the secret and secret binding to the garden cluster
1. Configure and deploy a cloud profile to the garden cluster
1. Configure and deploy a shoot to the garden cluster

#### Base Cluster and Garden Cluster

Deploying your base cluster is beyond the scope of these documents. Deploy it any way that works for you.
We recommend following the instructions at [garden-setup](https://github.com/gardener/garden-setup).
The documentation there describes, as well, how to deploy a seed.

#### Ensure the Seed is Deployed

You need at least one seed deployed. For simplicity, you should deploy the seed right into the garden cluster.

1. Connect to the "gardener API server", e.g. `KUBECONFIG=./export/kube-apiserver/kubeconfig`
1. `kubectl get seed` - this should return at least one functional seed

For example:

```sh
$ KUBECONFIG=./export/kube-apiserver/kubeconfig  kubectl get seed
NAME   STATUS   PROVIDER   REGION        AGE   VERSION   K8S VERSION
gcp    Ready    gcp        us-central1   60d   v1.17.1   v1.19.9-gke.1400
```

Note that the above garden cluster is running in GCP, and so the seed is using the `PROVIDER` named `gcp`. That is fine.
The Equinix Metal extension will be invoked when deploying a shoot to Equinix Metal, and can do so from a seed
in GCP.

#### Deploy a Project

A `Project` groups together shoots and infrastructure secrets in a namespace.
A sample `Project` is available at [example/23-project.yaml](./example/23-project.yaml). Copy it over to a temporary workspace,
modify it as needed, and then apply it.

```sh
kubectl apply -f 23-project.yaml
```

Unless you actually will be using the Gardener UI, most of the RBAC entries in the file do not matter for development.
The only really important elements are:

* `name`: pick a unique one for the `Project`
* `namespace`: you will need to be consistent in using the same namespace for multiple elements

#### Secret

You need two pieces of information for Gardener to do its job against the Equinix Metal API:

* project ID - the UUID of the project in which it should manage devices
* API key - your unique API key that has the rights to create/delete devices in that project

Both of these should be placed in a Kubernetes `Secret`, as well as the `SecretBinding` that enables the extension to use them.

A sample `Secret` is available at [example/25-secret.yaml](./example/25-secret.yaml). Copy it over to a temporary workspace,
modify it as needed, and then apply it.

**Important:** The `Secret` and `SecretBinding` must be in the same `namespace` as the `Project`.

```sh
kubectl apply -f 25-secret.yaml
```

#### Cloud Profile

The `CloudProfile` is a resource that contains the list of acceptable machine types and OS images. Each
`Shoot`, when deployed, uses a specific `CloudProfile` and picks elements from it.

A sample `CloudProfile` is available at [example/26-cloudprofile.yaml](./example/26-cloudprofile.yaml)/ Copy it over to a
temporary workspace, modify as needed, and then apply it.

```sh
kubectl apply -f 26-cloudprofile.yaml
```

#### Shoot

With all of the above, you are ready to create a `Shoot`. The `Shoot` is the instructions for actually deploying
a Kubernetes cluster, managed by Gardener, on machines deployed by Gardener. The interaction with Equinix Metal
occurs via this extension provider.

A sample `Shoot` is available at [example/90-shoot.yaml](./example/90-shoot.yaml). Copy it over to a temporary workspace,
modify as needed, and then apply it.

```sh
kubectl apply -f 90-shoot.yaml
````

The `Shoot` resource is long and complex, and beyond the scope of this document. The important parts to note are:

* `namespace` - must match the namespace for the `Project` and `Secret`/`SecretBinding`
* `secretBindingName` - must match the name of the `SecretBinding`
* `cloudProfileName` - must match the name of the `CloudProfile`
* `region` - must be one of the regions in the `CloudProfile`
* `infrastructureConfig` - must conform to the infrastructure config, see [here](./example/30-infrastructure.yaml#L48-L49)
* `controlPlaneConfig` - must conform to the control plane config, see [here](./example/30-controlplane.yaml#L58-L59)

### Extension

With the setup in place, you can run your local extension against the cluster.

1. Get the namespace in which your extensions are running; see below
1. Determine the port on which the hook will be listening; see below.
1. Run the webhook: `./hack/hook-me.sh <namespace> <port>`, where `<namespace>` and `<port>` are as determined above
1. Run the extension locally: `make start`

You need the webhook so that the "gardener server" can communicate with your locally running extension.

`<namespace>` is the namespace that the extension is running in on the "base cluster". _After_ you deployed the shoot,
connect to the "base cluster" in which the seed is running, and do `kubectl get ns`; the target namespace should have
your provider name in it.

For example:

```sh
$ kubectl get ns
NAME                                STATUS   AGE
cert-manager                        Active   61d
default                             Active   64d
extension-dns-external-j4psd        Active   6d21h
extension-networking-calico-xjshs   Active   6d21h
extension-os-ubuntu-nrhxn           Active   6d21h
extension-provider-gcp-bcdwd        Active   6d21h
extension-provider-packet-9r7xh     Active   6d21h
garden                              Active   61d
kube-node-lease                     Active   64d
kube-public                         Active   64d
kube-system                         Active   64d
shoot--em--em-test                  Active   6d21h
```

In this case, it is `extension-provider-packet-9r7xh`.

`<port>` is the port on which the provider is listening. By default, it is `8443`, but it is set in the `Makefile` as `WEBHOOK_CONFIG_PORT`.

Combining the above, we get:

```sh
./hack/hook-me.sh extension-provider-packet-9r7xh 8443
```

### Deploying an extension to a "garden" cluster

>NOTE: Some operations in the development workflow don't support Go modules. It is recommended to
>clone the extension's source directory into your GOPATH.

If you don't already have it, clone the extension's source repository:

```console
mkdir -p $GOPATH/src/github.com/gardener && cd $_
git clone git@github.com:gardener/gardener-extension-provider-packet.git
cd gardener-extension-provider-packet
```

Set your `KUBECONFIG` env var to point at the kubeconfig file belonging to the **Gardener** API
server (not the base API server). For a Gardener cluster deployed using
[garden-setup](https://github.com/gardener/garden-setup), the command should be the following:

```
export KUBECONFIG=$(pwd)/export/kube-apiserver/kubeconfig
```

Verify connectivity to the Gardener API server:

```
kubectl get seeds
```

Sample output:

```
NAME   STATUS   PROVIDER   REGION         AGE    VERSION   K8S VERSION
aws    Ready    aws        eu-central-1   104m   v1.14.0   v1.18.12
```

The extension is in [controller-registration.yaml](./example/controller-registration.yaml).
It contains the Helm chart from [here](./charts/gardener-extension-provider-packet/values.yaml), tarred, gzipped, and then
base64-encoded. By default, it deploys one replica of the extension pod from a registry.

>NOTE: The `controller-registration.yaml` file in the `master` branch contains a reference to an
>**unreleased** container image. Either overwrite the `tag` field in the file or check out a Git
>tag before registering the extension.

If you are developing locally, you don't actually want to run one replica; you want 0. Modify the values so that it has:

```yaml
values:
  replicaCount: 0
```

This will set up all of the prerequisites but not try to deploy a controller pod, leaving you free to run `make start`

Next, register the extension:

```
kubectl apply -f ./example/controller-registration.yaml
```

From this point the extension should be deployed on-demand, e.g. when a shoot cluster gets created.

### Making changes to an extension

Install the development requirements:

```
make install-requirements
```

Make code changes, then run the following:

```
# Run the tests
make test

# Re-generate any auto-generated files
make generate

# Verify auto-generated files are up to date
make check-generate

# Build a Docker image
make docker-images
```

You can now tag and push the resulting image to a Docker registry:

```
docker tag 9fa4c197cf04 quay.io/jlieb/gardener-extension-packet:testing
docker push quay.io/jlieb/gardener-extension-packet:testing
```

To use the new image, register or re-register the extension with the Gardener API server. Note that
when re-pushing an existing Docker tag, you will likely have to temporarily set the extension's
`imagePullPolicy` field to `Always` and delete the extension's pod to force the kubelet to re-pull
the image.

### Debugging

To check the logs of the extension controller, run the following command against the **base** API
server:

```
kubectl -n extension-provider-packet-xxxxx logs gardener-extension-provider-packet-xxxxx
```

### Miscellaneous and "gotchas"

The Helm chart that is used by Gardener to deploy the extension controller is packaged as a
**base64-encoded tgz archive** inline inside `example/controller-registration.yaml` under the
`chart` field. To view its contents, run the following command:

```
cat example/controller-registration.yaml | grep chart | awk {'print $2'} | base64 -d | tar zxvf -
```

Some generated files are **gitignored**. To avoid building images using outdated generated data, be
sure to **always** run `make generate` before running `make docker-images` (even when there is no
visible Git diff).
