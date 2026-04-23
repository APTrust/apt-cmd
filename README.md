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

The latest version is 3.0.4, released April 23, 2026. See the [change log](./CHANGELOG.md) for details.

| Platform | Architecture | Version | SHA-256 |
| -------- | ------------ | ------- | ------- |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.4/windows/amd64/apt-cmd.exe) | Intel 64-bit | v3.0.4 | 8f2785f5c88fe2839d2ec95ce31a758163d4cd6808e4a0600fbfa8fa76079489 |
| [Windows](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.4/windows/arm64/apt-cmd.exe) | ARM 64-bit | v3.0.4 | c25b677a6206fb59e9cd6dced4a272b375484c9b5c46faa87b77c37bb1a2cf23 |
| [Mac Intel](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.4/mac/amd64/apt-cmd)  | Intel 64-bit | v3.0.4 | 8c8a37b0f56d1a1618b773624586961a5c235dd527dca61b6eb9cb5570038860 |
| [Mac ARM](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.4/mac/arm64/apt-cmd) | Apple Silicon (M series) | v3.0.4 | f9f1e1f2a02eba3871c366b15ed75cd7b56f4ae654ae5ac7105b618e8c928da8 |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.4/linux/amd64/apt-cmd) | Intel 64-bit | v3.0.4 | 95fcc4cdc89a416c8fd472874e817e40b7d3d3344823606a5704bbfa4dfcca7c |
| [Linux](https://s3.amazonaws.com/aptrust.public.download/apt-cmd/v3.0.4/linux/arm64/apt-cmd) | ARM 64-bit | v3.0.4 | dc2d69d22f70abd339813cfc97afb6f84f5f402a0c62e8ab27c7cbb09a2bd4dd |


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
