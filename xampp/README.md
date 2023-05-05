# Installing xampp

This example was chosen because it installs to a location unsupported by rpm-ostree(/opt/)

In this example, we:

- store the installer in /tmp/
- create a folder at the "expected" location behind the symlink
- Install the package via provided installer
- move the files to /usr/lib/lampp from the temporary install location
- set up a tmpfile.d for automating the following at startup:
    (i) creating the symlinks pointing to dir /usr/lib/lampp and binary at /usr/lib/lampp/xampp
    (ii) creating an opt dir in /var/*
- remove files left in /tmp/ 

More links:
- RHCOS layering: https://docs.openshift.com/container-platform/4.13/post_installation_configuration/coreos-layering.html#coreos-layering