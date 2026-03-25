import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { apiFetch } from '../../hooks/useApi'
import './HomePage.css'

export default function HomePage() {
  const [recent, setRecent] = useState(null)

  useEffect(() => {
    document.title = 'Tech & Science Blog'
    apiFetch('/articles')
      .then((data) => setRecent(data.slice(0, 5)))
      .catch(() => setRecent([]))
  }, [])

  return (
    <div className="home">
      {/* ── 1. Intro ── */}
      <section className="home-intro">
        <h1 className="home-greeting">Estas en:</h1>
        <p className="home-tagline">
          Un espacio para documentar lo que aprendo construyendo software,
          escalando un negocio y tratando de entender el mundo.
        </p>
        <p className="home-rule">
          La regla: publica lo que te hubiera servido a ti hace 6 meses.
        </p>
      </section>

      <hr className="home-divider" />

      {/* ── 2. Recientes ── */}
      <section className="home-recent">
        <h2 className="home-section-header">Recientes</h2>
        <ul className="home-article-list">
          {recent === null ? (
            Array.from({ length: 5 }, (_, i) => (
              <li key={i} className="home-skeleton-item">
                <div className="home-skeleton-line home-skeleton-title" />
                <div className="home-skeleton-line home-skeleton-excerpt" />
              </li>
            ))
          ) : (
            recent.map((a) => (
              <li key={a.id}>
                <Link to={`/entradas/${a.slug}`} className="home-article-item">
                  <div className="home-article-info">
                    <span className="home-article-title">{a.title}</span>
                    <span className="home-article-excerpt">{a.excerpt}</span>
                  </div>
                  <div className="home-article-meta">
                    <span className="home-article-tag">{a.category}</span>
                    <time>{new Date(a.published_at).toLocaleDateString('es-MX', { day: 'numeric', month: 'short', year: 'numeric' })}</time>
                  </div>
                </Link>
              </li>
            ))
          )}
        </ul>
        <Link to="/entradas" className="home-all-link">Ver todos los artículos →</Link>
      </section>

      <hr className="home-divider" />

      {/* ── 3. Explora ── */}
      <section className="home-explore">
        <div className="home-categories">
          <h2 className="home-section-header">Categorías</h2>
          <div className="home-cat-grid">
            {[
              { name: 'Code', desc: 'Go, Node.js, arquitectura, DevOps' },
              { name: 'Business', desc: 'Emprendimiento, importación, operaciones' },
              { name: 'Ideas', desc: 'Geopolítica, ciencia, opinión' },
              { name: 'Stack de vida', desc: 'Productividad, herramientas' },
              { name: 'Aprendiendo en público', desc: 'Errores y descubrimientos' },
            ].map((cat) => (
              <Link key={cat.name} to="/entradas" className="home-cat-chip">
                <strong>{cat.name}</strong>
                <span>{cat.desc}</span>
              </Link>
            ))}
          </div>
        </div>

        <div className="home-formats">
          <h2 className="home-section-header">Formatos</h2>
          <div className="home-formats-row">
            <div className="home-format">
              <strong>Deep Dive</strong>
              <span><em>1000-1500 palabras.</em> Investigación a fondo.</span>
            </div>
            <div className="home-format">
              <strong>Nota Rápida</strong>
              <span><em>300-500 palabras.</em> Una lección, sin perfeccionismo.</span>
            </div>
            <div className="home-format">
              <strong>TIL</strong>
              <span><em>Párrafo corto.</em> Lo que aprendí hoy.</span>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
