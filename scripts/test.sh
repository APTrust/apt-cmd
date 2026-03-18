#!/usr/bin/env bash
# coding: utf-8

# Run unit and integration tests for apt-cmd.
# Bash 3.x compatible.

# Parallel arrays for service tracking (no associative arrays in bash 3.x)
SERVICE_NAMES=()
SERVICE_PIDS=()
REGISTRY_STARTED="false"
SERVICES_STOPPED="false"
TEST_NAME=""
START_TIME=$(date +%s)

os_name() {
    case "$OSTYPE" in
        darwin*)             echo "osx" ;;
        cygwin*|msys*|mingw*) echo "windows" ;;
        *)                   echo "linux" ;;
    esac
}

project_root() {
    local script_dir
    script_dir="$(cd "$(dirname "$0")" && pwd)"
    echo "$(cd "$script_dir/.." && pwd)"
}

bin_dir() {
    local os bin
    os=$(os_name)
    bin="$(project_root)/bin/$os"
    if [ "$os" = "osx" ]; then
        if [ "$(uname -m)" = "arm64" ]; then
            bin="$bin/arm64"
        else
            bin="$bin/amd64"
        fi
    fi
    echo "$bin"
}

log_file_path() {
    echo "$HOME/tmp/logs/$1.log"
}

clean_test_cache() {
    echo "Deleting test cache from last run"
    go clean -testcache
}

make_test_dirs() {
    local base dir full_dir bucket full_bucket
    base="$HOME/tmp"

    # Safety check: only delete if path ends with /tmp
    case "$base" in
        */tmp) echo "Deleting $base" ;;
    esac
    rm -rf "$base"

    for dir in bin logs minio nsq redis restore; do
        full_dir="$base/$dir"
        echo "Creating $full_dir"
        mkdir -p "$full_dir"
    done

    # S3 buckets for minio. We should ideally read these from the .env.test file.
    for bucket in test-bucket-1 test-bucket-2; do
        full_bucket="$base/minio/$bucket"
        echo "Creating local minio bucket $bucket"
        mkdir -p "$full_bucket"
    done
}

setup_env() {
    export MINIO_ACCESS_KEY="minioadmin"
    export MINIO_SECRET_KEY="minioadmin"
    if [ "$TEST_NAME" != "units" ] && [ "$(os_name)" != "windows" ]; then
        if [ -z "$REGISTRY_ROOT" ]; then
            echo "Set env var REGISTRY_ROOT" >&2
            exit 1
        fi
    fi
}

print_results() {
    local exit_status=$1
    local end_time elapsed
    end_time=$(date +%s)
    elapsed=$((end_time - START_TIME))
    echo ""
    echo "Elapsed time: $elapsed seconds"
    echo "Logs are in $HOME/tmp/logs"
    if [ "$exit_status" -eq 0 ]; then
        printf "\n\n    **** \xF0\x9F\x98\x81 PASS \xF0\x9F\x98\x81 **** \n\n"
    else
        printf "\n\n    **** \xF0\x9F\xA4\xAC FAIL \xF0\x9F\xA4\xAC **** \n\n"
        exit 1
    fi
}

# Stop a service on Linux by name using pidof.
# Note: this kills ALL processes with the given name.
stop_service_linux() {
    local name=$1
    local pids_found pid
    pids_found=$(pidof "$name" 2>/dev/null) || true
    for pid in $pids_found; do
        if kill -TERM "$pid" 2>/dev/null; then
            echo "(Linux) Killed $name with pid $pid"
        else
            echo "Hmm... Couldn't kill $name."
            echo "Check system processes to see if a version"
            echo "of that process is lingering from a previous test run."
        fi
    done
}

stop_service() {
    local name=$1
    local pid=$2
    local signal

    if [ -z "$pid" ] || [ "$pid" -eq 0 ] 2>/dev/null; then
        echo "Pid for $name is zero. Can't kill that..."
        return
    fi

    if [ "$(os_name)" = "linux" ]; then
        stop_service_linux "$name"
        return
    fi

    signal="TERM"
    if [ "$(os_name)" = "windows" ]; then
        signal="KILL"
    fi

    echo "Stopping $name service (pid $pid) with signal $signal"
    if ! kill -"$signal" "$pid" 2>/dev/null; then
        echo "Hmm... Couldn't kill $name."
        echo "Check system processes to see if a version"
        echo "of that process is lingering from a previous test run."
    fi
}

stop_all_services() {
    [ "$SERVICES_STOPPED" = "true" ] && return
    echo "Stopping all services"
    local i=0
    while [ "$i" -lt "${#SERVICE_NAMES[@]}" ]; do
        stop_service "${SERVICE_NAMES[$i]}" "${SERVICE_PIDS[$i]}"
        i=$((i + 1))
    done
    SERVICES_STOPPED="true"
}

start_service() {
    local name=$1
    local cmd=$2
    local msg=$3
    local log_file pid

    log_file=$(log_file_path "$name")

    # Start service in background, redirecting output to log file
    eval "$cmd" >> "$log_file" 2>&1 &
    pid=$!
    disown "$pid" 2>/dev/null || true

    echo ""
    echo "Started $name with command '$cmd' and pid $pid"
    echo "$msg"
    echo "Log file is $log_file"
    echo ""

    # Append to parallel arrays (bash 3.x compatible)
    SERVICE_NAMES=("${SERVICE_NAMES[@]}" "$name")
    SERVICE_PIDS=("${SERVICE_PIDS[@]}" "$pid")
}

registry_load_fixtures() {
    echo "Loading registry fixtures"
    local log_file
    log_file=$(log_file_path "registry_fixtures")
    (
        export APT_ENV="integration"
        cd "$REGISTRY_ROOT" || exit 1
        go run loader/load_fixtures.go >> "$log_file" 2>&1
    )
    echo "Registry fixtures loaded"
}

# Note: This assumes you have the registry repo source tree on your machine.
# See https://github.com/APTrust/registry
registry_start() {
    [ "$REGISTRY_STARTED" = "true" ] && return
    REGISTRY_STARTED="true"

    registry_load_fixtures

    local log_file registry_pid actual_pid registry_process
    log_file=$(log_file_path "registry")

    # Important! Adding -tags=test here turns on the special testing endpoints
    # prepare_file_delete and prepare_object_delete, which are disabled in all
    # non-test environments.
    (
        export APT_ENV="integration"
        cd "$REGISTRY_ROOT" || exit 1
        go run -tags=test registry.go >> "$log_file" 2>&1
    ) &
    registry_pid=$!
    disown "$registry_pid" 2>/dev/null || true
    sleep 3

    # go run compiles an executable, puts it in a temp directory, and runs it
    # as a new process. We need to get the pid of that process.
    # Note that the temp dir pattern will be different on Linux.
    # /var/folders works for Mac.
    registry_process=$(ps -ef | grep registry | grep "/var/folders" 2>/dev/null | head -1)
    actual_pid=$(echo "$registry_process" | awk '{print $2}')
    if [ -n "$actual_pid" ] && [ "$actual_pid" -gt 0 ] 2>/dev/null; then
        registry_pid=$actual_pid
    fi

    echo "Started Registry with command 'go run -tags=test registry.go' and pid $registry_pid"

    SERVICE_NAMES=("${SERVICE_NAMES[@]}" "registry")
    SERVICE_PIDS=("${SERVICE_PIDS[@]}" "$registry_pid")
}

# Initialize for integration, interactive tests, and end-to-end tests.
# This clears and rebuilds data directories, starts all services, and
# creates all NSQ topics.
init_for_integration() {
    local bin minio minio_data_dir

    clean_test_cache
    make_test_dirs

    if [ "$(os_name)" != "windows" ]; then
        registry_start
        sleep 8
    fi

    bin=$(bin_dir)
    minio="minio"
    if [ "$(os_name)" = "windows" ]; then
        minio="minio.exe"
    fi

    minio_data_dir="$HOME/tmp/minio"
    # For localhost testing, use 'localhost' instead of '127.0.0.1'
    # because Minio signed URLs use hostname, not IP.
    echo "$minio_data_dir"
    start_service \
        "minio" \
        "$bin/$minio server --address=localhost:9899 $minio_data_dir" \
        "Minio is running on localhost:9899. User/Pwd: minioadmin/minioadmin"

    sleep 5
}

run_unit_tests() {
    clean_test_cache
    run_go_unit_tests "$1"
}

run_go_unit_tests() {
    local arg="${1:-./...}"
    echo "Starting unit tests..."
    echo "go test $arg"
    setup_env
    (cd "$(project_root)" && go test $arg)
    print_results $?
}

run_integration_tests() {
    local arg="${1:-./...}"
    local tags exit_status

    init_for_integration

    tags="integration"
    if [ "$(os_name)" = "windows" ]; then
        tags="integration,windows"
        echo "*** NOT testing with local Registry on Windows. ***"
        echo "You must test Registry manually against an external server."
    fi

    echo "Starting integration tests..."
    echo "go test -tags=$tags $arg"

    setup_env

    if [ "$(os_name)" = "windows" ]; then
        # On Windows, capture stdout and print it explicitly (Open3 behavior)
        local stdout
        stdout=$(cd "$(project_root)" && go test -tags="$tags" $arg 2>/dev/null)
        exit_status=$?
        if [ ${#stdout} -gt 0 ]; then
            echo "$stdout"
        fi
    else
        (cd "$(project_root)" && go test -tags="$tags" $arg)
        exit_status=$?
    fi

    print_results "$exit_status"
}

print_help() {
    echo ""
    echo "APTrust partner tools tests"
    echo ""
    echo "Usage:"
    echo "  test.sh units                   # Run unit tests"
    echo "  test.sh integration             # Run integration tests"
    echo ""
    echo "Note that running integration tests also runs unit tests."
    echo ""
}

# ------------------------------------------------------------
# Main
# ------------------------------------------------------------
TEST_NAME="${1:-}"

if [ "$TEST_NAME" != "units" ] && [ "$TEST_NAME" != "integration" ]; then
    print_help
    exit 1
fi

trap stop_all_services EXIT

case "$TEST_NAME" in
    units)
        run_unit_tests "$2"
        ;;
    integration)
        run_integration_tests "$2"
        ;;
esac
