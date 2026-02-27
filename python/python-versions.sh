#!/bin/sh
#
# python-versions.sh — Discover available CPython versions and emit a
#                       docker-bake override JSON for PYTHON_VERSIONS_LIST.
#
# Usage:
#   ./python-versions.sh              # prints JSON to stdout
#   ./python-versions.sh > override.json
#
# Dependencies: uv, jq
#
set -eu

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------
readonly PYTHON_IMPL="cpython"          # implementation prefix from `uv python list`
readonly EXCLUDE_PATTERN="freethreaded" # variants to skip

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

log() {
    printf '[python-versions] %s\n' "$1" >&2
}

die() {
    printf '[python-versions] ERROR: %s\n' "$1" >&2
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
require_cmd jq
log "dependencies OK (uv=$(uv --version 2>/dev/null || echo '?'), jq=$(jq --version 2>/dev/null || echo '?'))"

# ---------------------------------------------------------------------------
# Step 1 — Fetch all managed CPython versions from uv
# ---------------------------------------------------------------------------
log "fetching managed python versions from uv..."
raw_versions=$(uv python list --python-preference only-managed) \
    || die "failed to list python versions via uv"
raw_count=$(printf '%s\n' "${raw_versions}" | wc -l | tr -d ' ')
log "found ${raw_count} total entries from uv"

# ---------------------------------------------------------------------------
# Step 2 — Filter to standard (non-freethreaded) CPython entries
# ---------------------------------------------------------------------------
filtered_versions=$(
    printf '%s\n' "${raw_versions}" \
        | grep "^${PYTHON_IMPL}-" \
        | grep -v "+${EXCLUDE_PATTERN}"
)

if [ -z "${filtered_versions}" ]; then
    die "no ${PYTHON_IMPL} versions found (excluding ${EXCLUDE_PATTERN})"
fi
filtered_count=$(printf '%s\n' "${filtered_versions}" | wc -l | tr -d ' ')
log "filtered to ${filtered_count} ${PYTHON_IMPL} entries (excluded ${EXCLUDE_PATTERN})"

# ---------------------------------------------------------------------------
# Step 3 — Extract bare version numbers  (e.g. "3.13.2")
# ---------------------------------------------------------------------------
version_numbers=$(printf '%s\n' "${filtered_versions}" | cut -d- -f2)
log "extracted version numbers (newest first): $(printf '%s\n' "${version_numbers}" | head -3 | tr '\n' ' ')..."

# ---------------------------------------------------------------------------
# Step 4 — Deduplicate to latest patch per minor, sort ascending, add "latest"
# ---------------------------------------------------------------------------
# uv outputs newest first, so the first occurrence of each minor is the
# latest patch.  We keep that, reverse into ascending order, and append
# a synthetic "latest" entry that the bake file maps to the default.
# ---------------------------------------------------------------------------
bake_json=$(
    printf '%s\n' "${version_numbers}" | jq -Rs '
        # Split lines and drop empties
        split("\n") | map(select(length > 0))

        # Keep only the first (newest) patch per minor series
        | reduce .[] as $v ({seen: {}, out: []};
            ($v | split(".")[0:2] | join(".")) as $minor |
            if .seen[$minor] then .
            else .seen[$minor] = true | .out += [$v]
            end
          )
        | .out

        # Reverse to ascending order and append the "latest" tag
        | reverse
        | . + ["latest"]

        # Wrap in the docker-bake variable override structure
        | { variable: { PYTHON_VERSIONS_LIST: { default: . } } }
    '
) || die "jq processing failed"

selected_versions=$(printf '%s\n' "${bake_json}" | jq -r '.variable.PYTHON_VERSIONS_LIST.default[]' 2>/dev/null | tr '\n' ' ')
log "selected versions: ${selected_versions}"

# ---------------------------------------------------------------------------
# Step 5 — Output
# ---------------------------------------------------------------------------
log "emitting docker-bake override JSON to stdout"
printf '%s\n' "${bake_json}"
