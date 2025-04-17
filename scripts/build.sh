#!/usr/bin/env bash

OS=`uname`
APPNAME="apt-cmd"
TEST_CMD="./scripts/test.rb units"
VERSION=`git describe --tags`
COMMIT_ID=`git rev-parse --short HEAD`
DATE=`date +"%Y-%m-%d"`
CGO_FLAG=""

BUILD_TAGS="posix"

if [[ $OS =~ "MINGW" || $OS =~ "Windows" ]]; then
    APPNAME="apt-cmd.exe"
    BUILD_TAGS="windows"
fi

BUILD_DIR="dist"
OUTPUT_FILE="${BUILD_DIR}/${APPNAME}"

# echo "Running unit tests"
# $TEST_CMD
# if [[ $? != 0 ]]; then
#     echo "Aborting because unit tests failed."
#     exit 1
# else
#     echo "Tests passed. Proceeding with build."
# fi

echo "Commit ID: ${COMMIT_ID}, Version: ${VERSION}"
echo "Building ${OUTPUT_FILE}"
echo ""
mkdir -p $BUILD_DIR

if [[ $OS =~ "Linux" ]]; then
    CGO_ENABLED=0 go build -tags ${BUILD_TAGS} -o ${OUTPUT_FILE} -ldflags="-X github.com/APTrust/apt-cmd/cmd.CommitId=$COMMIT_ID -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION -X github.com/APTrust/apt-cmd/cmd.BuildDate=$DATE"
else
    go build -tags ${BUILD_TAGS} -o ${OUTPUT_FILE} -ldflags="-X github.com/APTrust/apt-cmd/cmd.CommitId=$COMMIT_ID -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION -X github.com/APTrust/apt-cmd/cmd.BuildDate=$DATE"
fi


echo ""
echo "Running ${OUTPUT_FILE} version..."
echo ""

${OUTPUT_FILE} version
echo ""
