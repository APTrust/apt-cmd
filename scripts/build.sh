#!/usr/bin/env bash

OS=`uname`
ARCH=`uname -m`
APPNAME="apt-cmd"
VERSION=`git describe --tags`
COMMIT_ID=`git rev-parse --short HEAD`
DATE=`date +"%Y-%m-%d"`

if [[ $OS =~ "MINGW" || $OS =~ "Windows" ]]; then
    APPNAME="apt-cmd.exe"
    ARCH="Windows"
fi

BUILD_DIR="dist/${OS}/${ARCH}/${VERSION}"

if [[ $VERSION == "" ]]; then
    VERSION="Internal-Test-Build"
    BUILD_DIR="dist/testbuild"
fi

OUTPUT_FILE="${BUILD_DIR}/${APP_NAME}"


echo "Commit ID: ${COMMIT_ID}, Version: ${VERSION}"
echo "Building ${OUTPUT_FILE}"
mkdir -p $BUILD_DIR

go build -o ${OUTPUT_FILE} -ldflags="-X github.com/APTrust/apt-cmd/cmd.CommitId=$COMMIT_ID -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION -X github.com/APTrust/apt-cmd/cmd.BuildDate=$DATE"


