import { Sora } from 'next/font/google'
import './../globals.css'

const sora = Sora({ subsets: ['latin'] })

export const metadata = {
    title: 'Event Grid',
    description: 'Sets up your events with easy steps.',
}

export default function RootLayout({ children }) {
    return (
        <html lang="en">
            <body className={sora.className}>
                <div className="w-screen h-screen bg-primary-background flex justify-center items-center">
                    <div className="w-1/3 flex flex-col gap-4">
                        <h1 className="text-2xl text-primary-text">
                            Welcome to EventGrid
                        </h1>
                        <p className="text-sm text-primary-text opacity-75">
                            Easily organize your schedule, set reminders, and
                            never miss an important event.
                        </p>
                        {children}
                    </div>
                </div>
            </body>
        </html>
    )
}
