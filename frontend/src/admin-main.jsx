import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import AdminPage from './pages/AdminPage/AdminPage'
import './styles/global.css'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <AdminPage />
  </StrictMode>,
)
