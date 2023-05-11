import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC UPLOAD INFO BELOW                                            #

# NAME OF RECEIVING BUCKET
bucket_name = "aptrust.receiving.test.acr7d.edu"

# AWS KEYS
aws_access_key = ""
aws_secret_key = ""

# LIST PATH OF BAGS TO BE UPLOADED
bags = [
        { "bag_path":"/home/emory/aptrust/testing_files/testing_output/test_bag_2.tar"}
]

###############################################################################################################################

os.environ['APTRUST_AWS_KEY'] = aws_access_key
os.environ['APTRUST_AWS_SECRET'] = aws_secret_key

for bag in bags:
        upload_command = 'apt-cmd s3 upload --host=s3.amazonaws.com --bucket="' + str(bucket_name) + '" ' + bag["bag_path"]
        upload = subprocess.run(upload_command, shell=True, capture_output=True, text=True)
        if upload.returncode: print("ERROR UPLOADING: {}".format(bag["bag_path"]))

        print("Uploaded: {}".format(bag["bag_path"]))