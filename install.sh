#!/usr/bin/env bash

Help() {
  cat <<USAGE

Usage: $0 [--out]

Options:
    -o, --out:            location to install (default ~/bin)
    -h, --help:           this help screen
USAGE
    exit 1
}

out="$HOME/bin"

while [ "$1" != "" ]; do
  case $1 in
    -o | --out) shift
                out="$1"
                ;;
    -h | --help) Help
                 ;;
    *) Help
       ;;
  esac
  shift
done

os=$(uname)
arch=$(uname -m)
os_arch=$(echo "$os"'_'"$arch")

curl -sSL https://api.github.com/repos/eiladin/k8s-dotenv/releases/latest \
  | grep -E "browser_download_url" \
  | grep "$os_arch" \
  | cut -d '"' -f 4 \
  | wget -q -O /tmp/k8s-dotenv.tar.gz  -i -

tar -zxf /tmp/k8s-dotenv.tar.gz -C "$out"

rm /tmp/k8s-dotenv.tar.gz
