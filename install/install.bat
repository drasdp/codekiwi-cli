@echo off
setlocal EnableDelayedExpansion

:: Configuration
set GITHUB_REPO=drasdp/codekiwi-cli
set INSTALL_DIR=%USERPROFILE%\.codekiwi
set BIN_NAME=codekiwi.exe
set BINARY_NAME=codekiwi-windows-amd64.exe

:: Colors (not supported in cmd, but we'll use text markers)
echo =======================================
echo     CodeKiwi CLI Installation
echo =======================================
echo.

:: Check Docker
echo [INFO] Checking Docker installation...
docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not installed. Please install Docker Desktop first.
    echo         https://www.docker.com/products/docker-desktop
    exit /b 1
)

docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker daemon is not running. Please start Docker Desktop.
    exit /b 1
)
echo [OK] Docker is installed and running
echo.

:: Create install directory
echo [INFO] Creating installation directory...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
if not exist "%INSTALL_DIR%\bin" mkdir "%INSTALL_DIR%\bin"
if not exist "%INSTALL_DIR%\instances" mkdir "%INSTALL_DIR%\instances"

:: Get latest release info
echo [INFO] Fetching latest release information...
set TEMP_JSON=%TEMP%\codekiwi_release.json

:: Download release info
curl -s https://api.github.com/repos/%GITHUB_REPO%/releases/latest > %TEMP_JSON% 2>nul
if errorlevel 1 (
    echo [ERROR] Failed to fetch release information
    exit /b 1
)

:: Extract download URL and version using PowerShell
echo [INFO] Processing release information...
for /f "delims=" %%i in ('powershell -NoProfile -Command "try { $json = Get-Content '%TEMP_JSON%' -Raw | ConvertFrom-Json; $asset = $json.assets | Where-Object {$_.name -eq '%BINARY_NAME%'}; if($asset) { Write-Host $asset.browser_download_url } else { Write-Host 'ERROR' } } catch { Write-Host 'ERROR' }"') do set DOWNLOAD_URL=%%i

if "%DOWNLOAD_URL%"=="ERROR" (
    echo [ERROR] Failed to find Windows release
    echo         Please check https://github.com/%GITHUB_REPO%/releases
    del %TEMP_JSON% 2>nul
    exit /b 1
)

if "%DOWNLOAD_URL%"=="" (
    echo [ERROR] Failed to parse release information
    del %TEMP_JSON% 2>nul
    exit /b 1
)

:: Get version
for /f "delims=" %%i in ('powershell -NoProfile -Command "try { $json = Get-Content '%TEMP_JSON%' -Raw | ConvertFrom-Json; Write-Host $json.tag_name } catch { Write-Host 'unknown' }"') do set VERSION=%%i
echo [INFO] Latest version: %VERSION%

:: Download binary
echo [INFO] Downloading CodeKiwi binary...
curl -fsSL %DOWNLOAD_URL% -o "%INSTALL_DIR%\bin\%BIN_NAME%" 2>nul
if errorlevel 1 (
    echo [ERROR] Failed to download CodeKiwi binary
    del %TEMP_JSON% 2>nul
    exit /b 1
)
echo [OK] Binary downloaded successfully

:: Download configuration files
echo [INFO] Downloading configuration files...

:: Download docker-compose.yaml
curl -fsSL https://raw.githubusercontent.com/%GITHUB_REPO%/main/docker-compose.yaml -o "%INSTALL_DIR%\docker-compose.yaml" 2>nul
if errorlevel 1 (
    echo [WARN] Failed to download docker-compose.yaml
)

:: Download config.env
curl -fsSL https://raw.githubusercontent.com/%GITHUB_REPO%/main/config.env -o "%INSTALL_DIR%\config.env" 2>nul
if errorlevel 1 (
    echo [WARN] Failed to download config.env
)

echo [OK] Configuration files downloaded
echo.

:: Add to PATH
echo [INFO] Adding CodeKiwi to PATH...
set NEW_PATH=%INSTALL_DIR%\bin

:: Get current user PATH
for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v PATH 2^>nul') do set USER_PATH=%%b

:: Check if already in PATH
echo !USER_PATH! | findstr /C:"%NEW_PATH%" >nul
if errorlevel 1 (
    :: Add to PATH
    if defined USER_PATH (
        setx PATH "%USER_PATH%;%NEW_PATH%" >nul 2>&1
    ) else (
        setx PATH "%NEW_PATH%" >nul 2>&1
    )
    echo [OK] Added to PATH: %NEW_PATH%

    :: Also update current session PATH
    set PATH=%PATH%;%NEW_PATH%
) else (
    echo [INFO] Already in PATH: %NEW_PATH%
)

:: Clean up temp files
del %TEMP_JSON% 2>nul

:: Create a batch file wrapper (optional, for better cmd integration)
echo @echo off > "%INSTALL_DIR%\bin\codekiwi.bat"
echo "%INSTALL_DIR%\bin\codekiwi.exe" %%* >> "%INSTALL_DIR%\bin\codekiwi.bat"

:: Success message
echo.
echo =======================================
echo [OK] CodeKiwi installed successfully!
echo =======================================
echo.
echo Installation directory: %INSTALL_DIR%
echo.
echo To start using CodeKiwi:
echo   1. Open a NEW Command Prompt window
echo   2. Run: codekiwi --help
echo.
echo To start a project:
echo   codekiwi start [project-directory]
echo.
echo Note: You must open a new Command Prompt for PATH changes to take effect.
echo.

:: If running from temp (one-line install), schedule self-deletion
if "%~dp0"=="%TEMP%\" (
    timeout /t 2 /nobreak >nul
    del "%~f0" >nul 2>&1
)

pause