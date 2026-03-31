import { useMsal } from '@azure/msal-react'
import { apiRequest } from './msalConfig'

export function useApiToken() {
  const { instance, accounts } = useMsal()

  return async (): Promise<string> => {
    const response = await instance.acquireTokenSilent({
      ...apiRequest,
      account: accounts[0],
    })
    return response.accessToken
  }
}
