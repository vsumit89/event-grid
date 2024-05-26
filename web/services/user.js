export const getProfile = async () => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/users/profile`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
    })

    const data = await response.json()

    if (!response.ok) {
        throw new Error(data.error.message || 'Something went wrong')
    }

    return data
}