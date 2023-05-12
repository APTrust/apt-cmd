import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC VALIDATE INFO BELOW                                          #

# BAG PROFILE TO VALIDATE WITH
profile = ""

# LIST PATH OF BAGS TO BE VALIDATED
bags = [
        { "bag_path":"/path/to/bag.tar"},
        { "bag_path":"/path/to/bag.tar"}
]

###############################################################################################################################

for bag in bags:
    validate_command = "apt-cmd bag validate -p " + profile + " " + bag["bag_path"]
    validate = subprocess.call(validate_command, shell=True, capture_output=True, text=True)
    if validate: 
        print("ERROR VALIDATING: {}".format(bag["bag_path"]))
    else:
        print("Validated: {}".format(bag["bag_path"]))