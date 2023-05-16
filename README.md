# Example containers that derive from RHEL CoreOS

See https://docs.openshift.com/container-platform/4.12/post_installation_configuration/coreos-layering.html#coreos-layering

This repository contains example container builds.

## Examples

TODO

## Running an example

Build an image using an example from this repo and push it to an image registry,
using any container build tooling (`podman build`, OpenShift builds, Tekton, etc.)
and push it to a registry.

### Testing outside of an OpenShift cluster

One can directly boot a RHEL CoreOS image using e.g. an Ignition (or Butane) config
that just sets up an SSH key

### Testing on an OpenShift cluster

See the [documentation](https://docs.openshift.com/container-platform/4.12/post_installation_configuration/coreos-layering.html#coreos-layering).

## Test
Adding it to test ci run on a PR.