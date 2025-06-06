# Change Log for apt-cmd

For the user guide, see https://aptrust.github.io/userguide/partner_tools/.

## [3.0.2] 2025-06-05

* Storage-Option tag is no longer required in APTrust BagIt profile v2.3, to match v2.2.

## [3.0.1] 2025-06-05

* Fixed omissions in the algorithms lists for allowed manifests and allowed tag manifests in version 2.3 of the APTrust BagIt profile. This list now correctly includes sha1 and sha512 as allowed algorithms.

## [3.0.0] 2025-04-23

* Added support for APTrust BagIt Profile [version 2.3](https://github.com/APTrust/preservation-services/blob/master/profiles/aptrust-v2.3.json), which is identical to the [2.2 profile](https://github.com/APTrust/preservation-services/blob/master/profiles/aptrust-v2.2.json) except for the addition of Wasabi-TX as a valid storage option.
* Version 2.3 is now the default APTrust profile for apt-cmd. All bags that are valid according to version 2.2 are also valid according to version 2.3.
* While the flag `--profile=aptrust` now indicates the version 2.3 profile, you can force apt-cmd to use the version 2.2 profile by specifying `--profile=aptrust-2.2`

## [3.0.0-beta] 2023-04-01

* Initial release supporting APTrust BagIt Profile version 2.2 and BTR BagIt profile.