#!/bin/bash
# Build script untuk Wails desktop app
# Script ini akan copy frontend/dist ke backend/frontend-dist sebelum build

echo "Copying frontend/dist to backend/frontend-dist..."
rm -rf frontend-dist
mkdir -p frontend-dist
cp -r ../frontend/dist/* frontend-dist/

echo "Copying .env file to build/bin for production..."
if [ -f ".env" ]; then
    mkdir -p build/bin
    cp ".env" "build/bin/.env"
    echo ".env file copied to build/bin"
else
    echo "Warning: .env file not found in backend folder"
fi

echo "Building Wails application..."
wails build

echo "Done!"

