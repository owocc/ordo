import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './global.css'

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <div>
            <h1 className="text-lg font-bold">Ordo</h1>
            <i className="i-lucide-app-window-mac text-2xl"></i>
        </div>
    </StrictMode>
)
