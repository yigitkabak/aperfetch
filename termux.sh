#!/data/data/com.termux/files/usr/bin/sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

check_dep() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf "${RED}ERROR: The command '%s' is required for installation but was not found on your system.${NC}\n" "$1"
        printf "${YELLOW}Please install it with 'pkg install %s' and try again.${NC}\n" "$1"
        exit 1
    fi
}

main() {
    printf "${GREEN}Starting pixfetch Termux installation script...${NC}\n\n"

    printf "${YELLOW}Checking for required commands...${NC}\n"
    check_dep "go"
    check_dep "git"
    printf "${GREEN}Basic requirements are met.${NC}\n\n"

    if ! command -v "termux-battery-status" >/dev/null 2>&1; then
        printf "${YELLOW}WARNING: 'termux-battery-status' command not found.${NC}\n"
        printf "Termux:API app needs to be installed to display battery information.\n"
        printf "To install: 'pkg install termux-api'\n\n"
    fi

    printf "${YELLOW}Downloading and preparing Go modules...${NC}\n"
    if [ ! -f "go.mod" ]; then
        go mod init pixfetch
    fi
    go mod tidy
    printf "${GREEN}Go modules prepared successfully.${NC}\n\n"

    printf "${YELLOW}Compiling pixfetch... (This may take some time depending on your device's speed)${NC}\n"
    go build -o pixfetch
    printf "${GREEN}Compilation complete.${NC}\n\n"

    printf "${YELLOW}Installing pixfetch for Termux...${NC}\n"
    if mv pixfetch "$PREFIX/bin/"; then
        printf "\n${GREEN}---- INSTALLATION SUCCESSFUL! ----${NC}\n"
        printf "You can now run the 'pixfetch' command from anywhere in Termux.\n"
    else
        printf "\n${RED}---- INSTALLATION FAILED! ----${NC}\n"
        printf "Could not move the file to '$PREFIX/bin/'.\n"
        printf "Please check your Termux permissions.\n"
        exit 1
    fi
}

main
