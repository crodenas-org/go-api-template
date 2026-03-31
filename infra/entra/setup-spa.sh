#!/usr/bin/env bash
# Entra ID app registration for the hello-world-go React SPA.
# Creates a public client SPA registration and grants it permission
# to call the hello-world-go API.
#
# Prerequisites:
#   az login (already done)
#   setup.sh must have been run first (API app reg must exist)
#
# To re-run cleanly:
#   ./infra/entra/teardown-spa.sh && ./infra/entra/setup-spa.sh

set -euo pipefail

TENANT_ID="a38e45e7-9d8c-49c2-b524-4f1ece71c53f"
SPA_APP_NAME="hello-world-go-spa"
API_APP_NAME="hello-world-go"

echo "==> Looking up API app registration: $API_APP_NAME"
API_CLIENT_ID=$(az ad app list --display-name "$API_APP_NAME" --query "[0].appId" -o tsv)

if [ -z "$API_CLIENT_ID" ]; then
  echo "    ERROR: API app registration not found. Run setup.sh first."
  exit 1
fi

echo "    API client ID: $API_CLIENT_ID"

echo "==> Looking up API service principal"
API_SP_OBJECT_ID=$(az ad sp show --id "$API_CLIENT_ID" --query id -o tsv)

echo "==> Looking up access_as_user scope ID from API app"
SCOPE_ID=$(az ad app show --id "$API_CLIENT_ID" \
  --query "api.oauth2PermissionScopes[?value=='access_as_user'].id" -o tsv)

if [ -z "$SCOPE_ID" ]; then
  echo "    ERROR: access_as_user scope not found on API app. Re-run setup.sh."
  exit 1
fi

echo "    Scope ID: $SCOPE_ID"

echo "==> Creating SPA app registration: $SPA_APP_NAME"
SPA_CLIENT_ID=$(az ad app create \
  --display-name "$SPA_APP_NAME" \
  --sign-in-audience "AzureADMyOrg" \
  --query appId -o tsv)

echo "    SPA client ID: $SPA_CLIENT_ID"

echo "==> Configuring SPA platform (PKCE, no implicit grant)"
az rest --method PATCH \
  --uri "https://graph.microsoft.com/v1.0/applications(appId='$SPA_CLIENT_ID')" \
  --body '{
    "spa": {
      "redirectUris": [
        "http://localhost:5173"
      ]
    },
    "isFallbackPublicClient": true
  }'

echo "==> Granting API permission: api://$API_CLIENT_ID/access_as_user"
az rest --method PATCH \
  --uri "https://graph.microsoft.com/v1.0/applications(appId='$SPA_CLIENT_ID')" \
  --body "{
    \"requiredResourceAccess\": [{
      \"resourceAppId\": \"$API_CLIENT_ID\",
      \"resourceAccess\": [{
        \"id\": \"$SCOPE_ID\",
        \"type\": \"Scope\"
      }]
    }]
  }"

echo "==> Creating service principal for SPA"
SPA_SP_OBJECT_ID=$(az ad sp create --id "$SPA_CLIENT_ID" --query id -o tsv)
echo "    SPA service principal object ID: $SPA_SP_OBJECT_ID"

echo "==> Granting admin consent for SPA to access the API"
az rest --method POST \
  --uri "https://graph.microsoft.com/v1.0/oauth2PermissionGrants" \
  --body "{
    \"clientId\": \"$SPA_SP_OBJECT_ID\",
    \"consentType\": \"AllPrincipals\",
    \"resourceId\": \"$API_SP_OBJECT_ID\",
    \"scope\": \"access_as_user\"
  }"

echo "==> Adding current user as app owner"
USER_OBJECT_ID=$(az ad signed-in-user show --query id -o tsv)
az ad app owner add --id "$SPA_CLIENT_ID" --owner-object-id "$USER_OBJECT_ID"

echo ""
echo "==> Done. Add the following to web/.env.local:"
echo ""
echo "    VITE_AZURE_TENANT_ID=$TENANT_ID"
echo "    VITE_AZURE_CLIENT_ID=$SPA_CLIENT_ID"
echo "    VITE_AZURE_API_CLIENT_ID=$API_CLIENT_ID"
echo "    VITE_API_BASE_URL=http://localhost:8080"
echo ""
echo "==> Users still need the items.read / items.write roles assigned via setup.sh."
echo "    The SPA acquires tokens on behalf of the user — roles come from the API app reg."
