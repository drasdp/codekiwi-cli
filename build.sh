#!/bin/bash
# Build script for CodeKiwi Launcher (macOS/Linux)

echo "Installing PyInstaller..."
uv pip install pyinstaller

echo "Building CodeKiwi executable..."
uv run pyinstaller --onefile --windowed \
  --name CodeKiwi \
  --add-data "core:core" \
  main.py

echo ""
echo "Build complete!"
echo "Executable location: dist/CodeKiwi"
