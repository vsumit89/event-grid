export const getEvents = async (startDate, endDate) => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/events?` + new URLSearchParams({
        start_date: startDate,
        end_date: endDate,
    }), {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
    })

    const data = await response.json()

    console.log("data", data)
    if (!response.ok) {
        throw new Error(data.error.message || 'Something went wrong')
    }

    return data
}