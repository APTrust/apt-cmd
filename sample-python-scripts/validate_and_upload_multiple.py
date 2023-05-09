import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC UPLOAD AND VALIDATE INFO BELOW                               #

# NAME OF RECEIVING BUCKET
bucket_name = ""

# AWS KEYS
aws_access_key = ""
aws_secret_key = ""

# THIS IS THE PATH TO WHERE THE PARTNER TOOLS apt-cmd EXECUTABLE IS
path_to_where_apt_cmd_is = ""

# LIST PATH OF BAGS TO BE UPLOADED
path_to_bags = ""

# BAG PROFILE TO VALIDATE WITH
bag_profile = ""

###############################################################################################################################

files_in_directory = os.listdir(path_to_bags)
os.environ['APTRUST_AWS_KEY'] = aws_access_key
os.environ['APTRUST_AWS_SECRET'] = aws_secret_key
os.chdir(path_to_where_apt_cmd_is)

for file in files_in_directory:
    if file.endswith('.tar'):
        print("\nValidating {}".format(file))
        validate_command = './apt-cmd bag validate -p ' + str(bag_profile) + ' ' + path_to_bags + "/" + file
        proc = subprocess.Popen(validate_command, shell=True, stdout=subprocess.PIPE)
        validate_result = 'error' in str(proc.stdout.readline())
        if validate_result:
            print("Error validating " + str(file) + " bag will not be uploaded.")
            continue
        print("Bag Validated, Uploading {}".format(file))
        upload_command = './apt-cmd s3 upload --host=s3.amazonaws.com --bucket="' + str(bucket_name) + '" ' + path_to_bags + "/" + file
        subprocess.run(upload_command, shell=True)

print("\nDone uploading bags")
