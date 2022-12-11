#!/bin/bash
# Copyright 1999-2022. Plesk International GmbH.

PREFIX=/usr/local/bin/

BIN_URLS=$(curl -fsSL https://api.github.com/repos/plesk/pleskapp/releases/latest | grep browser_download_url | cut -d '"' -f 4)
OS_NAME=$(uname -s)
OS_ARCH=$(uname -m)

if [ "Linux" = "$OS_NAME" -o "Darwin" = "$OS_NAME" ]; then
  if [ "Linux" = "$OS_NAME" ]; then
    if [ $EUID -ne 0 ]; then
      echo "This script must be run as root user"
      exit 1
    fi
    LINUX_ARCHIVE=$(echo "$BIN_URLS" | grep linux-x86_64.tar.gz)
    curl -fsSL "$LINUX_ARCHIVE" --output plesk-latest.tar.gz
  fi

  if [ "Darwin" = "$OS_NAME" ]; then
    if [ "arm64" = "$OS_ARCH" ]; then
      MAC_ARCHIVE=$(echo "$BIN_URLS" | grep mac-arm64.tar.gz)
    else
      MAC_ARCHIVE=$(echo "$BIN_URLS" | grep mac-x86_64.tar.gz)
    fi
    curl -fsSL "$MAC_ARCHIVE" --output plesk-latest.tar.gz
  fi

  tar xzf plesk-latest.tar.gz --directory=$PREFIX plesk
  rm plesk-latest.tar.gz
  [ -f $PREFIX/plesk ] && echo "The utility 'plesk' has been successfully installed to $PREFIX"
else
  echo "Unsupported OS."
  exit 1
fi
