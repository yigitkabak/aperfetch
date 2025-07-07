---
# pixfetch: Sistem Bilgisi Aracı

pixfetch, sisteminizin temel bilgilerini, hoş ve sade bir **ASCII logosu** eşliğinde gösteren basit ve hızlı bir komut satırı aracıdır. Windows, macOS ve Linux (Termux dahil Android) sistemlerinde çalışacak şekilde tasarlanmıştır.

## Özellikler

* **Çoklu Platform Desteği:** Windows, macOS, Linux (Termux dahil Android) üzerinde çalışır.
* **Minimalist Tasarım:** Tek bir şık ASCII logosu ile temiz ve düzenli bir çıktı sunar.
* **Kapsamlı Sistem Bilgileri:** Aşağıdaki bilgileri görüntüler:
    * Kullanıcı ve Host adı
    * İşletim Sistemi bilgisi
    * Kernel sürümü
    * Çalışma süresi (Uptime)
    * Kullanılan Shell
    * GPU bilgisi
    * Bellek kullanımı
    * Disk kullanımı
    * Yerel IP adresi
    * Ekran çözünürlüğü
    * Yüklü paket sayısı (dpkg, pacman, rpm, brew destekli)
    * Cihaz üreticisi ve modeli
    * Pil durumu (Android ve Linux için)
* **Renkli Çıktı:** Bilgileri daha okunaklı hale getirmek için ANSI renk kodlarını kullanır.

## Nasıl Çalışır?

pixfetch, sistem bilgilerini toplamak için çeşitli komut satırı araçlarından (örneğin `uname`, `lspci`, `getprop`, `wmic`, `sysctl`, `xrandr`, `dpkg`, vb.) ve Go'nun `gopsutil` kütüphanesinden faydalanır. Topladığı bu bilgileri, zarif bir ASCII logosuyla birlikte hizalı bir şekilde terminalinize yazdırır.

NOT: `pixfetch ascii_distro` komudunu kullanarak istediğiniz bir distronun ascii yazısını kullanabilirsiniz!
örnek:
```
❯ pixfetch --ascii_distro macos

                      c.'          yigitkabak@yigitkabak
                   ,xNMM.          -----------------
                 .OMMMMo           OS: CachyOS
                 lMM"              Kernel: 6.15.5-2-cachyos
       .;loddo:.  .olloddol;.      Uptime: 0 days, 0 hours, 46 mins
     cKMMMMMMMMMMNWMMMMMMMMMM0:    Packages: 1128 (pacman)
   .KMMMMMMMMMMMMMMMMMMMMMMMWd.    DE: KDE
   XMMMMMMMMMMMMMMMMMMMMMMMX.      Shell: fish
  ;MMMMMMMMMMMMMMMMMMMMMMMM:       Terminal: xterm-256color
  :MMMMMMMMMMMMMMMMMMMMMMMM:       Resolution: 1366x768
  .MMMMMMMMMMMMMMMMMMMMMMMMX.      Theme: Breeze
   kMMMMMMMMMMMMMMMMMMMMMMMMWd.    Icons: breeze-dark
   'XMMMMMMMMMMMMMMMMMMMMMMMMk     CPU: Intel(R) Core(TM) i7-5500U CPU @ 2.40GHz (4) [0.75%]
    'XMMMMMMMMMMMMMMMMMMMMMMMMK.   GPU: Intel Corporation HD Graphics 5500 (rev 09)
      kMMMMMMMMMMMMMMMMMMMMMMd     Memory: 1.71 GiB / 7.64 GiB
       ;KMMMMMMMWXXWMMMMMMMk.      Disk: 9.83 GiB / 931.22 GiB
         "cooc*"    "*coo'"        Battery: 100% [Full]
                                   Locale: tr_TR
                                   Local IP (wlan0): *********

```

## Desteklenen Sistemler
* Ubuntu
* Arch
* NixOS
* Android
* Debian
* Fedora
* Pop!_OS
* CachyOS
* ZorinOS
* OpenSUSE
* Elemantary OS
* Endeavour OS
* Manjaro
* Mint
* Gentoo
* Pardus
* Windows
* MacOS

## Kurulum

İlk olarak aşağıdaki komutla aperfetch deposunu klonlayın:

```bash
aperium clone -github yigitkabak/pixfetch
```

Ardından son olarak şunları yapın:

```bash
cd pixfetch
bash linux.sh (termux ise termux.sh yazın.)
```
