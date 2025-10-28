@echo off
echo ================================
echo CodeKiwi Launcher Build Script
echo ================================
echo.

REM Check Python
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Python is not installed or not in PATH
    pause
    exit /b 1
)

REM Check/Install PyInstaller
echo Checking PyInstaller...
pip show pyinstaller >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing PyInstaller...
    pip install pyinstaller
    if %errorlevel% neq 0 (
        echo [ERROR] Failed to install PyInstaller
        pause
        exit /b 1
    )
)

REM Change to src directory
cd src
if %errorlevel% neq 0 (
    echo [ERROR] src directory not found
    pause
    exit /b 1
)

REM Clean previous build
echo.
echo Cleaning previous build...
if exist build rmdir /s /q build
if exist dist rmdir /s /q dist
if exist __pycache__ rmdir /s /q __pycache__
if exist *.pyc del /q *.pyc

REM Build with PyInstaller
echo.
echo Building executable...
echo This may take a few minutes...
echo.

pyinstaller launcher.spec

if %errorlevel% neq 0 (
    echo.
    echo [ERROR] Build failed!
    pause
    exit /b 1
)

echo.
echo ================================
echo + Build completed successfully!
echo ================================
echo.
echo Output location: launcher\src\dist\
echo Executable: src\dist\CodeKiwiLauncher.exe
echo.
echo To run:
echo   src\dist\CodeKiwiLauncher.exe
echo.
pause