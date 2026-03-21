import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import './HomePage.css'

export default function HomePage() {
  useEffect(() => { document.title = 'Alejandro G. Mota — Blog' }, [])

  return (
    <div className="home">
      <section className="hero">
        <div className="hero-content">
          <h1>Alejandro G. Mota</h1>
          <p>Construyo cosas y aprendo en el proceso. Code, Business, Ideas y todo lo que vale la pena documentar.</p>
          <Link to="/entradas" className="hero-cta">Ver Artículos</Link>
        </div>
      </section>

      <section className="about">
        <div className="container">
          <h2>Sobre este blog</h2>
          <p>
            Publico lo que me hubiera servido a mí hace 6 meses.
            Desde lecciones técnicas reales hasta análisis de lo que pasa en el mundo
            — sin filtro y aprendiendo en público.
          </p>
        </div>
      </section>

      <section className="services">
        <div className="container">
          <h2>Categorías</h2>
          <div className="services-grid">
            <div className="service-card">
              <h3>Code</h3>
              <p>Arquitectura, microservicios, Go, Node.js, DevOps y lecciones técnicas del trabajo.</p>
            </div>
            <div className="service-card">
              <h3>Business</h3>
              <p>Emprendimiento, importación, operaciones y lecciones de escalar un negocio en México.</p>
            </div>
            <div className="service-card">
              <h3>Ideas</h3>
              <p>Geopolítica, ciencia y análisis de cosas que te hacen pensar.</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
