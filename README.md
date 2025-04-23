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

The latest version is 3.0.0, released April 23, 2025. See the [change log](./CHANGELOG.md) for details.

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0/windows/amd64/apt-cmd.exe) | Intel 64-bit | v3.0.0 | 0aa6629c275c4780ad031568c4a6f8f5e929e1100dca40dd40a3937262d924ef |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0/windows/arm64/apt-cmd.exe) | ARM 64-bit | v3.0.0 | 216812412f027e8cff4919da8c2f10ab8ca56aa53f1f5f284cf514d29d804a3a |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0/mac/amd64/apt-cmd)  | Intel 64-bit | v3.0.0 | c2356d9486a77530c561011d34717cd91895316788e474b82a5a674d144b5330 |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0/mac/arm64/apt-cmd) | Apple Silicon (M series) | v3.0.0 | 4eb060b1669887837801d70a42856e17a948cc3fb9db58dc0a8463ab575fe8a0 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0/linux/amd64/apt-cmd) | Intel 64-bit | v3.0.0 | a3c78535987026080512a9626f8d7a7fd47afc7e9178d03334edb205086780cf |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0/linux/arm64/apt-cmd) | ARM 64-bit | v3.0.0 | 2fff16d7ab063fbfaee8a47662c77ff0420525de986f375f6532b844fa9df35b |

# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.rb units`

Integration tests: `./scripts/test.rb integration`

Note that when running integration tests, Registry tests do not run on Windows.

# Building

`./scripts/build.sh`
