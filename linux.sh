#!/bin/sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

check_dep() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf "${RED}HATA: Kurulum için '%s' komutu gerekli ama sistemde bulunamadı.${NC}\n" "$1"
        printf "${YELLOW}Lütfen '%s' yükleyip tekrar deneyin.${NC}\n" "$1"
        exit 1
    fi
}

main() {
    printf "${GREEN}glintfetch kurulum betiği başlatılıyor...${NC}\n\n"

    printf "${YELLOW}Gerekli komutlar kontrol ediliyor...${NC}\n"
    check_dep "go"
    check_dep "git"
    printf "${GREEN}Tüm gereksinimler mevcut.${NC}\n\n"

    printf "${YELLOW}Go modülleri indiriliyor ve hazırlanıyor...${NC}\n"
    if [ ! -f "go.mod" ]; then
        go mod init glintfetch
    fi
    go mod tidy
    printf "${GREEN}Go modülleri başarıyla hazırlandı.${NC}\n\n"

    printf "${YELLOW}glintfetch derleniyor...${NC}\n"
    go build -o glintfetch
    printf "${GREEN}Derleme tamamlandı.${NC}\n\n"

    printf "${YELLOW}glintfetch sisteme kurulacak. Yönetici (sudo) parolası istenebilir.${NC}\n"
    if sudo mv glintfetch /usr/local/bin/; then
        printf "\n${GREEN}---- KURULUM BAŞARILI! ----${NC}\n"
        printf "Artık terminalde herhangi bir yerden 'glintfetch' komutunu çalıştırabilirsiniz.\n"
    else
        printf "\n${RED}---- KURULUM BAŞARISIZ! ----${NC}\n"
        printf "Dosya '/usr/local/bin/' dizinine taşınamadı.\n"
        printf "İzinleri kontrol edin veya manuel olarak 'sudo mv glintfetch /usr/local/bin/' komutunu çalıştırın.\n"
        exit 1
    fi
}

main
