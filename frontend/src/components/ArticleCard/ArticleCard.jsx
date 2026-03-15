import { Link } from 'react-router-dom'
import './ArticleCard.css'

export default function ArticleCard({ article }) {
  return (
    <article className="article-card">
      <div className="article-card-body">
        <span className="article-card-category">{article.category}</span>
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
