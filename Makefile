NAME = mhc
VERSION = 1.0
RPM = $(NAME)-$(VERSION)-1.el8.x86_64.rpm
TAR = $(NAME)-$(VERSION).tar.gz
EXEC_FILE = $(NAME)
SRC_FILE = cmd/mhc.go
JUNIT_REPORT = report.xml
COVERAGE_REPORT = cover.out
LD_FLAGS=-ldflags=-linkmode=external

all: build

rpm: clean
	rpmdev-setuptree
	touch ./$(TAR)
	tar -czf $(TAR) --exclude=$(TAR) --transform 's,^\.,$(NAME)-$(VERSION),' .
	mv $(TAR) ~/rpmbuild/SOURCES
	rpmbuild -ba $(NAME).spec
	cp ~/rpmbuild/RPMS/x86_64/$(RPM) $(RPM)

build: $(EXEC_FILE)

$(EXEC_FILE):
	go build -v $(LD_FLAGS) -o $(NAME) $(SRC_FILE)

test:
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=$(COVERAGE_REPORT) ./... 2>&1 | go-junit-report > $(JUNIT_REPORT)

clean:
	rm -f ~/rpmbuild/RPMS/x86_64/$(RPM) || true
	rm -f ./artifacts/$(RPM) || true
	rm -f ./$(EXEC_FILE) || true
	rm -f ./$(TAR) || true
	rm -f ~/rpmbuild/SOURCES/$(TAR) || true

build-image:
	docker build -t mid-health-check-svc .

build-centos: build-image
	docker run -w /src -v "$(PWD):/src" mid-health-check-svc bash -c "make rpm"
	mkdir artifacts || true
	mv $(RPM) ./artifacts

.PHONY: all build test clean rpm build-centos