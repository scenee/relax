#!/bin/bash -eu

BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BOLD='\033[1m'
NC='\033[0m'
ARROW="${BLUE}==>${NC}${BOLD}"
WARN="${RED}Warning${NC}${BOLD}"
ERROR="${RED}Error${NC}${BOLD}"

print() {
	echo -e "$1${NC}"
}

os_name=$(uname)
dest=~/.relax

print "$ARROW Check xcode-select"
xcode-select -p > /dev/null
if test 0 -ne $?; then
	print "==> Run: xcode-select --install"
	if xcode-select --install; then
		print "$ERROR Please install the command line developer tools as accoording to a user interface dialog"
		exit 1
	fi
fi

print "$ARROW Fetch relax..."
[[ -d "$dest" ]] && rm -fr "$dest"
git clone https://github.com/SCENEE/relax.git --depth 1 -b master "$dest"
rm -rf "$dest"/sample

print "${BOLD}Done!${NC}"
cat <<-EOM
	    Add Relax binary path to your bash profile:
	        export PATH=\$HOME/.relax/bin:\$PATH
	    To enable completion, add the following to your bash profile:
	        if which relax > /dev/null; then source "\$(relax init completion)"; fi
EOM

