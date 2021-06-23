FROM centos:7

RUN yum install golang golang-pkg-linux-amd64 golang-godoc golang-vet golang-src make -y