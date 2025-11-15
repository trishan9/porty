# porty

### **Modern, interactive port manager for Linux & Mac**

<img width="863" height="655" alt="Porty" src="https://github.com/user-attachments/assets/f04767b1-4a91-4455-8fe9-53e803aa6c8d" />

## Overview

**porty** is a fast, real-time interactive port manager written in Go.
It reads Linux networking data directly from /proc for maximum accuracy, presenting TCP/UDP sockets in a polished, Tokyo Night–inspired TUI.

Built for developers, sysadmins, and anyone who needs to quickly inspect or kill processes bound to ports.

## ⚙️ Installation

### 📥 Curl Installer (Quickest & Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/trishan9/porty/main/install.sh | bash
```

### 🐧 Arch Linux (AUR)

```bash
yay -S porty-bin
# or
paru -S porty-bin
```

AUR Package: https://aur.archlinux.org/packages/porty-bin

### 🛠️ Build from source

```bash
git clone https://github.com/trishan9/porty
cd porty
go build -o porty .
sudo mv porty /usr/local/bin/
```

## Usage Examples

### TUI (Interactive Listing, Killing, and all):

```bash
porty list
```

### Kill a port:

```
porty kill --port 3000
```

### Kill multiple Ports:

```
porty kill --port 3000,8000
```

### Kill a PID:

```
porty kill --pid 1234
```

### Export all port details to JSON:

```bash
porty list --json
porty list --json > ports.json # Saving as a file
```

### Check version:

```
porty version
```

## Features

### Full Port Scanner

- TCP + UDP port detection
- Kernel sockets (e.g. systemd-resolved, dhclient, avahi)
- PID → process name mapping
- UID → user detection
- Tags for:

  - **USER**
  - **SYSTEM**
  - **KERNEL**
  - **SELF**

### Real-Time Updates

- Auto-refresh every 2 seconds
- Manual refresh with `r`

### Keyboard Shortcuts

| Key           | Action        |
| ------------- | ------------- |
| ↑ / ↓ / j / k | Move cursor   |
| Space         | Select port   |
| Enter / x     | Kill process  |
| r             | Refresh ports |
| q             | Quit          |

## Architecture

Porty uses:

- `/proc/net/tcp` and `/proc/net/udp`
- `/proc/<pid>/fd` to resolve inode → PID
- `/proc/<pid>/comm` for process names
- `/proc/<pid>/status` for UID
- BubbleTea (TUI framework)
- LipGloss for UI styling

This avoids external dependencies like `lsof` or `netstat`, making Porty extremely fast.

## Contributing

Open issues for feature ideas or bugs.
