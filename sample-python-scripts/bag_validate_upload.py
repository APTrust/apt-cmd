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

output_dir = "/home/emory/aptrust/scripts/bag_dir"

# BAG PROFILE TO USE WHEN CREATING
profile = "aptrust"

# WHAT MAINIFEST AND TAG MANIFESTS TO INCLUDE: i.e. 'md5,sha256'
manifest_algs = "md5,sha256"

# WHAT ORGANIZATION TO PUT IN THE BAG INFO
source_organization = "College"

# BAG STORAGE TYPE: i.e. 'Standard'
storage_option = "Standard"

# THIS IS THE PATH TO WHERE THE PARTNER TOOLS apt-cmd EXECUTABLE IS
path_to_where_apt_cmd_is = "/home/emory/aptrust/scripts"


jobs = [
        { "source_dir":"/home/emory/aptrust/scripts/test_bag", "title": "Bag 1", "access": "Institution" },
        { "source_dir":"/home/emory/aptrust/scripts/test2_bag", "title": "Bag 2", "access": "Consortia" }
        ]

###############################################################################################################################

profile_full = " --profile=" + profile
manifest_algs = " --manifest-algs=" + manifest_algs
source_organization = " --tags='bag-info.txt/Source-Organization=" + source_organization + "'"
output_file = " --output-file=" + output_dir + "/"
storage_option = " --tags='aptrust-info.txt/Storage-Option=" + storage_option + "'"

os.environ['APTRUST_AWS_KEY'] = aws_access_key
os.environ['APTRUST_AWS_SECRET'] = aws_secret_key

for bag in jobs:
    os.chdir(path_to_where_apt_cmd_is)
    title = " --tags='aptrust-info.txt/Title=" + bag["title"] + "'"
    access = " --tags='aptrust-info.txt/Access=" + bag["access"] + "'"
    bag_dir = " --bag-dir=" + bag["source_dir"]
    index_slash = bag["source_dir"][::-1].find('/')
    bag_name = bag["source_dir"][-index_slash::]
    create_command = "./apt-cmd bag create" + profile_full + manifest_algs + output_file + bag_name + ".tar" + bag_dir + source_organization + title + access + storage_option
    subprocess.run(create_command, shell=True)
    validate_command = "./apt-cmd bag validate -p " + profile + " " + output_dir + "/" + bag_name + ".tar"
    subprocess.run(validate_command, shell=True)
    upload_command = './apt-cmd s3 upload --host=s3.amazonaws.com --bucket="' + str(bucket_name) + '" ' + output_dir + "/" + bag_name + ".tar"
    subprocess.run(upload_command, shell=True)
