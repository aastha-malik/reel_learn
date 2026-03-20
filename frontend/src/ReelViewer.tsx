import { useRef, useState, useEffect, useCallback } from 'react'
import './ReelViewer.css'

interface Reel {
  path: string
  start_time: number
  end_time: number
}

interface Props {
  reels: Reel[]
  onBack: () => void
}

export default function ReelViewer({ reels, onBack }: Props) {
  const [index, setIndex] = useState(0)
  const [paused, setPaused] = useState(false)
  const [transitioning, setTransitioning] = useState(false)
  const [buffering, setBuffering] = useState(true)
  const [direction, setDirection] = useState<'next' | 'prev'>('next')

  const videoRef = useRef<HTMLVideoElement>(null)
  const touchStartY = useRef<number | null>(null)

  // tracks seconds watched per reel — ready for future use
  const watchedSeconds = useRef<number[]>(new Array(reels.length).fill(0))

  const videoUrl = (reel: Reel) =>
    `http://localhost:8080/video?path=${encodeURIComponent(reel.path)}`

  const saveProgress = () => {
    const video = videoRef.current
    if (!video) return
    watchedSeconds.current[index] = Math.floor(video.currentTime)
  }

  const navigate = useCallback((dir: 'next' | 'prev') => {
    if (transitioning) return
    if (dir === 'next' && index >= reels.length - 1) return
    if (dir === 'prev' && index <= 0) return

    saveProgress()
    setDirection(dir)
    setTransitioning(true)

    setTimeout(() => {
      setIndex(i => (dir === 'next' ? i + 1 : i - 1))
      setPaused(false)
      setBuffering(true)
      setTransitioning(false)
    }, 250)
  }, [transitioning, index, reels.length])

  useEffect(() => {
    const video = videoRef.current
    if (!video) return
    video.load()
    video.play().catch(() => {})
  }, [index])

  const togglePlay = () => {
    const video = videoRef.current
    if (!video) return
    if (video.paused) {
      video.play()
      setPaused(false)
    } else {
      video.pause()
      setPaused(true)
    }
  }

  const handleWheel = useCallback((e: WheelEvent) => {
    e.preventDefault()
    if (e.deltaY > 30) navigate('next')
    else if (e.deltaY < -30) navigate('prev')
  }, [navigate])

  useEffect(() => {
    window.addEventListener('wheel', handleWheel, { passive: false })
    return () => window.removeEventListener('wheel', handleWheel)
  }, [handleWheel])

  const handleTouchStart = (e: React.TouchEvent) => {
    touchStartY.current = e.touches[0].clientY
  }

  const handleTouchEnd = (e: React.TouchEvent) => {
    if (touchStartY.current === null) return
    const diff = touchStartY.current - e.changedTouches[0].clientY
    if (diff > 60) navigate('next')
    else if (diff < -60) navigate('prev')
    touchStartY.current = null
  }

  const animClass = transitioning
    ? direction === 'next' ? 'slide-out-up' : 'slide-out-down'
    : direction === 'next' ? 'slide-in-up' : 'slide-in-down'

  const progress = ((index + 1) / reels.length) * 100
  const isFirst = index === 0
  const isLast = index === reels.length - 1

  return (
    <div
      className="reel-viewer"
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
    >
      <div className="progress-bar">
        <div className="progress-fill" style={{ width: `${progress}%` }} />
      </div>

      <div className="reel-top">
        <button className="back-btn" onClick={onBack}>✕</button>
        <span className="reel-counter">{index + 1} / {reels.length}</span>
      </div>

      <div className={`reel-wrap ${animClass}`} onClick={togglePlay}>
        <video
          ref={videoRef}
          className="reel-video"
          src={videoUrl(reels[index])}
          autoPlay
          playsInline
          onEnded={() => navigate('next')}
          onWaiting={() => setBuffering(true)}
          onPlaying={() => setBuffering(false)}
          onCanPlay={() => setBuffering(false)}
          onTimeUpdate={saveProgress}
        />

        {buffering && !paused && (
          <div className="overlay">
            <div className="spinner" />
          </div>
        )}

        {paused && !buffering && (
          <div className="overlay">
            <span className="play-icon">▶</span>
          </div>
        )}
      </div>

      <div className="reel-hint">
        {isFirst && isLast
          ? 'Only one reel'
          : isFirst
          ? 'Scroll down for next'
          : isLast
          ? 'Scroll up for previous · End of video'
          : 'Scroll up or down to navigate'}
      </div>
    </div>
  )
}
