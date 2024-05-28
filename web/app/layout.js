import './globals.css'

export const metadata = {
    title: 'Event Grid',
    description: 'Sets up your events with easy steps.',
}

export default function RootLayout({ children }) {
    return (
        <html lang="en">
            <head>
                <link rel="preconnect" href="https://fonts.googleapis.com"></link>
                <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin></link> 
                <link href="https://fonts.googleapis.com/css2?family=Sora:wght@100..800&display=swap" rel="stylesheet"></link> 
            </head>
            <body 
            >{children}</body>
        </html>
    )
}
