export function useSearchUsers() {
    const [loading, setLoading] = useState(false)
    const [users, setUsers] = useState([])
    const [error, setError] = useState(null)

    const searchUsers = async (query) => {
        setLoading(true)
        setError(null)
        setUsers([])
        const response = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL}/api/v1/users/search?` +
                new URLSearchParams({ query, limit: 5 }),
            {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            }
        )

        const data = await response.json()

        if (!response.ok) {
            setLoading(false)
            setError(data.error.message || 'Something went wrong')
            return
        }

        setLoading(false)
        setUsers(data)
    }

    return { loading, users, error, searchUsers }
}
