# APTrust Command Line Tool

The aptrust command-line utility (apt-cmd) enables you to create and validate
bags, manage S3 files, and query data in the APTrust Registry. They replace
the older version 2.x partner tools which did not have bag creation features
and worked only with Pharos.

You can create workflows with apt-cmd that include:

* creating a bag
* validating the bag
* uploading the bag to an S3 bucket
* checking the APTrust registry to see when the bag was ingested
* retrieving object and file details, including checksums and PREMIS events,
from the APTrust registry

For bag creating and validation, the current release supports tarred bags
only. Though access to the APTrust registry is limited to APTrust depositors,
anyone can use apt-cmd's bagging and S3 features.

# Downloads

The latest version is 3.0.0, released April 23, 2025.

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/windows/apt-cmd.exe) | ARM 64-bit | v3.0.0 |  |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/windows/apt-cmd.exe) | Intel 64-bit | v3.0.0 |  |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/mac-intel/apt-cmd)  | Intel 64-bit | v3.0.0 |  |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/mac-arm/apt-cmd) | Apple Silicon (M series) | v3.0.0 |  |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/linux/apt-cmd) | ARM 64-bit | v3.0.0 |  |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/linux/apt-cmd) | Intel 64-bit | v3.0.0 |  |

# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.rb units`

Integration tests: `./scripts/test.rb integration`

Note that when running integration tests, Registry tests do not run on Windows.

# Building

`./scripts/build.sh`
