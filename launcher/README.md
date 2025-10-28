# CodeKiwi Launcher

Simple GUI launcher for CodeKiwi CLI - Start, stop, and manage your CodeKiwi projects with ease!

## Features

- 🚀 **Simple & Clean UI** - Minimal, straightforward interface
- 📁 **Project Browser** - Easy project folder selection
- 🤖 **Auto-Detection** - CodeKiwi automatically detects project type
- 📊 **Live Output** - Real-time command output display with UTF-8 support
- ⚡ **Quick Actions** - Start, Stop, List, Update with one click

## Requirements

- Python 3.6+
- CodeKiwi CLI installed and in PATH
- Docker Desktop running

## Development

### Run from source
```bash
cd launcher/src
python launcher.py
```

### Build executable

#### Windows
```batch
cd launcher
build.bat
```

#### Mac/Linux
```bash
cd launcher
chmod +x build.sh
./build.sh
```

## Usage

1. **Launch the application**
   - Run `CodeKiwiLauncher.exe` (Windows) or `CodeKiwiLauncher` (Mac/Linux)

2. **Select a project**
   - Click "Browse" to select your project folder

3. **Start your project**
   - Click "START" to launch CodeKiwi
   - CodeKiwi will automatically detect your project type
   - View real-time output in the console area

4. **Manage instances**
   - Click "LIST" to see all running projects
   - Click "STOP" to stop the selected project (or all if none selected)
   - Click "UPDATE" to update CodeKiwi CLI and Docker images

## Keyboard Shortcuts

- `Ctrl+O` - Browse for project folder
- `Ctrl+S` - Start project
- `Ctrl+Q` - Stop project
- `Ctrl+L` - List instances
- `Ctrl+C` - Clear output

## Architecture

```
launcher/
├── src/
│   ├── launcher.py        # Main application (single file)
│   └── launcher.spec      # PyInstaller configuration
├── build.bat              # Windows build script
├── build.sh              # Mac/Linux build script
├── requirements.txt      # Python dependencies
└── README.md            # This file
```

## Philosophy: KISS (Keep It Simple, Stupid)

This launcher follows the KISS principle:
- **Single file** - All code in one file for simplicity
- **No dependencies** - Only uses Python standard library (tkinter)
- **Minimal UI** - Just the essentials, no fluff
- **Direct integration** - Simply wraps CLI commands

## Troubleshooting

### "codekiwi command not found"
- Make sure CodeKiwi CLI is installed
- Check that it's in your system PATH

### "Docker is not running"
- Start Docker Desktop
- Wait for Docker to fully start

### Build fails
- Ensure Python 3.6+ is installed
- Install PyInstaller: `pip install pyinstaller`

## License

Part of CodeKiwi project