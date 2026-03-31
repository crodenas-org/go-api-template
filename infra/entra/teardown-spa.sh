#!/usr/bin/env bash
# Removes the hello-world-go-spa app registration and service principal.
# Use before re-running setup-spa.sh from scratch.

set -euo pipefail

SPA_APP_NAME="hello-world-go-spa"

echo "==> Looking up SPA app registration: $SPA_APP_NAME"
SPA_CLIENT_ID=$(az ad app list --display-name "$SPA_APP_NAME" --query "[0].appId" -o tsv)

if [ -z "$SPA_CLIENT_ID" ]; then
  echo "    Not found, nothing to do."
  exit 0
fi

echo "    SPA client ID: $SPA_CLIENT_ID"

echo "==> Deleting SPA app registration (also removes service principal)"
az ad app delete --id "$SPA_CLIENT_ID"

echo "==> Done."
