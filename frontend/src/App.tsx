import { useState } from 'react'
import './App.css'

function App() {
  const [url, setUrl] = useState('')
  const [status, setStatus] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.SyntheticEvent) => {
    e.preventDefault()
    setStatus(null)
    setError(null)
    setLoading(true)

    try {
      const res = await fetch('http://localhost:8080/process', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      })

      const data = await res.json()

      if (!res.ok) {
        setError(data.error ?? 'Something went wrong')
      } else {
        setStatus(`Received: ${data.url}`)
      }
    } catch {
      setError('Could not reach the backend. Is it running on :8080?')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container">
      <h1>ReelLearn</h1>
      <p>Paste a YouTube URL to start watching in reel mode.</p>

      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="https://youtube.com/watch?v=..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Processing...' : 'Start'}
        </button>
      </form>

      {status && <p className="success">{status}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  )
}

export default App
