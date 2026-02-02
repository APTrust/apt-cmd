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

The latest version is 3.0.3, released February 3, 2026. See the [change log](./CHANGELOG.md) for details.

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.3/windows/amd64/apt-cmd.exe) | Intel 64-bit | v3.0.3 | f039bb5ed01a4dd53a9d07ebab6b3440b13042a6ae71ab73361b2c27d4b3269e |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.3/windows/arm64/apt-cmd.exe) | ARM 64-bit | v3.0.3 | f0a04c2a276d6aca5dfc94805f0fbb499176edb9e3c60243d3b4e561e86e1c54 |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.3/mac/amd64/apt-cmd)  | Intel 64-bit | v3.0.3 | 648c3dc7bc67126e7a7aae108f3808d9d6881632989e0efc032418ca469ddb2d |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.3/mac/arm64/apt-cmd) | Apple Silicon (M series) | v3.0.3 | 7559817ba00da694d01501c53f4226f4c6e5b792f7750bc2b5294ecaeb0b5069 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.3/linux/amd64/apt-cmd) | Intel 64-bit | v3.0.3 | f7ed19dc5d2e06226baa3a733667090020994f184149d44a8016ad2e36b3efeb |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.3/linux/arm64/apt-cmd) | ARM 64-bit | v3.0.3 | 30a2380b10e5e5d4799982229bd345d4f2855ed0dcdaa00ce1d2c040e482b1dd |


# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.rb units`

Integration tests: `./scripts/test.rb integration`

Note that when running integration tests, Registry tests do not run on Windows.

# Building

`./scripts/build.sh`
