import { Link } from 'react-router-dom'
import './HomePage.css'

export default function HomePage() {
  return (
    <div className="home">
      <section className="hero">
        <div className="hero-content">
          <h1>Tech & Science Blog</h1>
          <p>Explorando el futuro de la tecnología, inteligencia artificial y ciencia</p>
          <Link to="/entradas" className="hero-cta">Ver Artículos</Link>
        </div>
      </section>

      <section className="about">
        <div className="container">
          <h2>Quiénes Somos</h2>
          <p>
            Somos un blog dedicado a explorar y compartir los avances más recientes
            en tecnología, inteligencia artificial, ciencia de datos y agronomía tecnológica.
          </p>
        </div>
      </section>

      <section className="services">
        <div className="container">
          <h2>Qué Hacemos</h2>
          <div className="services-grid">
            <div className="service-card">
              <h3>Inteligencia Artificial</h3>
              <p>Análisis y explicaciones sobre los avances en IA y machine learning.</p>
            </div>
            <div className="service-card">
              <h3>Ciencia de Datos</h3>
              <p>Aplicaciones prácticas del data science en diversas industrias.</p>
            </div>
            <div className="service-card">
              <h3>Agronomía Tech</h3>
              <p>Tecnología aplicada al campo y la agricultura moderna.</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
