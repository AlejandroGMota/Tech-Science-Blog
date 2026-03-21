import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { apiFetch } from '../../hooks/useApi'
import ArticleRating from '../../components/ArticleRating/ArticleRating'
import './ArticlePage.css'

export default function ArticlePage() {
  const { slug } = useParams()
  const [article, setArticle] = useState(null)
  const [error, setError] = useState(null)

  useEffect(() => {
    apiFetch(`/articles/${slug}`)
      .then((data) => {
        setArticle(data)
        document.title = `${data.title} — Alejandro G. Mota`
      })
      .catch((err) => setError(err.message))
  }, [slug])

  if (error) {
    return (
      <div className="article-page">
        <div className="container">
          <p className="error">{error}</p>
          <Link to="/entradas">Volver a artículos</Link>
        </div>
      </div>
    )
  }

  if (!article) {
    return (
      <div className="article-page">
        <div className="container">
          <p className="loading">Cargando...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="article-page">
      <div className="container">
        <Link to="/entradas" className="back-link">← Volver a artículos</Link>

        <article className="article-content">
          <header>
            <span className="article-category">{article.category}</span>
            <h1>{article.title}</h1>
            <div className="article-meta">
              <time>{new Date(article.published_at).toLocaleDateString('es-MX', { year: 'numeric', month: 'long', day: 'numeric' })}</time>
              <span>Por {article.author}</span>
            </div>
          </header>

          <ArticleRating slug={slug} />

          <div
            className="article-body"
            dangerouslySetInnerHTML={{ __html: article.content }}
          />

          <ArticleRating slug={slug} />
        </article>
      </div>
    </div>
  )
}
