[![Go](https://github.com/ARMmaster17/atc-mid-cache-health-check/actions/workflows/build.yml/badge.svg)](https://github.com/ARMmaster17/atc-mid-cache-health-check/actions/workflows/build.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/2b02133a5f8bd7909fb3/maintainability)](https://codeclimate.com/github/ARMmaster17/atc-mid-cache-health-check/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2b02133a5f8bd7909fb3/test_coverage)](https://codeclimate.com/github/ARMmaster17/atc-mid-cache-health-check/test_coverage)

# atc-mid-cache-health-check

Automated service to disable unreachable mid cache servers in an
Apache Traffic Control stack.

# Installing

Navigate to the GitHub Actions tab of this repository and download
the RPM zip. Run the following to install:
```bash
unzip "RPM Packags.rpm"
sudo dnf install mhc-1.0-1.el8.x86_64.rpm
sudo systemctl start mhc
```

# Local Development

Build the project with `make build`. Test with `make test`. Build
the RPM with `make build-centos` (requires Docker).  
