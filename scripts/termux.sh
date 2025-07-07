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
    printf "${GREEN}Starting omnifetch Termux installation script...${NC}\n\n"

    # Store the current directory of the script
    SCRIPT_DIR="$(dirname "$(realpath "$0")")"
    # Determine the project root (one level up from scripts/)
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

    printf "${YELLOW}Checking for required commands...${NC}\n"
    check_dep "go"
    check_dep "git"
    printf "${GREEN}Basic requirements are met.${NC}\n\n"

    if ! command -v "termux-battery-status" >/dev/null 2>&1; then
        printf "${YELLOW}WARNING: 'termux-battery-status' command not found.${NC}\n"
        printf "Termux:API app needs to be installed to display battery information.\n"
        printf "To install: 'pkg install termux-api'\n\n"
    fi

    printf "${YELLOW}Navigating to project root: %s${NC}\n" "$PROJECT_ROOT"
    # Change to the project root directory
    cd "$PROJECT_ROOT"

    printf "${YELLOW}Downloading and preparing Go modules...${NC}\n"
    # Check if go.mod exists in the project root, if not, initialize
    if [ ! -f "go.mod" ]; then
        go mod init omnifetch
    fi
    go mod tidy
    printf "${GREEN}Go modules prepared successfully.${NC}\n\n"

    printf "${YELLOW}Compiling omnifetch... (This may take some time depending on your device's speed)${NC}\n"
    # Build from the src directory
    go build -o omnifetch ./src
    printf "${GREEN}Compilation complete.${NC}\n\n"

    printf "${YELLOW}Installing omnifetch for Termux...${NC}\n"
    # The executable is now in the project root
    if mv omnifetch "$PREFIX/bin/"; then
        printf "\n${GREEN}---- INSTALLATION SUCCESSFUL! ----${NC}\n"
        printf "You can now run the 'omnifetch' command from anywhere in Termux.\n"
    else
        printf "\n${RED}---- INSTALLATION FAILED! ----${NC}\n"
        printf "Could not move the file to '$PREFIX/bin/'.\n"
        printf "Please check your Termux permissions.\n"
        exit 1
    fi

    # Navigate back to the script's original directory (optional, but good practice)
    # cd "$SCRIPT_DIR"
}

main
