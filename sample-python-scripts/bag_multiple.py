import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC BAG INFO BELOW                                               #

# THIS IS THE PATH TO WHERE THE PARTNER TOOLS apt-cmd EXECUTABLE IS
path_to_where_apt_cmd_is = ""

# LIST PATH OF BAGS TO BE UPLOADED
path_to_raw_bags = ""

# LIST PATH OF WHERE TO PLACE CREATED BAGS
path_to_created_bags = ""

# BAG PROFILE TO USE WHEN CREATING
profile = ""

# WHAT MAINIFEST AND TAG MANIFESTS TO INCLUDE: i.e. 'md5,sha256'
manifest_algs = ""

# WHAT ORGANIZATION TO PUT IN THE BAG INFO
source_organization = ""

# NAME OF BAG
title = ""

# BAG ACCESS TYPE: i.e. 'Institution'
access = ""

# BAG STORAGE TYPE: i.e. 'Standard'
storage_option = ""

###############################################################################################################################

profile = " --profile=" + profile
manifest_algs = " --manifest-algs=" + manifest_algs
source_organization = " --tags='bag-info.txt/Source-Organization=" + source_organization + "'"
output_file = " --output-file=" + path_to_created_bags + "/"
bag_dir = " --bag-dir=" + path_to_raw_bags + "/"
title = " --tags='aptrust-info.txt/Title=" + title + "'"
access = " --tags='aptrust-info.txt/Access=" + access + "'"
storage_option = " --tags='aptrust-info.txt/Storage-Option=" + storage_option + "'"

files_in_directory = os.listdir(path_to_raw_bags)
os.chdir(path_to_where_apt_cmd_is)

for file in files_in_directory:
    os.chdir(path_to_raw_bags)
    if os.path.isdir(file):
        os.chdir(path_to_where_apt_cmd_is)
        create_command = "./apt-cmd bag create" + profile + manifest_algs + output_file + file + ".tar" + bag_dir + file + source_organization + title + access + storage_option
        subprocess.run(create_command, shell=True)
    
print("\nDone creating bags.")