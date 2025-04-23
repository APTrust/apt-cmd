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

LDFLAGS="-X github.com/APTrust/apt-cmd/cmd.CommitId=$COMMIT_ID -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION -X github.com/APTrust/apt-cmd/cmd.BuildDate=$DATE"

echo "Running unit tests"
$TEST_CMD
if [[ $? != 0 ]]; then
    echo "Aborting because unit tests failed."
    exit 1
else
    echo "Tests passed. Proceeding with build."
fi

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



echo "Commit ID: ${COMMIT_ID}, Version: ${VERSION}"
echo "Building ${OUTPUT_FILE}"
echo ""
mkdir -p $BUILD_DIR

platforms=("darwin" "linux" "windows")
architectures=("amd64" "arm64")

for platform in "${platforms[@]}"
do
    for architecture in "${architectures[@]}"
    do
        echo "Creating directory $BUILD_DIR/$platform/$architecture"
        mkdir -p $BUILD_DIR/$platform/$architecture
        echo "Building $platform/$architecture"
        GOOS=$platform GOARCH=$architecture CGO_ENABLED=0 go build -tags=$(build_tags $platform) -o $BUILD_DIR/$platform/$architecture/$(app_name $platform) -ldflags="$LDFLAGS"
    done
done

# if [[ $OS =~ "Linux" ]]; then
#     CGO_ENABLED=0 go build -tags ${BUILD_TAGS} -o ${OUTPUT_FILE} -ldflags="-X github.com/APTrust/apt-cmd/cmd.CommitId=$COMMIT_ID -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION -X github.com/APTrust/apt-cmd/cmd.BuildDate=$DATE"
# else
#     go build -tags ${BUILD_TAGS} -o ${OUTPUT_FILE} -ldflags="-X github.com/APTrust/apt-cmd/cmd.CommitId=$COMMIT_ID -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION -X github.com/APTrust/apt-cmd/cmd.BuildDate=$DATE"
# fi


# echo ""
# echo "Running ${OUTPUT_FILE} version..."
# echo ""

# ${OUTPUT_FILE} version
# echo ""
