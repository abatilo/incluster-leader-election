# incluster-leader-election

This is an implementation of utilizing the Kubernetes `coordination.k8s.io` API
paired with a [tilt.dev](https://tilt.dev/) environment for testing.

This project was immensely inspired by [this blog
post](https://carlosbecker.dev/posts/k8s-leader-election) and most of the code
is copy pasta.

## Getting started

This project installs and handles dependencies using
[asdf-vm](https://asdf-vm.com/#/) and minimum version `v0.8.0` is required.
Once `asdf-vm` is installed, the `Makefile` will install any necessary
dependencies for you.

Otherwise, please refer to the Makefile:

```bash
⇒  make
help                           View help information
bootstrap                      Perform all bootstrapping to start your project
clean                          Delete local dev environment
up                             Run a local development environment
```
