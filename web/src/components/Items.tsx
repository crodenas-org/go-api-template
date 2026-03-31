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
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editName, setEditName] = useState('')

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
        headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
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

  const handleEditSave = async (id: number) => {
    if (!editName.trim()) return
    setError(null)
    try {
      const token = await getToken()
      const res = await fetch(`${API_BASE}/items/${id}`, {
        method: 'PUT',
        headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: editName.trim() }),
      })
      if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
      setEditingId(null)
      await fetchItems()
    } catch (e) {
      setError(String(e))
    }
  }

  const handleDelete = async (id: number) => {
    setError(null)
    try {
      const token = await getToken()
      const res = await fetch(`${API_BASE}/items/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
      await fetchItems()
    } catch (e) {
      setError(String(e))
    }
  }

  const startEdit = (item: Item) => {
    setEditingId(item.id)
    setEditName(item.name)
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
              <th></th>
            </tr>
          </thead>
          <tbody>
            {items.map(item => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>
                  {editingId === item.id ? (
                    <input
                      type="text"
                      value={editName}
                      onChange={e => setEditName(e.target.value)}
                      onKeyDown={e => {
                        if (e.key === 'Enter') handleEditSave(item.id)
                        if (e.key === 'Escape') setEditingId(null)
                      }}
                      autoFocus
                    />
                  ) : (
                    item.name
                  )}
                </td>
                <td>{new Date(item.created_at).toLocaleString()}</td>
                <td className="actions">
                  {editingId === item.id ? (
                    <>
                      <button onClick={() => handleEditSave(item.id)} disabled={!editName.trim()}>Save</button>
                      <button onClick={() => setEditingId(null)}>Cancel</button>
                    </>
                  ) : (
                    <>
                      <button onClick={() => startEdit(item)}>Edit</button>
                      <button onClick={() => handleDelete(item.id)}>Delete</button>
                    </>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
