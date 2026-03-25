import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom'
import { useEffect, lazy, Suspense } from 'react'
import Navbar from './components/Navbar/Navbar'
import Footer from './components/Footer/Footer'

const HomePage = lazy(() => import('./pages/HomePage/HomePage'))
const EntradasPage = lazy(() => import('./pages/EntradasPage/EntradasPage'))
const ArticlePage = lazy(() => import('./pages/ArticlePage/ArticlePage'))

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
        <Suspense fallback={<div style={{ minHeight: '100vh' }} />}>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/entradas" element={<EntradasPage />} />
            <Route path="/entradas/:slug" element={<ArticlePage />} />
          </Routes>
        </Suspense>
      </main>
      <Footer />
    </BrowserRouter>
  )
}
