@echo off
REM Build script untuk Wails desktop app
REM Script ini akan copy frontend/dist ke backend/frontend-dist sebelum build

echo Copying frontend/dist to backend/frontend-dist...
if exist "frontend-dist" rmdir /s /q "frontend-dist"
mkdir "frontend-dist"
xcopy /E /I /Y "..\frontend\dist\*" "frontend-dist\"

echo Copying .env file to build/bin for production...
if exist ".env" (
    if not exist "build\bin" mkdir "build\bin"
    copy ".env" "build\bin\.env"
    echo .env file copied to build/bin
) else (
    echo Warning: .env file not found in backend folder
)

echo Building Wails application...
wails build

echo Done!

