/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        surface: {
          DEFAULT: '#0A0A0F',
          raised: '#13131A',
          overlay: '#1A1A24',
          border: '#22222E',
        },
        accent: {
          indigo: '#6366F1',
          cyan: '#06B6D4',
          emerald: '#10B981',
        },
        text: {
          primary: '#F8FAFC',
          secondary: '#94A3B8',
          muted: '#64748B',
        },
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      animation: {
        'glow': 'glow 2s ease-in-out infinite alternate',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        glow: {
          '0%': { boxShadow: '0 0 5px rgba(99,102,241,0.5)' },
          '100%': { boxShadow: '0 0 20px rgba(99,102,241,0.8)' }
        }
      },
      backdropBlur: {
        xs: '2px',
      }
    },
  },
  plugins: [require('tailwindcss-animate')],
}
