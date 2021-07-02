Name:           mhc
Version:        1.0
Release:        1%{?dist}
Summary:        Disables upstream mids in Traffic Server based on data from Apache Traffic Control.

License:        MIT
URL:            https://github.com/ARMmaster17/atc-mid-cache-health-check
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang
Requires:       systemd-rpm-macros

%description

%global debug_package %{nil}

%prep
%autosetup


%build
make build

%install
rm -rf $RPM_BUILD_ROOT
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}
install -Dpm 0755 %{name}.conf %{buildroot}%{_sysconfdir}/%{name}/%{name}.conf
install -Dpm 644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service


%post
%systemd_post %{name}.service


%preun
%systemd_preun %{name}.service


%files
%dir %{_sysconfdir}/%{name}
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%config(noreplace) %{_sysconfdir}/%{name}/%{name}.conf


%changelog
* Thu Jul  1 2021 root
- 
