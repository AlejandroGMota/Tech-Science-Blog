import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar/Navbar'
import Footer from './components/Footer/Footer'
import HomePage from './pages/HomePage/HomePage'
import EntradasPage from './pages/EntradasPage/EntradasPage'
import ArticlePage from './pages/ArticlePage/ArticlePage'

export default function App() {
  return (
    <BrowserRouter>
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
