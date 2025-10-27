#!/bin/bash
# Build script for CodeKiwi Installer (macOS/Linux)

set -e

echo "================================"
echo "CodeKiwi Installer Build Script"
echo "================================"
echo ""

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "Error: Python 3 is not installed"
    echo "Please install Python 3.8 or later"
    exit 1
fi

echo "Python version: $(python3 --version)"
echo ""

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv venv
    echo "✓ Virtual environment created"
    echo ""
fi

# Activate virtual environment
echo "Activating virtual environment..."
source venv/bin/activate
echo ""

# Install dependencies
echo "Installing dependencies..."
pip install -q --upgrade pip
pip install -q -r requirements.txt
echo "✓ Dependencies installed"
echo ""

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf src/build src/dist src/__pycache__
echo "✓ Cleaned"
echo ""

# Build executable
echo "Building executable with PyInstaller..."
cd src
pyinstaller installer.spec
cd ..
echo ""

# Show results
if [ -f "src/dist/CodeKiwiInstaller" ] || [ -d "src/dist/CodeKiwiInstaller.app" ]; then
    echo "================================"
    echo "✓ Build completed successfully!"
    echo "================================"
    echo ""
    echo "Output location: installer/src/dist/"
    echo ""

    if [ "$(uname)" == "Darwin" ]; then
        echo "macOS App Bundle: src/dist/CodeKiwiInstaller.app"
        echo ""
        echo "To run:"
        echo "  open src/dist/CodeKiwiInstaller.app"
    else
        echo "Executable: src/dist/CodeKiwiInstaller"
        echo ""
        echo "To run:"
        echo "  ./src/dist/CodeKiwiInstaller"
    fi
    echo ""
else
    echo "✗ Build failed!"
    exit 1
fi

# Deactivate virtual environment
deactivate
