import os
import subprocess
import sys

###############################################################################################################################
#                                               FILL IN SPECIFIC VALIDATE INFO BELOW                                          #

# BAG PROFILE TO VALIDATE WITH
profile = "aptrust"

# LIST PATH OF BAGS TO BE VALIDATED
bags = [
        { "bag_path":"C:/Users/etd4sv/aptrust/test_bags/input_folder/test_bag_1"}
        
]

###############################################################################################################################

for bag in bags:
    validate = subprocess.Popen(['apt-cmd.exe', 'bag', 'validate', '-p', profile, bag["bag_path"]], shell=True, stdout=subprocess.DEVNULL)
    if validate.returncode:
        print("ERROR VALIDATING: {}".format(bag["bag_path"]))
    else:
        print("Validated: {}".format(bag["bag_path"]))