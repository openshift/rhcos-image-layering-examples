# Get RHCOS base image of target cluster `oc adm release info --image-for rhel-coreos`
# hadolint ignore=DL3006
FROM quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256...

# This example installs to /usr/local/bin, so to work around this, create a dir behind the symlink
# After install, move the binary to /usr/bin/ for CLI use.
RUN mkdir -p /var/usrlocal/bin && \ 
    rpm -Uvh https://github.com/mozilla/sops/releases/download/v3.7.3/sops-3.7.3-1.x86_64.rpm && \
    mv /var/usrlocal/bin/sops /usr/bin/sops && \
    ostree container commit
