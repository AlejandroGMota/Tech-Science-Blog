import './Footer.css'

export default function Footer() {
  return (
    <footer className="footer">
      <div className="footer-container">
        <p>&copy; {new Date().getFullYear()} Tech & Science Blog — Alejandro G. Mota</p>
      </div>
    </footer>
  )
}
