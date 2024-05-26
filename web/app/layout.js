import { Sora } from 'next/font/google'
import './globals.css'

const sora = Sora({ subsets: ['latin'] })

export const metadata = {
    title: 'Event Grid',
    description: 'Sets up your events with easy steps.',
}

export default function RootLayout({ children }) {
    return (
        <html lang="en">
            <body className={sora.className}>{children}</body>
        </html>
    )
}
