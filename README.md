# SSD Indexer

High-performance SSD analyzer with real-time scanning and smart file categorization built in Go.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Supported File Types](#supported-file-types)
- [Installation](#installation)
- [Usage](#usage)
- [Dependencies](#-dependencies)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

**SSD Indexer** is a fast and efficient tool designed to scan and index the contents of your SSD.  
It provides insights into your storage usage, including the percentage of memory occupied by different file types (photos, videos, documents, compressed files, and others).  
With a modern terminal user interface (TUI), it offers live progress updates while scanning large storage volumes.

---

## Features

- Fast SSD scanning using parallel processing
- Real-time TUI updates with live progress
- Categorizes files: Photos, Videos, Documents, Compressed, Other
- Color-coded progress bars and storage percentages
- Displays top file extensions for "Other" category
- Shows total storage used
- Keyboard navigation: `q` to quit
- Safe scanning without modifying any files

---

## Supported File Types

| Category    | Extensions                                                                 |
|------------|----------------------------------------------------------------------------|
| Photos      | .jpg, .jpeg, .png, .gif, .bmp, .heic                                       |
| Videos      | .mp4, .mov, .avi, .mkv                                                     |
| Documents   | .pdf, .txt, .doc, .docx, .md                                              |
| Compressed  | .zip, .rar, .7z, .tar, .gz, .bz2, .xz, .tar.gz, .tar.bz2, .tar.xz, .z, .lzma, .lz4, .zst, .arj, .cab, .deb, .rpm, .dmg, .iso, .img |
| Other       | All other file extensions                                                  |

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/SanujaRubasinghe/ssd-indexer.git
cd ssd-indexer
```

2. Build the project:

```bash
go build -o ssdindexer main.go
```

3. Make sure dependencies are installed:

```bash
go mod tidy
```

---

## Usage

Run the tool (requires Go installed):

```bash
./ssdindexer /path/to/scan
```

- By default, it scans the root directory `/`.
- Press `q` to quit the TUI at any time.

### Example Final Output

```
SSD Analyzer — DONE

Memory Composition:

Photos     [=========---------] 45.67% 2.3 GB
Videos     [====--------------] 20.15% 1.0 GB
Docs       [===---------------] 15.23% 768.5 MB
Compressed [====--------------] 18.95% 950.2 MB
Other      [=-----------------] 0.00% 12.3 KB

Top Extensions (Other):
- .log 8.91% 423.7 MB
- .tmp 5.67% 269.4 MB
- .cache 3.45% 163.8 MB
- .bak 2.34% 111.2 MB
- .old 1.23% 58.4 MB

Total Size: 4.1 GB

Press Q to quit
```

---

## 🔧 Dependencies

- [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) – TUI framework  
- [github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) – UI components  
- [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) – Styling

---

## Project Structure

```
ssd-indexer/
├── cmd/
│   └── analyzer/       # Main executable
├── internal/
│   ├── scanner/        # Scanning logic
│   └── ui/             # Terminal UI logic
├── go.mod
├── go.sum
└── README.md
```

---

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

1. Fork the repo
2. Create your feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m 'Add feature'`
4. Push to branch: `git push origin feature-name`
5. Open a Pull Request

---
