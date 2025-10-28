#!/usr/bin/env python3
"""
CodeKiwi Launcher - Simple GUI for CodeKiwi CLI
"""

import tkinter as tk
from tkinter import ttk, filedialog, scrolledtext, messagebox
import subprocess
import threading
import os
import sys

class CodeKiwiLauncher:
    def __init__(self):
        self.root = tk.Tk()
        self.root.title("CodeKiwi Launcher")
        self.root.geometry("600x500")
        self.root.resizable(False, False)

        # Variables
        self.project_path = tk.StringVar()
        self.process = None

        # Set style
        self.style = ttk.Style()
        self.style.theme_use('clam')

        self.create_widgets()

    def create_widgets(self):
        """Create all UI widgets"""
        # Header
        header = ttk.Label(self.root, text="CodeKiwi Launcher",
                          font=('Arial', 16, 'bold'))
        header.pack(pady=10)

        # Project path frame
        frame1 = ttk.Frame(self.root, padding="10")
        frame1.pack(fill=tk.X)

        ttk.Label(frame1, text="Project:").pack(side=tk.LEFT)
        ttk.Entry(frame1, textvariable=self.project_path, width=45).pack(side=tk.LEFT, padx=5)
        ttk.Button(frame1, text="Browse", command=self.browse_folder).pack(side=tk.LEFT)

        # Action buttons frame
        frame2 = ttk.Frame(self.root, padding="10")
        frame2.pack(fill=tk.X)

        self.start_btn = ttk.Button(frame2, text="▶ START", command=self.start_project, width=15)
        self.start_btn.pack(side=tk.LEFT, padx=5)

        self.stop_btn = ttk.Button(frame2, text="■ STOP", command=self.stop_project, width=15)
        self.stop_btn.pack(side=tk.LEFT, padx=5)

        self.list_btn = ttk.Button(frame2, text="☰ LIST", command=self.list_instances, width=15)
        self.list_btn.pack(side=tk.LEFT, padx=5)

        self.update_btn = ttk.Button(frame2, text="↻ UPDATE", command=self.update_codekiwi, width=15)
        self.update_btn.pack(side=tk.LEFT, padx=5)

        # Output frame
        output_frame = ttk.LabelFrame(self.root, text="Output", padding="10")
        output_frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=5)

        # Create output text widget with scrollbar
        self.output = scrolledtext.ScrolledText(
            output_frame,
            wrap=tk.WORD,
            height=18,
            font=("Consolas", 9),
            bg="#1e1e1e",
            fg="#ffffff",
            insertbackground="#ffffff"
        )
        self.output.pack(fill=tk.BOTH, expand=True)

        # Configure text tags for colors
        self.output.tag_config("command", foreground="#00ff00")
        self.output.tag_config("error", foreground="#ff6666")
        self.output.tag_config("success", foreground="#66ff66")
        self.output.tag_config("info", foreground="#66b3ff")

        # Control buttons
        control_frame = ttk.Frame(output_frame)
        control_frame.pack(fill=tk.X, pady=5)

        ttk.Button(control_frame, text="Clear", command=self.clear_output).pack(side=tk.LEFT)
        ttk.Button(control_frame, text="Copy All", command=self.copy_output).pack(side=tk.LEFT, padx=5)

        # Status bar
        self.status = ttk.Label(self.root, text="Ready", relief=tk.SUNKEN, anchor=tk.W)
        self.status.pack(side=tk.BOTTOM, fill=tk.X)

        # Add welcome message
        self.output.insert(tk.END, "=== CodeKiwi Launcher ===\n", "info")
        self.output.insert(tk.END, "Simple GUI for CodeKiwi CLI\n\n", "info")
        self.output.insert(tk.END, "1. Select a project folder\n")
        self.output.insert(tk.END, "2. Click START to launch\n")
        self.output.insert(tk.END, "   (CodeKiwi will auto-detect project type)\n\n")

    def browse_folder(self):
        """Open folder browser dialog"""
        folder = filedialog.askdirectory(
            title="Select Project Folder",
            initialdir=os.path.expanduser("~")
        )
        if folder:
            self.project_path.set(folder)
            self.status.config(text=f"Selected: {folder}")

    def start_project(self):
        """Start CodeKiwi project"""
        path = self.project_path.get()
        if not path:
            messagebox.showerror("Error", "Please select a project folder")
            return

        if not os.path.exists(path):
            messagebox.showerror("Error", "Selected folder does not exist")
            return

        # Add -d for detached mode to run in background
        cmd = ["codekiwi", "start", path, "-d"]
        self.run_command(cmd, command_type="start")
        self.status.config(text="Starting project...")

    def stop_project(self):
        """Stop CodeKiwi project"""
        path = self.project_path.get()
        if not path:
            # If no path selected, stop all
            response = messagebox.askyesno("Stop All", "No project selected. Stop all running instances?")
            if response:
                cmd = ["codekiwi", "kill", "--all", "--force"]
                self.run_command(cmd, command_type="stop")
                self.status.config(text="Stopping all instances...")
            return

        cmd = ["codekiwi", "kill", path, "--force"]
        self.run_command(cmd, command_type="stop")
        self.status.config(text="Stopping project...")

    def list_instances(self):
        """List running CodeKiwi instances"""
        cmd = ["codekiwi", "list", "--all"]
        self.run_command(cmd, command_type="list")
        self.status.config(text="Listing instances...")

    def update_codekiwi(self):
        """Update CodeKiwi CLI and images"""
        response = messagebox.askyesno("Update", "Update CodeKiwi CLI and Docker images?")
        if response:
            cmd = ["codekiwi", "update"]
            self.run_command(cmd, command_type="update")
            self.status.config(text="Updating CodeKiwi...")

    def run_command(self, cmd, command_type=None):
        """Execute command asynchronously"""
        # Display command
        self.output.insert(tk.END, f"\n$ {' '.join(cmd)}\n", "command")
        self.output.see(tk.END)

        # Selective button disabling based on command type
        if command_type == "start":
            # Only disable START button for start commands
            self.start_btn.config(state=tk.DISABLED)
        elif command_type == "stop":
            # Only disable STOP button for stop commands
            self.stop_btn.config(state=tk.DISABLED)
        elif command_type == "list":
            # Only disable LIST button for list commands
            self.list_btn.config(state=tk.DISABLED)
        elif command_type == "update":
            # Only disable UPDATE button for update commands
            self.update_btn.config(state=tk.DISABLED)
        else:
            # Default: disable all buttons
            self.set_buttons_state(tk.DISABLED)

        # Start command in thread
        thread = threading.Thread(target=self._execute_command, args=(cmd, command_type))
        thread.daemon = True
        thread.start()

    def _execute_command(self, cmd, command_type=None):
        """Execute command and capture output"""
        try:
            # Start process with UTF-8 encoding to avoid cp949 codec errors
            self.process = subprocess.Popen(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                encoding='utf-8',  # Force UTF-8 encoding
                errors='replace',   # Replace any undecodable bytes
                bufsize=1,
                shell=(sys.platform == "win32")  # Use shell on Windows
            )

            # Read output line by line
            for line in iter(self.process.stdout.readline, ''):
                if line:
                    # Determine line type and color
                    tag = None
                    if "error" in line.lower() or "failed" in line.lower():
                        tag = "error"
                    elif "success" in line.lower() or "✓" in line:
                        tag = "success"
                    elif "warning" in line.lower() or "!" in line:
                        tag = "info"

                    # Insert with appropriate tag
                    if tag:
                        self.output.insert(tk.END, line, tag)
                    else:
                        self.output.insert(tk.END, line)

                    self.output.see(tk.END)
                    self.root.update()

            # Wait for process to complete
            self.process.wait()

            # Display result
            if self.process.returncode == 0:
                if command_type == "start":
                    self.output.insert(tk.END, "✓ Project started in background (detached mode)\n", "success")
                    self.output.insert(tk.END, "   Use LIST to see running instances\n", "info")
                    self.output.insert(tk.END, "   Use STOP to stop the project\n", "info")
                    self.status.config(text="Project started in background")
                else:
                    self.output.insert(tk.END, "✓ Command completed successfully\n", "success")
                    self.status.config(text="Command completed")
            else:
                self.output.insert(tk.END, f"✗ Command failed (exit code: {self.process.returncode})\n", "error")
                self.status.config(text="Command failed")

            self.output.see(tk.END)

        except FileNotFoundError:
            error_msg = "Error: 'codekiwi' command not found. Please install CodeKiwi CLI first.\n"
            self.output.insert(tk.END, error_msg, "error")
            self.status.config(text="Error: CodeKiwi CLI not found")
            messagebox.showerror("Error", "CodeKiwi CLI is not installed.\nPlease install it first.")

        except Exception as e:
            self.output.insert(tk.END, f"Error: {str(e)}\n", "error")
            self.status.config(text="Error occurred")

        finally:
            self.process = None
            # Re-enable buttons based on command type
            if command_type == "start":
                self.root.after(0, lambda: self.start_btn.config(state=tk.NORMAL))
            elif command_type == "stop":
                self.root.after(0, lambda: self.stop_btn.config(state=tk.NORMAL))
            elif command_type == "list":
                self.root.after(0, lambda: self.list_btn.config(state=tk.NORMAL))
            elif command_type == "update":
                self.root.after(0, lambda: self.update_btn.config(state=tk.NORMAL))
            else:
                self.root.after(0, lambda: self.set_buttons_state(tk.NORMAL))

    def set_buttons_state(self, state):
        """Enable/disable buttons"""
        self.start_btn.config(state=state)
        self.stop_btn.config(state=state)
        self.list_btn.config(state=state)
        self.update_btn.config(state=state)

    def clear_output(self):
        """Clear output text"""
        self.output.delete(1.0, tk.END)
        self.status.config(text="Output cleared")

    def copy_output(self):
        """Copy all output to clipboard"""
        text = self.output.get(1.0, tk.END)
        self.root.clipboard_clear()
        self.root.clipboard_append(text)
        self.status.config(text="Output copied to clipboard")

    def run(self):
        """Start the GUI application"""
        # Center window on screen
        self.root.update_idletasks()
        width = self.root.winfo_width()
        height = self.root.winfo_height()
        x = (self.root.winfo_screenwidth() // 2) - (width // 2)
        y = (self.root.winfo_screenheight() // 2) - (height // 2)
        self.root.geometry(f'{width}x{height}+{x}+{y}')

        # Start main loop
        self.root.mainloop()

def main():
    """Main entry point"""
    app = CodeKiwiLauncher()
    app.run()

if __name__ == "__main__":
    main()