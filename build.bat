@echo off
REM Build script for CodeKiwi Launcher (Windows)

echo Installing PyInstaller...
uv pip install pyinstaller

echo Building CodeKiwi executable...
uv run pyinstaller --onefile --windowed --name CodeKiwi --add-data "core;core" main.py

echo.
echo Build complete!
echo Executable location: dist\CodeKiwi.exe
pause
