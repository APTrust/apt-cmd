import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC VALIDATE INFO BELOW                                          #

# THIS IS THE PATH TO WHERE THE PARTNER TOOLS apt-cmd EXECUTABLE IS
path_to_where_apt_cmd_is = ""

# LIST PATH OF BAGS TO BE VALIDATED
path_to_bags = ""

# BAG PROFILE TO VALIDATE WITH
bag_profile = ""

###############################################################################################################################

files_in_directory = os.listdir(path_to_bags)
os.chdir(path_to_where_apt_cmd_is)

for file in files_in_directory:
    if file.endswith('.tar'):
        print("\nValidating {}".format(file))
        validate_command = './apt-cmd bag validate -p ' + str(bag_profile) + ' ' + path_to_bags + "/" + file
        subprocess.run(validate_command, shell=True)

print("\nDone validating bags")
