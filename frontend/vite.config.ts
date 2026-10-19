import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Target API backend. Default-nya memakai nama service Docker ("backend").
// Di luar Docker (mis. `make frontend`), set VITE_API_TARGET=http://localhost:8080.
const apiTarget =
  (globalThis as { process?: { env?: Record<string, string> } }).process?.env
    ?.VITE_API_TARGET ?? 'http://backend:8080'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    // Wajib: Cloudflare tunnel mengirim header Host berupa domain publik,
    // dan Vite >= 5.4 menolak host yang tidak terdaftar.
    allowedHosts: ['ujiemusi.orasaurus.my.id', 'localhost', '127.0.0.1'],
    proxy: {
      '/api': {
        target: apiTarget,
        changeOrigin: true,
      },
      // File statis (gambar, working instruction PDF) dilayani backend,
      // tanpa proxy ini Vite mengembalikan index.html sehingga PDF tidak tampil.
      '/uploads': {
        target: apiTarget,
        changeOrigin: true,
      },
    },
  },
})
