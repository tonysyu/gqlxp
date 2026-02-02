#!/bin/sh
# Install script for gqlxp
# Usage: curl -sSfL https://raw.githubusercontent.com/tonysyu/gqlxp/main/install.sh | sh
#    or: curl -sSfL https://raw.githubusercontent.com/tonysyu/gqlxp/main/install.sh | sh -s -- -b /custom/install/path

set -eu

REPO="tonysyu/gqlxp"
BINARY_NAME="gqlxp"
INSTALL_DIR="${INSTALL_DIR:-$HOME/bin}"

usage() {
    cat <<EOF
Install script for gqlxp

Usage:
    install.sh [options]

Options:
    -b, --bin-dir DIR    Installation directory (default: $HOME/bin)
    -h, --help           Show this help message

Environment variables:
    INSTALL_DIR          Installation directory (overridden by -b flag)
EOF
}

# Parse arguments
while [ $# -gt 0 ]; do
    case "$1" in
        -b|--bin-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

detect_os() {
    os=$(uname -s)
    case "$os" in
        Linux)  echo "Linux" ;;
        Darwin) echo "Darwin" ;;
        *)
            echo "Error: Unsupported operating system: $os" >&2
            exit 1
            ;;
    esac
}

detect_architecture() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)  echo "x86_64" ;;
        arm64|aarch64) echo "arm64" ;;
        i386|i686)     echo "i386" ;;
        *)
            echo "Error: Unsupported architecture: $arch" >&2
            exit 1
            ;;
    esac
}

get_latest_app_version() {
    version=$(curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        echo "Error: Could not determine latest version" >&2
        exit 1
    fi
    echo "$version"
}

# Download file with error handling
#
# @param url     Url of file to download
# @param output  Filepath for
download_file() {
    url="$1"
    output="$2"
    if ! curl -sSfL -o "$output" "$url"; then
        echo "Error: Failed to download $url" >&2
        exit 1
    fi
}

# Verify checksum
#
# @param archive         Filepath to verify
# @param checksums_file  Filepath to file containing checksums
#
# Each line of the checksums file should contain:
#     {CHECKSUM_HASH}  {FILENAME}
verify_checksum() {
    archive="$1"
    checksums_file="$2"
    archive_name=$(basename "$archive")

    expected=$(grep "$archive_name" "$checksums_file" | awk '{print $1}')
    if [ -z "$expected" ]; then
        echo "Error: Checksum not found for $archive_name" >&2
        exit 1
    fi

    if command -v sha256sum >/dev/null 2>&1; then
        actual=$(sha256sum "$archive" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        actual=$(shasum -a 256 "$archive" | awk '{print $1}')
    else
        echo "Error: No sha256sum or shasum command found" >&2
        exit 1
    fi

    if [ "$expected" != "$actual" ]; then
        echo "Error: Checksum verification failed" >&2
        echo "  Expected: $expected" >&2
        echo "  Actual:   $actual" >&2
        exit 1
    fi

    echo "Checksum verified successfully"
}

main() {
    os=$(detect_os)
    arch=$(detect_architecture)
    version=$(get_latest_app_version)
    # Version string with leading "v" trimmed off
    version_number="${version#v}"

    echo "Installing ${BINARY_NAME} ${version} for ${os}/${arch}..."

    # Construct download URLs (goreleaser naming convention)
    archive_name="${BINARY_NAME}_${os}_${arch}.tar.gz"
    base_url="https://github.com/${REPO}/releases/download/${version}"
    archive_url="${base_url}/${archive_name}"
    checksums_url="${base_url}/gqlxp_${version_number}_checksums.txt"

    # Create temporary directory
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT

    archive_path="${tmp_dir}/${archive_name}"
    checksums_path="${tmp_dir}/checksums.txt"

    # Download archive and checksums
    echo "Downloading ${archive_url}..."
    download_file "$archive_url" "$archive_path"

    echo "Downloading checksums..."
    download_file "$checksums_url" "$checksums_path"

    # Verify checksum
    verify_checksum "$archive_path" "$checksums_path"

    # Extract archive
    echo "Extracting archive..."
    tar -xzf "$archive_path" -C "$tmp_dir"

    # Install binary
    echo "Installing to ${INSTALL_DIR}..."
    if [ ! -d "$INSTALL_DIR" ]; then
        mkdir -p "$INSTALL_DIR"
    fi

    if [ -w "$INSTALL_DIR" ]; then
        mv "${tmp_dir}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo "Elevated permissions required to install to ${INSTALL_DIR}"
        sudo mv "${tmp_dir}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    echo ""
    echo "${BINARY_NAME} ${version} installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"

    # Check if install dir is in PATH
    case ":$PATH:" in
        *":${INSTALL_DIR}:"*) ;;
        *)
            echo ""
            echo "Note: ${INSTALL_DIR} is not in your PATH."
            echo "Add it with: export PATH=\"\$PATH:${INSTALL_DIR}\""
            ;;
    esac
}

main
