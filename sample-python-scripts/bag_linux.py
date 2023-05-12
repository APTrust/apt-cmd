import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC BAG INFO BELOW                                               #

# LIST PATH OF WHERE TO PLACE CREATED BAGS
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

profile_full = " --profile=" + profile
manifest_algs = " --manifest-algs=" + manifest_algs
source_organization = " --tags='bag-info.txt/Source-Organization=" + source_organization + "'"
output_file = " --output-file=" + output_dir + "/"
storage_option = " --tags='aptrust-info.txt/Storage-Option=" + storage_option + "'"

for bag in jobs:
    title = " --tags='aptrust-info.txt/Title=" + bag["title"] + "'"
    access = " --tags='aptrust-info.txt/Access=" + bag["access"] + "'"
    bag_dir = " --bag-dir=" + bag["source_dir"]
    index_slash = bag["source_dir"][::-1].find('/')
    bag_name = bag["source_dir"][-index_slash::]

    create_command = "apt-cmd bag create" + profile_full + manifest_algs + output_file + bag_name + ".tar" + bag_dir + source_organization + title + access + storage_option
    create = subprocess.call(create_command, shell=True, capture_output=True, text=True)
    if create:
        print("ERROR CREATING: {}".format(bag_name))
    else:
        print("Bagged: {}".format(bag_name))