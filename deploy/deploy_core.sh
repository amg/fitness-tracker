#!/bin/bash

helpFunction()
{
   echo "Core functionality such as SQL DB and domain mapping takes a long time to recreate"
   echo "This only likely needs to be done once"
   echo "Using global TF_VARs to align some of the dependencies"
   echo "Usage: $0 -m plan/apply/destroy"
   echo -e "\t-m Mode for Terraform: plan or apply"
   echo -e "\t-a authenticate docker to gcloud"
   exit 1 # Exit script after printing help
}

authenticate()
{
    echo "Authenticating:"
    gcloud auth configure-docker
    exit 1 # Exit script after authenticating
}

while getopts "m:a:" opt
do
   case "$opt" in
      m ) mode="$OPTARG" ;;
      a ) authenticate ;;
      ? ) helpFunction ;; # Print helpFunction in case parameter is non-existent
   esac
done

# Print helpFunction in case parameters are empty
if [ -z "$mode" ]
then
   echo "Mode needs to be specified";
   helpFunction
fi

# Begin script in case all parameters are correct
echo "Using mode: $mode"

# Jump into core/ subfolder
cd core/

eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' ../../.secrets/.deploy.env) \
&& eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' ../common/.variables_resources) \
&& terraform $mode
