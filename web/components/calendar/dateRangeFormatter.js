'use client'
import { formatDateRange } from '@/commons/dateTime'
import { memo, useEffect, useState } from 'react'

function DateRangeFormatter({ dateRange }) {
    // const [isLoading, setIsLoading] = useState(true)
    const [formattedDateRange, setFormattedDateRange] = useState('')

    useEffect(() => {
        // setIsLoading(true)
        const timer = setTimeout(() => {
            const formatted = formatDateRange(dateRange)
            setFormattedDateRange(formatted)
            // setIsLoading(false)
        }, 100) // Simulating a 1-second delay

        return () => clearTimeout(timer)
    }, [dateRange])

    return (
        <div>
            <span>{formattedDateRange}</span>
        </div>
    )
}

const areSameDate = (prevProps, nextProps) => {
    return (
        prevProps.initialDate !== nextProps.initialDate &&
        prevProps.finalDate !== nextProps.finalDate
    )
}

export default memo(DateRangeFormatter, areSameDate)
