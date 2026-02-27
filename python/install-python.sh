#!/bin/sh
#
# install-python.sh — Install a CPython version via uv into a distroless rootfs.
#
# Usage:
#   ./install-python.sh [OPTIONS]
#
# Options:
#   --python-version=VERSION  Python version to install (default: latest)
#   --path=PATH               rootfs install target     (default: /rootfs/usr/local)
#   --strip=BOOL              remove pip/tk/idle etc    (default: true)
#
# Expects:
#   - uv in PATH
#   - ${path}/{bin,include,lib,share} directories to exist
#
# Dependencies: uv
#
# Examples:
#
#   # Install latest version with defaults
#   ./install-python.sh
#
#   # Install a specific version (distroless, stripped)
#   ./install-python.sh --python-version=3.13.2
#
#   # Install into a custom rootfs path
#   ./install-python.sh --python-version=3.12.9 --path=/opt/python
#
#   # Keep pip, tkinter, idle, etc. (dev/debug image)
#   ./install-python.sh --python-version=3.13.2 --strip=false
#
#   # Dockerfile usage
#   RUN sh /install-python.sh --python-version=3.14 --path=/rootfs/usr/local --strip=true
#
set -eu

# ---------------------------------------------------------------------------
# Configuration (defaults)
# ---------------------------------------------------------------------------
INSTALL_DIR="/python"
ROOTFS_PREFIX="/rootfs/usr/local"
STRIP_EXTRAS="true"
PYTHON_VERSION=""

# ---------------------------------------------------------------------------
# Parse flags
# ---------------------------------------------------------------------------
for arg in "$@"; do
    case "${arg}" in
        --python-version=*) PYTHON_VERSION="${arg#*=}" ;;
        --path=*)           ROOTFS_PREFIX="${arg#*=}" ;;
        --strip=*)          STRIP_EXTRAS="${arg#*=}" ;;
        *) printf '[install-python] ERROR: unknown flag: %s\n' "${arg}" >&2; exit 1 ;;
    esac
done

readonly INSTALL_DIR ROOTFS_PREFIX STRIP_EXTRAS PYTHON_VERSION

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

log() {
    printf '[install-python] %s\n' "$1"
}

die() {
    printf '[install-python] ERROR: %s\n' "$1" >&2
    exit 1
}

require_cmd() {
    command -v "$1" >/dev/null 2>&1 || die "'$1' is required but not found in PATH"
}

# ---------------------------------------------------------------------------
# Step 0 — Preflight checks
# ---------------------------------------------------------------------------
log "checking dependencies..."
require_cmd uv

if [ ! -d "${ROOTFS_PREFIX}/bin" ]; then
    die "rootfs not found at ${ROOTFS_PREFIX} — run apko build first"
fi

log "target version: ${PYTHON_VERSION:-latest (default)}"
log "rootfs prefix:  ${ROOTFS_PREFIX}"
log "strip extras:   ${STRIP_EXTRAS}"

# ---------------------------------------------------------------------------
# Step 1 — Install Python via uv
# ---------------------------------------------------------------------------
log "installing python via uv into ${INSTALL_DIR}..."
mkdir -p "${INSTALL_DIR}"
# shellcheck disable=SC2086  # intentional word-splitting when PYTHON_VERSION is empty
uv python install ${PYTHON_VERSION} --install-dir="${INSTALL_DIR}" --no-bin
log "uv install complete"

# ---------------------------------------------------------------------------
# Step 2 — Locate the installed cpython directory
# ---------------------------------------------------------------------------
python_dir=$(ls -d "${INSTALL_DIR}"/cpython-* 2>/dev/null | head -n1)

if [ -z "${python_dir}" ]; then
    die "no cpython directory found under ${INSTALL_DIR} — uv install may have failed"
fi

log "found installation at ${python_dir}"

# ---------------------------------------------------------------------------
# Step 3 — Move runtime files into the rootfs
# ---------------------------------------------------------------------------
for dir in bin include lib share; do
    if [ -d "${python_dir}/${dir}" ]; then
        mv "${python_dir}/${dir}"/* "${ROOTFS_PREFIX}/${dir}/"
        log "moved ${dir}/ into rootfs"
    else
        log "skipping ${dir}/ (not present in installation)"
    fi
done

# ---------------------------------------------------------------------------
# Step 4 — Strip unnecessary components to minimise image size
# ---------------------------------------------------------------------------
if [ "${STRIP_EXTRAS}" = "true" ]; then
    log "removing unnecessary components..."

    # Unneeded CLI tools
    rm -rf "${ROOTFS_PREFIX}"/bin/{pip*,idle*,pydoc*,python*-config}

    # Tcl/Tk libraries (no GUI in distroless)
    rm -rf "${ROOTFS_PREFIX}"/lib/{libtcl*.so*,tcl*,tk*,itcl*,thread*}

    # Python stdlib modules not needed at runtime
    rm -rf "${ROOTFS_PREFIX}"/lib/python*/{tkinter,idlelib,ensurepip,pydoc_data,turtle.py,turtledemo,__phello__}

    # Bundled package managers (uv replaces these)
    rm -rf "${ROOTFS_PREFIX}"/lib/python*/site-packages/{pip,pip-*,setuptools,setuptools-*,wheel,wheel-*}

    # Man pages and terminfo (already in base layer)
    rm -rf "${ROOTFS_PREFIX}"/share/{man,terminfo}

    log "cleanup complete"
else
    log "skipping cleanup (STRIP_EXTRAS=${STRIP_EXTRAS})"
fi

# ---------------------------------------------------------------------------
# Step 5 — Verify the installation
# ---------------------------------------------------------------------------
python_bin=$(ls "${ROOTFS_PREFIX}"/bin/python3.* 2>/dev/null | head -n1)

if [ -z "${python_bin}" ]; then
    die "no python3 binary found in rootfs after install"
fi

log "installed binary: ${python_bin}"
log "rootfs size: $(du -sh "${ROOTFS_PREFIX}" | cut -f1)"
log "done"
