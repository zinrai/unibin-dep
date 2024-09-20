# unibin-dep: Universal Binary Deployment Tool

`unibin-dep` is a command-line tool written in Go that facilitates the safe download and deployment of binary files across different operating systems and architectures. It ensures that only compatible binaries are deployed to the specified directory.

## Features

- Downloads files from specified URLs
- Detects whether the downloaded file is a binary or text file
- Checks binary compatibility with the host system (supports Linux, macOS, and Windows)
- Safely handles downloads using temporary files
- Moves compatible binaries to the specified directory
- Automatically sets executable permissions for compatible binaries
- Cleans up incompatible binaries

## Installation

Build the binary:

```
$ go build
```

## Usage

The basic syntax for using `unibin-dep` is:

```
$ unibin-dep -u <download_url> -d <save_directory> [-f <custom_filename>]
```

Example:

```
$ unibin-dep -u https://github.com/kubernetes-sigs/kind/releases/download/v0.24.0/kind-linux-amd64 -d ~/bin -f kind
```

## How it works

1. The tool downloads the file from the specified URL to a temporary location.
2. It checks if the downloaded file is a binary.
3. If it's a binary, it verifies compatibility with the host system.
4. Compatible binaries are moved to the specified directory and given executable permissions.
5. Incompatible binaries are removed.
6. Text files are moved to the specified directory without additional checks.

## Supported Systems

`unibin-dep` can check compatibility for the following systems and architectures:

- Linux: x86_64 (amd64), ARM64
- macOS: x86_64 (amd64), ARM64
- Windows: x86_64 (amd64), x86

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
