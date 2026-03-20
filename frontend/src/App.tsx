import { useState } from 'react'
import './App.css'
import ReelViewer from './ReelViewer'

interface Reel {
  path: string
  start_time: number
  end_time: number
}

interface Metadata {
  title: string
  channel: string
  duration: number
  thumbnail: string
}

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  if (h > 0) return `${h}h ${m}m ${s}s`
  if (m > 0) return `${m}m ${s}s`
  return `${s}s`
}

function App() {
  const [url, setUrl] = useState('')
  const [reels, setReels] = useState<Reel[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchingMeta, setFetchingMeta] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [viewing, setViewing] = useState(false)
  const [meta, setMeta] = useState<Metadata | null>(null)

  const fetchMetadata = async (inputUrl: string) => {
    if (!inputUrl.includes('youtube.com') && !inputUrl.includes('youtu.be')) return
    setFetchingMeta(true)
    setMeta(null)
    try {
      const res = await fetch(
        `http://localhost:8080/metadata?url=${encodeURIComponent(inputUrl)}`
      )
      if (res.ok) {
        const data = await res.json()
        setMeta(data)
      }
    } catch {
      // silently ignore — metadata is optional
    } finally {
      setFetchingMeta(false)
    }
  }

  const handleBlur = () => {
    if (url) fetchMetadata(url)
  }

  const handleSubmit = async (e: React.SyntheticEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      const res = await fetch('http://localhost:8080/process', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ youtube_url: url }),
      })

      const data = await res.json()

      if (!res.ok) {
        setError(data.error ?? 'Something went wrong')
      } else {
        setReels(data.reels)
        setViewing(true)
      }
    } catch {
      setError('Could not reach the backend. Is it running on :8080?')
    } finally {
      setLoading(false)
    }
  }

  if (viewing && reels.length > 0) {
    return <ReelViewer reels={reels} onBack={() => setViewing(false)} />
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
          onChange={(e) => { setUrl(e.target.value); setMeta(null) }}
          onBlur={handleBlur}
          required
        />
        <button type="submit" disabled={loading || fetchingMeta}>
          {loading ? 'Processing...' : 'Start'}
        </button>
      </form>

      {fetchingMeta && <p className="hint">Fetching video info...</p>}

      {meta && !fetchingMeta && (
        <div className="meta-card">
          {meta.thumbnail && (
            <img className="meta-thumb" src={meta.thumbnail} alt={meta.title} />
          )}
          <div className="meta-info">
            <p className="meta-title">{meta.title}</p>
            <p className="meta-sub">{meta.channel} · {formatDuration(meta.duration)}</p>
          </div>
        </div>
      )}

      {loading && (
        <p className="hint">Downloading and processing — this takes 1-2 minutes...</p>
      )}
      {error && <p className="error">{error}</p>}
    </div>
  )
}

export default App
