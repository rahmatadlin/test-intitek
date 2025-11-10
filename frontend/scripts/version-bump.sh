#!/bin/bash
# Linux/Mac shell script untuk version bump
# Usage: ./version-bump.sh [patch|minor|major]

VERSION_TYPE=${1:-patch}

node scripts/version-bump.js "$VERSION_TYPE"

