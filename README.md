# scallop

Scallop is a command-line tool for scanning Docker images for vulnerabilities and analyzing their structure.

- [scallop](#scallop)
  - [Features](#features)
  - [Installation](#installation)
    - [From Source](#from-source)
  - [Usage](#usage)
    - [Basic Usage](#basic-usage)
    - [Output Formats](#output-formats)
    - [Examples](#examples)
  - [Project Structure](#project-structure)
  - [Development](#development)
  - [References](#references)

## Features

- **Docker Image Vulnerability Scanning**: Scans Docker images for security vulnerabilities
- **Directory Structure Analysis**: Analyzes the directory structure of Docker images
- **Size Analysis**: Analyzes the size of Docker images and identifies large files and directories
- **Security Analysis**: Identifies sensitive files, hardcoded secrets, and vulnerable packages
- **Recommendations**: Provides recommendations for improving the security and size of Docker images

## Installation

### From Source

1. Clone the repository:

```bash
git clone https://github.com/yourusername/scallop.git
cd scallop
```

2. Build the project:

```bash
go build -o scallop ./cmd/scallop
```

3. Install the binary (optional):

```bash
sudo mv scallop /usr/local/bin/
```

## Usage

### Basic Usage

```bash
# Scan a Docker image from Docker daemon
scallop --image nginx:latest

# Scan a Docker image from a tar file
scallop --image path/to/image.tar

# Enable verbose output
scallop --image ubuntu:20.04 --verbose
```

### Output Formats

Scallop supports two output formats: text (default) and JSON.

```bash
# Output in text format (default)
scallop --image nginx:latest

# Output in JSON format
scallop --image nginx:latest --format json
```

### Examples

**Scan a Docker image from Docker daemon:**

```bash
scallop --image nginx:latest
```

**Scan a Docker image from a tar file:**

```bash
# First, save a Docker image to a tar file
docker save -o nginx.tar nginx:latest

# Then scan the tar file
scallop --image nginx.tar
```

**Output in JSON format:**

```bash
scallop --image nginx:latest --format json
```

**Enable verbose output:**

```bash
scallop --image nginx:latest --verbose
```

## Project Structure

```
scallop/
├── cmd/
│   └── scallop/
│       ├── main.go             # Entry point (CLI)
├── internal/
│   ├── analyzer/               # Analysis logic
│   │   ├── analyzer.go         # Main analyzer
│   │   ├── directory.go        # Directory structure analysis
│   │   ├── security.go         # Security risk analysis
│   │   ├── size.go             # Size analysis
│   ├── docker/                 # Docker image operations
│   │   ├── docker.go           # Docker image extraction and analysis
│   └── utils/                  # Common utilities
│       └── file.go             # File operation utilities
├── pkg/
│   ├── config/                 # Configuration file management
│   │   └── config.go
├── scripts/                    # Scripts (CI scripts, build scripts, etc.)
├── go.mod
├── go.sum
└── README.md
```

## Development

1. Clone the repository:

```bash
git clone https://github.com/yourusername/scallop.git
cd scallop
```

2. Install dependencies:

```bash
go mod tidy
```

3. Build the project:

```bash
go build -o scallop ./cmd/scallop
```

4. Run tests:

```bash
go test ./...
```

## References

- [Docker Engine API](https://docs.docker.com/engine/api/)
- [Container Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [Docker Image Specification](https://github.com/moby/moby/blob/master/image/spec/v1.md)
