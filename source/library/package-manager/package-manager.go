package packagemanager

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
	"yathahi/source/library"
)

func PackageManager(pkgName string) {
	slog.Info("Package Manager initialized", "package", pkgName)

	switch pkgName {
	case library.CorePackageName,
		library.AdministratorPackageName,
		library.ClientPackageName,
		library.WebPackageName,
		library.DnsPackageName:
		if err := CheckComponentAvailability(pkgName); err != nil {
			slog.Error("Component setup failed",
				"category", "yathahi_component",
				"package", pkgName,
				"error", err,
			)
		}

	default:
		if err := linuxPackage(pkgName); err != nil {
			slog.Error("Linux package installation failed",
				"category", "linux_package",
				"package", pkgName,
				"error", err,
			)
		}
	}
}

//--------------------------------------------------------------------------------
// Linux Package - Installs Official Linux Packages.
//--------------------------------------------------------------------------------

func linuxPackage(pkgName string) error {
	pm := getPm()
	if pm == "unknown" {
		return errors.New("failed to identify a supported package manager")
	}

	var cmd *exec.Cmd
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
	case "pkg":
		cmd = exec.Command("pkg", "install", "-y", pkgName)
	case "yum":
		cmd = exec.Command("yum", "install", "-y", pkgName)
	case "brew":
		cmd = exec.Command("brew", "install", pkgName)
	case "swupd":
		cmd = exec.Command("swupd", "bundle-add", pkgName)
	case "nix-env":
		cmd = exec.Command("nix-env", "-iA", pkgName)
	default:
		return fmt.Errorf("unsupported package manager: %s", pm)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install %s using %s: %w\nOutput: %s", pkgName, pm, err, string(output))
	}

	slog.Info("Successfully processed package", "package", pkgName, "manager", pm)
	return nil
}

//--------------------------------------------------------------------------------
// Get Package Manager - Find The Supported Package Manager.
//--------------------------------------------------------------------------------

func getPm() string {
	slog.Info("Analysing the package manager")
	managers := []string{"apt", "dnf", "pacman", "zypper", "apk", "pkg", "yum", "brew", "swupd", "nix-env"}
	for _, pm := range managers {
		if _, err := exec.LookPath(pm); err == nil {
			return pm
		}
	}
	return "unknown"
}

//--------------------------------------------------------------------------
// Checking the availability of the component
//--------------------------------------------------------------------------

func CheckComponentAvailability(component string) error {
	binaryPath := filepath.Join(library.BinaryFolder, component)

	if _, err := os.Stat(binaryPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to check component %s: %w", component, err)
		}
		slog.Info("Component is not installed, starting installation", "component", component)
		return ComponentInstaller(component)
	}

	slog.Info("Component binary is already installed", "component", component)

	serviceName := component + ".service"
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	if err != nil || status != "active" {
		slog.Info("Service is not active, enabling now", "service", serviceName, "status", status)
		return enableComponent(component)
	}

	slog.Info("Component service is already active", "service", serviceName)
	return nil
}

//--------------------------------------------------------------------------
// Install the component
//--------------------------------------------------------------------------

func ComponentInstaller(component string) error {
	downloadURL := library.BinaryURL + component + ".tar.zst"
	archivePath, err := urlPkg(downloadURL)
	if err != nil {
		return fmt.Errorf("unable to download %s: %w", component, err)
	}
	// Clean up the downloaded archive when finished
	defer os.Remove(archivePath)

	extractCmd := exec.Command("tar", "--zstd", "-xf", archivePath, "-C", "/tmp/")
	if output, err := extractCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unable to extract %s: %w\nOutput: %s", archivePath, err, string(output))
	}

	extractedBinary := filepath.Join("/tmp", component)
	destBinary := filepath.Join(library.BinaryFolder, component)

	if err := os.MkdirAll(library.BinaryFolder, 0700); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	mvCmd := exec.Command("mv", extractedBinary, destBinary)
	if output, err := mvCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unable to move %s to %s: %w\nOutput: %s", component, library.BinaryFolder, err, string(output))
	}

	mvCmd = exec.Command("mv", "/tmp/"+component+".service", "/etc/systemd/system/"+component+".service")
	if output, err := mvCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unable to move %s to %s: %w\nOutput: %s", component, "/etc/systemd/system/", err, string(output))
	}

	slog.Info("Binary successfully installed", "component", component, "destination", destBinary)
	return enableComponent(component)
}

//--------------------------------------------------------------------------------
// URL Package - Downloads A Package From A URL
//--------------------------------------------------------------------------------

func urlPkg(fileURL string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL structure: %w", err)
	}

	fileName := path.Base(parsedURL.Path)
	if fileName == "" || fileName == "." || fileName == "/" {
		return "", errors.New("the URL string does not terminate with a valid file name segment")
	}

	targetDir := "/tmp"
	finalSavePath := filepath.Join(targetDir, fileName)

	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create directory structure: %w", err)
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad server status: %s", resp.Status)
	}

	out, err := os.Create(finalSavePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	slog.Info("Downloading and saving package", "path", finalSavePath)
	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed while writing file: %w", err)
	}

	return finalSavePath, nil
}

//--------------------------------------------------------------------------
// Enabling the component
//--------------------------------------------------------------------------

func enableComponent(component string) error {
	reloadCmd := exec.Command("systemctl", "daemon-reload")
	if output, err := reloadCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unable to reload systemd daemon: %w\nOutput: %s", err, string(output))
	}

	serviceName := component + ".service"
	slog.Info("Enabling and starting service", "service", serviceName)

	enableCmd := exec.Command("systemctl", "enable", "--now", serviceName)
	if output, err := enableCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unable to enable service %s: %w\nOutput: %s", serviceName, err, string(output))
	}

	checkCmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := checkCmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	if err != nil || status != "active" {
		return fmt.Errorf("service %s failed to start (status: %s): %w", serviceName, status, err)
	}

	slog.Info("Component service is running actively", "service", serviceName)
	return nil
}
