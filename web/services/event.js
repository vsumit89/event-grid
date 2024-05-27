export const getEvents = async (startDate, endDate) => {
    const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/events?` +
            new URLSearchParams({
                start_date: startDate,
                end_date: endDate,
            }),
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
        throw new Error(data.error.message || 'Something went wrong')
    }

    return data
}

export const getEvent = async (id) => {
    const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/events/${id}`,
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
        throw new Error(data.error.message || 'Something went wrong')
    }

    return data
}

export const deleteEvent = async (id) => {
    const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/events/${id}`,
        {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        }
    )

    const data = await response.json()

    if (!response.ok) {
        throw new Error(data.error.message || 'Something went wrong')
    }

    return data
}

export const createEvent = async (data) => {
    try {
        const response = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL}/api/v1/events`,
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
                credentials: 'include',
            }
        )

        const responseData = await response.json()

        if (!response.ok) {
            throw new Error(
                responseData.error.message || 'Something went wrong'
            )
        }
    } catch (error) {
        throw new Error(error)
    }

    return responseData
}

export const updateEvent = async (id, data) => {
    try {
        const response = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL}/api/v1/events/${id}`,
            {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
                credentials: 'include',
            }
        )

        const responseData = await response.json()

        if (!response.ok) {
            throw new Error(
                responseData.error.message || 'Something went wrong'
            )
        }
    } catch (error) {
        throw new Error(error)
    }

    return responseData
}
