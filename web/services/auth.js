export const login = async (email, password) => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            email,
            password,
        }),
        credentials: 'include',
    })

    const data = await response.json()

    if (!response.ok) {
        throw new Error(data.error.message || 'Something went wrong')
    }

    return data
}

