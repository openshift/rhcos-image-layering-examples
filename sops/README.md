# Installing sops

This example was chosen because it installs to a location unsupported by rpm-ostree(/usr/local/)

In this example, we:

- create a folder at the "expected" location behind the symlink
- Install the package via rpm
- move the binary to /usr/bin from the temporary install location

More links:
- About sops: https://github.com/mozilla/sops#sops-secrets-operations
- RHCOS layering: https://docs.openshift.com/container-platform/4.13/post_installation_configuration/coreos-layering.html#coreos-layering