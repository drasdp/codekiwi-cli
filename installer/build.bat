@echo off
REM Build script for CodeKiwi Installer (Windows)

echo ================================
echo CodeKiwi Installer Build Script
echo ================================
echo.

REM Check if Python is installed
python --version >nul 2>&1
if errorlevel 1 (
    echo Error: Python is not installed
    echo Please install Python 3.8 or later from https://www.python.org/
    pause
    exit /b 1
)

python --version
echo.

REM Create virtual environment if it doesn't exist
if not exist "venv" (
    echo Creating virtual environment...
    python -m venv venv
    echo + Virtual environment created
    echo.
)

REM Activate virtual environment
echo Activating virtual environment...
call venv\Scripts\activate.bat
echo.

REM Install dependencies
echo Installing dependencies...
python -m pip install -q --upgrade pip
python -m pip install -q -r requirements.txt
echo + Dependencies installed
echo.

REM Clean previous builds
echo Cleaning previous builds...
if exist "src\build" rmdir /s /q src\build
if exist "src\dist" rmdir /s /q src\dist
if exist "src\__pycache__" rmdir /s /q src\__pycache__
echo + Cleaned
echo.

REM Build executable
echo Building executable with PyInstaller...
cd src
pyinstaller installer.spec
cd ..
echo.

REM Show results
if exist "src\dist\CodeKiwiInstaller.exe" (
    echo ================================
    echo + Build completed successfully!
    echo ================================
    echo.
    echo Output location: installer\src\dist\
    echo Executable: src\dist\CodeKiwiInstaller.exe
    echo.
    echo To run:
    echo   src\dist\CodeKiwiInstaller.exe
    echo.
) else (
    echo x Build failed!
    pause
    exit /b 1
)

REM Deactivate virtual environment
call venv\Scripts\deactivate.bat

pause
