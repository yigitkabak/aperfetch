package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
)

const (
	RESET = "\x1b[0m"
	BOLD  = "\x1b[1m"
	CYAN  = "\x1b[36m"
)

const CUSTOM_LOGO = `

  ___
 / _ \
/ /_\ \_ __   ___ _ __
|  _  | '_ \ / _ \ '__|
| | | | |_) |  __/ |
\_| |_/ .__/ \___|_|
      | |
      |_|

`

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func main() {
	logoColor := CYAN
	logo := strings.Trim(CUSTOM_LOGO, "\n\r")
	logoLines := strings.Split(logo, "\n")
	logoWidth := 0
	for _, line := range logoLines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}

	var info []string

	user, err := os.UserHomeDir()
	if err != nil {
		user = "aperfetch"
	} else {
		user = filepath.Base(user)
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	info = append(info, fmt.Sprintf("%s%s%s@%s%s", BOLD, CYAN, user, hostname, RESET))
	info = append(info, fmt.Sprintf("%s%s%s", CYAN, "-----------------", RESET))

	info = append(info, fmt.Sprintf("%sOS: %s%s", BOLD, getOsInfo(), RESET))
	info = append(info, fmt.Sprintf("%sKernel: %s%s", BOLD, getKernelVersion(), RESET))
	info = append(info, fmt.Sprintf("%sUptime: %s%s", BOLD, getUptime(), RESET))
	info = append(info, fmt.Sprintf("%sShell: %s%s", BOLD, getShell(), RESET))
	info = append(info, fmt.Sprintf("%sCPU: %s (%d)%s", BOLD, getCpuModel(), runtime.NumCPU(), RESET))
	info = append(info, fmt.Sprintf("%sGPU: %s%s", BOLD, getGpuInfo(), RESET))
	info = append(info, fmt.Sprintf("%sMemory: %s%s", BOLD, getMemory(), RESET))
	info = append(info, fmt.Sprintf("%sDisk: %s%s", BOLD, getDiskInfo(), RESET))
	if ip := getLocalIP(); ip != "Unknown" {
		info = append(info, fmt.Sprintf("%s%s%s", BOLD, ip, RESET))
	}
	if resolution := getResolution(); resolution != "Unknown" {
		info = append(info, fmt.Sprintf("%sResolution: %s%s", BOLD, resolution, RESET))
	}
	if packages := getPackageCount(); packages != "Unknown" {
		info = append(info, fmt.Sprintf("%sPackages: %s%s", BOLD, packages, RESET))
	}
	if manufacturer := getDeviceManufacturer(); manufacturer != "Unknown" {
		info = append(info, fmt.Sprintf("%sManufacturer: %s%s", BOLD, manufacturer, RESET))
	}
	if model := getDeviceModel(); model != "Unknown" {
		info = append(info, fmt.Sprintf("%sModel: %s%s", BOLD, model, RESET))
	}
	if battery := getBatteryStatus(); battery != "Unknown" {
		info = append(info, fmt.Sprintf("%sBattery: %s%s", BOLD, battery, RESET))
	}

	maxLines := len(logoLines)
	if len(info) > maxLines {
		maxLines = len(info)
	}

	fmt.Println()
	for i := 0; i < maxLines; i++ {
		logoPart := ""
		if i < len(logoLines) {
			logoPart = logoLines[i]
		}

		infoPart := ""
		if i < len(info) {
			infoPart = info[i]
		}

		fmt.Printf("  %s%-*s%s   %s\n", logoColor, logoWidth, logoPart, RESET, infoPart)
	}
	fmt.Println()
}

func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(out.String())
}

func readFileTrim(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(content))
}

func isAndroid() bool {
	return os.Getenv("ANDROID_ROOT") != ""
}

func getOsInfo() string {
	if isAndroid() {
		ver := runCommand("getprop", "ro.build.version.release")
		if ver != "Unknown" && ver != "" {
			return "Android " + ver
		}
		return "Android"
	}
	if runtime.GOOS == "linux" {
		if content, err := os.ReadFile("/etc/os-release"); err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(content)))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					return strings.Trim(strings.SplitN(line, "=", 2)[1], "\"")
				}
			}
		}
	}
	info, _, _, err := host.PlatformInformation()
	if err == nil {
		return info
	}
	return "Unknown"
}

func getKernelVersion() string {
	if runtime.GOOS == "windows" {
		return "NT"
	}
	ver, err := host.KernelVersion()
	if err == nil {
		return ver
	}
	return "Unknown"
}

func getUptime() string {
	uptime, err := host.Uptime()
	if err != nil {
		return "Unknown"
	}
	days := uptime / 86400
	hours := (uptime % 86400) / 3600
	mins := (uptime % 3600) / 60
	return fmt.Sprintf("%d days, %d hours, %d mins", days, hours, mins)
}

func getShell() string {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		if runtime.GOOS == "windows" {
			return "PowerShell/CMD"
		}
		return "Unknown"
	}
	return filepath.Base(shellPath)
}

func getCpuModel() string {
	if runtime.GOOS == "linux" && isAndroid() {
		hw := runCommand("getprop", "ro.board.platform")
		if hw != "Unknown" && hw != "" {
			return hw
		}
	}
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		return cpuInfo[0].ModelName
	}
	return "Unknown"
}

func getMemory() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "Unknown"
	}
	usedGiB := float64(memInfo.Used) / (1024 * 1024 * 1024)
	totalGiB := float64(memInfo.Total) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2f GiB / %.2f GiB", usedGiB, totalGiB)
}

func getDiskInfo() string {
	path := "/"
	if isAndroid() {
		path = "/data"
	}
	diskInfo, err := disk.Usage(path)
	if err != nil {
		return "Unknown"
	}
	usedGiB := float64(diskInfo.Used) / (1024 * 1024 * 1024)
	totalGiB := float64(diskInfo.Total) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2f GiB / %.2f GiB", usedGiB, totalGiB)
}

func getGpuInfo() string {
	if isAndroid() {
		output := runCommand("getprop", "ro.board.platform")
		if output != "Unknown" && output != "" {
			return output
		}
	}
	switch runtime.GOOS {
	case "linux":
		output := runCommand("lspci")
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "VGA compatible controller") || strings.Contains(line, "3D controller") || strings.Contains(line, "Display controller") {
				parts := strings.Split(line, ": ")
				if len(parts) > 1 {
					return parts[len(parts)-1]
				}
			}
		}
	case "darwin":
		output := runCommand("system_profiler", "SPDisplaysDataType")
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "Chipset Model:") {
				return strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			}
		}
	case "windows":
		output := runCommand("wmic", "path", "win32_VideoController", "get", "name")
		lines := strings.Split(output, "\n")
		if len(lines) > 1 {
			return strings.TrimSpace(lines[1])
		}
	}
	return "Unknown"
}

func getLocalIP() string {
	ifaces, err := psnet.Interfaces()
	if err != nil {
		return "Unknown"
	}

	for _, iface := range ifaces {
		if !contains(iface.Flags, "up") || contains(iface.Flags, "loopback") {
			continue
		}

		for _, addr := range iface.Addrs {
			var ip net.IP
			if strings.Contains(addr.Addr, "/") {
				parsedIP, _, err := net.ParseCIDR(addr.Addr)
				if err == nil {
					ip = parsedIP
				}
			} else {
				ip = net.ParseIP(addr.Addr)
			}

			if ip != nil && ip.To4() != nil {
				return fmt.Sprintf("Local IP (%s): %s", iface.Name, ip.String())
			}
		}
	}
	return "Unknown"
}

func getResolution() string {
	switch runtime.GOOS {
	case "linux":
		if isAndroid() {
			output := runCommand("wm", "size")
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Physical size:") {
					return strings.TrimSpace(strings.Replace(line, "Physical size:", "", 1))
				}
			}
		} else {
			output := runCommand("sh", "-c", "xrandr 2>/dev/null | grep '*' | awk '{print $1}'")
			if output != "Unknown" && output != "" {
				return output
			}
		}
	case "darwin":
		output := runCommand("system_profiler", "SPDisplaysDataType")
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "Resolution:") {
				res := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				return strings.Replace(res, " x ", "x", 1)
			}
		}
	}
	return "Unknown"
}

func getPackageCount() string {
	if _, err := exec.LookPath("dpkg"); err == nil {
		output := runCommand("sh", "-c", "dpkg -l | grep -c '^ii'")
		return output + " (dpkg)"
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		output := runCommand("sh", "-c", "pacman -Q | wc -l")
		return output + " (pacman)"
	}
	if _, err := exec.LookPath("rpm"); err == nil {
		output := runCommand("sh", "-c", "rpm -qa | wc -l")
		return output + " (rpm)"
	}
	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("brew"); err == nil {
			output := runCommand("sh", "-c", "brew list --formula | wc -l")
			return output + " (brew)"
		}
	}
	return "Unknown"
}

func getDeviceManufacturer() string {
	if isAndroid() {
		return runCommand("getprop", "ro.product.manufacturer")
	}
	if runtime.GOOS == "linux" {
		return readFileTrim("/sys/class/dmi/id/sys_vendor")
	}
	return "Unknown"
}

func getDeviceModel() string {
	if isAndroid() {
		return runCommand("getprop", "ro.product.model")
	}
	if runtime.GOOS == "linux" {
		model := readFileTrim("/sys/class/dmi/id/product_name")
		if model != "Unknown" {
			return model
		}
		return readFileTrim("/sys/class/dmi/id/product_version")
	}
	return "Unknown"
}

func getBatteryStatus() string {
	if isAndroid() {
		if _, err := exec.LookPath("termux-battery-status"); err == nil {
			output := runCommand("termux-battery-status")
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(output), &result); err == nil {
				percentage := result["percentage"].(float64)
				status := result["status"].(string)
				health := result["health"].(string)
				return fmt.Sprintf("%.0f%% [%s, %s]", percentage, status, health)
			}
		}
	}
	if runtime.GOOS == "linux" {
		batDirs, err := os.ReadDir("/sys/class/power_supply")
		if err != nil {
			return "Unknown"
		}
		for _, dir := range batDirs {
			if strings.HasPrefix(dir.Name(), "BAT") {
				capacityPath := filepath.Join("/sys/class/power_supply", dir.Name(), "capacity")
				statusPath := filepath.Join("/sys/class/power_supply", dir.Name(), "status")
				capacity := readFileTrim(capacityPath)
				status := readFileTrim(statusPath)
				if capacity != "Unknown" && status != "Unknown" {
					return fmt.Sprintf("%s%% [%s]", capacity, status)
				}
			}
		}
	}
	return "Unknown"
}
