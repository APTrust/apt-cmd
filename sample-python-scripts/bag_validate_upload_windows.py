import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC UPLOAD AND VALIDATE INFO BELOW                                            #

# NAME OF RECEIVING BUCKET
bucket_name = "aptrust.receiving.test.acr7d.edu"

# AWS KEYS
aws_access_key = ""
aws_secret_key = ""

output_dir = 'C:/Users/etd4sv/aptrust/test_bags/output_folder'

# BAG PROFILE TO USE WHEN CREATING
profile = "aptrust"

# WHAT MAINIFEST AND TAG MANIFESTS TO INCLUDE: i.e. 'md5,sha256'
manifest_algs = "md5,sha256"

# WHAT ORGANIZATION TO PUT IN THE BAG INFO
source_organization = "College"

# BAG STORAGE TYPE: i.e. 'Standard'
storage_option = "Standard"

jobs = [
        { "source_dir":"C:/Users/etd4sv/aptrust/test_bags/input_folder/test_bag_1", "title": "Bag 1", "access": "Institution" },
        { "source_dir":"C:/Users/etd4sv/aptrust/test_bags/input_folder/test_bag_2", "title": "Bag 2", "access": "Consortia" },
        { "source_dir":"C:/Users/etd4sv/aptrust/test_bags/input_folder/test_bag_3", "title": "Bag 2", "access": "Consortia" }
        ]

###############################################################################################################################

profile_full = "--profile=" + profile
manifest_algs = "--manifest-algs=" + manifest_algs
source_organization = '--tags=bag-info.txt/Source-Organization=' + source_organization
output_file = "--output-file=" + output_dir + "/"
storage_option = '--tags=aptrust-info.txt/Storage-Option=' + storage_option

os.environ['APTRUST_AWS_KEY'] = aws_access_key
os.environ['APTRUST_AWS_SECRET'] = aws_secret_key

for bag in jobs:
    title = '--tags=aptrust-info.txt/Title=' + bag["title"]
    access = '--tags=aptrust-info.txt/Access=' + bag["access"]
    bag_dir = "--bag-dir=" + bag["source_dir"]
    index_slash = bag["source_dir"][::-1].find('/')
    bag_name = bag["source_dir"][-index_slash::]

    create = subprocess.Popen(['apt-cmd.exe', 'bag', 'create', profile_full, manifest_algs, output_file + str(bag_name) + '.tar', bag_dir, source_organization, title, access, storage_option], shell=True, stdout=subprocess.DEVNULL)
    if create.returncode: 
        print("ERROR CREATING: {}".format(bag_name))
    else:
        print("Bagged: {}".format(bag_name))

    validate = subprocess.Popen(['apt-cmd.exe', 'bag', 'validate', '-p', profile, output_dir + "/" + bag_name + ".tar"], shell=True, stdout=subprocess.DEVNULL)
    if validate.returncode:
        print("ERROR VALIDATING: {}".format(bag_name))
        continue
    else:
        print("Validated: {}".format(bag_name))
    
    upload = subprocess.run(['apt-cmd.exe', 's3', 'upload', '--host=s3.amazonaws.com', '--bucket=' + str(bucket_name), output_dir + "/" + bag_name + ".tar"], shell=True, stdout=subprocess.DEVNULL)
    if upload.returncode: 
        print("ERROR UPLOADING: {}".format(bag_name))
    else:
        print("Uploaded: {}".format(bag_name))

    print("Bagged, Validated, Uploaded: {}\n".format(bag_name))