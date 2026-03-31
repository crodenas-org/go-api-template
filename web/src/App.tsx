import { useIsAuthenticated, useMsal } from '@azure/msal-react'
import { apiRequest } from './auth/msalConfig'
import Items from './components/Items'
import './App.css'

function App() {
  const { instance, accounts } = useMsal()
  const isAuthenticated = useIsAuthenticated()

  const handleLogin = () => {
    instance.loginRedirect(apiRequest)
  }

  const handleLogout = () => {
    instance.logoutRedirect()
  }

  return (
    <div className="app">
      <header>
        <h1>
          <a href="https://github.com/crodenas-org/go-api-template" target="_blank" rel="noreferrer">
            hello-world-go
          </a>
        </h1>
        <div>
          {isAuthenticated ? (
            <>
              <span>{accounts[0]?.username}</span>
              <button onClick={handleLogout}>Sign out</button>
            </>
          ) : (
            <button onClick={handleLogin}>Sign in</button>
          )}
        </div>
      </header>
      <main>
        {isAuthenticated ? <Items /> : <p>Sign in to continue.</p>}
      </main>
    </div>
  )
}

export default App
