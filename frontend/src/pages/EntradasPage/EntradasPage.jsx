import { useState, useEffect } from 'react'
import { apiFetch } from '../../hooks/useApi'
import ArticleCard from '../../components/ArticleCard/ArticleCard'
import './EntradasPage.css'

export default function EntradasPage() {
  const [articles, setArticles] = useState([])
  const [category, setCategory] = useState('')
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const params = new URLSearchParams()
    if (category) params.set('category', category)
    if (search) params.set('search', search)
    const query = params.toString()

    apiFetch(`/articles${query ? `?${query}` : ''}`)
      .then(setArticles)
      .catch(() => setArticles([]))
      .finally(() => setLoading(false))
  }, [category, search])

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
          <select
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            className="category-select"
          >
            <option value="">Todas las categorías</option>
            <option value="IA">Inteligencia Artificial</option>
            <option value="Blockchain">Blockchain</option>
            <option value="Agronomía">Agronomía</option>
            <option value="Data Science">Data Science</option>
          </select>
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
