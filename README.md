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

The latest version is 3.0.2, released June 6, 2025. See the [change log](./CHANGELOG.md) for details.

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.2/windows/amd64/apt-cmd.exe) | Intel 64-bit | v3.0.1 | ff998da3ac9dd555e3a58ce800f905b42ecc338bbb370cda3226b3ce1a19ddcb |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.2/windows/arm64/apt-cmd.exe) | ARM 64-bit | v3.0.1 | 419e41f40c485c61b4584ea28d91a02bf317fe54f94bbbab12156473ae7ace18 |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.2/mac/amd64/apt-cmd)  | Intel 64-bit | v3.0.1 | 886a0cb98a0f7e3374fa07b11e6646d45ce76bedaca911db64371d7f88aadba9 |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.2/mac/arm64/apt-cmd) | Apple Silicon (M series) | v3.0.1 | c24f34295ec2113fdc91db7b1c1c1b3faff0199bca8e0a6603541ec65e967204 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.2/linux/amd64/apt-cmd) | Intel 64-bit | v3.0.1 | ef53916c63a57470db264406323740cfbd23c353dbda18b43dbf8cf65d418cb7 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.2/linux/arm64/apt-cmd) | ARM 64-bit | v3.0.1 | 46b9d592a119633a686b995cda48f24b2302eadf3dabde003ecd7b0e1690b605 |

# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.rb units`

Integration tests: `./scripts/test.rb integration`

Note that when running integration tests, Registry tests do not run on Windows.

# Building

`./scripts/build.sh`
