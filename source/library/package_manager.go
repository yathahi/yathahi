package library

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Abbrevations
// Package Manager = pm
// Package = pkg

/*	Download Types
	1. Linux Packages
	2. URL
*/

//--------------------------------------------------------------------------------
// Package Manager
//--------------------------------------------------------------------------------

func PackageManager(pkg string) error {

	slog.Info("Package manager has been initialized! for installing the package", slog.String("pkg", pkg))
	pkgType := pkgtypeparser(pkg)
	if pkgType == "1" {
		err := linuxPackage(pkg)
		return err
	} else {
		err := urlPkg(pkg)
		return err
	}

}

//--------------------------------------------------------------------------------
// Package Type Parser - Identify the Package Is a linux Package or a URL.
//--------------------------------------------------------------------------------

func pkgtypeparser(input string) string {
	// Clean up whitespace
	input = strings.TrimSpace(input)

	// Attempt to parse the string as a URL
	u, err := url.Parse(input)

	// If it fails basic parsing, it's definitely just a string
	if err != nil {
		return "1"
	}

	// Crucial Check: A valid web URL MUST have a scheme (http/https) and a host.
	// Without this, plain text like "my-file.txt" or "hello" gets misidentified as a URL path.
	if (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" {
		return "2"
	}

	return "1"
}

//--------------------------------------------------------------------------------
// Linux Package - Installs Officical Linux Packages.
//--------------------------------------------------------------------------------

func linuxPackage(pkgName string) error {

	var cmd *exec.Cmd

	pm := getPm()
	if pm == "unknown" {
		return fmt.Errorf("failed to identify the package manager.")
	}

	switch pm {
	case "apt":
		cmd = exec.Command("apt-get", "install", "-y", pkgName)
	case "dnf":
		cmd = exec.Command("dnf", "install", "-y", pkgName)
	case "pacman":
		cmd = exec.Command("pacman", "-S", "--noconfirm", "--needed", pkgName)
	case "zypper":
		cmd = exec.Command("zypper", "--non-interactive", "install", pkgName)
	case "apk":
		cmd = exec.Command("apk", "add", pkgName)
	case "pkg": // FreeBSD style package manager used on some Linux/Termux systems
		cmd = exec.Command("pkg", "install", "-y", pkgName)
	case "yum":
		cmd = exec.Command("yum", "install", "-y", pkgName)

	// --- Bonus Package Managers You Missed ---
	case "brew": // Homebrew for Linux (increasingly common on Ubuntu/Fedora)
		cmd = exec.Command("brew", "install", pkgName)
	case "swupd": // Intel's Clear Linux OS bundle manager
		cmd = exec.Command("swupd", "bundle-add", pkgName)
	case "nix-env": // Nix package manager (used on NixOS and standalone)
		cmd = exec.Command("nix-env", "-iA", pkgName)

	default:
		return fmt.Errorf("unsupported package manager: %s", pm)
	}

	// Execute the command and capture output/errors
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install %s using %s: %v\nOutput: %s", pkgName, pm, err, string(output))

	}

	slog.Info("Successfully processed %s via %s\n", pkgName, pm)
	return nil
}

//--------------------------------------------------------------------------------
//	Get Package Manager - Find The Supported Package Manager.
//--------------------------------------------------------------------------------

func getPm() string {
	slog.Info("Analysing the pacakage manager.")
	for _, pm := range []string{"apt", "dnf", "pacman", "zypper", "apk", "pkg", "yum"} {
		if _, err := exec.LookPath(pm); err == nil {
			return pm
		}
	}
	return "unknown"
}

//--------------------------------------------------------------------------------
// URL Package - Downloads A package From A URL
//--------------------------------------------------------------------------------

func urlPkg(fileURL string) error {
	// 1. Parse the URL string to cleanly extract its path component
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("invalid URL structure: %w", err)
	}

	// 2. Extract the absolute last part of the path string
	fileName := path.Base(parsedURL.Path)
	if fileName == "" || fileName == "." || fileName == "/" {
		return errors.New("the URL string does not terminate with a valid file name segment")
	}

	// 3. Define the mandatory target directory and join with the filename
	targetDir := "/tmp"
	finalSavePath := filepath.Join(targetDir, fileName)

	// 4. Ensure the parent directory structure exists with correct permissions
	err = os.MkdirAll(targetDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// 5. Send the HTTP GET request
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad server status: %s", resp.Status)
	}

	// 6. Create the local file at the absolute target path
	out, err := os.Create(finalSavePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 7. Stream the download straight to disk
	slog.Info("Downloading and auto-saving to: ", "path", finalSavePath)
	_, err = io.Copy(out, resp.Body)
	return err
}
