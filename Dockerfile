FROM centos:8

RUN dnf install git make -y
RUN dnf module install go-toolset -y
RUN dnf groupinstall "RPM Development Tools" -y
RUN rpmdev-setuptree

WORKDIR /src