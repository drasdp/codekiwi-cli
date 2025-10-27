# CodeKiwi GUI Installer

A cross-platform graphical installer for CodeKiwi CLI, built with Python, tkinter, and PyInstaller.

## Features

- **Cross-platform**: Supports Windows (.exe), macOS (.app), and Linux
- **Wizard-style installation**: Step-by-step guided installation process
- **Prerequisite checking**: Validates Docker and curl installation
- **Automatic setup**: Downloads and configures CodeKiwi CLI automatically
- **User-friendly**: Clean GUI with progress indicators and detailed logs

## Installation Steps

The installer guides users through the following steps:

1. **Welcome**: System information and prerequisites overview
2. **Docker Check**: Validates Docker installation and daemon status
   - Must have Docker 20.10+ installed and running
   - Must have Docker Compose 2.0+
   - Provides download links for Docker Desktop
3. **curl Check**: Validates curl installation
   - Automatically checks for curl
   - Provides installation instructions if missing
4. **Installation**: Downloads and installs CodeKiwi
   - Downloads latest release binary
   - Downloads configuration files
   - Sets up PATH configuration
5. **Finish**: Installation complete with next steps

## Prerequisites Checked

The installer validates the following prerequisites based on analysis of the CodeKiwi codebase:

### System Requirements
- **Docker**: Version 20.10+ (required)
- **Docker Compose**: Version 2.0+ (required)
- **curl**: For downloading files
- **Memory**: 2GB minimum
- **Disk Space**: ~1GB for Docker runtime image

### Platform Support
- **macOS**: All versions (x86_64 and ARM64/M1+)
- **Linux**: All distributions (x86_64 and ARM64)
- **Windows**: Windows 10+ (x86_64)

## Building the Installer

### Requirements
- Python 3.8 or later
- pip (Python package manager)

### Build on macOS/Linux

```bash
cd installer
./build.sh
```

The script will:
1. Create a virtual environment
2. Install dependencies (PyInstaller)
3. Build the executable

Output:
- **macOS**: `src/dist/CodeKiwiInstaller.app`
- **Linux**: `src/dist/CodeKiwiInstaller`

### Build on Windows

```cmd
cd installer
build.bat
```

The script will:
1. Create a virtual environment
2. Install dependencies (PyInstaller)
3. Build the executable

Output:
- **Windows**: `src\dist\CodeKiwiInstaller.exe`

## Running the Installer

### macOS
```bash
open src/dist/CodeKiwiInstaller.app
```

Or double-click the app in Finder.

### Linux
```bash
./src/dist/CodeKiwiInstaller
```

Or double-click the file if your desktop environment supports it.

### Windows
```cmd
src\dist\CodeKiwiInstaller.exe
```

Or double-click the .exe file.

## Manual Build (Advanced)

If you prefer to build manually:

```bash
# Create virtual environment
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Build
cd src
pyinstaller installer.spec
cd ..
```

## What the Installer Does

1. **Creates installation directory**: `~/.codekiwi/`
2. **Downloads CodeKiwi binary**: Platform-specific binary from GitHub releases
3. **Downloads configuration files**:
   - `docker-compose.yaml`
   - `config.env`
4. **Configures PATH**:
   - **macOS/Linux**: Updates shell config files (.bashrc, .zshrc, etc.)
   - **Windows**: Updates user PATH environment variable
5. **Creates symlink** (macOS/Linux): `/usr/local/bin/codekiwi` (if permissions allow)

## Project Structure

```
installer/
├── README.md              # This file
├── requirements.txt       # Python dependencies (PyInstaller)
├── build.sh              # Build script for macOS/Linux
├── build.bat             # Build script for Windows
└── src/
    ├── installer.py      # Main installer application
    └── installer.spec    # PyInstaller configuration
```

## Installer Features Detail

### Docker Validation
- Checks if Docker is installed (`docker --version`)
- Verifies Docker daemon is running (`docker info`)
- Validates Docker Compose availability (v2 or v1)
- Provides download links for Docker Desktop
- **Blocks progression** until Docker is properly installed and running

### curl Validation
- Automatically checks curl installation
- Provides platform-specific installation instructions
- Required for downloading CodeKiwi files

### Installation Process
- Shows real-time progress with progress bar
- Displays detailed installation log
- Handles errors gracefully with clear error messages
- Downloads from official GitHub releases
- Supports all platform/architecture combinations

### PATH Configuration
- **Windows**: Uses PowerShell to update user PATH environment variable
- **macOS/Linux**: Updates common shell config files
- Creates symlinks for easy access
- Provides instructions for applying changes

## Troubleshooting

### Docker Not Detected
- Ensure Docker Desktop is installed and running
- Verify Docker daemon is started: `docker info`
- Restart Docker Desktop if needed

### Permission Errors (macOS/Linux)
- The installer may need `sudo` for creating symlinks in `/usr/local/bin`
- Alternatively, add `~/.codekiwi/bin` to your PATH manually

### Network Issues
- Ensure stable internet connection
- GitHub API and release downloads require network access
- Check firewall settings if downloads fail

### Windows PATH Not Updated
- Restart PowerShell/Command Prompt after installation
- Manually add `C:\Users\YourName\.codekiwi\bin` to PATH if needed

## Technical Details

### Dependencies
- **tkinter**: GUI framework (included with Python)
- **PyInstaller**: Creates standalone executables
- **urllib**: For downloading files
- **subprocess**: For running system commands
- **threading**: For non-blocking operations

### Binary Detection
The installer automatically detects:
- Operating system (Windows, macOS, Linux)
- Architecture (x86_64/amd64, ARM64/aarch64)
- Downloads the appropriate binary for the platform

### GitHub Integration
- Fetches latest release from GitHub API
- Downloads platform-specific binaries
- Downloads configuration files from main branch

## License

This installer is part of the CodeKiwi CLI project.

## Related Links

- **CodeKiwi CLI Repository**: https://github.com/drasdp/codekiwi-cli
- **Docker Desktop**: https://www.docker.com/products/docker-desktop
- **Python**: https://www.python.org/
- **PyInstaller**: https://pyinstaller.org/

## Support

For issues or questions:
- Open an issue at: https://github.com/drasdp/codekiwi-cli/issues
- Check the main README: https://github.com/drasdp/codekiwi-cli
