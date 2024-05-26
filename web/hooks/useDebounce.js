import { useEffect, useRef } from 'react'

export const useDebouncedCallback = (callback, delay) => {
    const timeoutRef = useRef(null)

    const debouncedCallback = (...args) => {
        clearTimeout(timeoutRef.current)
        timeoutRef.current = setTimeout(() => {
            callback(...args)
        }, delay)
    }

    useEffect(() => {
        return () => {
            clearTimeout(timeoutRef.current) // Cleanup timeout on unmount
        }
    }, [])

    return debouncedCallback
}
