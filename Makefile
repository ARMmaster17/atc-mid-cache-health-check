RPM = mid-health-check.rpm
EXEC_FILE = mid-health-check
SRC_FILE = main.go
JUNIT_REPORT = report.xml

all: $(RPM)

$(RPM): clean $(EXEC_FILE)
	# TODO: Make the RPM

$(EXEC_FILE): test
	go build $(SRC_FILE)

test:
	go get -u github.com/jstemmer/go-junit-report
	go test -v 2>&1 | go-junit-report > $(JUNIT_REPORT)

clean:
	rm -f $(RPM)
	rm -f $(EXEC_FILE)

.PHONY: all test clean