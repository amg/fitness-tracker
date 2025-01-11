#!/bin/bash

# https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
R='\033[1;31m'
G='\033[0;32m'
O='\033[1;33m'
NC='\033[0m' # No Color

helpFunction()
{
   echo ""
   echo "Usage: $0 -a build -e dev -m normal -p arm"
   echo -e "\t-a Action: build, watch changes on disk (local) or logs (local)"
   echo -e "\t-e Environment: dev (local), staging (gcp)"
   echo -e "\t-m Mode: normal or dry run"
   echo -e "\t-p Platform: arm/amd. M1 uses arm, GCP uses amd"
   echo -e "\t-u Upload to registry true/false (push)"
   exit 1 # Exit script after printing help
}

setEnvPath()
{
    if [ "$env" == "dev" ]; then
        echo -e "${G}Using $env environment${NC}";
        envPath="./.secrets/.env.development"
    elif [ "$env" == "staging" ]; then
        echo -e  "${G}Using $env environment${NC}";
        envPath="./.secrets/.env.staging"
    else
        echo -e  "${R}Incorrect environment parameter${NC}";
        helpFunction    
    fi
}

# Ignores secrets used:
#  dev build is just local so secrets are passed using env variables (not secure but local)
#  staging build is declaring same dockerfile but actually doesn't get passed in any secure params
#   instead it is using GCP secrets store, but docker is still complaining
buildArgs="--set "*.args.BUILDKIT_DOCKERFILE_CHECK=skip=SecretsUsedInArgOrEnv""

# Global variable since shell doesn't allow proper return from functions
envPath=""

while getopts "a:e:m:p:u:" opt
do
   case "$opt" in
      a ) action="$OPTARG" ;;
      e ) env="$OPTARG" ;;
      m ) mode="$OPTARG" ;;
      p ) platform="$OPTARG" ;;
      u ) push=$OPTARG ;;
      ? ) helpFunction ;; # Print helpFunction in case parameter is non-existent
   esac
done

# Print helpFunction in case parameters are empty
if [ -z "$action" ]; then
   echo -e  "${R}Action is required${NC}";
   helpFunction
fi

if [ "$action" != "build" ] && [ "$action" != "watch" ] && [ "$action" != "logs" ]; then
   echo -e  "${R}Incorrect action parameter${NC}";
   helpFunction
fi

if [ "$action" == "build" ]; then
    setEnvPath
    echo -e  "${G}Building${NC}"
    platformOverride=""
    if [ "$platform" == "arm" ]; then
        platformOverride="--set *.platform=linux/arm64"
        echo -e  "${G}Platform override 'arm'${NC}"
    elif [ "$platform" == "amd" ]; then
        platformOverride="--set *.platform=linux/amd64"
        echo -e  "${G}Platform override 'amd'${NC}"
    else
        echo -e  "${O}Platform override ignored. Building for both arm and amd${NC}"
    fi

    dryRun=""
    if [ "$mode" == "dry" ]; then
        dryRun="--print"
        echo -e  "${G}Dry run, will only print config${NC}"
    fi

    pushCommand=""
    if [ $push == true ]; then
        pushCommand="--push"
        echo -e  "${G}Push to registry after building${NC}"
    fi

    echo -e "\n"
    eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' $envPath) && docker buildx bake $platformOverride $dryRun $pushCommand $buildArgs
elif [ "$action" == "watch" ] ; then
    setEnvPath
    echo -e  "${G}Watching${NC}"
    eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' $envPath) && docker compose up --watch
elif [ "$action" == "logs" ] ; then
    setEnvPath
    echo -e  "${G}Logs${NC}"
    eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' $envPath) && docker compose logs
else
    echo -e  "${R}Incorrect action parameter${NC}";
    helpFunction
fi