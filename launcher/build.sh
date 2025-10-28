#!/bin/bash

echo "================================"
echo "CodeKiwi Launcher Build Script"
echo "================================"
echo

# Check Python
if ! command -v python3 &> /dev/null; then
    echo "[ERROR] Python3 is not installed"
    exit 1
fi

# Check/Install PyInstaller
echo "Checking PyInstaller..."
if ! pip3 show pyinstaller &> /dev/null; then
    echo "Installing PyInstaller..."
    pip3 install pyinstaller
    if [ $? -ne 0 ]; then
        echo "[ERROR] Failed to install PyInstaller"
        exit 1
    fi
fi

# Change to src directory
cd src || {
    echo "[ERROR] src directory not found"
    exit 1
}

# Clean previous build
echo
echo "Cleaning previous build..."
rm -rf build dist __pycache__ *.pyc

# Build with PyInstaller
echo
echo "Building executable..."
echo "This may take a few minutes..."
echo

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    pyinstaller launcher.spec --onefile
else
    # Linux
    pyinstaller launcher.spec
fi

if [ $? -ne 0 ]; then
    echo
    echo "[ERROR] Build failed!"
    exit 1
fi

echo
echo "================================"
echo "+ Build completed successfully!"
echo "================================"
echo
echo "Output location: launcher/src/dist/"

if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Executable: src/dist/CodeKiwiLauncher"
else
    echo "Executable: src/dist/CodeKiwiLauncher"
fi

echo
echo "To run:"
echo "  ./src/dist/CodeKiwiLauncher"
echo