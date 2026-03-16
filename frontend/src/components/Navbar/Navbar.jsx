import { Link } from 'react-router-dom'
import './Navbar.css'

export default function Navbar() {
  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          <img src="/logo.png" alt="Tech & Science Blog" className="navbar-logo-img" />
          Tech & Science Blog
        </Link>
        <div className="navbar-links">
          <Link to="/">Inicio</Link>
          <Link to="/entradas">Artículos</Link>
        </div>
      </div>
    </nav>
  )
}
