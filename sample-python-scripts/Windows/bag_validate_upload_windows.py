import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC UPLOAD AND VALIDATE INFO BELOW                                            #

# NAME OF RECEIVING BUCKET
bucket_name = ""

# AWS KEYS
aws_access_key = ""
aws_secret_key = ""

output_dir = ""

# BAG PROFILE TO USE WHEN CREATING
profile = ""

# WHAT MAINIFEST AND TAG MANIFESTS TO INCLUDE: i.e. 'md5,sha256'
manifest_algs = ""

# WHAT ORGANIZATION TO PUT IN THE BAG INFO
source_organization = ""

# BAG STORAGE TYPE: i.e. 'Standard'
storage_option = ""

jobs = [
        { "source_dir":"", "title": "", "access": "" },
        { "source_dir":"", "title": "", "access": "" }
]

###############################################################################################################################

profile_full = "--profile=" + profile
manifest_algs = "--manifest-algs=" + manifest_algs
source_organization = '--tags=bag-info.txt/Source-Organization=' + source_organization
output_file = "--output-file=" + output_dir + "\\"
storage_option = '--tags=aptrust-info.txt/Storage-Option=' + storage_option

os.environ['APTRUST_AWS_KEY'] = aws_access_key
os.environ['APTRUST_AWS_SECRET'] = aws_secret_key

for bag in jobs:
    title = '--tags=aptrust-info.txt/Title=' + bag["title"]
    access = '--tags=aptrust-info.txt/Access=' + bag["access"]
    bag_dir = "--bag-dir=" + bag["source_dir"]
    index_slash = bag["source_dir"][::-1].find('\\')
    bag_name = bag["source_dir"][-index_slash::]

    create = subprocess.call(['apt-cmd.exe', 'bag', 'create', profile_full, manifest_algs, output_file + str(bag_name) + '.tar', bag_dir, source_organization, title, access, storage_option], shell=True, stdout=subprocess.DEVNULL)
    if create: 
        print("ERROR CREATING: {}".format(bag_name))
        continue
    else:
        print("Bagged: {}".format(bag_name))

    validate = subprocess.call(['apt-cmd.exe', 'bag', 'validate', '-p', profile, output_dir + "\\" + bag_name + ".tar"], shell=True, stdout=subprocess.DEVNULL)
    if validate:
        print("ERROR VALIDATING: {}".format(bag_name))
        continue
    else:
        print("Validated: {}".format(bag_name))
    
    upload = subprocess.call(['apt-cmd.exe', 's3', 'upload', '--host=s3.amazonaws.com', '--bucket=' + str(bucket_name), output_dir + "\\" + bag_name + ".tar"], shell=True, stdout=subprocess.DEVNULL)
    if upload: 
        print("ERROR UPLOADING: {}".format(bag_name))
    else:
        print("Uploaded: {}".format(bag_name))

    print("Bagged, Validated, Uploaded: {}\n".format(bag_name))