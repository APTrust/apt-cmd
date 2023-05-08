# Sample Python Scripts for Using apt-cmd

This folder contains several examples of python3 scripts that utilize the apt-cmd command line tool. Examples are present for bagging multiple items, validating multiple bags, uploading multiple bags, validating and uploading multiple bags in the same script.

## bag_multiple.py

Takes in:
* Path to where the command line tool is: ```path_to_where_apt_cmd_is```
* Path the raw folders to be bagged: ```path_to_raw_bags```
* Location to place the created bags: ```path_to_created_bags```
* Bag profile to use: ```profile```
* What manifests and tag manifests to use: ```manifest_algs```
* What organization belongs in bag info: ```source_organization```
* Name of the bag: ```title```
* Access type for the bag: ```access```
* Storage type of the bag: ```storage_option```

Upon running this script, any folder within the chosen directory will be bagged and placed in the directory specified.

## validate_multiple.py

Takes in:
* Path to where the command line tool is: ```path_to_where_apt_cmd_is```
* Path to the bags to be validated: ```path_to_bags```
* Bag profile to use: ```profile```

Upon running this script, any file ending in .tar within the chosen directory will be validated using the profile specified.

## upload_multiple.py

Takes in:
* Name of the S3 receiving bucket to upload to: ```bucket_name```
* AWS access key of the user: ```aws_access_key```
* AWS secret key of the user: ```aws_secret_key```
* Path to where the command line tool is: ```path_to_where_apt_cmd_is```
* Path to the bags to be uploaded: ```path_to_bags```

Upon running this script, any file ending in .tar within the chosen directory will be uploaded to the specified bucket using the credentials specified.

## validate_and_upload_multiple.py

Takes in:
* Name of the S3 receiving bucket to upload to: ```bucket_name```
* AWS access key of the user: ```aws_access_key```
* AWS secret key of the user: ```aws_secret_key```
* Path to where the command line tool is: ```path_to_where_apt_cmd_is```
* Path to the bags to be validated and uploaded: ```path_to_bags```
* Bag profile to validate with: ```profile```

Upon running this script, any file ending in .tar within the chosen directory will be validated with the given profile. If the file is validated successfully, it will then be uploaded to the specified S3 receiving bucket using the AWS keys provided.