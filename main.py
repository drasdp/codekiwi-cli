#!/usr/bin/env python3
"""
CodeKiwi Launcher - Docker-based Web Development Environment
A tkinter GUI application for managing CodeKiwi Docker containers
"""

import os
import sys
import threading
import webbrowser
from pathlib import Path
from tkinter import Tk, Frame, Label, Button, Text, Scrollbar, StringVar, messagebox, filedialog
from tkinter import ttk
import docker
from docker.errors import DockerException, APIError, BuildError, ContainerError, NotFound, ImageNotFound
import psutil


class CodeKiwiLauncher:
    """Docker SDK wrapper for CodeKiwi container management"""

    def __init__(self):
        self.client = None
        self.container = None

    def check_docker(self):
        """Check if Docker is installed and running"""
        try:
            self.client = docker.from_env()
            self.client.ping()
            return True, "Docker is running"
        except DockerException as e:
            return False, f"Docker error: {str(e)}"

    def build_image(self, callback=None):
        """
        Build Docker image from core/ directory

        Args:
            callback: Function to receive build log messages

        Returns:
            (success: bool, result: str)
        """
        try:
            # Get absolute path to core directory
            core_dir = Path(__file__).parent / "core"

            if not core_dir.exists():
                return False, f"core/ directory not found at {core_dir}"

            if callback:
                callback(f"Building from {core_dir}...")

            # Build image with real-time logs
            image, build_logs = self.client.images.build(
                path=str(core_dir),
                tag="codekiwi:latest",
                rm=True,  # Remove intermediate containers
                forcerm=True,  # Always remove intermediate containers
            )

            # Stream build logs
            for log in build_logs:
                if 'stream' in log:
                    message = log['stream'].strip()
                    if message and callback:
                        callback(message)
                elif 'error' in log:
                    return False, f"Build error: {log['error']}"

            return True, image.short_id

        except BuildError as e:
            return False, f"Build failed: {e.msg}"
        except APIError as e:
            return False, f"API error: {e.explanation}"
        except Exception as e:
            return False, f"Unexpected error: {str(e)}"

    def check_existing_container(self, container_name):
        """Check if a container with the given name exists"""
        try:
            container = self.client.containers.get(container_name)
            return True, container
        except NotFound:
            return False, None
        except APIError as e:
            return False, str(e)

    def remove_container(self, container_name):
        """Stop and remove existing container"""
        try:
            container = self.client.containers.get(container_name)

            # Stop if running
            if container.status == 'running':
                container.stop(timeout=10)

            # Remove
            container.remove()
            return True, "Container removed"

        except NotFound:
            return True, "Container not found"
        except APIError as e:
            return False, f"Failed to remove: {e.explanation}"

    def start_container(self, directory, container_name):
        """
        Start a new container with specified directory mounted

        Args:
            directory: Local directory to mount to /workspace
            container_name: Name for the container

        Returns:
            (success: bool, result: str)
        """
        try:
            self.container = self.client.containers.run(
                image="codekiwi:latest",
                name=container_name,
                detach=True,  # Run in background
                ports={
                    '80/tcp': 8080,
                    '3000/tcp': 3000
                },
                volumes={
                    str(directory): {
                        'bind': '/workspace',
                        'mode': 'rw'  # Read-write
                    }
                },
                restart_policy={
                    "Name": "unless-stopped"
                },
                remove=False,  # Don't auto-remove on stop
            )

            return True, self.container.short_id

        except ContainerError as e:
            return False, f"Container error: {e.stderr.decode() if e.stderr else str(e)}"
        except ImageNotFound:
            return False, "Image 'codekiwi:latest' not found. Please build first."
        except APIError as e:
            return False, f"API error: {e.explanation}"
        except Exception as e:
            return False, f"Unexpected error: {str(e)}"

    def get_container_status(self, container_name):
        """Get container status information"""
        try:
            container = self.client.containers.get(container_name)
            container.reload()  # Refresh status

            return {
                'status': container.status,  # running, exited, etc.
                'id': container.short_id,
                'created': container.attrs.get('Created', 'Unknown'),
            }

        except NotFound:
            return None
        except APIError as e:
            return {'error': str(e)}

    def stop_container(self, container_name):
        """Stop a running container"""
        try:
            container = self.client.containers.get(container_name)
            container.stop(timeout=10)
            return True, "Container stopped"
        except NotFound:
            return False, "Container not found"
        except APIError as e:
            return False, f"Failed to stop: {e.explanation}"

    def restart_container(self, container_name):
        """Restart a container"""
        try:
            container = self.client.containers.get(container_name)
            container.restart(timeout=10)
            return True, "Container restarted"
        except NotFound:
            return False, "Container not found"
        except APIError as e:
            return False, f"Failed to restart: {e.explanation}"


class CodeKiwiApp(Tk):
    """Main tkinter GUI application"""

    def __init__(self):
        super().__init__()

        self.title("CodeKiwi Launcher")
        self.geometry("700x600")
        self.resizable(True, True)

        # Initialize launcher
        self.launcher = CodeKiwiLauncher()

        # Check Docker
        success, message = self.launcher.check_docker()
        if not success:
            messagebox.showerror(
                "Docker Error",
                f"{message}\n\nPlease install and start Docker Desktop."
            )
            self.destroy()
            return

        # State variables
        self.selected_dir = StringVar()
        self.container_name = StringVar()
        self.status_text = StringVar(value="« Not running")
        self.monitoring = False

        self.setup_ui()

    def setup_ui(self):
        """Setup the GUI layout"""

        # Main container with padding
        main_frame = Frame(self, padx=20, pady=20)
        main_frame.pack(fill='both', expand=True)

        # Title
        title = Label(
            main_frame,
            text=">] CodeKiwi Launcher",
            font=('Arial', 18, 'bold')
        )
        title.pack(pady=(0, 20))

        # Directory selection section
        dir_frame = Frame(main_frame)
        dir_frame.pack(fill='x', pady=(0, 10))

        Label(dir_frame, text="=Á Work Directory:", font=('Arial', 11)).pack(anchor='w')

        dir_select_frame = Frame(dir_frame)
        dir_select_frame.pack(fill='x', pady=(5, 0))

        dir_entry = Label(
            dir_select_frame,
            textvariable=self.selected_dir,
            relief='sunken',
            anchor='w',
            bg='white',
            padx=10,
            pady=8,
            font=('Arial', 10)
        )
        dir_entry.pack(side='left', fill='x', expand=True, padx=(0, 5))

        browse_btn = Button(
            dir_select_frame,
            text="Browse...",
            command=self.select_directory,
            padx=15,
            pady=5
        )
        browse_btn.pack(side='right')

        # Container info section
        info_frame = Frame(main_frame)
        info_frame.pack(fill='x', pady=(10, 10))

        Label(info_frame, text="Container:", font=('Arial', 10)).pack(side='left')
        Label(info_frame, textvariable=self.container_name, font=('Arial', 10, 'bold')).pack(side='left', padx=(5, 20))

        Label(info_frame, text="Status:", font=('Arial', 10)).pack(side='left')
        self.status_label = Label(info_frame, textvariable=self.status_text, font=('Arial', 10, 'bold'))
        self.status_label.pack(side='left', padx=(5, 0))

        # Control buttons
        btn_frame = Frame(main_frame)
        btn_frame.pack(fill='x', pady=(10, 10))

        self.start_btn = Button(
            btn_frame,
            text="¶ Start",
            command=self.on_start_click,
            bg='#4CAF50',
            fg='white',
            font=('Arial', 11, 'bold'),
            padx=20,
            pady=8
        )
        self.start_btn.pack(side='left', padx=(0, 5))

        self.stop_btn = Button(
            btn_frame,
            text="ù Stop",
            command=self.on_stop_click,
            bg='#f44336',
            fg='white',
            font=('Arial', 11, 'bold'),
            padx=20,
            pady=8,
            state='disabled'
        )
        self.stop_btn.pack(side='left', padx=5)

        self.restart_btn = Button(
            btn_frame,
            text="= Restart",
            command=self.on_restart_click,
            bg='#2196F3',
            fg='white',
            font=('Arial', 11, 'bold'),
            padx=20,
            pady=8,
            state='disabled'
        )
        self.restart_btn.pack(side='left', padx=5)

        self.open_browser_btn = Button(
            btn_frame,
            text="< Open Browser",
            command=self.open_browser,
            bg='#FF9800',
            fg='white',
            font=('Arial', 11, 'bold'),
            padx=20,
            pady=8
        )
        self.open_browser_btn.pack(side='left', padx=5)

        # Log section
        log_label = Label(main_frame, text="Logs:", font=('Arial', 11, 'bold'))
        log_label.pack(anchor='w', pady=(10, 5))

        log_frame = Frame(main_frame)
        log_frame.pack(fill='both', expand=True)

        scrollbar = Scrollbar(log_frame)
        scrollbar.pack(side='right', fill='y')

        self.log_text = Text(
            log_frame,
            wrap='word',
            yscrollcommand=scrollbar.set,
            bg='#1e1e1e',
            fg='#d4d4d4',
            font=('Courier', 10),
            padx=10,
            pady=10
        )
        self.log_text.pack(side='left', fill='both', expand=True)
        scrollbar.config(command=self.log_text.yview)

        self.log("CodeKiwi Launcher initialized")
        self.log("Docker is running ")

    def select_directory(self):
        """Open directory selection dialog"""
        directory = filedialog.askdirectory(
            title="Select Work Directory",
            initialdir=os.path.expanduser("~")
        )

        if directory:
            self.selected_dir.set(directory)
            dir_name = os.path.basename(directory)
            self.container_name.set(f"codekiwi-{dir_name}")
            self.log(f"Selected directory: {directory}")

    def log(self, message):
        """Append message to log text widget"""
        self.log_text.insert('end', f"{message}\n")
        self.log_text.see('end')  # Auto-scroll to bottom

    def check_ports(self, ports=[8080, 3000]):
        """Check if ports are in use"""
        conflicts = []

        for port in ports:
            for conn in psutil.net_connections():
                if conn.laddr.port == port and conn.status == 'LISTEN':
                    try:
                        proc = psutil.Process(conn.pid)
                        conflicts.append({
                            'port': port,
                            'pid': conn.pid,
                            'process': proc.name()
                        })
                    except (psutil.NoSuchProcess, psutil.AccessDenied):
                        pass

        return conflicts

    def handle_port_conflicts(self):
        """Handle port conflicts with user confirmation"""
        conflicts = self.check_ports()

        if not conflicts:
            return True

        # Build warning message
        message = "The following ports are already in use:\n\n"
        for c in conflicts:
            message += f"Port {c['port']}: {c['process']} (PID: {c['pid']})\n"
        message += "\nDo you want to terminate these processes?"

        # Show confirmation dialog
        response = messagebox.askyesno(
            "Port Conflict",
            message,
            icon='warning'
        )

        if response:
            # Terminate processes
            for c in conflicts:
                try:
                    proc = psutil.Process(c['pid'])
                    proc.terminate()
                    proc.wait(timeout=5)
                    self.log(f"Terminated process on port {c['port']}")
                except Exception as e:
                    self.log(f"Failed to terminate process: {e}")
                    messagebox.showerror(
                        "Error",
                        f"Failed to terminate process on port {c['port']}\n"
                        "You may need administrator/sudo privileges."
                    )
                    return False
            return True
        else:
            return False

    def on_start_click(self):
        """Handle start button click"""

        # Check directory selection
        if not self.selected_dir.get():
            messagebox.showwarning("Warning", "Please select a work directory first")
            return

        directory = Path(self.selected_dir.get())
        if not directory.exists():
            messagebox.showerror("Error", f"Directory does not exist: {directory}")
            return

        # Check port conflicts
        if not self.handle_port_conflicts():
            return

        # Check existing container
        container_name = self.container_name.get()
        exists, container = self.launcher.check_existing_container(container_name)

        if exists and container:
            response = messagebox.askyesno(
                "Container Exists",
                f"Container '{container_name}' already exists.\n"
                "Do you want to stop and remove it, then start fresh?",
                icon='warning'
            )

            if response:
                self.log(f"Removing existing container: {container_name}")
                success, msg = self.launcher.remove_container(container_name)
                if not success:
                    messagebox.showerror("Error", f"Failed to remove container: {msg}")
                    return
                self.log(msg)
            else:
                return

        # Start container asynchronously
        self.start_container_async()

    def start_container_async(self):
        """Start container in background thread to avoid GUI blocking"""

        def worker():
            # Disable start button
            self.after(0, lambda: self.start_btn.config(state='disabled'))

            # Build image
            self.after(0, lambda: self.log("Building Docker image..."))
            success, result = self.launcher.build_image(
                callback=lambda msg: self.after(0, lambda m=msg: self.log(m))
            )

            if not success:
                self.after(0, lambda: self.log(f"L Build failed: {result}"))
                self.after(0, lambda: messagebox.showerror("Build Error", result))
                self.after(0, lambda: self.start_btn.config(state='normal'))
                return

            self.after(0, lambda: self.log(f" Image built: {result}"))

            # Start container
            directory = self.selected_dir.get()
            container_name = self.container_name.get()

            self.after(0, lambda: self.log(f"Starting container: {container_name}..."))
            success, result = self.launcher.start_container(directory, container_name)

            if success:
                self.after(0, lambda: self.log(f" Container started: {result}"))
                self.after(0, lambda: self.log(""))
                self.after(0, lambda: self.log("< Access at http://localhost:8080"))
                self.after(0, lambda: self.log("=' Dev server at http://localhost:3000"))
                self.after(0, lambda: self.start_monitoring())
            else:
                self.after(0, lambda: self.log(f"L Failed: {result}"))
                self.after(0, lambda: messagebox.showerror("Start Error", result))

            # Re-enable start button
            self.after(0, lambda: self.start_btn.config(state='normal'))

        # Run in background thread
        thread = threading.Thread(target=worker, daemon=True)
        thread.start()

    def start_monitoring(self):
        """Start periodic container status monitoring"""

        def check_status():
            if not self.monitoring:
                return

            container_name = self.container_name.get()
            if not container_name:
                return

            status = self.launcher.get_container_status(container_name)

            if status and 'status' in status:
                self.update_status_ui(status['status'])
            else:
                self.update_status_ui('stopped')

            # Check again in 5 seconds
            if self.monitoring:
                self.after(5000, check_status)

        self.monitoring = True
        check_status()

    def stop_monitoring(self):
        """Stop status monitoring"""
        self.monitoring = False

    def update_status_ui(self, status):
        """Update UI based on container status"""
        if status == 'running':
            self.status_text.set("=â Running")
            self.status_label.config(fg='green')
            self.stop_btn.config(state='normal')
            self.restart_btn.config(state='normal')
        else:
            self.status_text.set("« Stopped")
            self.status_label.config(fg='gray')
            self.stop_btn.config(state='disabled')
            self.restart_btn.config(state='disabled')

    def on_stop_click(self):
        """Handle stop button click"""
        container_name = self.container_name.get()

        self.log(f"Stopping container: {container_name}...")
        success, msg = self.launcher.stop_container(container_name)

        if success:
            self.log(f" {msg}")
        else:
            self.log(f"L {msg}")
            messagebox.showerror("Error", msg)

    def on_restart_click(self):
        """Handle restart button click"""
        container_name = self.container_name.get()

        self.log(f"Restarting container: {container_name}...")
        success, msg = self.launcher.restart_container(container_name)

        if success:
            self.log(f" {msg}")
        else:
            self.log(f"L {msg}")
            messagebox.showerror("Error", msg)

    def open_browser(self):
        """Open browser to localhost:8080"""
        webbrowser.open("http://localhost:8080")
        self.log("Opened browser at http://localhost:8080")


def main():
    """Main entry point"""
    app = CodeKiwiApp()
    app.mainloop()


if __name__ == "__main__":
    main()
