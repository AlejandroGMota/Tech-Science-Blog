import { useState, useEffect } from 'react'
import { apiFetch } from '../../hooks/useApi'
import './AdminPage.css'

export default function AdminPage() {
  const [token, setToken] = useState(localStorage.getItem('token') || '')
  const [articles, setArticles] = useState([])
  const [editing, setEditing] = useState(null)
  const [form, setForm] = useState({ title: '', slug: '', content: '', excerpt: '', category: 'IA', author: 'Alejandro G. Mota' })
  const [loginForm, setLoginForm] = useState({ user: '', pass: '' })
  const [error, setError] = useState('')

  useEffect(() => {
    if (token) loadArticles()
  }, [token])

  const loadArticles = async () => {
    try {
      const data = await apiFetch('/articles')
      setArticles(data)
    } catch {
      setToken('')
      localStorage.removeItem('token')
    }
  }

  const handleLogin = async (e) => {
    e.preventDefault()
    setError('')
    try {
      const data = await apiFetch('/admin/login', {
        method: 'POST',
        body: JSON.stringify(loginForm),
      })
      localStorage.setItem('token', data.token)
      setToken(data.token)
    } catch (err) {
      setError(err.message)
    }
  }

  const handleLogout = async () => {
    await apiFetch('/admin/logout', { method: 'POST' }).catch(() => {})
    localStorage.removeItem('token')
    setToken('')
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    try {
      if (editing) {
        await apiFetch(`/articles/${editing}`, { method: 'PUT', body: JSON.stringify(form) })
      } else {
        await apiFetch('/articles', { method: 'POST', body: JSON.stringify(form) })
      }
      setForm({ title: '', slug: '', content: '', excerpt: '', category: 'IA', author: 'Alejandro G. Mota' })
      setEditing(null)
      loadArticles()
    } catch (err) {
      setError(err.message)
    }
  }

  const handleEdit = (article) => {
    setEditing(article.slug)
    setForm({ ...article })
  }

  const handleDelete = async (slug) => {
    if (!confirm('Eliminar este artículo?')) return
    await apiFetch(`/articles/${slug}`, { method: 'DELETE' })
    loadArticles()
  }

  const generateSlug = (title) => {
    return title.toLowerCase()
      .normalize('NFD').replace(/[\u0300-\u036f]/g, '')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '')
  }

  if (!token) {
    return (
      <div className="admin-page">
        <div className="login-card">
          <h1>Admin Login</h1>
          {error && <p className="error">{error}</p>}
          <form onSubmit={handleLogin}>
            <input type="text" placeholder="Usuario" value={loginForm.user}
              onChange={(e) => setLoginForm({ ...loginForm, user: e.target.value })} />
            <input type="password" placeholder="Contraseña" value={loginForm.pass}
              onChange={(e) => setLoginForm({ ...loginForm, pass: e.target.value })} />
            <button type="submit">Ingresar</button>
          </form>
        </div>
      </div>
    )
  }

  return (
    <div className="admin-page">
      <div className="admin-container">
        <header className="admin-header">
          <h1>Panel de Administración</h1>
          <button onClick={handleLogout} className="btn-logout">Cerrar Sesión</button>
        </header>

        {error && <p className="error">{error}</p>}

        <section className="admin-form-section">
          <h2>{editing ? 'Editar Artículo' : 'Nuevo Artículo'}</h2>
          <form onSubmit={handleSubmit} className="article-form">
            <input type="text" placeholder="Título" value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value, slug: editing ? form.slug : generateSlug(e.target.value) })} required />
            <input type="text" placeholder="Slug" value={form.slug}
              onChange={(e) => setForm({ ...form, slug: e.target.value })} required disabled={!!editing} />
            <select value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })}>
              <option value="IA">Inteligencia Artificial</option>
              <option value="Blockchain">Blockchain</option>
              <option value="Agronomía">Agronomía</option>
              <option value="Data Science">Data Science</option>
            </select>
            <input type="text" placeholder="Autor" value={form.author}
              onChange={(e) => setForm({ ...form, author: e.target.value })} />
            <textarea placeholder="Extracto / Resumen" value={form.excerpt}
              onChange={(e) => setForm({ ...form, excerpt: e.target.value })} rows={3} />
            <textarea placeholder="Contenido (HTML)" value={form.content}
              onChange={(e) => setForm({ ...form, content: e.target.value })} rows={15} />
            <div className="form-actions">
              <button type="submit">{editing ? 'Guardar Cambios' : 'Publicar'}</button>
              {editing && <button type="button" onClick={() => { setEditing(null); setForm({ title: '', slug: '', content: '', excerpt: '', category: 'IA', author: 'Alejandro G. Mota' }) }}>Cancelar</button>}
            </div>
          </form>
        </section>

        <section className="admin-articles">
          <h2>Artículos ({articles.length})</h2>
          <div className="articles-list">
            {articles.map((a) => (
              <div key={a.id} className="article-item">
                <div>
                  <strong>{a.title}</strong>
                  <span className="article-item-category">{a.category}</span>
                  <small>{new Date(a.published_at).toLocaleDateString('es-MX')}</small>
                </div>
                <div className="article-item-actions">
                  <button onClick={() => handleEdit(a)}>Editar</button>
                  <button onClick={() => handleDelete(a.slug)} className="btn-danger">Eliminar</button>
                </div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  )
}
