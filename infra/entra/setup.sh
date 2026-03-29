#!/usr/bin/env bash
# Entra ID app registration setup for hello-world-go.
# Creates the app registration, app roles, service principal,
# and assigns both roles to the current logged-in user.
#
# Prerequisites:
#   az login (already done)
#
# To re-run cleanly:
#   ./infra/entra/teardown.sh && ./infra/entra/setup.sh
#
# To reproduce in Terraform later, see infra/entra/terraform/ (TODO)

set -euo pipefail

TENANT_ID="a38e45e7-9d8c-49c2-b524-4f1ece71c53f"
APP_NAME="hello-world-go"

# UUIDs for app roles — fixed so the script is idempotent
ROLE_ITEMS_READ_ID="30bdff6c-ec7d-47cb-9d20-11689d8fcddc"
ROLE_ITEMS_WRITE_ID="637bb7fb-c4e0-4d7f-8b2d-48649b74f720"

echo "==> Creating app registration: $APP_NAME"
APP_ID=$(az ad app create \
  --display-name "$APP_NAME" \
  --sign-in-audience "AzureADMyOrg" \
  --app-roles "[
    {
      \"allowedMemberTypes\": [\"User\"],
      \"displayName\": \"Items Reader\",
      \"description\": \"Can read items\",
      \"value\": \"items.read\",
      \"id\": \"$ROLE_ITEMS_READ_ID\",
      \"isEnabled\": true
    },
    {
      \"allowedMemberTypes\": [\"User\"],
      \"displayName\": \"Items Writer\",
      \"description\": \"Can create and update items\",
      \"value\": \"items.write\",
      \"id\": \"$ROLE_ITEMS_WRITE_ID\",
      \"isEnabled\": true
    }
  ]" \
  --query appId -o tsv)

echo "    App ID (client ID): $APP_ID"

echo "==> Creating service principal"
SP_OBJECT_ID=$(az ad sp create --id "$APP_ID" --query id -o tsv)
echo "    Service principal object ID: $SP_OBJECT_ID"

echo "==> Getting current user object ID"
USER_OBJECT_ID=$(az ad signed-in-user show --query id -o tsv)
echo "    User object ID: $USER_OBJECT_ID"

echo "==> Assigning items.read role"
az rest --method POST \
  --uri "https://graph.microsoft.com/v1.0/users/$USER_OBJECT_ID/appRoleAssignments" \
  --body "{
    \"principalId\": \"$USER_OBJECT_ID\",
    \"resourceId\": \"$SP_OBJECT_ID\",
    \"appRoleId\": \"$ROLE_ITEMS_READ_ID\"
  }" --output none

echo "==> Assigning items.write role"
az rest --method POST \
  --uri "https://graph.microsoft.com/v1.0/users/$USER_OBJECT_ID/appRoleAssignments" \
  --body "{
    \"principalId\": \"$USER_OBJECT_ID\",
    \"resourceId\": \"$SP_OBJECT_ID\",
    \"appRoleId\": \"$ROLE_ITEMS_WRITE_ID\"
  }" --output none

echo ""
echo "==> Done. Add the following to your .env:"
echo ""
echo "    AZURE_TENANT_ID=$TENANT_ID"
echo "    AZURE_CLIENT_ID=$APP_ID"
echo ""
echo "==> To get a test token:"
echo "    az account get-access-token --resource api://$APP_ID --query accessToken -o tsv"
