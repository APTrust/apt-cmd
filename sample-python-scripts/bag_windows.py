import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC BAG INFO BELOW                                               #

# LIST PATH OF WHERE TO PLACE CREATED BAGS
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