import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC UPLOAD INFO BELOW                                            #

# NAME OF RECEIVING BUCKET
bucket_name = ""

# AWS KEYS
aws_access_key = ""
aws_secret_key = ""

# THIS IS THE PATH TO WHERE THE PARTNER TOOLS apt-cmd EXECUTABLE IS
path_to_where_apt_cmd_is = ""

# LIST PATH OF BAGS TO BE UPLOADED
path_to_bags = ""

###############################################################################################################################

files_in_directory = os.listdir(path_to_bags)
os.environ['APTRUST_AWS_KEY'] = aws_access_key
os.environ['APTRUST_AWS_SECRET'] = aws_secret_key
os.chdir(path_to_where_apt_cmd_is)

for file in files_in_directory:
    if file.endswith('.tar'):
        print("\nUploading {}".format(file))
        upload_command = './apt-cmd s3 upload --host=s3.amazonaws.com --bucket="' + str(bucket_name) + '" ' + path_to_bags + "/" + file
        subprocess.run(upload_command, shell=True)

print("\nDone uploading bags")