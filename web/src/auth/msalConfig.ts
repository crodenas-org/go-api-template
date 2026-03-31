import type { Configuration, PopupRequest } from '@azure/msal-browser'

export const msalConfig: Configuration = {
  auth: {
    clientId: import.meta.env.VITE_AZURE_CLIENT_ID,
    authority: `https://login.microsoftonline.com/${import.meta.env.VITE_AZURE_TENANT_ID}`,
    redirectUri: window.location.origin,
  },
}

// Scopes requested when acquiring a token for the API
export const apiRequest: PopupRequest = {
  scopes: [`api://${import.meta.env.VITE_AZURE_API_CLIENT_ID}/access_as_user`],
}
