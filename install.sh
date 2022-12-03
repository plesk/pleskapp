#!/bin/bash
# Copyright 1999-2022. Plesk International GmbH.

PREFIX=/usr/local/bin/

BIN_URLS=$(curl -fsSL https://api.github.com/repos/plesk/pleskapp/releases/latest | grep browser_download_url | cut -d '"' -f 4)
LINUX_ARCHIVE=$(echo "$BIN_URLS" | grep linux.tgz)
MAC_ARCHIVE=$(echo "$BIN_URLS" | grep mac.tgz)
OS_NAME=$(uname -s)

if [ "Linux" = "$OS_NAME" -o "Darwin" = "$OS_NAME" ]; then
  if [ "Linux" = "$OS_NAME" ]; then
    if [ $EUID -ne 0 ]; then
      echo "This script must be run as root user"
      exit 1
    fi
    curl -fsSL "$LINUX_ARCHIVE" --output plesk-latest.tgz
  fi

  [ "Darwin" = "$OS_NAME" ] && curl -fsSL "$MAC_ARCHIVE" --output plesk-latest.tgz

  tar xzf plesk-latest.tgz --directory=$PREFIX
  rm plesk-latest.tgz
  [ -f $PREFIX/plesk ] && echo "The utility 'plesk' has been successfully installed to $PREFIX"
else
  echo "Unsupported OS."
  exit 1
fi
