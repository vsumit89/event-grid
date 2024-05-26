'use client'

import { addDaysToDate } from '@/commons/dateTime'

import { CaretLeft, CaretRight } from '@phosphor-icons/react'
import { useEffect } from 'react'

import DateRangeFormatter from '../calendar/dateRangeFormatter'

export function CalendarSettings({
    selectedDate,
    setSelectedDate,
    dateRange,
    setDateRange,
    numDays,
    setNumDays,
    setTimeFormat,
    timeFormat
}) {
    const updateDateRange = (increment) => {
        setDateRange((prev) => {
            const newInitialDate = addDaysToDate(prev.initialDate, increment)
            setSelectedDate(newInitialDate)
            return {
                initialDate: newInitialDate,
                finalDate: addDaysToDate(newInitialDate, increment),
            }
        })
    }

    useEffect(() => {
        setDateRange({
            initialDate: selectedDate,
            finalDate: addDaysToDate(selectedDate, numDays),
        })
    }, [selectedDate])

    return (
        <div className="px-4 py-5 flex justify-between text-white">
            <div className="flex items-center gap-4">
                <div className="flex gap-2">
                    <div
                        className="p-2 hover:bg-secondary-background rounded-md cursor-pointer"
                        onClick={() => updateDateRange(-(numDays + 1))}
                    >
                        <CaretLeft className="text-white font-bold w-4 h-4" />
                    </div>
                    <div
                        className="p-2 hover:bg-secondary-background rounded-md cursor-pointer"
                        onClick={() => updateDateRange(numDays + 1)}
                    >
                        <CaretRight className="text-white font-bold w-4 h-4" />
                    </div>
                </div>
                <DateRangeFormatter dateRange={dateRange} />
                <button
                    className="bg-button-primary ml-4 px-2 py-1 rounded-md border border-primary-border hover:border-white text-sm"
                    onClick={() => {
                        setSelectedDate(new Date())
                    }}
                >
                    Today
                </button>
            </div>
            <div className="flex items-center gap-2">
                {' '}
                <div className="flex gap-4 py-1 px-2 border border-primary-border rounded-md text-sm text-primary-text hover:border-white">
                    <button
                        className={`${timeFormat === 12 ? 'bg-secondary-background px-2 py-1 rounded-md text-white' : 'bg-transparent'}`}
                        onClick={() => {
                            setTimeFormat(12)
                        }}
                    >
                        12h
                    </button>
                    <button
                        className={`${timeFormat === 24? 'bg-secondary-background px-2 py-1 rounded-md text-white' : 'bg-transparent'}`}
                        onClick={() => {
                            setTimeFormat(24)
                        }}
                    >
                        24h
                    </button>
                </div>
                <div className="flex gap-4 py-1 px-2 border border-primary-border rounded-md text-sm text-primary-text hover:border-white">
                    <button
                        className={`${numDays + 1 === 1 ? 'bg-secondary-background px-2 py-1 rounded-md text-white' : 'bg-transparent'}`}
                        onClick={() => {
                            setNumDays(0)
                        }}
                    >
                        Daily
                    </button>
                    <button
                        className={`${numDays + 1 === 7 ? 'bg-secondary-background px-2 py-1 rounded-md text-white' : 'bg-transparent'}`}
                        onClick={() => {
                            setNumDays(6)
                        }}
                    >
                        Weekly
                    </button>
                </div>
            </div>
        </div>
    )
}
