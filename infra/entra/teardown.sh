#!/usr/bin/env bash
# Removes the hello-world-go app registration and service principal.
# Use before re-running setup.sh from scratch.

set -euo pipefail

APP_NAME="hello-world-go"

echo "==> Looking up app registration: $APP_NAME"
APP_ID=$(az ad app list --display-name "$APP_NAME" --query "[0].appId" -o tsv)

if [ -z "$APP_ID" ]; then
  echo "    Not found, nothing to do."
  exit 0
fi

echo "    App ID: $APP_ID"

echo "==> Deleting app registration (also removes service principal)"
az ad app delete --id "$APP_ID"

echo "==> Done."
