Name:           vxagent
Version:        $VERSION
Release:        $CI_PIPELINE_IID
Summary:        This service for work XDR agent
License:        -

Source0:        vxagent
Requires:       bash, systemd, initscripts, libstdc++, libgcc, glibc

BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)


%description
VX agent

%global _enable_debug_package 0
%global debug_package %{nil}
%global __os_install_post /usr/lib/rpm/brp-compress %{nil}

%install
export DONT_STRIP=1
install -D -pm 755 %{SOURCE0} %{buildroot}/opt/vxcontrol/vxagent/bin/vxagent

%post -p /bin/bash
ln -s /usr/lib64/librt.so.1 /usr/lib64/librt.so 2>/dev/null || true
ln -s /usr/lib64/libpthread.so.0 /usr/lib64/libpthread.so 2>/dev/null || true
chmod +x /opt/vxcontrol/vxagent/bin/vxagent && chown root:root /opt/vxcontrol/vxagent/bin/vxagent

if [ -n "\$$VXSERVER_CONNECT" ]; then
  /opt/vxcontrol/vxagent/bin/vxagent -connect \$$VXSERVER_CONNECT -command install
else
  /opt/vxcontrol/vxagent/bin/vxagent -command install
fi
/opt/vxcontrol/vxagent/bin/vxagent -command start

%preun -p /bin/bash
/opt/vxcontrol/vxagent/bin/vxagent -command stop || true
if pgrep -f vxagent &>/dev/null; then
  kill -9 $(cat /var/run/vxagent.pid 2>/dev/null) &>/dev/null || true
fi
/opt/vxcontrol/vxagent/bin/vxagent -command uninstall || true

%postun -p /bin/bash
if ! [ -d /opt/ ]; then
  mkdir "/opt" || true
fi
rm -rf /opt/vxcontrol/vxagent || true
rmdir /opt/vxcontrol/ || true

%files
/opt/vxcontrol/vxagent/*

%clean
rm -rf $RPM_BUILD_ROOT

%changelog
* Mon Nov 22 2021 Positive Technologies <help@vxcontrol.com>
- Create RPM build
