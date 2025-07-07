#!/bin/sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

check_dep() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf "${RED}ERROR: The command '%s' is required for installation but was not found on your system.${NC}\n" "$1"
        printf "${YELLOW}Please install '%s' and try again.${NC}\n" "$1"
        exit 1
    fi
}

main() {
    printf "${GREEN}Starting omnifetch installation script...${NC}\n\n"

    # Store the current directory of the script
    SCRIPT_DIR="$(dirname "$(realpath "$0")")"
    # Determine the project root (one level up from scripts/)
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

    printf "${YELLOW}Checking for required commands...${NC}\n"
    check_dep "go"
    check_dep "git"
    printf "${GREEN}All requirements are met.${NC}\n\n"

    printf "${YELLOW}Navigating to project root: %s${NC}\n" "$PROJECT_ROOT"
    # Change to the project root directory
    cd "$PROJECT_ROOT"

    printf "${YELLOW}Downloading and preparing Go modules...${NC}\n"
    if [ ! -f "go.mod" ]; then
        go mod init omnifetch
    fi
    go mod tidy
    printf "${GREEN}Go modules prepared successfully.${NC}\n\n"

    printf "${YELLOW}Compiling omnifetch...${NC}\n"
    # Build from the src directory
    go build -o omnifetch ./src
    printf "${GREEN}Compilation complete.${NC}\n\n"

    printf "${YELLOW}omnifetch will be installed to your system. You may be asked for your administrator (sudo) password.${NC}\n"
    if sudo mv omnifetch /usr/local/bin/; then
        printf "\n${GREEN}---- INSTALLATION SUCCESSFUL! ----${NC}\n"
        printf "You can now run the 'omnifetch' command from anywhere in your terminal.\n"
    else
        printf "\n${RED}---- INSTALLATION FAILED! ----${NC}\n"
        printf "Could not move the file to '/usr/local/bin/'.\n"
        printf "Please check permissions or manually run 'sudo mv omnifetch /usr/local/bin/'.\n"
        exit 1
    fi
}

main
