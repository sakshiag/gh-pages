#!/bin/bash

artifactory_registry=terraform-docker-local.artifactory.swg-devops.com

#Remove all exited containers which might include the UAT which has been run for the same commit
sudo docker rm $(sudo docker ps -a -f status=exited -q)

BUILD_ENV=$1
GIT_SHA=$2
REPORT_URL=$3
CODE_REPO=/tmp/${BUILD_ENV}/${GIT_SHA}
GIT_REPO=git@github.ibm.com:blueprint/bluemix-terraform-provider-dev.git
mkdir -p $CODE_REPO
cd $CODE_REPO

echo "Cloning $GIT_REPO at  $CODE_REPO"
git clone $GIT_REPO .

echo "Checking out $GIT_SHA to temp branch"
git checkout -b temp $GIT_SHA

#Prep the uatrunner environment with DOCKER_USER and DOCKER_PASSWORD env variables
sudo docker login -u "$DOCKER_USER" -p "$DOCKER_PASSWORD" -e "$DOCKER_EMAIL" $artifactory_registry

echo "Building the docker uatrunner:${BUILD_ENV}_${GIT_SHA}"
sudo docker build -t uatrunner:${BUILD_ENV}_${GIT_SHA} .

echo "Run the docker which will run the UAT by calling build.sh of the main repo"
#IBMID IBMID_PASSWORD and SL_ACCOUNT_NUMBER must be set in the UAT runner enviroments
sudo docker run -d  --name ${BUILD_ENV}_${GIT_SHA} \
-e UAT="true" \
-e TRAVIS_COMMIT="$GIT_SHA" \
-e REPORT_URL="$REPORT_URL" \
-e IBMID="${IBMID}" \
-e IBMCLOUD_VIRTUAL_GUEST_IMAGE_ID="${IBMCLOUD_VIRTUAL_GUEST_IMAGE_ID}" \
-e TESTARGS="${TESTARGS}" \
-e BUILD_ENV="${BUILD_ENV}" \
-e SOFTLAYER_TIMEOUT=120 \
-e IBMID_PASSWORD="${IBMID_PASSWORD}" \
-e SL_ACCOUNT_NUMBER="${SL_ACCOUNT_NUMBER}" \
uatrunner:${BUILD_ENV}_${GIT_SHA}
