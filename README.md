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

The latest version is 3.0.5, released October 28, 2026. See the [change log](./CHANGELOG.md) for details.

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.5/windows/amd64/apt-cmd.exe) | Intel 64-bit | v3.0.5 | 6526d0d24b65a285966c663eac64b8a626d2581bc64dc955ec725100ff8c1a95 |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.5/windows/arm64/apt-cmd.exe) | ARM 64-bit | v3.0.5 | 0d61d7fba302ebe4564a36a3d190f07de478e94d345167042fdcd367eedd65bb |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.5/mac/amd64/apt-cmd)  | Intel 64-bit | v3.0.5 | ec6cb5e0c0baa5693891417cbb5ee70466f6e66d5d8b3e566a9be36b3ff231e1 |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.5/mac/arm64/apt-cmd) | Apple Silicon (M series) | v3.0.5 | 214a1a615e347fcfbfa036c58f204ed482fa3e49a3458fae902c1e7eadef48b4 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.5/linux/amd64/apt-cmd) | Intel 64-bit | v3.0.5 | ff5a629f2f0de22b125af1328656892f0315e58c7761d4e332290dd1a0daa626 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.5/linux/arm64/apt-cmd) | ARM 64-bit | v3.0.5 | d4be8cf9fa5f29e2bf8bee78a3d7897218a4ea006617d699c8da430a12c27ff5 |

# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.sh units`

Integration tests: `./scripts/test.sh integration`

Note that when running integration tests, Registry tests do not run on Windows.

# Building

`./scripts/build.sh`

## Releasing

To create a new release:

1. Tag the release with a new version number using `git tag -a <version>`
2. Make sure your APTrust AWS credentials are in your environment. See ./scripts/release.sh for the environment variable names.
3. Run `./scripts/release.sh --pre-check <version>`
4. If all is good, run `./scripts/release.sh <version>`
5. Update links and checksums as instructed by release.sh.
6. Test the new release links to ensure that downloads work.
