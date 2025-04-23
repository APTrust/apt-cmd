#!/usr/bin/env bash

OS=`uname`
CURRENT_ARCHITECTURE=`uname -m`
APPNAME="apt-cmd"
TEST_CMD="./scripts/test.rb units"
VERSION=`git describe --tags`
COMMIT_ID=`git rev-parse --short HEAD`
DATE=`date +"%Y-%m-%d"`
CGO_FLAG=""
CURRENT_OS="linux"
BUILD_DIR="dist"
LDFLAGS="-X github.com/APTrust/apt-cmd/cmd.CommitId=$COMMIT_ID -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION -X github.com/APTrust/apt-cmd/cmd.BuildDate=$DATE"

if [[ $OS =~ "MINGW" || $OS =~ "Windows" ]]; then
    CURRENT_OS="windows"
elif [[ $OS =~ "Darwin" ]]; then
    CURRENT_OS="darwin"
fi

run_unit_tests() {
    echo "Running unit tests"
    $TEST_CMD
    if [[ $? != 0 ]]; then
        echo "Aborting because unit tests failed."
        exit 1
    else
        echo "Tests passed. Proceeding with build."
    fi
}

build_tags() {
    case "$1" in
        "windows")
            echo "windows"
            ;;
        *)
            echo "posix"
            ;;
    esac
}

app_name() {
    case "$1" in
        "windows")
            echo "apt-cmd.exe"
            ;;
        *)
            echo "apt-cmd"
            ;;
    esac
}

build_all() {
    echo "Building Commit ID: ${COMMIT_ID}, Version: ${VERSION}"

    platforms=("darwin" "linux" "windows")
    architectures=("amd64" "arm64")

    for platform in "${platforms[@]}"
    do
        for architecture in "${architectures[@]}"
        do
            echo " "
            echo "Creating directory $BUILD_DIR/$platform/$architecture"
            mkdir -p $BUILD_DIR/$platform/$architecture
            echo "Building $BUILD_DIR/$platform/$architecture/$(app_name $platform)"
            GOOS=$platform GOARCH=$architecture CGO_ENABLED=0 go build -tags=$(build_tags $platform) -o $BUILD_DIR/$platform/$architecture/$(app_name $platform) -ldflags="$LDFLAGS"
        done
    done
}

run_version() {
    EXECUTABLE=$BUILD_DIR/$CURRENT_OS/$CURRENT_ARCHITECTURE/apt-cmd
    if [[ "$CURRENT_OS" == "windows" ]]; then
        EXECUTABLE=$BUILD_DIR/$CURRENT_OS/$CURRENT_ARCHITECTURE/apt-cmd.exe
    fi

    echo ""
    echo "Running ${EXECUTABLE} version..."
    echo ""

    ${EXECUTABLE} version
    echo ""
}

#
# Run the whole show.
#
# TODO:
#    - Run this as a GitLab workflow.
#    - If build succeeds, upload executables to our S3 public download bucket.
#
run_unit_tests
build_all
run_version
