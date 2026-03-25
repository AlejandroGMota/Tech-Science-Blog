import { useState, useEffect, useRef } from 'react'
import { Link } from 'react-router-dom'
import { useTheme } from '../../hooks/useTheme'
import './Navbar.css'

export default function Navbar() {
  const { theme, toggle } = useTheme()
  const [hidden, setHidden] = useState(false)
  const lastY = useRef(0)

  useEffect(() => {
    const onScroll = () => {
      const y = window.scrollY
      setHidden(y > 60 && y > lastY.current)
      lastY.current = y
    }
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  return (
    <nav className={`navbar${hidden ? ' navbar-hidden' : ''}`}>
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          <img src="/logo-44.png" alt="Tech & Science Blog" className="navbar-logo-img" width="22" height="22" />
          Tech & Science Blog
        </Link>
        <div className="navbar-links">
          <Link to="/entradas" className="navbar-main-link">Artículos</Link>
          <button className="theme-toggle" onClick={toggle} aria-label="Cambiar tema">
            <span className="theme-toggle-track">
              <span className="theme-toggle-icon theme-toggle-sun">☀️</span>
              <span className="theme-toggle-icon theme-toggle-moon">🌙</span>
              <span className="theme-toggle-thumb" />
            </span>
          </button>
        </div>
      </div>
    </nav>
  )
}
