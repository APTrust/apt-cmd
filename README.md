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

| Platform | Architecture | Download URL |
| -------- | ------------ | ------------ |
| Windows | Intel 64-bit | Coming soon... |
| Mac  | Intel 64-bit | Coming soon... |
| Mac  | Apple Silicon (M series) | Coming soon... |
| Linux | Intel 64-bit | Coming soon... |

# Documentation

You'll find extensive documentation and usage examples at https://aptrust.github.io/userguide/partner_tools/

# Testing

Unit tests: `./scripts/test.rb units`

Integration tests: `./scripts/test.rb integration`

Note that when running integration tests, Registry tests do not run on Windows.
