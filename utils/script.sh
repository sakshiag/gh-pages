#!/bin/bash

artifactory_registry=terraform-docker-local.artifactory.swg-devops.com

#Remove all exited containers which might include the e2e which has been run for the same commit
sudo docker rm $(sudo docker ps -a -f status=exited -q)

BUILD_ENV=$1
GIT_SHA=$2
E2E_SHA=$3
REPORT_URL=$4
CODE_REPO=/tmp/${BUILD_ENV}/${GIT_SHA}
GIT_REPO=git@github.ibm.com:terraform-devops-tools/e2etest.git
TIMESTAMP=`date +'%y%m%d%H%M%S'`

echo $TIMESTAMP

mkdir -p $CODE_REPO

cd $CODE_REPO

echo "Cloning $GIT_REPO at  $CODE_REPO"
git clone $GIT_REPO .

echo "Checking out $E2E_SHA to temp branch"
git checkout -b temp $E2E_SHA


#Prep the e2erunner environment with DOCKER_USER and DOCKER_PASSWORD env variables
sudo docker login -u "$DOCKER_USER" -p "$DOCKER_PASSWORD" -e "$DOCKER_EMAIL" $artifactory_registry

echo "Building the docker e2erunner:${BUILD_ENV}_${GIT_SHA}"
sudo docker build --build-arg FTP_USERNAME="${FTP_USERNAME}" --build-arg FTP_PASSWORD="${FTP_PASSWORD}" -t e2erunner:${BUILD_ENV}_${GIT_SHA} .

echo "Run the docker which will run the e2e by calling build.sh of the main repo"
#SL_USERNAME and SL_API_KEY must be set in the e2e runner enviroments
#FTP_USERNMAME and FTP_PASSWORD must be set in the e2e runner environments
sudo docker run -d  --name ${BUILD_ENV}_${GIT_SHA} \
-e e2e="true" \
-e TRAVIS_COMMIT="$GIT_SHA" \
-e E2E_COMMIT="$E2E_SHA" \
-e REPORT_URL="$REPORT_URL" \
-e TESTARGS="${TESTARGS}" \
-e BUILD_ENV="${BUILD_ENV}" \
-e SOFTLAYER_TIMEOUT=300 \
-e SL_USERNAME="${SL_USERNAME}" \
-e SL_API_KEY="${SL_API_KEY}" \
-e IBMID="${IBMID}" \
-e IBMID_PASSWORD="${IBMID_PASSWORD}" \
-e SL_ACCOUNT_NUMBER="${SL_ACCOUNT_NUMBER}" \
-e TIMESTAMP="${TIMESTAMP}" \
-e FTP_USERNAME="${FTP_USERNAME}" \
-e FTP_PASSWORD="${FTP_PASSWORD}" \
e2erunner:${BUILD_ENV}_${GIT_SHA}
