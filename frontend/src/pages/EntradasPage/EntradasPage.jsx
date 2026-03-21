import { useState, useEffect } from 'react'
import { apiFetch } from '../../hooks/useApi'
import ArticleCard from '../../components/ArticleCard/ArticleCard'
import './EntradasPage.css'

const CATEGORIES = [
  { value: 'Code', label: 'Code' },
  { value: 'Business', label: 'Business' },
  { value: 'Ideas', label: 'Ideas' },
  { value: 'Stack de vida', label: 'Stack de vida' },
  { value: 'Aprendiendo en público', label: 'Aprendiendo en público' },
]

const ARTICLE_TYPES = [
  { value: 'deep_dive', label: 'Deep Dive' },
  { value: 'nota_rapida', label: 'Nota Rápida' },
  { value: 'til', label: 'TIL' },
]

export default function EntradasPage() {
  const [articles, setArticles] = useState([])
  const [category, setCategory] = useState('')
  const [articleType, setArticleType] = useState('')
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const params = new URLSearchParams()
    if (category) params.set('category', category)
    if (search) params.set('search', search)
    if (articleType) params.set('article_type', articleType)
    const query = params.toString()

    apiFetch(`/articles${query ? `?${query}` : ''}`)
      .then(setArticles)
      .catch(() => setArticles([]))
      .finally(() => setLoading(false))
  }, [category, search, articleType])

  return (
    <div className="entradas-page">
      <div className="container">
        <h1>Artículos</h1>

        <div className="filters">
          <input
            type="text"
            placeholder="Buscar artículos..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="search-input"
          />
          <div className="type-filters">
            {ARTICLE_TYPES.map((t) => (
              <button
                key={t.value}
                className={`type-chip${articleType === t.value ? ' active' : ''}`}
                onClick={() => setArticleType(articleType === t.value ? '' : t.value)}
              >
                {t.label}
              </button>
            ))}
          </div>
        </div>

        <div className="category-filters">
          {CATEGORIES.map((cat) => (
            <button
              key={cat.value}
              className={`category-chip${category === cat.value ? ' active' : ''}`}
              onClick={() => setCategory(category === cat.value ? '' : cat.value)}
            >
              {cat.label}
            </button>
          ))}
        </div>

        {loading ? (
          <p className="loading">Cargando artículos...</p>
        ) : articles.length === 0 ? (
          <p className="no-results">No se encontraron artículos.</p>
        ) : (
          <div className="articles-grid">
            {articles.map((a) => (
              <ArticleCard key={a.id} article={a} />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
