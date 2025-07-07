package main

import (
        "bufio"
        "bytes"
        "encoding/json"
        "flag"
        "fmt"
        "net"
        "os"
        "os/exec"
        "path/filepath"
        "runtime"
        "strings"
        "time"

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

var distroLogos = map[string]string{
        "CachyOS": `
          .-------------------------:
          .+=========================.
         :++===++==================-       :++-
        :*++====+++++=============-        .==:
       -*+++=====+***++==========:
      =*++++========------------:
     =*+++++=====-                     ...
   .+*+++++=-===:                    .=+++=:
  :++++=====-==:                     -*****+
 :++========-=.                      .=+**+.
.+==========-.                          .
 :+++++++====-                                .--==-.
  :++==========.                             :+++++++:
    .-===========:                           .+*****+:
      -=======++++:::::::::::::::::::::::::-:  .---:
       :======++++====+++******************=.
        :=====+++==========++++++++++++++*-
         .====++==============++++++++++*-
          .===+==================+++++++:
           .-=======================+++:
             ..........................
`,
        "Zorin OS": `
       "osssssssssssssssssssso"
       .osssssssssssssssssssssso.
      .+oooooooooooooooooooooooo+.


  "::::::::::::::::::::::.         .:"
 "+ssssssssssssssssss+:."     "".:+ssso"
.ossssssssssssssso/.       "-+ossssssso.
ssssssssssssso/-"      "-/osssssssssssss
.ossssssso/-"      .-/ossssssssssssssso.
 "+sss+:.      ".:+ssssssssssssssssss+"
  "":.         .::::::::::::::::::::::"


      .+oooooooooooooooooooooooo+.
       -osssssssssssssssssssssso-
        "osssssssssssssssssssso"
    `,

        "OpenSUSE": `
           .;ldkO0000Okdl;.
       .;d00xl:^''''''^:ok00d;.
     .d00l'                'o00d.
   .d0Kd'  Okxol:;,.          :O0d.
  .OKKKK0kOKKKKKKKKKKOxo:,      lKO.
 ,0KKKKKKKKKKKKKKKK0P^,,,^dx:    ;00,
.OKKKKKKKKKKKKKKKKk'.oOPPb.'0k.   cKO.
:KKKKKKKKKKKKKKKKK: kKx..dd lKd   'OK:
dKKKKKKKKKKKOx0KKKd ^0KKKO' kKKc   dKd
dKKKKKKKKKKKK;.;oOKx,..^..;kKKK0.  dKd
:KKKKKKKKKKKK0o;...^cdxxOK0O/^^'  .0K:
 kKKKKKKKKKKKKKKK0x;,,......,;od  lKk
 '0KKKKKKKKKKKKKKKKKKKKK00KKOo^  c00'
  'kKKKOxddxkOO00000Okxoc;''   .dKk'
    l0Ko.                    .c00l'
     'l0Kk:.              .;xK0l'
        'lkK0xl:;,,,,;:ldO0kl'
            '^:ldxkkkkxdl:^'
    `,
        "elementary OS": `
         eeeeeeeeeeeeeeeee
      eeeeeeeeeeeeeeeeeeeeeee
    eeeee  eeeeeeeeeeee   eeeee
  eeee   eeeee       eee     eeee
 eeee   eeee          eee     eeee
eee    eee            eee       eee
eee   eee            eee        eee
ee    eee           eeee       eeee
ee    eee         eeeee      eeeeee
ee    eee       eeeee      eeeee ee
eee   eeee   eeeeee      eeeee  eee
eee    eeeeeeeeee     eeeeee    eee
 eeeeeeeeeeeeeeeeeeeeeeee    eeeee
  eeeeeeee eeeeeeeeeeee      eeee
    eeeee                 eeeee
      eeeeeee         eeeeeee
         eeeeeeeeeeeeeeeee
    `,
        "Arch": `
                   -"
                  .o+"
                 "ooo/
                "+oooo:
               "+oooooo:
               -+oooooo+:
             "/:-:++oooo+:
            "/++++/+++++++:
           "/++++++++++++++:
          "/+++ooooooooooooo/"
         ./ooosssso++osssssso+"
        .oossssso-""""/ossssss+"
       -osssssso.      :ssssssso.
      :osssssss/        osssso+++.
     /ossssssss/        +ssssooo/-
   "/ossssso+/:-        -:/+osssso+-
  "+sso+:-"                 ".-/+oso:
 "++:.                           "-/+/
 ."                                 "/
`,
        "EndeavourOS": `
                     ./o.
                   ./sssso-
                 ":osssssss+-
               ":+sssssssssso/.
             "-/ossssssssssssso/.
           "-/+sssssssssssssssso+:
         "-:/+sssssssssssssssssso+/.
       ".://osssssssssssssssssssso++-
      .://+ssssssssssssssssssssssso++:
    .:///ossssssssssssssssssssssssso++:
  ":////ssssssssssssssssssssssssssso+++.
"-////+ssssssssssssssssssssssssssso++++-
 "..-+oosssssssssssssssssssssssso+++++/"
   ./++++++++++++++++++++++++++++++/:.
  ":::::::::::::::::::::::::------""
    `,
        "Pop!_OS": `
            /////////////
         /////////////////////
      ///////*767////////////////
    //////7676767676*//////////////
   /////76767//7676767//////////////
  /////767676///*76767///////////////
 ///////767676///76767.///7676*///////
/////////767676//76767///767676////////
//////////76767676767////76767/////////
///////////76767676//////7676//////////
////////////,7676,///////767///////////
/////////////*7676///////76////////////
///////////////7676////////////////////
 ///////////////7676///767////////////
  //////////////////////'////////////
   //////.7676767676767676767,//////
    /////767676767676767676767/////
      ///////////////////////////
         /////////////////////
             /////////////
    `,
        "Fedora": `
             .',;::::;,'.
         .';:cccccccccccc:;,.
      .;cccccccccccccccccccccc;.
    .:cccccccccccccccccccccccccc:.
  .;ccccccccccccc;.:dddl:.;ccccccc;.
 .:ccccccccccccc;OWMKOOXMWd;ccccccc:.
.:ccccccccccccc;KMMc;cc;xMMc;ccccccc:.
,cccccccccccccc;MMM.;cc;;WW:;cccccccc,
:cccccccccccccc;MMM.;cccccccccccccccc:
:ccccccc;oxOOOo;MMM0OOk.;cccccccccccc:
cccccc;0MMKxdd:;MMMkddc.;cccccccccccc;
ccccc;XM0';cccc;MMM.;cccccccccccccccc'
ccccc;MMo;ccccc;MMW.;ccccccccccccccc;
ccccc;0MNc.ccc.xMMd;ccccccccccccccc;
cccccc;dNMWXXXWM0:;cccccccccccccc:,
cccccccc;.:odl:.;cccccccccccccc:,.
:cccccccccccccccccccccccccccc:'.
.:cccccccccccccccccccccc:;,..
  '::cccccccccccccc::;,.
    `,
        "Android": `
         -o          o-            
          +hydNNNNdyh+           
        +mMMMMMMMMMMMMm+           
      "dMMm:NMMMMMMN:mMMd"         
      hMMMMMMMMMMMMMMMMMMh        
  ..  yyyyyyyyyyyyyyyyyyyy  ..    
.mMMm"MMMMMMMMMMMMMMMMMMMM"mMMm.   
:MMMM-MMMMMMMMMMMMMMMMMMMM-MMMM:   
:MMMM-MMMMMMMMMMMMMMMMMMMM-MMMM:  
:MMMM-MMMMMMMMMMMMMMMMMMMM-MMMM:   
:MMMM-MMMMMMMMMMMMMMMMMMMM-MMMM:  
-MMMM-MMMMMMMMMMMMMMMMMMMM-MMMM-   
 +yy+ MMMMMMMMMMMMMMMMMMMM +yy+     
      mMMMMMMMMMMMMMMMMMMm        
      "/++MMMMh++hMMMM++/"       
          MMMMo  oMMMM              
          MMMMo  oMMMM            
          oNMm-  -mMNs

        `,
        "Ubuntu": `
            .-/+oossssoo+\-.              
        ":+ssssssssssssssssss+:"          
      -+ssssssssssssssssssyyssss+-        
    .ossssssssssssssssssdMMMNysssso.      
   /ssssssssssshdmmNNmmyNMMMMhssssss\     
  +ssssssssshmydMMMMMMMNddddyssssssss+     
 /sssssssshNMMMyhhyyyyhmNMMMNhssssssss\    
.ssssssssdMMMNhsssssssssshNMMMdssssssss.   
+sssshhhyNMMNyssssssssssssyNMMMysssssss+  
ossyNMMMNyMMhsssssssssssssshmmmhssssssso   
ossyNMMMNyMMhsssssssssssssshmmmhssssssso   
+sssshhhyNMMNyssssssssssssyNMMMysssssss+   
.ssssssssdMMMNhsssssssssshNMMMdssssssss.   
 \sssssssshNMMMyhhyyyyhdNMMMNhssssssss/     
  +sssssssssdmydMMMMMMMMddddyssssssss+    
   \ssssssssssshdmNNNNmyNMMMMhssssss/       
    .ossssssssssssssssssdMMMNysssso.       
      -+sssssssssssssssssyyyssss+-
        ":+ssssssssssssssssss+:"  
            .-\+oossssoo+/-.
`,
        "Debian": `
       _,met$$$$$gg.
    ,g$$$$$$$$$$$$$$$P.
  ,g$$P"        """Y$$."
 ,$$P'              .$$$."
',$$P       ,ggs.     "$$b:
"d$$'     ,$P"'   .    $$$
 $$P      d$'     ,    $$P
 $$:      $$.   -    ,d$$'
 $$;      Y$b._   _,d$P'
 Y$$.    ".""Y$$$$P"'
 "$$b      "-.__
  "Y$$
   "Y$$.
     "$$b.
       "Y$$b.
          ""Y$b._
              """"

`,
        "Manjaro": `
██████████████████  ████████
██████████████████  ████████
██████████████████  ████████
██████████████████  ████████
████████            ████████
████████  ████████  ████████
████████  ████████  ████████
████████  ████████  ████████
████████  ████████  ████████
████████  ████████  ████████
████████  ████████  ████████
████████  ████████  ████████
████████  ████████  ████████
████████  ████████  ████████
`,
        "NixOS": `
          ▗▄▄▄       ▗▄▄▄▄    ▄▄▄▖
          ▜███▙       ▜███▙  ▟███▛
           ▜███▙       ▜███▙▟███▛
            ▜███▙       ▜██████▛
     ▟█████████████████▙ ▜████▛     ▟▙
    ▟███████████████████▙ ▜███▙    ▟██▙
           ▄▄▄▄▖           ▜███▙  ▟███▛
          ▟███▛             ▜██▛ ▟███▛
         ▟███▛               ▜▛ ▟███▛
▟███████████▛                  ▟██████████▙
▜██████████▛                  ▟███████████▛
      ▟███▛ ▟▙               ▟███▛
     ▟███▛ ▟██▙             ▟███▛
    ▟███▛  ▜███▙           ▝▀▀▀▀
    ▜██▛    ▜███▙ ▜██████████████████▛
     ▜▛     ▟████▙ ▜████████████████▛
           ▟██████▙       ▜███▙
          ▟███▛▜███▙       ▜███▙
         ▟███▛  ▜███▙       ▜███▙
         ▝▀▀▀    ▀▀▀▀▘       ▀▀▀▘
`,
        "Mint": `
             ...-:::::-...
          .-MMMMMMMMMMMMMMM-.
      .-MMMM"..-:::::::-.."MMMM-.
    .:MMMM.:MMMMMMMMMMMMMMM:.MMMM:.
   -MMM-M---MMMMMMMMMMMMMMMMMMM.MMM-
 ":MMM:MM"  :MMMM:....::-...-MMMM:MMM:"
 :MMM:MMM"  :MM:"  ""    ""  ":MMM:MMM:
.MMM.MMMM"  :MM.  -MM.  .MM-  "MMMM.MMM.
:MMM:MMMM"  :MM.  -MM-  .MM:  "MMMM-MMM:
:MMM:MMMM"  :MM.  -MM-  .MM:  "MMMM:MMM:
:MMM:MMMM"  :MM.  -MM-  .MM:  "MMMM-MMM:
.MMM.MMMM"  :MM:--:MM:--:MM:  "MMMM.MMM.
 :MMM:MMM-  "-MMMMMMMMMMMM-"  -MMM-MMM:
  :MMM:MMM:"                ":MMM:MMM:
   .MMM.MMMM:--------------:MMMM.MMM.
     '-MMMM.-MMMMMMMMMMMMMMM-.MMMM-'
       '.-MMMM""--:::::--""MMMM-.'
            '-MMMMMMMMMMMMM-'
               ""-:::::-""
`,
        "Gentoo": `
         -/oyddmdhs+:.
     -odNMMMMMMMMNNmhy+-"
   -yNMMMMMMMMMMMNNNmmdhy+-
 "omMMMMMMMMMMMMNmdmmmmddhhy/"
 omMMMMMMMMMMMNhhyyyohmdddhhhdo"
.ydMMMMMMMMMMdhs++so/smdddhhhhdm+"
 oyhdmNMMMMMMMNdyooydmddddhhhhyhNd.
  :oyhhdNNMMMMMMMNNNmmdddhhhhhyymMh
    .:+sydNMMMMMNNNmmmdddhhhhhyymMmy
       /mMMMMMMNNNmmmdddhhhhhmMhs:
    "oNMMMMMMMNNNmmmddddhhdmMNhs+"
  "sNMMMMMMMMNNNmmmdddddmNMmhs/.
 /NMMMMMMMMNNNNmmmdddmNMNdso:"
+MMMMMMMNNNNNmmmmdmNMNdso/-
yMMNNNNNNNmmmmmNNMmhs+/-"
/hMMNNNNNNNNMNdhs++/-"
"/ohdmmddhys+++/:."
  "-//////:--.

`,
        "Pardus": `
 .smNdy+-    ".:/osyyso+:."    -+ydmNs.
/Md- -/ymMdmNNdhso/::/oshdNNmdMmy/. :dM/
mN.     oMdyy- -y          "-dMo     .Nm
.mN+"  sMy hN+ -:             yMs  "+Nm.
 "yMMddMs.dy "+"               sMddMMy"
   +MMMo  ."  .                 oMMM+
   "NM/    """""."    "."""""    +MN"
   yM+   ".-:yhomy    ymohy:-."   +My
   yM:          yo    oy          :My
   +Ms         .N"    "N.      +h sM+
   "MN      -   -::::::-   : :o:+"NM"
    yM/    sh   -dMMMMd-   ho  +y+My
    .dNhsohMh-//: /mm/ ://-yMyoshNd"
      "-ommNMm+:/. oo ./:+mMNmmo:"
     "/o+.-somNh- :yy: -hNmos-.+o/"
    ./" .s/"s+sMdd+""+ddMs+s"/s. "/.
        : -y.  -hNmddmNy.  .y- :
         -+       ".."       +-

`,
        "Windows": `
        ,.=:!!t3Z3z.,
       :tt:::tt333EE3
       Et:::ztt33EEEL @Ee.,      ..,
      ;tt:::tt333EE7 ;EEEEEEttttt33#
     :Et:::zt333EEQ. $EEEEEttttt33QL
     it::::tt333EEF @EEEEEEttttt33F
    ;3=*^"""*4EEV :EEEEEEttttt33@.
    ,.=::::!t=., " @EEEEEEtttz33QF
   ;::::::::zt33)   "4EEEtttji3P*
  :t::::::::tt33.:Z3z..  "" ,..g.
  i::::::::zt33F AEEEtttt::::ztF
 ;:::::::::t33V ;EEEttttt::::t3
 E::::::::zt33L @EEEtttt::::z3F
{3=*^"""*4E3) ;EEEtttt:::::tZ"
             " :EEEEtttt::::z7
                 "VEzjt:;;z>*"
`,
        "macOS": `
                    c.'
                 ,xNMM.
               .OMMMMo
               lMM"
     .;loddo:.  .olloddol;.
   cKMMMMMMMMMMNWMMMMMMMMMM0:
 .KMMMMMMMMMMMMMMMMMMMMMMMWd.
 XMMMMMMMMMMMMMMMMMMMMMMMX.
;MMMMMMMMMMMMMMMMMMMMMMMM:
:MMMMMMMMMMMMMMMMMMMMMMMM:
.MMMMMMMMMMMMMMMMMMMMMMMMX.
 kMMMMMMMMMMMMMMMMMMMMMMMMWd.
 'XMMMMMMMMMMMMMMMMMMMMMMMMk
  'XMMMMMMMMMMMMMMMMMMMMMMMMK.
    kMMMMMMMMMMMMMMMMMMMMMMd
     ;KMMMMMMMWXXWMMMMMMMk.
       "cooc*"    "*coo'"
`,
        "aper": `
  ___
 / _ \
/ /_\ \_ __   ___ _ __
|  _  | '_ \ / _ \ '__|
| | | | |_) |  __/ |
\_| |_| .__/ \___|_|
      | |
      |_|
        `,
        "Linux": `
    .--.
   |o_o |
   |:_/ |
  //   \ \
 (|     | )
/'\_   _/` + "`" + `\
\___)=(___/
`,
}

const CUSTOM_LOGO = `
  ___
 / _ \
/ /_\ \_ __   ___ _ __
|  _  | '_ \ / _ \ '__|
| | | | |_) |  __/ |
\_| |_| .__/ \___|_|
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

func getLogo(osName string) string {
        osNameLower := strings.ToLower(osName)
        for distro, logo := range distroLogos {
                if strings.Contains(osNameLower, strings.ToLower(distro)) {
                        return strings.Trim(logo, "\n\r")
                }
        }
        if runtime.GOOS == "linux" {
                return strings.Trim(distroLogos["Linux"], "\n\r")
        }
        if runtime.GOOS == "darwin" {
                return strings.Trim(distroLogos["macOS"], "\n\r")
        }
        if runtime.GOOS == "windows" {
                return strings.Trim(distroLogos["Windows"], "\n\r")
        }
        return strings.Trim(CUSTOM_LOGO, "\n\r")
}

func main() {
        asciiDistro := flag.String("ascii_distro", "", "Force a specific ASCII distro logo (e.g., macos, arch, ubuntu)")
        flag.Parse()

        osInfo := getOsInfo()
        var logo string

        if *asciiDistro != "" {

                found := false
                for distroName, distroLogo := range distroLogos {
                        if strings.EqualFold(distroName, *asciiDistro) {
                                logo = strings.Trim(distroLogo, "\n\r")
                                found = true
                                break
                        }
                }
                if !found {
                        fmt.Printf("Warning: ASCII logo for '%s' not found. Falling back to default.\n", *asciiDistro)
                        logo = getLogo(osInfo)
                }
        } else {

                logo = getLogo(osInfo)
        }

        logoColor := CYAN
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
                user = "fetcher"
        } else {
                user = filepath.Base(user)
        }
        hostname, err := os.Hostname()
        if err != nil {
                hostname = "localhost"
        }
        info = append(info, fmt.Sprintf("%s%s%s@%s%s", BOLD, CYAN, user, hostname, RESET))
        info = append(info, fmt.Sprintf("%s%s%s", CYAN, "-----------------", RESET))

        info = append(info, fmt.Sprintf("%sOS: %s%s", BOLD, osInfo, RESET))
        info = append(info, fmt.Sprintf("%sKernel: %s%s", BOLD, getKernelVersion(), RESET))
        info = append(info, fmt.Sprintf("%sUptime: %s%s", BOLD, getUptime(), RESET))
        info = append(info, fmt.Sprintf("%sPackages: %s%s", BOLD, getPackageCount(), RESET))
        if de := getDesktopEnvironment(); de != "Unknown" {
                info = append(info, fmt.Sprintf("%sDE: %s%s", BOLD, de, RESET))
        }
        info = append(info, fmt.Sprintf("%sShell: %s%s", BOLD, getShell(), RESET))
        if term := getTerminal(); term != "Unknown" {
                info = append(info, fmt.Sprintf("%sTerminal: %s%s", BOLD, term, RESET))
        }
        if resolution := getResolution(); resolution != "Unknown" {
                info = append(info, fmt.Sprintf("%sResolution: %s%s", BOLD, resolution, RESET))
        }
        if theme := getTheme(); theme != "Unknown" {
                info = append(info, fmt.Sprintf("%sTheme: %s%s", BOLD, theme, RESET))
        }
        if icons := getIcons(); icons != "Unknown" {
                info = append(info, fmt.Sprintf("%sIcons: %s%s", BOLD, icons, RESET))
        }
        info = append(info, fmt.Sprintf("%sCPU: %s (%d) [%s]%s", BOLD, getCpuModel(), runtime.NumCPU(), getCPUUsage(), RESET))
        info = append(info, fmt.Sprintf("%sGPU: %s%s", BOLD, getGpuInfo(), RESET))
        info = append(info, fmt.Sprintf("%sMemory: %s%s", BOLD, getMemory(), RESET))
        info = append(info, fmt.Sprintf("%sDisk: %s%s", BOLD, getDiskInfo(), RESET))
        if battery := getBatteryStatus(); battery != "Unknown" {
                info = append(info, fmt.Sprintf("%sBattery: %s%s", BOLD, battery, RESET))
        }
        if locale := getLocale(); locale != "Unknown" {
                info = append(info, fmt.Sprintf("%sLocale: %s%s", BOLD, locale, RESET))
        }
        if ip := getLocalIP(); ip != "Unknown" {
                info = append(info, fmt.Sprintf("%s%s%s", BOLD, ip, RESET))
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
        if runtime.GOOS == "darwin" {
                prodName := runCommand("sw_vers", "-productName")
                prodVer := runCommand("sw_vers", "-productVersion")
                if prodName != "Unknown" && prodVer != "Unknown" {
                        return fmt.Sprintf("%s %s", prodName, prodVer)
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
                output := runCommand("ver")
                if output != "Unknown" {
                        return output
                }
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
        } else if runtime.GOOS == "windows" {
                path = "C:"
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
                lines := strings.Split(strings.TrimSpace(output), "\r\n")
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
                        output := runCommand("sh", "-c", "xrandr --current 2>/dev/null | grep '*' | uniq | awk '{print $1}'")
                        if output != "Unknown" && output != "" {
                                return strings.Split(output, "\n")[0]
                        }
                }
        case "darwin":
                output := runCommand("system_profiler", "SPDisplaysDataType")
                scanner := bufio.NewScanner(strings.NewReader(output))
                for scanner.Scan() {
                        line := strings.TrimSpace(scanner.Text())
                        if strings.HasPrefix(line, "Resolution:") {
                                res := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
                                if strings.Contains(res, "UI") {
                                        res = strings.Split(res, " ")[0]
                                }
                                return strings.Replace(res, " x ", "x", 1)
                        }
                }
        case "windows":
                output := runCommand("powershell", "-Command", "Get-WmiObject -Class Win32_VideoController | Select-Object -Property VideoModeDescription, CurrentHorizontalResolution, CurrentVerticalResolution | Format-List")
                lines := strings.Split(output, "\n")
                var h, v string
                for _, line := range lines {
                        if strings.HasPrefix(strings.TrimSpace(line), "CurrentHorizontalResolution") {
                                parts := strings.Split(line, ":")
                                if len(parts) > 1 {
                                        h = strings.TrimSpace(parts[1])
                                }
                        }
                        if strings.HasPrefix(strings.TrimSpace(line), "CurrentVerticalResolution") {
                                parts := strings.Split(line, ":")
                                if len(parts) > 1 {
                                        v = strings.TrimSpace(parts[1])
                                }
                        }
                }
                if h != "" && v != "" {
                        return fmt.Sprintf("%sx%s", h, v)
                }
        }
        return "Unknown"
}

func getPackageCount() string {
        if _, err := exec.LookPath("dpkg"); err == nil {
                output := runCommand("sh", "-c", "dpkg -l | grep -c '^ii'")
                if output != "Unknown" {
                        return output + " (dpkg)"
                }
        }
        if _, err := exec.LookPath("pacman"); err == nil {
                output := runCommand("sh", "-c", "pacman -Q | wc -l")
                if output != "Unknown" {
                        return output + " (pacman)"
                }
        }
        if _, err := exec.LookPath("rpm"); err == nil {
                output := runCommand("sh", "-c", "rpm -qa | wc -l")
                if output != "Unknown" {
                        return output + " (rpm)"
                }
        }
        if runtime.GOOS == "darwin" {
                if _, err := exec.LookPath("brew"); err == nil {
                        output := runCommand("sh", "-c", "brew list --formula | wc -l")
                        if output != "Unknown" {
                                return output + " (brew)"
                        }
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
        if runtime.GOOS == "darwin" {
                return "Apple"
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
        if runtime.GOOS == "darwin" {
                return runCommand("sysctl", "-n", "hw.model")
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
        if runtime.GOOS == "darwin" {
                output := runCommand("pmset", "-g", "batt")
                scanner := bufio.NewScanner(strings.NewReader(output))
                for scanner.Scan() {
                        line := scanner.Text()
                        if strings.Contains(line, "InternalBattery") {
                                parts := strings.Split(line, "\t")
                                if len(parts) > 1 {
                                        statusParts := strings.Split(parts[1], ";")
                                        if len(statusParts) > 2 {
                                                charge := statusParts[0]
                                                state := statusParts[1]
                                                return fmt.Sprintf("%s [%s]", strings.TrimSpace(charge), strings.TrimSpace(state))
                                        }
                                }
                        }
                }
        }
        return "Unknown"
}

func getDesktopEnvironment() string {
        if de := os.Getenv("XDG_CURRENT_DESKTOP"); de != "" {
                return de
        }
        if de := os.Getenv("DESKTOP_SESSION"); de != "" {
                return de
        }
        return "Unknown"
}

func getTerminal() string {
        if term := os.Getenv("TERM_PROGRAM"); term != "" {
                return term
        }
        if term := os.Getenv("TERM"); term != "" {
                return term
        }
        return "Unknown"
}

func getCPUUsage() string {
        percentages, err := cpu.Percent(time.Second, false)
        if err != nil || len(percentages) == 0 {
                return "N/A"
        }
        return fmt.Sprintf("%.2f%%", percentages[0])
}

func getGtkInfo(key string) string {
        if _, err := exec.LookPath("gsettings"); err == nil {
                schema := "org.gnome.desktop.interface"
                output := runCommand("gsettings", "get", schema, key)
                if output != "Unknown" {
                        return strings.Trim(output, "'")
                }
        }
        return "Unknown"
}

func getTheme() string {
        return getGtkInfo("gtk-theme")
}

func getIcons() string {
        return getGtkInfo("icon-theme")
}

func getLocale() string {
        locale := os.Getenv("LANG")
        if locale != "" {
                return strings.Split(locale, ".")[0]
        }
        return "Unknown"
}