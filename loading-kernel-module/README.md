# Loading Kernel Module

### Full demo

[![asciicast](https://asciinema.org/a/522655.svg)](https://asciinema.org/a/522655)

### Prerequisits

* OCP 4.12+
* `rpm-ostree` version `2022.10`+

### Background

We are going to build a container image that contain a driver that we wish to
load into the node's kernel.

We will use the [driver-toolkit](https://github.com/openshift/driver-toolkit)
as a base image in order to build the [simple-kmod](https://github.com/openshift-psap/simple-kmod)
driver and then add the driver's binary to a new layer on top of the REHL-CoreOS
image that runs on the node.

### Get the correct driver-toolkit image

The driver-toolkit contains the kernel package and headers in order to build
kernel modules for a specific kernel, therefore, we need to make sure we use
the correct digest of the driver-toolkit image to fit our specific node.

Let's get the driver-toolkit image for our 4.12.18 cluster
```
$ oc adm release info quay.io/openshift-release-dev/ocp-release:4.12.18-x86_64 --image-for=driver-toolkit
quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:d8f228d47ebf0a5ec0c343f7ed3c5aec1d62d441d671703e5e7f40d78b8d97b9
```

For any other version/arch you can use
```
$ oc adm release info quay.io/openshift-release-dev/ocp-release:<cluste version>-x86_64 --image-for=driver-toolkit
$ oc adm release info quay.io/openshift-release-dev/ocp-release:<cluster version>-aarch64 --image-for=driver-toolkit
```

### Get the correct rhel-coreos-8 image

We want to make sure we are using the same image that runs on the nodes to be
used as the base image of our driver-container image for minimal side effects.
We only want to add our new `.ko` file (and required configuration files if such
exists) to the node. Nothing else should change.

Note:
OCP 4.13 is using `rhel-coreos-9` while OCP 4.12 is using `rhel-coreos-8`.
Be sure to use the correct `--image-for` based on your cluster version.

Let's get the `rhel-coreos` image for our cluster version.
```
$ oc adm release info quay.io/openshift-release-dev/ocp-release:4.12.18-x86_64 --image-for=rhel-coreos-8
quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:127670aaa6073e97bd6d919e1b57bd16b60435410b096aeefce7cd64f6c83d24
```

### Get the correct kernel-version

We also need the correct kernel version in order to build the driver-container.

Since after reboot, the kernel of the host will be the kernel RPM installed on
the new image, aka, the `rhel-coreos-8` image that we are using as the last
layer, then the correct way to get the kernel version is by getting it from the
image itself.
```
$ podman run -it quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:127670aaa6073e97bd6d919e1b57bd16b60435410b096aeefce7cd64f6c83d24 \
    rpm -qa | grep kernel
kernel-modules-4.18.0-372.53.1.el8_6.x86_64
kernel-core-4.18.0-372.53.1.el8_6.x86_64
kernel-4.18.0-372.53.1.el8_6.x86_64
kernel-modules-extra-4.18.0-372.53.1.el8_6.x86_64
```

If you don't have access to the image, make sure you are using your [pull-secret](https://console.redhat.com/openshift/install/pull-secret) with Podman.

We can see the kernel RPM in the image is `4.18.0-372.53.1.el8_6.x86_64`.

### Set up the environment

For simplicity, let's define some env variables to hold our previous results
```
$ export DTK_IMAGE=quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:d8f228d47ebf0a5ec0c343f7ed3c5aec1d62d441d671703e5e7f40d78b8d97b9
$ export RHCOS_IMAGE=quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:127670aaa6073e97bd6d919e1b57bd16b60435410b096aeefce7cd64f6c83d24
$ export KERNEL_VERSION=4.18.0-372.53.1.el8_6.x86_64
```

### Build the container image

We are going to use the following [Containerfile](./loading-kernel-module/Containerfile)
for building
```
$ podman build \
    --build-arg KERNEL_VERSION=${KERNEL_VERSION} \
    --build-arg DTK_IMAGE=${DTK_IMAGE} \
    --build-arg RHCOS_IMAGE=${RHCOS_IMAGE} \
    -t quay.io/ybettan/coreos-layering:simple-kmod \
    loading-kernel-module
[1/2] STEP 1/6: FROM quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:d8f228d47ebf0a5ec0c343f7ed3c5aec1d62d441d671703e5e7f40d78b8d97b9 AS builder
[1/2] STEP 2/6: ARG KERNEL_VERSION
--> Using cache 0a1ec1eb4b3eb2accba3475b2e347ac81c0f01b607d25c5e81c0d01471092d21
--> 0a1ec1eb4b3
[1/2] STEP 3/6: WORKDIR /build/
--> Using cache 02fd2db894dbf7f2c59b7dbabf64c2114eb283ddd230163fb9bd8df892045e02
--> 02fd2db894d
[1/2] STEP 4/6: RUN git clone https://github.com/openshift-psap/simple-kmod.git
--> Using cache 419fda4f52482891200333bfa165ab40d66c40bd4180be1c207dd48ba6fa7d27
--> 419fda4f524
[1/2] STEP 5/6: WORKDIR /build/simple-kmod
--> Using cache 80d50b672deb18ba0f3876e32d487eaaa5832765c87f724c079130576c8cfd6d
--> 80d50b672de
[1/2] STEP 6/6: RUN make all install KVER=4.18.0-372.53.1.el8_6.x86_64
make -C /lib/modules/4.18.0-372.53.1.el8_6.x86_64/build M=/build/simple-kmod EXTRA_CFLAGS=-DKMODVER=\\\"e852852\\\" modules
make[1]: Entering directory '/usr/src/kernels/4.18.0-372.53.1.el8_6.x86_64'
  CC [M]  /build/simple-kmod/simple-kmod.o
  CC [M]  /build/simple-kmod/simple-procfs-kmod.o
  Building modules, stage 2.
  MODPOST 2 modules
  CC      /build/simple-kmod/simple-kmod.mod.o
  LD [M]  /build/simple-kmod/simple-kmod.ko
  CC      /build/simple-kmod/simple-procfs-kmod.mod.o
  LD [M]  /build/simple-kmod/simple-procfs-kmod.ko
make[1]: Leaving directory '/usr/src/kernels/4.18.0-372.53.1.el8_6.x86_64'
gcc -o spkut ./simple-procfs-kmod-userspace-tool.c
install -v -m 755 spkut /bin/
'spkut' -> '/bin/spkut'
install -v -m 755 -d /lib/modules/4.18.0-372.53.1.el8_6.x86_64/
install -v -m 644 simple-kmod.ko        /lib/modules/4.18.0-372.53.1.el8_6.x86_64/simple-kmod.ko
'simple-kmod.ko' -> '/lib/modules/4.18.0-372.53.1.el8_6.x86_64/simple-kmod.ko'
install -v -m 644 simple-procfs-kmod.ko /lib/modules/4.18.0-372.53.1.el8_6.x86_64/simple-procfs-kmod.ko
'simple-procfs-kmod.ko' -> '/lib/modules/4.18.0-372.53.1.el8_6.x86_64/simple-procfs-kmod.ko'
depmod -F /lib/modules/4.18.0-372.53.1.el8_6.x86_64/System.map 4.18.0-372.53.1.el8_6.x86_64
--> 26dfedc51be
[2/2] STEP 1/5: FROM quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:127670aaa6073e97bd6d919e1b57bd16b60435410b096aeefce7cd64f6c83d24
[2/2] STEP 2/5: ARG KERNEL_VERSION
--> 2879805a870
[2/2] STEP 3/5: COPY --from=builder /etc/driver-toolkit-release.json /etc/
--> 8dea5bb2bd8
[2/2] STEP 4/5: COPY --from=builder /lib/modules/4.18.0-372.53.1.el8_6.x86_64/*.ko /usr/lib/modules/4.18.0-372.53.1.el8_6.x86_64/
--> 6249995c281
[2/2] STEP 5/5: RUN depmod -a 4.18.0-372.53.1.el8_6.x86_64 && echo simple_kmod > /etc/modules-load.d/simple_kmod.conf
[2/2] COMMIT quay.io/ybettan/coreos-layering:simple-kmod
--> 612c8c8a54f
Successfully tagged quay.io/ybettan/coreos-layering:simple-kmod
612c8c8a54f0800cf8dcb4c20d5d3d384dc47463711d6d3ad1bd93d52e982b1e
```

Now, let's push the image to some image registry in order for the image to be accesible
by the cluster.
```
$ podman push quay.io/ybettan/coreos-layering:simple-kmod
```

Make sure the image is public or it has the correct pull-secret to pull the image.
The pull secret in the cluster is
```
$ oc get secret/pull-secret -n openshift-config.
```
For more info, check the [how to update the global pull-secret](https://docs.openshift.com/container-platform/4.12/openshift_images/managing_images/using-image-pull-secrets.html#images-update-global-pull-secret_using-image-pull-secrets) doc.

### Loading the kernel module into the nodes

Make sure to change the `osImageURL` section in the [machineconfig.yaml](./loading-kernel-module/machineconfig.yaml),
and then apply it.
```
oc apply -f loading-kernel-module/machineconfig.yaml
```

The Machine Config Operator will start updating workers nodes one by one with
the new image specified in `osImageURL` which contain our kernel module.

Altough there are some optimizations in place to apply live-updates, in some
cases, we will experience a nodes reboot.

This may take some time and you can track the progress using
```
$ watch oc get nodes
```
until all workers nodes are `Ready` again.

Here is a snapshot of the node updating process
```
NAME                           STATUS                     ROLES                  AGE   VERSION
ip-10-0-146-108.ec2.internal   Ready                      control-plane,master   39m   v1.25.8+37a9a08
ip-10-0-147-43.ec2.internal    Ready,SchedulingDisabled   worker                 32m   v1.25.8+37a9a08
ip-10-0-166-128.ec2.internal   Ready                      control-plane,master   39m   v1.25.8+37a9a08
ip-10-0-178-9.ec2.internal     Ready                      worker                 30m   v1.25.8+37a9a08
ip-10-0-226-15.ec2.internal    Ready                      control-plane,master   39m   v1.25.8+37a9a08
ip-10-0-228-190.ec2.internal   Ready                      worker                 29m   v1.25.8+37a9a08
```

For some additional info regarding MCO in case it doesn't update the nodes check
the [machine-config-pull-status](https://docs.openshift.com/container-platform/4.13/post_installation_configuration/machine-configuration-tasks.html#checking-mco-status_post-install-machine-configuration-tasks).

### Validate that the kernel module was loaded

We can start with validating that the node is indeed running a customer image
```
$ oc debug node/ip-10-0-147-43.ec2.internal -- chroot host/ rpm-ostree status
Starting pod/ip-10-0-147-43.ec2.internal-debug ...
To use host binaries, run `chroot /host`
State: idle
Deployments:
* ostree-unverified-registry:quay.io/ybettan/coreos-layering:simple-kmod
                   Digest: sha256:9894d87312b03e04a01371864dd30b56d223f37d7b808ee08f4054eef92f5033
                  Version: 412.86.202305161131-0 (2023-06-04T10:55:24Z)

Removing debug pod ...
```

We can also make sure that the kernel-module was loaded successfuly.
```
$ oc debug node/ip-10-0-147-43.ec2.internal -- chroot host/ lsmod | grep simple_kmod
Starting pod/ip-10-0-147-43ec2internal-debug ...
To use host binaries, run `chroot /host`
simple_kmod            16384  0

Removing debug pod ...
```
