package main

import (
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func isBinary(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true
		}
	}
	return false
}

func getBinaryInfo(filePath string) (string, string) {
	// Check ELF (Linux)
	if elfFile, err := elf.Open(filePath); err == nil {
		defer elfFile.Close()
		switch elfFile.Machine {
		case elf.EM_X86_64:
			return "linux", "amd64"
		case elf.EM_AARCH64:
			return "linux", "arm64"
		}
	}

	// Check Mach-O (macOS)
	if machoFile, err := macho.Open(filePath); err == nil {
		defer machoFile.Close()
		switch machoFile.Cpu {
		case macho.CpuAmd64:
			return "darwin", "amd64"
		case macho.CpuArm64:
			return "darwin", "arm64"
		}
	}

	// Check PE (Windows)
	if peFile, err := pe.Open(filePath); err == nil {
		defer peFile.Close()
		switch peFile.Machine {
		case pe.IMAGE_FILE_MACHINE_AMD64:
			return "windows", "amd64"
		case pe.IMAGE_FILE_MACHINE_I386:
			return "windows", "386"
		}
	}

	return "", ""
}

func isCompatibleBinary(filePath string) bool {
	binaryOS, binaryArch := getBinaryInfo(filePath)
	return binaryOS == runtime.GOOS && binaryArch == runtime.GOARCH
}

func setExecutable(filePath string) error {
	return os.Chmod(filePath, 0755)
}

func main() {
	url := flag.String("u", "", "URL of the file to download")
	saveDir := flag.String("d", "", "Directory to save the downloaded file")
	filename := flag.String("f", "", "Specify a custom filename for the downloaded file")
	flag.Parse()

	if *url == "" || *saveDir == "" {
		fmt.Println("Both URL and save directory are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *filename == "" {
		*filename = filepath.Base(*url)
	}

	tempFile, err := os.CreateTemp("", "unibin-dep-*")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		os.Exit(1)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFilePath)

	fmt.Println("Downloading file...")
	if err := downloadFile(*url, tempFilePath); err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File downloaded to temporary location.")

	if isBinary(tempFilePath) {
		fmt.Println("The downloaded file is a binary.")
		if isCompatibleBinary(tempFilePath) {
			fmt.Println("The binary is compatible with the current system.")

			if err := os.MkdirAll(*saveDir, 0755); err != nil {
				fmt.Printf("Error creating save directory: %v\n", err)
				os.Exit(1)
			}

			finalPath := filepath.Join(*saveDir, *filename)
			if err := os.Rename(tempFilePath, finalPath); err != nil {
				fmt.Printf("Error moving file to final location: %v\n", err)
				os.Exit(1)
			}

			if err := setExecutable(finalPath); err != nil {
				fmt.Printf("Error setting executable permissions: %v\n", err)
			} else {
				fmt.Printf("File moved to: %s and execution permissions granted.\n", finalPath)
			}
		} else {
			fmt.Println("The binary is not compatible with the current system.")
			fmt.Println("Temporary file will be removed.")
		}
	} else {
		fmt.Println("The downloaded file is a text file.")
		finalPath := filepath.Join(*saveDir, *filename)
		if err := os.Rename(tempFilePath, finalPath); err != nil {
			fmt.Printf("Error moving text file to final location: %v\n", err)
		} else {
			fmt.Printf("Text file moved to: %s\n", finalPath)
		}
	}
}
