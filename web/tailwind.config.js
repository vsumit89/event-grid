/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        './pages/**/*.{js,ts,jsx,tsx,mdx}',
        './components/**/*.{js,ts,jsx,tsx,mdx}',
        './app/**/*.{js,ts,jsx,tsx,mdx}',
    ],
    theme: {
        extend: {
            backgroundImage: {
                'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
                'gradient-conic':
                    'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
            },
            colors: {
                'primary-background': '#1f1f1f',
                'primary-border': '#575757',
                'primary-text': '#a6a6a6',
                'secondary-background': '#404040',
                'primary-label': '#FCFCFD',
                'button-primary': '#1a1a1a',
                'modal-background': '#1a1a1a',
            },
        },
    },
    plugins: [],
}
