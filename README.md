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

# Configuration (`/etc/mhc/mhc.conf`)

| Key                    | Description                                                                                                                     | Default             |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------|---------------------|
| MHC_TRAFFIC_CTL_DIR    | Directory where Traffic Server is installed.                                                                                    | /opt/trafficserver  |
| MHC_TM_HOSTS           | Comma separated FQDN of Traffic Monitor instances to connect to.                                                                |                     |
| MHC_TM_API_PATH        | API path for pulling cache stats from TM.                                                                                       | /api/cache-statuses |
| MHC_LOG_LEVEL          | See https://github.com/rs/zerolog#leveled-logging                                                                               | 1                   |
| MHC_TCP_CHECK_INTERVAL | Interval in seconds to perform a TCP check on each cache server (not implemented).                                              | 2                   |
| MHC_TM_CHECK_INTERVAL  | Interval in seconds to check data in TM for each cache server.                                                                  | 10                  |
| MHC_TO_CHECK_INTERVAL  | Interval in seconds to check Traffic Ops for status of each cache server.                                                       |                     |
| MHC_TO_HOSTNAME        | FQDN of the Traffic Ops server to connect to.                                                                                   |                     |
| MHC_TO_USERNAME        | Username used to log in to the Traffic Ops API.                                                                                 |                     |
| MHC_TO_PASSWORD        | Password used to log in to the Traffic Ops API.                                                                                 |                     |
| MHC_TO_INSECURE        | Determines if strict certificate checking should be ignored when connecting to Traffic Ops.                                     | FALSE               |
| MHC_TO_API_TIMEOUT     | Timeout in seconds when attempting to establish contact with Traffic Ops.                                                       | 10                  |
| MHC_USE_LOGFILE        | Determines if STDOUT and STDERR should also be written to a file at `/var/log/mhc.log`. Directory must exist for writes to occur. | TRUE                |
