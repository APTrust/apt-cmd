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

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/windows/apt-cmd.exe) | Intel 64-bit | v3.0.0-beta | b12d7daf68ca2a2ea99ea208143e4571cf49fd8221ea1998a9b4f5db9774b631 |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/mac-intel/apt-cmd)  | Intel 64-bit | v3.0.0-beta | 1b5ceb015744e9ca818e5526f0940988fd4dad7f56c1bde105762bd89522265b |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/mac-arm/apt-cmd) | Apple Silicon (M series) | v3.0.0-beta | 0327e04b44137ce856b342542563133b9f8184364513394013bf200939dd6c8e |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.0-beta/linux/apt-cmd) | Intel 64-bit | v3.0.0-beta | 4c1937567c5a31752bad04147efd1577dc1e3995f8334f225b5b683709117f58 |

# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.rb units`

Integration tests: `./scripts/test.rb integration`

Note that when running integration tests, Registry tests do not run on Windows.

# Building

`./scripts/build.sh`
