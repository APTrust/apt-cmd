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
    validate = subprocess.call(['apt-cmd.exe', 'bag', 'validate', '-p', profile, bag["bag_path"]], shell=True, stdout=subprocess.DEVNULL)
    if validate:
        print("ERROR VALIDATING: {}".format(bag["bag_path"]))
    else:
        print("Validated: {}".format(bag["bag_path"]))