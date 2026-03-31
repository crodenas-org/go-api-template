import { useEffect, useState } from 'react'
import { useApiToken } from '../auth/useApiToken'

interface Item {
  id: number
  name: string
  created_at: string
}

const API_BASE = import.meta.env.VITE_API_BASE_URL

export default function Items() {
  const getToken = useApiToken()
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [newName, setNewName] = useState('')
  const [creating, setCreating] = useState(false)

  const fetchItems = async () => {
    setLoading(true)
    setError(null)
    try {
      const token = await getToken()
      const res = await fetch(`${API_BASE}/items`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
      setItems(await res.json())
    } catch (e) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newName.trim()) return
    setCreating(true)
    setError(null)
    try {
      const token = await getToken()
      const res = await fetch(`${API_BASE}/items`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: newName.trim() }),
      })
      if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
      setNewName('')
      await fetchItems()
    } catch (e) {
      setError(String(e))
    } finally {
      setCreating(false)
    }
  }

  useEffect(() => { fetchItems() }, [])

  return (
    <div>
      <h2>Items</h2>

      <form onSubmit={handleCreate}>
        <input
          type="text"
          placeholder="Item name"
          value={newName}
          onChange={e => setNewName(e.target.value)}
          disabled={creating}
        />
        <button type="submit" disabled={creating || !newName.trim()}>
          {creating ? 'Creating...' : 'Create'}
        </button>
      </form>

      {error && <p className="error">{error}</p>}

      {loading ? (
        <p>Loading...</p>
      ) : items.length === 0 ? (
        <p>No items yet.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Created</th>
            </tr>
          </thead>
          <tbody>
            {items.map(item => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.name}</td>
                <td>{new Date(item.created_at).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
