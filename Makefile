RPM = mid-health-check.rpm
EXEC_FILE = mid-health-check
SRC_FILE = main.go
JUNIT_REPORT = report.xml
COVERAGE_REPORT = cover.out

all: $(RPM)

$(RPM): clean $(EXEC_FILE)
	# TODO: Make the RPM

$(EXEC_FILE): test
	go build $(SRC_FILE)

test:
	#go get -u github.com/jstemmer/go-junit-report
	#go test -v -coverprofile=$(COVERAGE_REPORT) ./... 2>&1 | go-junit-report > $(JUNIT_REPORT)
	go test ./...

clean:
	rm -f $(RPM)
	rm -f $(EXEC_FILE)

build-image:
	docker build -t mid-health-check-svc .

build-centos: build-image
	docker run -w /src -v "$(pwd):/src" mid-health-check-svc sh -c "make mid-health-check.rpm"

test-centos: build-image
	rm ./report.xml || true
	docker run -w /src -v "$(pwd):/src" mid-health-check-svc sh -c "/usr/local/go/bin/go test ."

.PHONY: all test clean build-centos test-centos