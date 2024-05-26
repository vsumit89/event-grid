
import { NextResponse } from 'next/server'

// This function can be marked `async` if using `await` inside
export function middleware(request) {
    let cookie = request.cookies.get('access_token')
    
    if (request.nextUrl.pathname.startsWith('/_next')) {
        return NextResponse.next()
    }

    if (request.nextUrl.pathname == '/') {
        if (!cookie) {
            return NextResponse.redirect(new URL('/auth/login', request.url))
        }
    }

    return NextResponse.next()
}

// See "Matching Paths" below to learn more
export const config = {
    matcher: '/((?!api|_next/static|_next/image|favicon.ico).*)',
}
