import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom'
import { useEffect } from 'react'
import Navbar from './components/Navbar/Navbar'
import Footer from './components/Footer/Footer'
import HomePage from './pages/HomePage/HomePage'
import EntradasPage from './pages/EntradasPage/EntradasPage'
import ArticlePage from './pages/ArticlePage/ArticlePage'

function ScrollToTop() {
  const { pathname } = useLocation()
  useEffect(() => { window.scrollTo(0, 0) }, [pathname])
  return null
}

export default function App() {
  return (
    <BrowserRouter>
      <ScrollToTop />
      <Navbar />
      <main>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/entradas" element={<EntradasPage />} />
          <Route path="/entradas/:slug" element={<ArticlePage />} />
        </Routes>
      </main>
      <Footer />
    </BrowserRouter>
  )
}
