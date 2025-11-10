@echo off
REM Windows batch script untuk version bump
REM Usage: version-bump.bat [patch|minor|major]

set VERSION_TYPE=%1
if "%VERSION_TYPE%"=="" set VERSION_TYPE=patch

node scripts/version-bump.js %VERSION_TYPE%

