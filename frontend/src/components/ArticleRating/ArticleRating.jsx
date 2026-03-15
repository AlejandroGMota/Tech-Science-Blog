import { useState, useEffect } from 'react'
import { apiFetch } from '../../hooks/useApi'
import './ArticleRating.css'

export default function ArticleRating({ slug }) {
  const [rating, setRating] = useState({ average: 0, count: 0 })
  const [hover, setHover] = useState(0)
  const [voted, setVoted] = useState(false)

  useEffect(() => {
    apiFetch(`/articles/${slug}/rating`).then(setRating).catch(() => {})
    setVoted(localStorage.getItem(`rated-${slug}`) === 'true')
  }, [slug])

  const handleRate = async (score) => {
    if (voted) return
    try {
      const result = await apiFetch(`/articles/${slug}/rating`, {
        method: 'POST',
        body: JSON.stringify({ score }),
      })
      setRating(result)
      setVoted(true)
      localStorage.setItem(`rated-${slug}`, 'true')
    } catch {
      // Already rated or error
    }
  }

  return (
    <div className="article-rating">
      <div className="stars">
        {[1, 2, 3, 4, 5].map((star) => (
          <button
            key={star}
            className={`star ${star <= (hover || Math.round(rating.average)) ? 'active' : ''} ${voted ? 'voted' : ''}`}
            onClick={() => handleRate(star)}
            onMouseEnter={() => !voted && setHover(star)}
            onMouseLeave={() => setHover(0)}
            disabled={voted}
            aria-label={`Rate ${star} stars`}
          >
            &#9733;
          </button>
        ))}
      </div>
      <span className="rating-info">
        {rating.average.toFixed(1)} ({rating.count} {rating.count === 1 ? 'voto' : 'votos'})
      </span>
    </div>
  )
}
