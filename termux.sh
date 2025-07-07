#!/data/data/com.termux/files/usr/bin/sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

check_dep() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf "${RED}HATA: Kurulum için '%s' komutu gerekli ama sistemde bulunamadı.${NC}\n" "$1"
        printf "${YELLOW}Lütfen 'pkg install %s' komutu ile yükleyip tekrar deneyin.${NC}\n" "$1"
        exit 1
    fi
}

main() {
    printf "${GREEN}pixfetch Termux kurulum betiği başlatılıyor...${NC}\n\n"

    printf "${YELLOW}Gerekli komutlar kontrol ediliyor...${NC}\n"
    check_dep "go"
    check_dep "git"
    printf "${GREEN}Temel gereksinimler mevcut.${NC}\n\n"

    if ! command -v "termux-battery-status" >/dev/null 2>&1; then
        printf "${YELLOW}UYARI: 'termux-battery-status' komutu bulunamadı.${NC}\n"
        printf "Pil bilgilerinin gösterilmesi için Termux:API uygulamasının kurulu olması gerekir.\n"
        printf "Yüklemek için: 'pkg install termux-api'\n\n"
    fi

    printf "${YELLOW}Go modülleri indiriliyor ve hazırlanıyor...${NC}\n"
    if [ ! -f "go.mod" ]; then
        go mod init pixfetch
    fi
    go mod tidy
    printf "${GREEN}Go modülleri başarıyla hazırlandı.${NC}\n\n"

    printf "${YELLOW}pixfetch derleniyor... (Bu işlem cihazınızın hızına göre biraz zaman alabilir)${NC}\n"
    go build -o pixfetch
    printf "${GREEN}Derleme tamamlandı.${NC}\n\n"

    printf "${YELLOW}glintfetch, Termux için kuruluyor...${NC}\n"
    if mv glintfetch "$PREFIX/bin/"; then
        printf "\n${GREEN}---- KURULUM BAŞARILI! ----${NC}\n"
        printf "Artık Termux'ta herhangi bir yerden 'pixfetch' komutunu çalıştırabilirsiniz.\n"
    else
        printf "\n${RED}---- KURULUM BAŞARISIZ! ----${NC}\n"
        printf "Dosya '$PREFIX/bin/' dizinine taşınamadı.\n"
        printf "Lütfen Termux izinlerinizi kontrol edin.\n"
        exit 1
    fi
}

main
