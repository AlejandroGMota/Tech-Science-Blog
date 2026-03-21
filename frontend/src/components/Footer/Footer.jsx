import { Link } from 'react-router-dom'
import './Footer.css'

export default function Footer() {
  return (
    <footer className="footer">
      <div className="footer-container">
        <div className="footer-links">
          <Link to="/entradas">Artículos</Link>
          <a href="https://alejandrogmota.com" target="_blank" rel="noopener noreferrer">alejandrogmota.com</a>
        </div>
        <p>&copy; {new Date().getFullYear()} Tech & Science Blog</p>
      </div>
    </footer>
  )
}
