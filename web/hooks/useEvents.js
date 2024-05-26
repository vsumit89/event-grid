import { getDateString } from '@/commons/dateTime'
import { getEvent, getEvents } from '@/services/event'

import { useState, useEffect, useCallback } from 'react'

export const useEvents = (start, end) => {
    const [events, setEvents] = useState([])

    const [loading, setLoading] = useState(true)

    const [error, setError] = useState('')

    const abortController = new AbortController()

    const fetchEvents = async () => {
        try {
            let startDate = getDateString(start)
            let endDate = getDateString(end)

            const data = await getEvents(startDate, endDate)

            setEvents(data)
        } catch (error) {
            setError(error.message)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchEvents()
        return () => {
            abortController.abort()
        }
    }, [start, end])

    const refreshEvents = () => {
        fetchEvents()
    }

    return { events, loading, error, refreshEvents }
}

export default useEvents

export function useEvent({ id }) {
    const [event, setEvent] = useState(null)

    const [loading, setLoading] = useState(true)

    const [error, setError] = useState('')

    const fetchEvent = useCallback(async () => {
        try {
            const data = await getEvent(id)
            setEvent(data)
        } catch (error) {
            setError(error.message)
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        fetchEvent()
    }, [fetchEvent])

    return { event, loading, error }
}
