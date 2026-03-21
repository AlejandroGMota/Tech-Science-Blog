import { Link } from 'react-router-dom'
import './ArticleCard.css'

export default function ArticleCard({ article }) {
  const tags = article.tags ? article.tags.split(',').map(t => t.trim()).filter(Boolean) : []

  return (
    <article className="article-card">
      <div className="article-card-body">
        {tags.length > 0 && (
          <div className="article-card-tags">
            {tags.map((tag) => (
              <span key={tag} className="article-card-tag">{tag}</span>
            ))}
          </div>
        )}
        <h3>
          <Link to={`/entradas/${article.slug}`}>{article.title}</Link>
        </h3>
        <p>{article.excerpt}</p>
        <div className="article-card-meta">
          <time>{new Date(article.published_at).toLocaleDateString('es-MX')}</time>
          <span>{article.author}</span>
        </div>
      </div>
    </article>
  )
}
