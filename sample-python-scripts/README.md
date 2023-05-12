# Sample Python Scripts for Using apt-cmd

This folder contains several examples of python3 scripts that utilize the apt-cmd command line tool. Examples are present for bagging multiple items, validating multiple bags, uploading multiple bags, and performing all three actions in the same script.

## Prerequisite for scripts

To run any of theses scripts, ```apt-cmd``` must be added to your systems path so that the system can access the executable. This can be done in different ways but here are straightforward ways in unix/mac and windows:

### Unix/Mac

1. Open Terminal
2. Create directory ```${HOME}/bin``` --- ```mkdir -p ${HOME}/bin```
3. Copy apt-cmd into the directory of ```${HOME}/bin```
4. Make the binary executable ```chmod 755 ${HOME}/bin/apt-cmd```
5. Open the shell config file:
    * Mac - Open the shell config file 
    * Unix - Open file ```${HOME}/.bashrc```
6. Add this line to the file ```export PATH="${HOME}/bin:${PATH}"```
7. Restart your Terminal

You can verify the binary is on the path by running ```command -v apt-cmd```

### Windows

1. Create folder ```C:\bin```
2. Save apt-cmd.exe in ```C:\bin```
3. Depending on your Windows version
    * Windows 8 or 10/11 - press the Windows key, then search for and select System (Control Panel)
    * Windows 7 - right click the Computer icon on the desktop and click Properties
4. Click Advanced system settings
5. Click Environment Variables
6. Under System Variables, find the PATH variable, select it, and click Edit. If there is no PATH variable, click New
7. Add ```C:\bin``` to the start of the variable value, followed by a ```;```. For example, if the value was ```C:\Windows\System32```, change it to ```C:\bin;C:\Windows\System32```
8. Click OK

You can verify the binary is on the path by running ```where.exe apt-cmd.exe```

# Script Specific Information

### bag_(linux/windows).py

Takes in:
* Location to place the created bags: ```output_dir```
* Bag profile to use: ```profile```
    * Allowed Values: aptrust, btr, empty
* What manifests and tag manifests to use: ```manifest_algs```
    * Specify one, or use comma-separated list for multiple. Supported algorithms: md5, sha1, sha256, sha512. Default is sha256.
* What organization belongs in bag info: ```source_organization```
* Storage type of the bag: ```storage_option```
    * Allowed Values: Standard, Glacier-OH, Glacier-OR, Glacier-VA, Glacier-Deep-OH, Glacier-Deep-OR, Glacier-Deep-VA
* Bag specific details, in the form of a python dictionary with structure:  ```jobs = [{ "source_dir":"", "title": "", "access": "" }]```
    * ```source_dir``` - path to directory to be bagged
    * ```title``` - title of bag to create
    * ```access``` - access type of the bag
        * Allowed Values: Restricted, Institution, Consortia

Upon running this script the bags listed in the ```jobs``` directory will be created and placed in ```output_dir```.

### validate_(linux/windows).py

Takes in:
* Bag profile to use: ```profile```
    * Allowed Values: aptrust, btr, empty
* Bags to validate in the form of a python dictionary: ```bags = [{ "bag_path":"/path/to/bag.tar"}]```

Upon running this script files in the ```bags``` directory will be validated using the ```profile``` specified.

### upload_(linux/windows).py

Takes in:
* Name of the S3 receiving bucket to upload to: ```bucket_name```
* AWS access key of the user: ```aws_access_key```
* AWS secret key of the user: ```aws_secret_key```
* Bags to upload in the form of a python dictionary: ```bags = [{ "bag_path":"/path/to/bag.tar"}]```

Upon running this script files in the ```bags``` directory will be uploaded to the specified bucket using the credentials specified.

### bag_validate_upload_(linux/windows).py

Takes in:
* Name of the S3 receiving bucket to upload to: ```bucket_name```
* AWS access key of the user: ```aws_access_key```
* AWS secret key of the user: ```aws_secret_key```
* Location to place the created bags: ```output_dir```
* Bag profile to use: ```profile```
    * Allowed Values: aptrust, btr, empty
* What manifests and tag manifests to use: ```manifest_algs```
    * Specify one, or use comma-separated list for multiple. Supported algorithms: md5, sha1, sha256, sha512. Default is sha256.
* What organization belongs in bag info: ```source_organization```
* Storage type of the bag: ```storage_option```
    * Allowed Values: Standard, Glacier-OH, Glacier-OR, Glacier-VA, Glacier-Deep-OH, Glacier-Deep-OR, Glacier-Deep-VA
* Bag specific details, in the form of a python dictionary with structure:  ```jobs = [{ "source_dir":"", "title": "", "access": "" }]```
    * ```source_dir``` - path to directory to be bagged
    * ```title``` - title of bag to create
    * ```access``` - access type of the bag
        * Allowed Values: Restricted, Institution, Consortia

This script will create bags specified in the ```jobs``` dictionary and place them in ```output_dir```. It then validates these bags, and if they are valid it will upload them to the specified bucket.