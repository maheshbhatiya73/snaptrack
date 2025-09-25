import { defineNuxtConfig } from 'nuxt/config'
import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  build: {
    transpile: ['vue-chart-3', 'chart.js']
  },
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      backendUrl: process.env.NUXT_PUBLIC_BACKEND_URL || 'http://localhost:8080'
    }
  },
  vite: {
    optimizeDeps: {
      include: ['vue-chart-3', 'chart.js']
    },
    plugins: [
      tailwindcss(),
    ],
  },
})