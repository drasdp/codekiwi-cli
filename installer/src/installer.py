#!/usr/bin/env python3
"""
CodeKiwi GUI Installer
A cross-platform installation wizard for CodeKiwi CLI
Supports: Windows (.exe), macOS, Linux
"""

import tkinter as tk
from tkinter import ttk, messagebox, scrolledtext
import subprocess
import platform
import os
import sys
import urllib.request
import json
import shutil
import threading
import time
import ctypes


def get_clean_env():
    """
    Get a clean environment for subprocess calls in PyInstaller context.
    This fixes the issue where subprocess can't find system executables like curl.
    """
    env = dict(os.environ)

    # If running as PyInstaller bundle
    if getattr(sys, 'frozen', False):
        # Remove _MEIPASS paths from PATH to avoid DLL conflicts
        meipass = getattr(sys, '_MEIPASS', None)
        if meipass and 'PATH' in env:
            paths = env['PATH'].split(os.pathsep)
            # Filter out PyInstaller temporary paths
            paths = [p for p in paths if not p.startswith(meipass)]
            env['PATH'] = os.pathsep.join(paths)

        # On Windows, reset DLL directory to allow system programs to find their DLLs
        if sys.platform == 'win32':
            try:
                # Reset DLL search path to default
                ctypes.windll.kernel32.SetDllDirectoryW(None)
            except Exception:
                pass  # Ignore if this fails

    return env


def get_subprocess_kwargs():
    """
    Get proper subprocess kwargs for PyInstaller GUI mode.
    Fixes handle errors in console=False mode.
    """
    kwargs = {
        'stdin': subprocess.PIPE,    # Required to prevent handle errors
        'stdout': subprocess.PIPE,   # Required to capture output
        'stderr': subprocess.PIPE,   # Required to capture errors
    }

    # On Windows, hide console window for subprocess
    if sys.platform == 'win32':
        si = subprocess.STARTUPINFO()
        si.dwFlags |= subprocess.STARTF_USESHOWWINDOW
        si.wShowWindow = subprocess.SW_HIDE  # Hide window
        kwargs['startupinfo'] = si

    # Use clean environment
    kwargs['env'] = get_clean_env()

    return kwargs


def find_curl_path():
    """Find the absolute path to curl executable."""
    if sys.platform == 'win32':
        # Try common Windows locations for curl.exe
        possible_paths = [
            r'C:\Windows\System32\curl.exe',
            r'C:\Windows\SysWOW64\curl.exe',
            os.path.join(os.environ.get('WINDIR', 'C:\\Windows'), 'System32', 'curl.exe'),
        ]

        # Check each path
        for path in possible_paths:
            if os.path.exists(path) and os.path.isfile(path):
                return path

        # Try to find curl using where.exe
        try:
            kwargs = get_subprocess_kwargs()
            result = subprocess.run(['where.exe', 'curl'],
                                  timeout=5,
                                  **kwargs)
            if result.returncode == 0:
                output = result.stdout.decode('utf-8', errors='ignore').strip()
                paths = output.split('\n')
                for path in paths:
                    path = path.strip()
                    if path and os.path.exists(path):
                        return path
        except Exception:
            pass

        # Last resort: try just 'curl' and let Windows search PATH
        return 'curl'

    else:
        # Unix-like systems (Linux, macOS)
        # Try which command first
        try:
            kwargs = get_subprocess_kwargs()
            result = subprocess.run(['which', 'curl'],
                                  timeout=5,
                                  **kwargs)
            if result.returncode == 0:
                path = result.stdout.decode('utf-8', errors='ignore').strip()
                if path and os.path.exists(path):
                    return path
        except Exception:
            pass

        # Common Unix locations
        possible_paths = [
            '/usr/bin/curl',
            '/usr/local/bin/curl',
            '/opt/homebrew/bin/curl',  # macOS with Homebrew on Apple Silicon
            '/opt/local/bin/curl',     # macOS with MacPorts
        ]

        for path in possible_paths:
            if os.path.exists(path) and os.path.isfile(path):
                return path

        # Default to 'curl' and let the system find it
        return 'curl'


class CodeKiwiInstaller:
    def __init__(self, root):
        self.root = root
        self.root.title("CodeKiwi Installer")
        self.root.geometry("700x500")
        self.root.resizable(False, False)

        # Set the style
        self.style = ttk.Style()
        self.style.theme_use('clam')

        # Variables
        self.current_page = 0
        self.os_type = platform.system()
        self.arch = platform.machine().lower()
        self.install_dir = self.get_default_install_dir()

        # Create main container
        self.main_frame = ttk.Frame(root, padding="20")
        self.main_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))

        # Configure grid weights
        root.columnconfigure(0, weight=1)
        root.rowconfigure(0, weight=1)

        # Pages
        self.pages = [
            self.create_welcome_page,
            self.create_docker_check_page,
            self.create_curl_check_page,
            self.create_install_page,
            self.create_finish_page
        ]

        # Navigation buttons frame
        self.nav_frame = ttk.Frame(self.main_frame)
        self.nav_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(20, 0))

        self.back_btn = ttk.Button(self.nav_frame, text="< Back", command=self.prev_page, state=tk.DISABLED)
        self.back_btn.grid(row=0, column=0, padx=5)

        self.next_btn = ttk.Button(self.nav_frame, text="Next >", command=self.next_page)
        self.next_btn.grid(row=0, column=1, padx=5)

        self.cancel_btn = ttk.Button(self.nav_frame, text="Cancel", command=self.cancel_installation)
        self.cancel_btn.grid(row=0, column=2, padx=5)

        # Page counter label
        self.page_label = ttk.Label(self.nav_frame, text="")
        self.page_label.grid(row=0, column=3, padx=20)

        # Content frame (where pages will be displayed)
        self.content_frame = ttk.Frame(self.main_frame)
        self.content_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        self.main_frame.rowconfigure(0, weight=1)
        self.main_frame.columnconfigure(0, weight=1)

        # Show first page
        self.show_page()

    def get_default_install_dir(self):
        """Get default installation directory based on OS"""
        home = os.path.expanduser("~")
        return os.path.join(home, ".codekiwi")

    def clear_content(self):
        """Clear the content frame"""
        for widget in self.content_frame.winfo_children():
            widget.destroy()

    def show_page(self):
        """Display the current page"""
        self.clear_content()
        self.pages[self.current_page]()
        self.update_navigation()
        self.page_label.config(text=f"Step {self.current_page + 1} of {len(self.pages)}")

    def update_navigation(self):
        """Update navigation button states"""
        self.back_btn.config(state=tk.NORMAL if self.current_page > 0 else tk.DISABLED)

        if self.current_page == len(self.pages) - 1:
            self.next_btn.config(text="Finish", command=self.finish_installation)
        else:
            self.next_btn.config(text="Next >", command=self.next_page)

    def next_page(self):
        """Go to next page"""
        if self.current_page < len(self.pages) - 1:
            # Validate current page before moving
            if self.validate_current_page():
                self.current_page += 1
                self.show_page()

    def prev_page(self):
        """Go to previous page"""
        if self.current_page > 0:
            self.current_page -= 1
            self.show_page()

    def validate_current_page(self):
        """Validate current page requirements"""
        if self.current_page == 1:  # Docker check page
            return self.check_docker_validation()
        elif self.current_page == 2:  # Curl check page
            return self.check_curl_validation()
        return True

    def cancel_installation(self):
        """Cancel the installation"""
        if messagebox.askyesno("Cancel Installation", "Are you sure you want to cancel the installation?"):
            self.root.quit()

    def finish_installation(self):
        """Finish the installation"""
        messagebox.showinfo("Installation Complete",
                          "CodeKiwi has been successfully installed!\n\n"
                          "You can now use 'codekiwi' command from your terminal.")
        self.root.quit()

    # ==================== PAGE 1: Welcome ====================
    def create_welcome_page(self):
        """Create welcome page"""
        frame = ttk.Frame(self.content_frame)
        frame.pack(fill=tk.BOTH, expand=True)

        # Title
        title = ttk.Label(frame, text="Welcome to CodeKiwi Installer",
                         font=('Arial', 16, 'bold'))
        title.pack(pady=20)

        # Description
        desc_frame = ttk.Frame(frame)
        desc_frame.pack(fill=tk.BOTH, expand=True, padx=20, pady=10)

        desc_text = """CodeKiwi is a Docker-based AI-powered development environment CLI.

This wizard will guide you through the installation process.

Prerequisites:
• Docker 20.10+ (required)
• Docker Compose 2.0+
• curl (for downloading)
• 2GB RAM minimum
• ~1GB disk space

System Information:
• Operating System: {os}
• Architecture: {arch}
• Installation Directory: {install_dir}

Click 'Next' to begin the installation process.""".format(
            os=self.os_type,
            arch=self.arch,
            install_dir=self.install_dir
        )

        desc_label = ttk.Label(desc_frame, text=desc_text, justify=tk.LEFT)
        desc_label.pack(anchor=tk.W)

    # ==================== PAGE 2: Docker Check ====================
    def create_docker_check_page(self):
        """Create Docker check page"""
        frame = ttk.Frame(self.content_frame)
        frame.pack(fill=tk.BOTH, expand=True)

        # Title
        title = ttk.Label(frame, text="Docker Installation Check",
                         font=('Arial', 14, 'bold'))
        title.pack(pady=20)

        # Description
        desc = ttk.Label(frame, text="Docker is required to run CodeKiwi.\n"
                                     "Please ensure Docker is installed and running.",
                        justify=tk.CENTER)
        desc.pack(pady=10)

        # Check results frame
        self.docker_result_frame = ttk.LabelFrame(frame, text="Check Results", padding=10)
        self.docker_result_frame.pack(fill=tk.BOTH, expand=True, padx=20, pady=10)

        self.docker_status_label = ttk.Label(self.docker_result_frame,
                                            text="Click 'Check Docker' to verify installation",
                                            wraplength=600)
        self.docker_status_label.pack(pady=10)

        # Check button
        check_btn = ttk.Button(frame, text="Check Docker", command=self.check_docker)
        check_btn.pack(pady=10)

        # Install Docker instructions
        install_frame = ttk.LabelFrame(frame, text="Don't have Docker?", padding=10)
        install_frame.pack(fill=tk.X, padx=20, pady=10)

        docker_links = self.get_docker_download_link()
        install_text = f"Download Docker Desktop from:\n{docker_links}"

        install_label = ttk.Label(install_frame, text=install_text,
                                 foreground="blue", cursor="hand2")
        install_label.pack()
        install_label.bind("<Button-1>", lambda e: self.open_url(docker_links))

        # Store validation state
        self.docker_validated = False

    def get_docker_download_link(self):
        """Get Docker download link based on OS"""
        if self.os_type == "Darwin":  # macOS
            return "https://docs.docker.com/desktop/install/mac-install/"
        elif self.os_type == "Windows":
            return "https://docs.docker.com/desktop/install/windows-install/"
        else:  # Linux
            return "https://docs.docker.com/desktop/install/linux-install/"

    def open_url(self, url):
        """Open URL in default browser"""
        import webbrowser
        webbrowser.open(url)

    def check_docker(self):
        """Check if Docker is installed and running"""
        self.docker_status_label.config(text="Checking Docker installation...", foreground="blue")
        self.root.update()

        # Run checks in a thread to avoid blocking UI
        thread = threading.Thread(target=self._check_docker_thread)
        thread.daemon = True
        thread.start()

    def _check_docker_thread(self):
        """Thread function to check Docker"""
        results = []
        all_passed = True
        kwargs = get_subprocess_kwargs()

        # Check 1: Docker command exists
        try:
            result = subprocess.run(["docker", "--version"],
                                  timeout=5,
                                  **kwargs)
            if result.returncode == 0:
                version = result.stdout.decode('utf-8', errors='ignore').strip()
                results.append(f"✓ Docker installed: {version}")
            else:
                results.append("✗ Docker command not found")
                all_passed = False
        except (subprocess.TimeoutExpired, FileNotFoundError) as e:
            results.append(f"✗ Docker command not found: {str(e)}")
            all_passed = False

        # Check 2: Docker daemon running
        try:
            result = subprocess.run(["docker", "info"],
                                  timeout=10,
                                  **kwargs)
            if result.returncode == 0:
                results.append("✓ Docker daemon is running")
            else:
                error_msg = result.stderr.decode('utf-8', errors='ignore').strip()
                if 'permission denied' in error_msg.lower():
                    results.append("✗ Docker daemon requires permissions\n  Add user to docker group or run with sudo")
                else:
                    results.append("✗ Docker daemon is not running\n  Please start Docker Desktop")
                all_passed = False
        except (subprocess.TimeoutExpired, FileNotFoundError) as e:
            results.append(f"✗ Docker daemon is not running: {str(e)}")
            all_passed = False

        # Check 3: Docker Compose
        compose_found = False
        try:
            # Try docker compose (v2)
            result = subprocess.run(["docker", "compose", "version"],
                                  timeout=5,
                                  **kwargs)
            if result.returncode == 0:
                version = result.stdout.decode('utf-8', errors='ignore').strip()
                results.append(f"✓ Docker Compose installed: {version}")
                compose_found = True
        except (subprocess.TimeoutExpired, FileNotFoundError):
            pass

        if not compose_found:
            try:
                # Try docker-compose (v1)
                result = subprocess.run(["docker-compose", "--version"],
                                      timeout=5,
                                      **kwargs)
                if result.returncode == 0:
                    version = result.stdout.decode('utf-8', errors='ignore').strip()
                    results.append(f"✓ Docker Compose installed: {version}")
                    compose_found = True
            except (subprocess.TimeoutExpired, FileNotFoundError):
                pass

        if not compose_found:
            results.append("✗ Docker Compose not found")
            all_passed = False

        # Update UI from main thread
        self.docker_validated = all_passed
        result_text = "\n".join(results)
        color = "green" if all_passed else "red"

        self.root.after(0, lambda: self.docker_status_label.config(
            text=result_text, foreground=color))

        if all_passed:
            self.root.after(0, lambda: messagebox.showinfo(
                "Success", "All Docker checks passed!\nYou can proceed to the next step."))
        else:
            self.root.after(0, lambda: messagebox.showwarning(
                "Docker Check Failed",
                "Some Docker checks failed. Please install Docker and ensure it's running."))

    def check_docker_validation(self):
        """Validate Docker before moving to next page"""
        if not self.docker_validated:
            messagebox.showerror("Docker Required",
                               "Docker must be installed and running to proceed.\n"
                               "Please run the Docker check and ensure all checks pass.")
            return False
        return True

    # ==================== PAGE 3: Curl Check ====================
    def create_curl_check_page(self):
        """Create curl check page"""
        frame = ttk.Frame(self.content_frame)
        frame.pack(fill=tk.BOTH, expand=True)

        # Title
        title = ttk.Label(frame, text="curl Installation Check",
                         font=('Arial', 14, 'bold'))
        title.pack(pady=20)

        # Description
        desc = ttk.Label(frame, text="curl is required to download CodeKiwi files.\n"
                                     "Checking if curl is installed...",
                        justify=tk.CENTER)
        desc.pack(pady=10)

        # Check results frame
        self.curl_result_frame = ttk.LabelFrame(frame, text="Check Results", padding=10)
        self.curl_result_frame.pack(fill=tk.BOTH, expand=True, padx=20, pady=10)

        self.curl_status_label = ttk.Label(self.curl_result_frame,
                                          text="Checking...",
                                          wraplength=600)
        self.curl_status_label.pack(pady=10)

        # Button frame for install/recheck
        self.curl_button_frame = ttk.Frame(self.curl_result_frame)
        self.curl_button_frame.pack(pady=10)

        # Auto-check curl
        self.root.after(100, self.check_curl)

        # Store validation state
        self.curl_validated = False

    def check_curl(self):
        """Check if curl is installed"""
        # Clear previous buttons
        for widget in self.curl_button_frame.winfo_children():
            widget.destroy()

        # Find curl path
        curl_path = find_curl_path()
        curl_found = False
        version_info = ""

        # Method 1: Try with absolute path and proper subprocess kwargs
        try:
            kwargs = get_subprocess_kwargs()
            result = subprocess.run([curl_path, '--version'],
                                  timeout=5,
                                  **kwargs)

            if result.returncode == 0:
                curl_found = True
                version_info = result.stdout.decode('utf-8', errors='ignore').split('\n')[0]
        except (subprocess.TimeoutExpired, FileNotFoundError, Exception):
            pass

        # Method 2: If not found, try with shell=True (Windows fallback)
        if not curl_found and sys.platform == 'win32':
            try:
                kwargs = get_subprocess_kwargs()
                # Use shell=True with quoted path
                cmd = f'"{curl_path}" --version' if ' ' in curl_path else f'{curl_path} --version'
                result = subprocess.run(cmd,
                                      shell=True,
                                      timeout=5,
                                      **kwargs)

                if result.returncode == 0:
                    curl_found = True
                    version_info = result.stdout.decode('utf-8', errors='ignore').split('\n')[0]
            except (subprocess.TimeoutExpired, Exception):
                pass

        # Method 3: Last resort - try just 'curl' command
        if not curl_found:
            try:
                kwargs = get_subprocess_kwargs()
                result = subprocess.run(['curl', '--version'],
                                      timeout=5,
                                      **kwargs)

                if result.returncode == 0:
                    curl_found = True
                    version_info = result.stdout.decode('utf-8', errors='ignore').split('\n')[0]
            except (subprocess.TimeoutExpired, FileNotFoundError, Exception):
                pass

        # Update UI based on result
        if curl_found:
            self.curl_status_label.config(
                text=f"✓ curl is installed:\n{version_info}",
                foreground="green")
            self.curl_validated = True
        else:
            self.curl_status_label.config(
                text="✗ curl is not installed",
                foreground="red")
            self.show_curl_install_instructions()

    def show_curl_install_instructions(self):
        """Show instructions for installing curl"""
        # Add automatic install button
        install_btn = ttk.Button(self.curl_button_frame,
                                text="Install curl Automatically",
                                command=self.install_curl_automatically)
        install_btn.pack(side=tk.LEFT, padx=5)

        # Add recheck button
        recheck_btn = ttk.Button(self.curl_button_frame,
                                text="Recheck",
                                command=self.check_curl)
        recheck_btn.pack(side=tk.LEFT, padx=5)

        instructions = ""
        if self.os_type == "Darwin":
            instructions = "curl is usually pre-installed on macOS.\nIf missing, we can install Xcode Command Line Tools for you."
        elif self.os_type == "Windows":
            instructions = "curl is included in Windows 10+.\nIf missing, we can try to install it using winget or chocolatey."
        else:  # Linux
            instructions = "We can try to install curl using your system's package manager."

        install_frame = ttk.LabelFrame(self.curl_result_frame, text="Installation Options", padding=10)
        install_frame.pack(fill=tk.X, pady=10)

        install_label = ttk.Label(install_frame, text=instructions, justify=tk.LEFT)
        install_label.pack()

    def check_curl_validation(self):
        """Validate curl before moving to next page"""
        if not self.curl_validated:
            messagebox.showerror("curl Required",
                               "curl must be installed to proceed.\n"
                               "Please install curl and restart the installer.")
            return False
        return True

    def install_curl_automatically(self):
        """Install curl automatically after user confirmation"""
        # Ask for confirmation
        message = "This will attempt to install curl on your system.\n\n"

        if self.os_type == "Darwin":
            message += "This will run: xcode-select --install\n"
            message += "You may need to confirm the installation in a popup dialog.\n\n"
        elif self.os_type == "Windows":
            message += "This will try to install curl using winget or chocolatey.\n"
            message += "Administrator privileges may be required.\n\n"
        else:  # Linux
            message += "This will run the appropriate package manager command.\n"
            message += "Administrator privileges (sudo) may be required.\n\n"

        message += "Do you want to continue?"

        if not messagebox.askyesno("Install curl", message):
            return

        # Update status
        self.curl_status_label.config(text="Installing curl...", foreground="blue")
        self.root.update()

        # Disable install button during installation
        for widget in self.curl_button_frame.winfo_children():
            widget.config(state=tk.DISABLED)

        # Run installation in thread
        thread = threading.Thread(target=self._install_curl_thread)
        thread.daemon = True
        thread.start()

    def _install_curl_thread(self):
        """Thread function to install curl"""
        success = False
        error_msg = ""
        kwargs = get_subprocess_kwargs()

        try:
            if self.os_type == "Darwin":
                # macOS: Install Xcode Command Line Tools
                result = subprocess.run(["xcode-select", "--install"],
                                      timeout=60,
                                      **kwargs)
                # xcode-select --install returns 1 even on success (prompts user)
                # We'll consider it successful if it doesn't throw an error
                success = True
                error_msg = "Xcode installation dialog should have appeared.\nPlease complete the installation and click 'Recheck'."

            elif self.os_type == "Windows":
                # Try winget first
                try:
                    result = subprocess.run(["winget", "install", "curl.curl"],
                                          timeout=120,
                                          **kwargs)
                    if result.returncode == 0:
                        success = True
                    else:
                        # Try chocolatey
                        result = subprocess.run(["choco", "install", "curl", "-y"],
                                              timeout=120,
                                              **kwargs)
                        if result.returncode == 0:
                            success = True
                        else:
                            stderr_msg = result.stderr.decode('utf-8', errors='ignore')
                            error_msg = f"Failed to install using winget or chocolatey.\n{stderr_msg}\nPlease install curl manually."
                except FileNotFoundError:
                    error_msg = "Neither winget nor chocolatey found.\nPlease install curl manually from https://curl.se/windows/"

            else:  # Linux
                # Try to detect package manager and install
                package_managers = [
                    (["apt-get", "install", "-y", "curl"], "apt-get"),
                    (["dnf", "install", "-y", "curl"], "dnf"),
                    (["yum", "install", "-y", "curl"], "yum"),
                    (["pacman", "-S", "--noconfirm", "curl"], "pacman"),
                    (["zypper", "install", "-y", "curl"], "zypper"),
                ]

                for cmd, pm_name in package_managers:
                    try:
                        # Check if package manager exists
                        check_result = subprocess.run(["which", cmd[0]],
                                                     timeout=5,
                                                     **kwargs)
                        if check_result.returncode == 0:
                            # Try to install with sudo
                            full_cmd = ["sudo"] + cmd
                            result = subprocess.run(full_cmd,
                                                  timeout=120,
                                                  **kwargs)
                            if result.returncode == 0:
                                success = True
                                break
                            else:
                                stderr_msg = result.stderr.decode('utf-8', errors='ignore')
                                error_msg = f"Failed to install using {pm_name}.\n{stderr_msg}"
                    except FileNotFoundError:
                        continue

                if not success and not error_msg:
                    error_msg = "Could not detect package manager.\nPlease install curl manually."

        except subprocess.TimeoutExpired:
            error_msg = "Installation timed out. Please try installing manually."
        except Exception as e:
            error_msg = f"Installation failed: {str(e)}"

        # Update UI from main thread
        if success:
            self.root.after(0, lambda: self.curl_status_label.config(
                text=error_msg if error_msg else "✓ curl installation initiated.\nClick 'Recheck' to verify.",
                foreground="blue"))
            self.root.after(0, lambda: messagebox.showinfo(
                "Installation Complete",
                "curl installation has been initiated.\nClick 'Recheck' to verify the installation."))
        else:
            self.root.after(0, lambda: self.curl_status_label.config(
                text=f"✗ Installation failed:\n{error_msg}",
                foreground="red"))
            self.root.after(0, lambda: messagebox.showerror(
                "Installation Failed",
                f"Failed to install curl automatically:\n\n{error_msg}"))

        # Re-enable buttons
        self.root.after(0, lambda: [widget.config(state=tk.NORMAL)
                                     for widget in self.curl_button_frame.winfo_children()])

    # ==================== PAGE 4: Installation ====================
    def create_install_page(self):
        """Create installation page"""
        frame = ttk.Frame(self.content_frame)
        frame.pack(fill=tk.BOTH, expand=True)

        # Title
        title = ttk.Label(frame, text="Installing CodeKiwi",
                         font=('Arial', 14, 'bold'))
        title.pack(pady=20)

        # Progress bar
        self.progress = ttk.Progressbar(frame, mode='indeterminate', length=400)
        self.progress.pack(pady=20)

        # Log frame
        log_frame = ttk.LabelFrame(frame, text="Installation Log", padding=10)
        log_frame.pack(fill=tk.BOTH, expand=True, padx=20, pady=10)

        self.log_text = scrolledtext.ScrolledText(log_frame, height=15, width=70,
                                                  state=tk.DISABLED)
        self.log_text.pack(fill=tk.BOTH, expand=True)

        # Disable navigation during installation
        self.back_btn.config(state=tk.DISABLED)
        self.next_btn.config(state=tk.DISABLED)
        self.cancel_btn.config(state=tk.DISABLED)

        # Start installation
        self.root.after(500, self.start_installation)

    def log(self, message):
        """Add message to installation log"""
        self.log_text.config(state=tk.NORMAL)
        self.log_text.insert(tk.END, message + "\n")
        self.log_text.see(tk.END)
        self.log_text.config(state=tk.DISABLED)
        self.root.update()

    def start_installation(self):
        """Start the installation process"""
        self.progress.start(10)

        # Run installation in thread
        thread = threading.Thread(target=self._install_thread)
        thread.daemon = True
        thread.start()

    def _install_thread(self):
        """Thread function for installation"""
        try:
            self.log("Starting CodeKiwi installation...")
            self.log(f"Operating System: {self.os_type}")
            self.log(f"Architecture: {self.arch}")
            self.log(f"Installation Directory: {self.install_dir}")
            self.log("")

            # Create installation directory
            self.log("Creating installation directory...")
            os.makedirs(self.install_dir, exist_ok=True)
            bin_dir = os.path.join(self.install_dir, "bin")
            os.makedirs(bin_dir, exist_ok=True)
            self.log(f"✓ Created directory: {self.install_dir}")
            self.log("")

            # Get latest release info
            self.log("Fetching latest release information...")
            release_info = self.get_latest_release()
            version = release_info.get("tag_name", "unknown")
            self.log(f"✓ Latest version: {version}")
            self.log("")

            # Determine binary name
            binary_name = self.get_binary_name()
            self.log(f"Target binary: {binary_name}")

            # Download binary
            self.log("Downloading CodeKiwi binary...")
            binary_url = self.get_binary_url(release_info, binary_name)
            binary_path = os.path.join(bin_dir, "codekiwi")
            if self.os_type == "Windows":
                binary_path += ".exe"

            self.download_file(binary_url, binary_path)
            self.log(f"✓ Downloaded to: {binary_path}")

            # Make binary executable (Unix-like systems)
            if self.os_type != "Windows":
                os.chmod(binary_path, 0o755)
                self.log("✓ Set executable permissions")
            self.log("")

            # Download configuration files
            self.log("Downloading configuration files...")
            config_files = [
                ("docker-compose.yaml", "docker-compose.yaml"),
                ("config.env", "config.env")
            ]

            for remote_file, local_file in config_files:
                url = f"https://raw.githubusercontent.com/drasdp/codekiwi-cli/main/{remote_file}"
                local_path = os.path.join(self.install_dir, local_file)
                self.download_file(url, local_path)
                self.log(f"✓ Downloaded: {local_file}")
            self.log("")

            # Add to PATH
            self.log("Configuring PATH...")
            self.configure_path(bin_dir)
            self.log("")

            self.log("Installation completed successfully!")
            self.log("")
            self.log("To use CodeKiwi, please:")
            if self.os_type == "Windows":
                self.log("1. Restart your terminal/PowerShell")
            else:
                self.log("1. Restart your terminal or run: source ~/.bashrc (or ~/.zshrc)")
            self.log("2. Run: codekiwi start")

            # Stop progress bar and enable next button
            self.root.after(0, self.progress.stop)
            self.root.after(0, lambda: self.next_btn.config(state=tk.NORMAL))
            self.root.after(0, lambda: self.cancel_btn.config(state=tk.NORMAL))

        except Exception as e:
            self.log(f"\n✗ ERROR: {str(e)}")
            self.root.after(0, self.progress.stop)
            self.root.after(0, lambda: messagebox.showerror(
                "Installation Failed",
                f"An error occurred during installation:\n{str(e)}"))
            self.root.after(0, lambda: self.back_btn.config(state=tk.NORMAL))
            self.root.after(0, lambda: self.cancel_btn.config(state=tk.NORMAL))

    def get_latest_release(self):
        """Get latest release info from GitHub API"""
        url = "https://api.github.com/repos/drasdp/codekiwi-cli/releases/latest"
        with urllib.request.urlopen(url, timeout=10) as response:
            data = response.read()
            return json.loads(data)

    def get_binary_name(self):
        """Get the binary name for current platform"""
        os_map = {
            "Darwin": "darwin",
            "Linux": "linux",
            "Windows": "windows"
        }

        arch_map = {
            "x86_64": "amd64",
            "amd64": "amd64",
            "arm64": "arm64",
            "aarch64": "arm64"
        }

        os_name = os_map.get(self.os_type, "linux")
        arch_name = arch_map.get(self.arch, "amd64")

        binary = f"codekiwi-{os_name}-{arch_name}"
        if self.os_type == "Windows":
            binary += ".exe"

        return binary

    def get_binary_url(self, release_info, binary_name):
        """Get download URL for binary"""
        for asset in release_info.get("assets", []):
            if asset["name"] == binary_name:
                return asset["browser_download_url"]

        raise Exception(f"Binary not found for {binary_name}")

    def download_file(self, url, dest_path):
        """Download file from URL to destination"""
        with urllib.request.urlopen(url, timeout=30) as response:
            with open(dest_path, 'wb') as out_file:
                out_file.write(response.read())

    def configure_path(self, bin_dir):
        """Add installation directory to PATH"""
        if self.os_type == "Windows":
            self.configure_path_windows(bin_dir)
        else:
            self.configure_path_unix(bin_dir)

    def configure_path_windows(self, bin_dir):
        """Configure PATH for Windows"""
        try:
            # Add to user PATH using PowerShell
            ps_command = f'''
            $oldPath = [Environment]::GetEnvironmentVariable('Path', 'User')
            if ($oldPath -notlike "*{bin_dir}*") {{
                $newPath = $oldPath + ";{bin_dir}"
                [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
            }}
            '''

            kwargs = get_subprocess_kwargs()
            result = subprocess.run(["powershell", "-Command", ps_command],
                                  timeout=10,
                                  **kwargs)

            if result.returncode == 0:
                self.log(f"✓ Added to PATH: {bin_dir}")
                self.log("  (Restart terminal for changes to take effect)")
            else:
                stderr_msg = result.stderr.decode('utf-8', errors='ignore')
                self.log(f"! Could not automatically add to PATH: {stderr_msg}")
                self.log(f"  Please manually add this directory to PATH: {bin_dir}")
        except Exception as e:
            self.log(f"! Could not automatically add to PATH: {e}")
            self.log(f"  Please manually add this directory to PATH: {bin_dir}")

    def configure_path_unix(self, bin_dir):
        """Configure PATH for Unix-like systems"""
        home = os.path.expanduser("~")
        shell_configs = [
            os.path.join(home, ".bashrc"),
            os.path.join(home, ".zshrc"),
            os.path.join(home, ".bash_profile"),
            os.path.join(home, ".profile")
        ]

        export_line = f'\n# CodeKiwi\nexport PATH="$PATH:{bin_dir}"\n'

        updated = False
        for config_file in shell_configs:
            if os.path.exists(config_file):
                try:
                    # Check if already added
                    with open(config_file, 'r') as f:
                        content = f.read()

                    if bin_dir not in content:
                        with open(config_file, 'a') as f:
                            f.write(export_line)
                        self.log(f"✓ Updated: {config_file}")
                        updated = True
                except Exception as e:
                    self.log(f"! Could not update {config_file}: {e}")

        # Try to create symlink to /usr/local/bin
        try:
            local_bin = "/usr/local/bin/codekiwi"
            source_bin = os.path.join(bin_dir, "codekiwi")

            if not os.path.exists(local_bin):
                os.symlink(source_bin, local_bin)
                self.log(f"✓ Created symlink: {local_bin}")
        except Exception as e:
            self.log(f"! Could not create symlink (may require sudo): {e}")

        if updated:
            self.log("✓ PATH configured")
            self.log("  (Run 'source ~/.bashrc' or restart terminal)")

    # ==================== PAGE 5: Finish ====================
    def create_finish_page(self):
        """Create finish page"""
        frame = ttk.Frame(self.content_frame)
        frame.pack(fill=tk.BOTH, expand=True)

        # Title
        title = ttk.Label(frame, text="Installation Complete!",
                         font=('Arial', 16, 'bold'),
                         foreground="green")
        title.pack(pady=30)

        # Success message
        msg_frame = ttk.Frame(frame)
        msg_frame.pack(fill=tk.BOTH, expand=True, padx=40)

        success_msg = f"""CodeKiwi has been successfully installed!

Installation Directory: {self.install_dir}

Next Steps:
1. {'Restart your terminal/PowerShell' if self.os_type == 'Windows' else 'Restart your terminal or run: source ~/.bashrc'}
2. Verify installation: codekiwi --version
3. Start CodeKiwi: codekiwi start

For more information:
• Documentation: https://github.com/drasdp/codekiwi-cli
• Issues: https://github.com/drasdp/codekiwi-cli/issues

Thank you for installing CodeKiwi!"""

        msg_label = ttk.Label(msg_frame, text=success_msg, justify=tk.LEFT)
        msg_label.pack(anchor=tk.W)

        # Change next button to finish
        self.next_btn.config(text="Finish", command=self.finish_installation)


def main():
    """Main entry point"""
    root = tk.Tk()
    app = CodeKiwiInstaller(root)
    root.mainloop()


if __name__ == "__main__":
    main()
