#!/bin/bash

helpFunction()
{
   echo ""
   echo "Usage: $0 -m plan/apply"
   echo -e "\t-m Mode for Terraform: plan/apply/destroy"
   echo ""
   echo -e "\t(auth) gcloud auth configure-docker"
   exit 1 # Exit script after printing help
}

while getopts "m:" opt
do
   case "$opt" in
      m ) mode="$OPTARG" ;;
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

eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' ../.secrets/.deploy.env) \
&& eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' ./common/.variables_resources) \
&& terraform $mode
