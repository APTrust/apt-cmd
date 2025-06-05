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

The latest version is 3.0.1, released June 5, 2025. See the [change log](./CHANGELOG.md) for details.

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.1/windows/amd64/apt-cmd.exe) | Intel 64-bit | v3.0.1 | f22212ce64245a167af265d46e07e86c6ac70af34e116bf788e2480bee87e23c |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.1/windows/arm64/apt-cmd.exe) | ARM 64-bit | v3.0.1 | c3a1f335acdbac5723c0fc6ada3cecc23fb53b5b98fda477e073e416318759bd |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.1/mac/amd64/apt-cmd)  | Intel 64-bit | v3.0.1 | 7d7ce0d9685361685a261e2ba729bc058d59b5f5c48ef19b4f59ce9f6dfd7655 |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.1/mac/arm64/apt-cmd) | Apple Silicon (M series) | v3.0.1 | 766d7a2d765d53064dcd0b4db3b6d1d320b7f681da50e8e7ab561333f4fa546d |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.1/linux/amd64/apt-cmd) | Intel 64-bit | v3.0.1 | 83ecf887f253c6253c149500cc843fd40cce303a134213ce01f123dd0f3c4740 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.1/linux/arm64/apt-cmd) | ARM 64-bit | v3.0.1 | b40a8f6f7ba783e3019c00a47ea30e238b26a4cfa7c22c5559800db6f4b0895f |

# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.rb units`

Integration tests: `./scripts/test.rb integration`

Note that when running integration tests, Registry tests do not run on Windows.

# Building

`./scripts/build.sh`
