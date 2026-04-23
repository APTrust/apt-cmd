#!/usr/bin/env bash
#
# release.sh - Build and release apt-cmd binaries to S3.
#
# Usage:
#   release.sh --pre-check <version>   Run pre-release checks only.
#   release.sh <version>               Build and release binaries.
#
# Example:
#   release.sh --pre-check v3.0.4
#   release.sh v3.0.4
#
# Binaries are uploaded to:
#   s3://aptrust.public.download/apt-cmd/<version>/linux/<arch>/apt-cmd
#   s3://aptrust.public.download/apt-cmd/<version>/mac/<arch>/apt-cmd
#   s3://aptrust.public.download/apt-cmd/<version>/windows/<arch>/apt-cmd.exe
#
# ----------------------------------------------------------------------------

set -euo pipefail

# ---------------------------------------------------------------------------
# Determine script and project directories
# ---------------------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

APPNAME="apt-cmd"
S3_BUCKET="aptrust.public.download"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
die() {
    echo "ERROR: $*" >&2
    exit 1
}

# ---------------------------------------------------------------------------
# Parse arguments
# ---------------------------------------------------------------------------
PRE_CHECK=false
VERSION=""

if [[ $# -eq 0 ]]; then
    die "Usage: release.sh [--pre-check] <version>  (e.g. release.sh v3.0.4)"
fi

if [[ "$1" == "--pre-check" ]]; then
    PRE_CHECK=true
    shift
fi

if [[ $# -eq 0 ]]; then
    die "Usage: release.sh [--pre-check] <version>  (e.g. release.sh v3.0.4)"
fi

VERSION="$1"

# ---------------------------------------------------------------------------
# Check required environment variables
# ---------------------------------------------------------------------------
check_env() {
    [[ -n "${AWS_ACCESS_KEY_ID:-}"     ]] || die "AWS_ACCESS_KEY_ID is not set in the environment."
    [[ -n "${AWS_SECRET_ACCESS_KEY:-}" ]] || die "AWS_SECRET_ACCESS_KEY is not set in the environment."
}

# ---------------------------------------------------------------------------
# Check that the current git tag matches the requested version
# ---------------------------------------------------------------------------
check_git_tag() {
    local current_tag
    current_tag="$(git -C "$PROJECT_DIR" describe --tags --exact-match HEAD 2>/dev/null || true)"
    if [[ "$current_tag" != "$VERSION" ]]; then
        die "Current git tag ('${current_tag:-<none>}') does not match the requested version '$VERSION'. " \
            "Please create or check out tag '$VERSION' before releasing."
    fi
}

# ---------------------------------------------------------------------------
# Verify CHANGELOG.md contains an h2 entry for this version with a date
# ---------------------------------------------------------------------------
check_changelog() {
    local changelog_file="$PROJECT_DIR/CHANGELOG.md"
    [[ -f "$changelog_file" ]] || die "CHANGELOG.md not found in $PROJECT_DIR."

    # Match the first ## headline that contains the version string.
    # The headline may look like "## v3.0.4 - 2026-04-23" or "## v3.0.4 - April 23, 2026".
    local match
    local version_no_v="${VERSION#v}"
    match="$(grep -m1 "^## .*${version_no_v}" "$changelog_file" || true)"

    if [[ -z "$match" ]]; then
        die "CHANGELOG.md does not contain an h2 (##) section for '$VERSION'. " \
            "Please add a changelog entry and release date for this version before releasing."
    fi

    # Check that the matched line contains a date.
    # Accept ISO dates (2026-04-23) or written dates (April 23, 2026 / Apr 23, 2026).
    if ! echo "$match" | grep -qE '([0-9]{4}-[0-9]{2}-[0-9]{2}|[A-Za-z]+ [0-9]{1,2},? [0-9]{4})'; then
        die "The '## $VERSION' headline in CHANGELOG.md does not include a release date. " \
            "Please add a date (e.g. '2026-04-23' or 'April 23, 2026') to the headline before releasing."
    fi

    echo "CHANGELOG.md: found entry for '$VERSION' with a release date."
}

# ---------------------------------------------------------------------------
# Build binaries for all platforms
# ---------------------------------------------------------------------------
build_binaries() {
    echo "Building $APPNAME binaries for all platforms..."
    cd "$PROJECT_DIR"

    local commit_id date ldflags
    commit_id="$(git rev-parse --short HEAD)"
    date="$(date +"%Y-%m-%d")"
    ldflags="-X github.com/APTrust/apt-cmd/cmd.CommitId=$commit_id \
             -X github.com/APTrust/apt-cmd/cmd.Version=$VERSION \
             -X github.com/APTrust/apt-cmd/cmd.BuildDate=$date"

    local platforms=("darwin" "linux" "windows")
    local architectures=("amd64" "arm64")

    for platform in "${platforms[@]}"; do
        for arch in "${architectures[@]}"; do
            local outname="$APPNAME"
            local build_tag="posix"
            if [[ "$platform" == "windows" ]]; then
                outname="${APPNAME}.exe"
                build_tag="windows"
            fi
            local outdir="$PROJECT_DIR/dist/$platform/$arch"
            mkdir -p "$outdir"
            echo "Building $outdir/$outname"
            GOOS="$platform" GOARCH="$arch" CGO_ENABLED=0 \
                go build -tags="$build_tag" \
                         -o "$outdir/$outname" \
                         -ldflags="$ldflags"
        done
    done

    echo "Build complete."
}

# ---------------------------------------------------------------------------
# Upload binaries to S3
# ---------------------------------------------------------------------------
upload_binaries() {
    echo "Uploading binaries to s3://$S3_BUCKET/$APPNAME/$VERSION/..."

    # Parallel arrays: local dist path -> arch argument for s3_helper.go
    local -a SRC_FILES=(
        "$PROJECT_DIR/dist/linux/amd64/$APPNAME"
        "$PROJECT_DIR/dist/linux/arm64/$APPNAME"
        "$PROJECT_DIR/dist/darwin/amd64/$APPNAME"
        "$PROJECT_DIR/dist/darwin/arm64/$APPNAME"
        "$PROJECT_DIR/dist/windows/amd64/${APPNAME}.exe"
        "$PROJECT_DIR/dist/windows/arm64/${APPNAME}.exe"
    )
    local -a ARCHES=(
        "linux/amd64"
        "linux/arm64"
        "mac/amd64"
        "mac/arm64"
        "windows/amd64"
        "windows/arm64"
    )

    cd "$SCRIPT_DIR"
    for i in "${!SRC_FILES[@]}"; do
        local src="${SRC_FILES[$i]}"
        local arch="${ARCHES[$i]}"
        [[ -f "$src" ]] || die "Expected binary not found at: $src"
        echo "Uploading $src (arch: $arch)..."
        go run s3_helper.go -upload -version "$VERSION" -arch "$arch" "$src" || \
            die "Failed to upload $src to S3."
    done

    echo "All binaries uploaded successfully."
}

# ---------------------------------------------------------------------------
# Print download links
# ---------------------------------------------------------------------------
print_download_links() {
    echo ""
    echo "Download links for $APPNAME $VERSION:"
    cd "$SCRIPT_DIR"
    go run s3_helper.go -get-links "$VERSION" || \
        die "Failed to retrieve download links for version '$VERSION'."
}

# ---------------------------------------------------------------------------
# Pre-check mode
# ---------------------------------------------------------------------------
run_pre_check() {
    echo "Running pre-release checks for version '$VERSION'..."
    check_changelog
    check_git_tag
    check_env
    echo "All pre-release checks passed."
}

# ---------------------------------------------------------------------------
# Full release
# ---------------------------------------------------------------------------
run_release() {
    # 1. Verify git tag, changelog, and AWS credentials
    check_git_tag
    check_changelog
    check_env

    # 2. Build binaries for all platforms
    build_binaries

    # 3. Create S3 folders for this version
    echo "Creating S3 folders for version '$VERSION'..."
    cd "$SCRIPT_DIR"
    go run s3_helper.go -make-folders -version "$VERSION" || \
        die "Failed to create S3 folders for version '$VERSION'."

    # 4. Upload binaries to S3
    upload_binaries

    # 5. Print download links
    print_download_links

    # 6. Print post-release checklist
    cat <<EOF

Checklist:

1. Ensure builds for all platforms (linux/amd64, linux/arm64, mac/amd64,
   mac/arm64, windows/amd64, windows/arm64) are present in:
   s3://$S3_BUCKET/$APPNAME/$VERSION/
2. Using the download links above, update the apt-cmd version number,
   release date, and download links in the following places:
      - https://github.com/APTrust/apt-cmd/blob/master/README.md
      - https://aptrust.github.io/userguide/partner_tools/#downloads
3. Push any README or documentation changes to GitHub. This will require
   you running "mkdocs gh-deploy" to publish the user guide.
EOF
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
if [[ "$PRE_CHECK" == "true" ]]; then
    run_pre_check
else
    run_release
fi
