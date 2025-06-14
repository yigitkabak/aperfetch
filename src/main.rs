use std::env;
use std::fs;
use std::process::Command;
use systemstat::{System, Platform, saturating_sub_bytes, IpAddr};
use serde_json::Value;

const RESET: &str = "\x1b[0m";
const BOLD: &str = "\x1b[1m";
const GREEN: &str = "\x1b[32m";
const YELLOW: &str = "\x1b[33m";
const BLUE: &str = "\x1b[34m";
const CYAN: &str = "\x1b[36m";

const CUSTOM_LOGO: &str = r#"

  ___                  
 / _ \                 
/ /_\ \_ __   ___ _ __ 
|  _  | '_ \ / _ \ '__|
| | | | |_) |  __/ |   
\_| |_/ .__/ \___|_|   
      | |              
      |_|              

"#;


fn main() {
    let sys = System::new();

    let logo_color = CYAN;
    let logo = CUSTOM_LOGO;
    let logo_lines: Vec<&str> = logo.trim().lines().collect();
    let logo_width = logo_lines.iter().map(|s| s.chars().count()).max().unwrap_or(0);

    let mut info = Vec::new();

    let user = env::var("USER").or_else(|_| env::var("USERNAME")).unwrap_or_else(|_| "aperfetch".to_string());
    let hostname = get_hostname();
    info.push(format!("{}{}{}@{}{}{}", BOLD, GREEN, user, hostname, RESET, BOLD));
    info.push(format!("{}{}{}{}", CYAN, "-----------------", RESET, BOLD));

    let os_info = get_os_info();
    info.push(format!("{}OS: {}{}", BOLD, os_info, RESET));

    let kernel_info = get_kernel_version();
    info.push(format!("{}Kernel: {}{}", BOLD, kernel_info, RESET));

    if let Ok(uptime) = sys.uptime() {
        let secs = uptime.as_secs();
        let days = secs / 86400;
        let hours = (secs % 86400) / 3600;
        let mins = (secs % 3600) / 60;
        info.push(format!("{}Uptime: {} days, {} hours, {} mins{}", BOLD, days, hours, mins, RESET));
    } else {
        info.push(format!("{}Uptime: Unknown{}", BOLD, RESET));
    }

    if let Ok(shell_path) = env::var("SHELL") {
        let shell_name = std::path::Path::new(&shell_path)
            .file_name()
            .and_then(|s| s.to_str())
            .unwrap_or("sh");
        info.push(format!("{}Shell: {}{}", BOLD, shell_name, RESET));
    } else {
        info.push(format!("{}Shell: Unknown{}", BOLD, RESET));
    }

    let cpu_model = get_cpu_model();
    let cpu_count = get_cpu_count();
    
    let gpu_info = get_gpu_info();
    info.push(format!("{}GPU: {}{}", BOLD, gpu_info, RESET));

    if let Ok(mem) = sys.memory() {
        info.push(format!(
            "{}Memory: {:.2} GiB / {:.2} GiB{}",
            BOLD,
            (mem.total.as_u64() - mem.free.as_u64()) as f64 / 1024.0 / 1024.0 / 1024.0,
            mem.total.as_u64() as f64 / 1024.0 / 1024.0 / 1024.0,
            RESET
        ));
    } else {
        info.push(format!("{}Memory: Unknown{}", BOLD, RESET));
    }

    let disk_info = get_disk_info(&sys);
    info.push(format!("{}Disk: {}{}", BOLD, disk_info, RESET));

    if let Ok(nets) = sys.networks() {
        for net in nets.values() {
            if net.name == "wlan0" || net.name == "eth0" || net.name.starts_with("en") || net.name.starts_with("wl") || net.name.starts_with("usb") {
                if let Some(network_addr) = net.addrs.iter().find(|a| matches!(a.addr, IpAddr::V4(_))) {
                    let ip_str = match network_addr.addr {
                        IpAddr::V4(ip) => ip.to_string(),
                        IpAddr::V6(ip) => ip.to_string(),
                        _ => "N/A".to_string(),
                    };
                    info.push(format!("{}Local IP ({}): {}{}", BOLD, net.name, ip_str, RESET));
                    break;
                }
            }
        }
    }

    let resolution = get_resolution();
    if resolution != "Unknown" {
        info.push(format!("{}Resolution: {}{}", BOLD, resolution, RESET));
    }

    let packages = get_package_count();
    if packages != "Unknown" {
        info.push(format!("{}Packages: {}{}", BOLD, packages, RESET));
    }
    
    let manufacturer = get_device_manufacturer();
    if manufacturer != "Unknown" {
        info.push(format!("{}Manufacturer: {}{}", BOLD, manufacturer, RESET));
    }
    
    let device_model = get_device_model();
    if device_model != "Unknown" {
        info.push(format!("{}Model: {}{}", BOLD, device_model, RESET));
    }

    let battery_status = get_battery_status();
    if battery_status != "Unknown" {
        info.push(format!("{}Battery: {}{}", BOLD, battery_status, RESET));
    }

    let max_lines = std::cmp::max(logo_lines.len(), info.len());

    println!();
    for i in 0..max_lines {
        let logo_part = logo_lines.get(i).cloned().unwrap_or("");
        let info_part = info.get(i).cloned().unwrap_or_default();

        println!("  {}{:<width$}{}   {}", logo_color, logo_part, RESET, info_part, width = logo_width);
    }
    println!();
}

fn is_android() -> bool {
    env::var("ANDROID_ROOT").is_ok()
}

fn get_hostname() -> String {
    Command::new("hostname")
        .output()
        .ok()
        .and_then(|output| if output.status.success() { Some(String::from_utf8_lossy(&output.stdout).trim().to_string()) } else { None })
        .unwrap_or_else(|| "localhost".to_string())
}

fn get_os_info() -> String {
    if is_android() {
        if let Ok(output) = Command::new("getprop").arg("ro.build.version.release").output() {
            if output.status.success() {
                let android_version = String::from_utf8_lossy(&output.stdout).trim().to_string();
                if !android_version.is_empty() {
                    return format!("Android {}", android_version);
                }
            }
        }
        return "Android".to_string();
    }
    if cfg!(target_os = "linux") {
        if let Ok(output) = Command::new("cat").arg("/etc/os-release").output() {
            if output.status.success() {
                for line in String::from_utf8_lossy(&output.stdout).lines() {
                    if line.starts_with("PRETTY_NAME=") {
                        return line.strip_prefix("PRETTY_NAME=").unwrap_or("").trim_matches('"').to_string();
                    }
                }
            }
        }
    } else if cfg!(target_os = "macos") {
        if let Ok(output) = Command::new("sw_vers").arg("-productName").output() {
            if output.status.success() {
                let name = String::from_utf8_lossy(&output.stdout).trim().to_string();
                if let Ok(ver_out) = Command::new("sw_vers").arg("-productVersion").output() {
                    let version = String::from_utf8_lossy(&ver_out.stdout).trim().to_string();
                    return format!("{} {}", name, version);
                }
                return name;
            }
        }
    } else if cfg!(target_os = "windows") {
        if let Ok(output) = Command::new("cmd").args(&["/C", "ver"]).output() {
            if output.status.success() {
                return String::from_utf8_lossy(&output.stdout).trim().to_string();
            }
        }
    }
    "Unknown".to_string()
}

fn get_kernel_version() -> String {
    if cfg!(target_os = "windows") { return "NT".to_string(); }
    Command::new("uname")
        .arg("-r")
        .output()
        .ok()
        .and_then(|output| if output.status.success() { Some(String::from_utf8_lossy(&output.stdout).trim().to_string()) } else { None })
        .unwrap_or_else(|| "Unknown".to_string())
}

fn get_cpu_model() -> String {
    if cfg!(target_os = "linux") {
        if let Ok(output) = Command::new("cat").arg("/proc/cpuinfo").output() {
            let mut hardware = "".to_string();
            for line in String::from_utf8_lossy(&output.stdout).lines() {
                if line.starts_with("model name") || line.starts_with("Processor") {
                    return line.split(':').nth(1).unwrap_or("").trim().to_string();
                }
                if line.starts_with("Hardware") {
                    if let Some(hw) = line.split(':').nth(1) {
                       hardware = hw.trim().to_string();
                    }
                }
            }
            if !hardware.is_empty() {
                return hardware;
            }
        }
    } else if cfg!(target_os = "macos") {
        if let Ok(output) = Command::new("sysctl").arg("-n").arg("machdep.cpu.brand_string").output() {
            return String::from_utf8_lossy(&output.stdout).trim().to_string();
        }
    } else if cfg!(target_os = "windows") {
        if let Ok(output) = Command::new("wmic").args(&["cpu", "get", "name"]).output() {
            return String::from_utf8_lossy(&output.stdout).lines().nth(1).unwrap_or("").trim().to_string();
        }
    }
    "Unknown".to_string()
}

fn get_cpu_count() -> String {
    num_cpus::get().to_string()
}

fn get_disk_info(sys: &System) -> String {
    let path_to_check = if is_android() { "/data" } else { "/" };
    if let Ok(mounts) = sys.mounts() {
        if let Some(mount) = mounts.iter().find(|m| m.fs_mounted_on == path_to_check) {
            let used_bytes = saturating_sub_bytes(mount.total, mount.avail);
            return format!(
                "{:.2} GiB / {:.2} GiB",
                used_bytes.as_u64() as f64 / 1024.0 / 1024.0 / 1024.0,
                mount.total.as_u64() as f64 / 1024.0 / 1024.0 / 1024.0
            );
        }
    }
    "Unknown".to_string()
}

fn get_gpu_info() -> String {
    if is_android() {
        if let Ok(output) = Command::new("dumpsys").args(&["SurfaceFlinger", "--list"]).output() {
            if output.status.success() {
                let stdout = String::from_utf8_lossy(&output.stdout);
                for line in stdout.lines() {
                    if line.contains("GL_RENDERER") {
                        if let Some(renderer) = line.split("GL_RENDERER = ").nth(1) {
                            let gpu_name = renderer.trim().to_string();
                            if !gpu_name.is_empty() {
                                return gpu_name;
                            }
                        }
                    }
                }
            }
        }
        if let Ok(output) = Command::new("getprop").arg("ro.board.platform").output() {
            if output.status.success() {
                let platform = String::from_utf8_lossy(&output.stdout).trim().to_string();
                if !platform.is_empty() {
                    return platform;
                }
            }
        }
    } else if cfg!(target_os = "linux") {
        if let Ok(output) = Command::new("lspci").output() {
            for line in String::from_utf8_lossy(&output.stdout).lines() {
                if line.contains("VGA compatible controller") || line.contains("3D controller") || line.contains("Display controller") {
                    return line.split(": ").last().unwrap_or("").trim().to_string();
                }
            }
        }
    } else if cfg!(target_os = "macos") {
        if let Ok(output) = Command::new("system_profiler").arg("SPDisplaysDataType").output() {
            for line in String::from_utf8_lossy(&output.stdout).lines() {
                if line.trim().starts_with("Chipset Model:") {
                    return line.split(':').nth(1).unwrap_or("").trim().to_string();
                }
            }
        }
    } else if cfg!(target_os = "windows") {
        if let Ok(output) = Command::new("wmic").args(&["path", "win32_VideoController", "get", "name"]).output() {
            return String::from_utf8_lossy(&output.stdout).lines().nth(1).unwrap_or("").trim().to_string();
        }
    }
    "Unknown".to_string()
}

fn get_resolution() -> String {
    if is_android() {
        if let Ok(output) = Command::new("wm").arg("size").output() {
            if let Some(line) = String::from_utf8_lossy(&output.stdout).lines().find(|l| l.contains("Physical size:")) {
                return line.replace("Physical size: ", "").trim().to_string();
            }
        }
    } else if cfg!(target_os = "linux") {
        if let Ok(output) = Command::new("sh").args(&["-c", "xrandr | grep \\* | cut -d' ' -f4"]).output() {
            let res = String::from_utf8_lossy(&output.stdout).trim().to_string();
            if !res.is_empty() { return res; }
        }
    } else if cfg!(target_os = "macos") {
         if let Ok(output) = Command::new("system_profiler").arg("SPDisplaysDataType").output() {
            let stdout = String::from_utf8_lossy(&output.stdout);
            for line in stdout.lines() {
                if line.trim().starts_with("Resolution:") {
                    return line.split(':').nth(1).unwrap_or("").trim().replace(" x ", "x").to_string();
                }
            }
        }
    }
    "Unknown".to_string()
}

fn get_package_count() -> String {
    let mut count_str = "Unknown".to_string();
    if let Ok(output) = Command::new("dpkg").arg("-l").output() {
        if output.status.success() {
            count_str = String::from_utf8_lossy(&output.stdout).lines().filter(|l| l.starts_with("ii")).count().to_string() + " (dpkg)";
        }
    }
    else if let Ok(output) = Command::new("pacman").arg("-Q").output() {
        if output.status.success() {
            count_str = String::from_utf8_lossy(&output.stdout).lines().count().to_string() + " (pacman)";
        }
    }
    else if let Ok(output) = Command::new("rpm").arg("-qa").output() {
         if output.status.success() {
            count_str = String::from_utf8_lossy(&output.stdout).lines().count().to_string() + " (rpm)";
        }
    }
    else if cfg!(target_os = "macos") && Command::new("brew").arg("--version").output().is_ok() {
         if let Ok(output) = Command::new("brew").args(&["list", "--formula"]).output() {
             if output.status.success() {
                count_str = String::from_utf8_lossy(&output.stdout).lines().count().to_string() + " (brew)";
            }
         }
    }
    count_str
}


fn get_battery_status() -> String {
    if is_android() {
        if let Ok(output) = Command::new("termux-battery-status").output() {
            if let Ok(json) = serde_json::from_slice::<Value>(&output.stdout) {
                if let (Some(percentage), Some(status), Some(health)) = (json["percentage"].as_u64(), json["status"].as_str(), json["health"].as_str()) {
                    return format!("{}% [{}, {}]", percentage, status, health);
                }
            }
        }
    } else if cfg!(target_os = "linux") {
        if let Ok(paths) = fs::read_dir("/sys/class/power_supply/") {
            for path in paths.flatten() {
                let dir_name = path.file_name().into_string().unwrap_or_default();
                if dir_name.starts_with("BAT") {
                    let capacity_path = path.path().join("capacity");
                    let status_path = path.path().join("status");
                    if let (Ok(capacity_str), Ok(status_str)) = (fs::read_to_string(capacity_path), fs::read_to_string(status_path)) {
                        return format!("{}% [{}]", capacity_str.trim(), status_str.trim());
                    }
                }
            }
        }
    }
    "Unknown".to_string()
}

fn read_file_trim(path: &str) -> String {
    fs::read_to_string(path).map(|s| s.trim().to_string()).unwrap_or_else(|_| "Unknown".to_string())
}


fn get_device_model() -> String {
    if is_android() {
        return read_from_command("getprop", "ro.product.model");
    } else if cfg!(target_os = "linux") {
        let model = read_file_trim("/sys/class/dmi/id/product_name");
        if model != "Unknown" { return model; }
        let version = read_file_trim("/sys/class/dmi/id/product_version");
        if version != "Unknown" { return version; }
    }
    "Unknown".to_string()
}

fn get_device_manufacturer() -> String {
    if is_android() {
        return read_from_command("getprop", "ro.product.manufacturer");
    } else if cfg!(target_os = "linux") {
        return read_file_trim("/sys/class/dmi/id/sys_vendor");
    }
    "Unknown".to_string()
}

fn read_from_command(cmd: &str, arg: &str) -> String {
     if let Ok(output) = Command::new(cmd).arg(arg).output() {
        if output.status.success() {
            let value = String::from_utf8_lossy(&output.stdout).trim().to_string();
            if !value.is_empty() { return value; }
        }
    }
    "Unknown".to_string()
}
