import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC VALIDATE INFO BELOW                                          #

# BAG PROFILE TO VALIDATE WITH
profile = "aptrust"

# LIST PATH OF BAGS TO BE VALIDATED
bags = [
        { "bag_path":"/home/emory/aptrust/testing_files/testing_output/test_bag_1.tar"}
]

###############################################################################################################################

for bag in bags:
    validate_command = "apt-cmd bag validate -p " + profile + " " + bag["bag_path"]
    validate = subprocess.run(validate_command, shell=True, capture_output=True, text=True)
    if validate.returncode: print("ERROR VALIDATING: {}".format(bag["bag_path"]))

    print("Bagged: {}".format(bag["bag_path"]))