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
        { "bag_path":"C:/Users/etd4sv/aptrust/test_bags/output_folder/test_bag_1.tar"}
]

###############################################################################################################################

os.environ['APTRUST_AWS_KEY'] = aws_access_key
os.environ['APTRUST_AWS_SECRET'] = aws_secret_key

for bag in bags:
        upload = subprocess.run(['apt-cmd.exe', 's3', 'upload', '--host=s3.amazonaws.com', '--bucket=' + str(bucket_name), bag["bag_path"]], shell=True, stdout=subprocess.DEVNULL)
        if upload.returncode: 
            print("ERROR UPLOADING: {}".format(bag["bag_path"]))
        else:
            print("Uploaded: {}".format(bag["bag_path"]))