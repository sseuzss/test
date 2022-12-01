#!/bin/bash
set -ex

apt update && apt install -y hashdeep fakeroot git wget
export VERSION=$( git describe --tags `git rev-list --tags --max-count=1` ).$CI_PIPELINE_IID
export VERSION=${VERSION/v/}
DEBIAN_FRONTEND=noninteractive apt install -y  rpm
mkdir install_linux/

mv DEBIAN/control TMP_control
mv DEBIAN/changelog TMP_changelog
mkdir -p vxagent/opt/vxcontrol/vxagent/bin && mkdir vxagent/DEBIAN && cp _tmp/linux/386/vxagent vxagent/opt/vxcontrol/vxagent/bin
arch="i386"
echo; echo "Generating vxagent/DEBIAN/control..."

eval "echo \"$(cat TMP_control)\"" > vxagent/DEBIAN/control
cat vxagent/DEBIAN/control || exit 1

echo; echo "Generating vxagent/DEBIAN/changelog..."
eval "echo \"$(cat TMP_changelog)\"" > vxagent/DEBIAN/changelog
cat vxagent/DEBIAN/changelog || exit 1

cp DEBIAN/* vxagent/DEBIAN

md5deep -r vxagent/opt/vxcontrol/vxagent > vxagent/DEBIAN/md5sums
chmod -R 755 vxagent/DEBIAN

fakeroot dpkg-deb --build vxagent vxagent-${VERSION}_${arch}.deb || exit 1

echo "Done create deb $arch"

rm -rf vxagent

mkdir -p vxagent/opt/vxcontrol/vxagent/bin && mkdir vxagent/DEBIAN && cp _tmp/linux/amd64/vxagent vxagent/opt/vxcontrol/vxagent/bin
arch="amd64"
echo; echo "Generating vxagent/DEBIAN/control..."

eval "echo \"$(cat TMP_control)\"" > vxagent/DEBIAN/control
cat vxagent/DEBIAN/control || exit 1

echo; echo "Generating vxagent/DEBIAN/changelog..."
eval "echo \"$(cat TMP_changelog)\"" > vxagent/DEBIAN/changelog
cat vxagent/DEBIAN/changelog || exit 1

cp DEBIAN/* vxagent/DEBIAN

md5deep -r vxagent/opt/vxcontrol/vxagent > vxagent/DEBIAN/md5sums
chmod -R 755 vxagent/DEBIAN

fakeroot dpkg-deb --build vxagent vxagent-${VERSION}_${arch}.deb || exit 1

echo "Done create deb $arch"

cp *.deb install_linux/
export VERSION=$( git describe --tags `git rev-list --tags --max-count=1` ).$CI_PIPELINE_IID
export VERSION=${VERSION/v/}
export VXSERVER_CONNECT="VXSERVER_CONNECT"
mkdir -p /root/rpmbuild/SOURCES/
arch="386"
eval "echo \"$(cat RPM/rpm.spec)\"" > rpm_$arch.spec
cp _tmp/linux/386/vxagent /root/rpmbuild/SOURCES/
rpmbuild -bb ./rpm_$arch.spec --target i386
cp /root/rpmbuild/RPMS/i386/* install_linux/vxagent-${VERSION}_i386.rpm

arch="amd64"
rm -rf /root/rpmbuild/SOURCES/* || true
cp _tmp/linux/amd64/vxagent /root/rpmbuild/SOURCES/
eval "echo \"$(cat RPM/rpm.spec)\"" > rpm_$arch.spec
rpmbuild -bb ./rpm_$arch.spec --target amd64
cp /root/rpmbuild/RPMS/amd64/* install_linux/vxagent-${VERSION}_amd64.rpm

