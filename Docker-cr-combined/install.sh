apt-get install -y git mercurial golang sqlite libdevmapper-dev libdevmapper-dev libseccomp-dev libseccomp2
./hack/vendor.sh
env AUTO_GOPATH=1 DOCKER_EXPERIMENTAL=1 DOCKER_BUILDTAGS='exclude_graphdriver_btrfs exclude_graphdriver_devicemapper selinux seccomp'  ./hack/make.sh binary
cp -v bundles/1.10.0-dev/binary/* /usr/bin/
