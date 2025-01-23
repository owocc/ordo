import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import deno from '@deno/vite-plugin'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), deno(), tailwindcss()],
  build: {
    outDir: '../../dist/web'
  },
  server: {
    port: 1420
  }
})
