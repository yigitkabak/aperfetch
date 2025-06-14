---
# aperfetch: Minimalist Sistem Bilgisi Aracı

aperfetch, sisteminizin temel bilgilerini, hoş ve sade bir **ASCII logosu** eşliğinde gösteren basit ve hızlı bir komut satırı aracıdır. Windows, macOS ve Linux (Termux dahil Android) sistemlerinde çalışacak şekilde tasarlanmıştır.

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

aperfetch, sistem bilgilerini toplamak için çeşitli komut satırı araçlarından (örneğin `uname`, `lspci`, `getprop`, `wmic`, `sysctl`, `xrandr`, `dpkg`, vb.) ve Rust'ın `systemstat` kütüphanesinden faydalanır. Topladığı bu bilgileri, zarif bir ASCII logosuyla birlikte hizalı bir şekilde terminalinize yazdırır.

## Kurulum
**ilk olarak:**
```
aperium clone -github yigitkabak/aperfetch
```
**son olarak şunu yapmanız yeterlidir:**
```
bash install.sh
```
Sistemi başlatmak için terminalinize `aperfetch` yazmanız yeterlidir.
