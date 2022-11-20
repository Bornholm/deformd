#!/bin/sh
set -e

DEFORMD_RELEASES_URL="https://github.com/Bornholm/deformd/releases"
DEFORMD_DESTDIR="."
DEFORMD_FILE_BASENAME="frmd"
DEFORMD_PLATFORM_SUFFIX="$(uname -s)_$(uname -m)"

function main {
    test -z "${DEFORMD_VERSION}" && DEFORMD_VERSION="$(curl -sfL -o /dev/null -w %{url_effective} "${DEFORMD_RELEASES_URL}/latest" |
        rev |
        cut -f1 -d'/'|
        rev)"

    # Check version variable initialization
    test -z "${DEFORMD_VERSION}" && {
        echo "Unable to get DEFORMD version !" >&2
        exit 1
    }

    test -z "${DEFORMD_TMPDIR}" && DEFORMD_TMPDIR="$(mktemp -d)"
    export TAR_FILE="${DEFORMD_TMPDIR}/${DEFORMD_FILE_BASENAME}_${DEFORMD_VERSION}_${DEFORMD_PLATFORM_SUFFIX}.tar.gz"

    (
        cd "${DEFORMD_TMPDIR}"

        # Download DEFORMD
        echo "Downloading DEFORMD ${DEFORMD_VERSION}..."
        curl -sfLo "${TAR_FILE}" \
            "${DEFORMD_RELEASES_URL}/download/${DEFORMD_VERSION}/${DEFORMD_FILE_BASENAME}_${DEFORMD_VERSION}_${DEFORMD_PLATFORM_SUFFIX}.tar.gz" ||
            ( echo  "Error while downloading DEFORMD !" >&2 && exit 1 )
        
        # Download checksums
        curl -sfLo "checksums.txt" "${DEFORMD_RELEASES_URL}/download/${DEFORMD_VERSION}/checksums.txt"
        
        echo "Verifying checksum..."
        check_sum ||
            ( echo  "Error while verifying checksums !" >&2 && exit 1 )
    )

    # Extracting archive files
    tar -xf "${TAR_FILE}" -C "${DEFORMD_TMPDIR}"

    # Moving downloaded binary to destination directory
    mv -f "${DEFORMD_TMPDIR}/${DEFORMD_FILE_BASENAME}" "${DEFORMD_DESTDIR}/"

    echo "You can now use '${DEFORMD_DESTDIR}/${DEFORMD_FILE_BASENAME}', enjoy !"
}

function check_sum {
    set -o pipefail
    cat checksums.txt | grep frmd_*_${DEFORMD_PLATFORM_SUFFIX}.tar.gz | sha256sum -c
    set +o pipefail
}

main $@