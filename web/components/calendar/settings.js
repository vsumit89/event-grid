'use client'

import { addDaysToDate } from '@/commons/dateTime'
import { shortMonthNames } from '@/constants/date'
import { CaretLeft, CaretRight } from '@phosphor-icons/react'
import { useEffect, useState } from 'react'

export function CalendarSettings({
    selectedDate,
    setSelectedDate,
    dateRange,
    setDateRange,
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
            finalDate: addDaysToDate(selectedDate, 6),
        })
    }, [selectedDate])

    return (
        <div className="px-4 py-5 flex justify-between text-white">
            <div className="flex items-center gap-4">
                <div className="flex gap-2">
                    <div
                        className="p-2 hover:bg-secondary-background rounded-md cursor-pointer"
                        onClick={() => updateDateRange(-7)}
                    >
                        <CaretLeft className="text-white font-bold w-4 h-4" />
                    </div>
                    <div
                        className="p-2 hover:bg-secondary-background rounded-md cursor-pointer"
                        onClick={() => updateDateRange(7)}
                    >
                        <CaretRight className="text-white font-bold w-4 h-4" />
                    </div>
                </div>
                <span className="text-base text-white">
                    {dateRange.initialDate.getDate()}{' '}
                    {dateRange.initialDate.getMonth() !==
                    dateRange.finalDate.getMonth()
                        ? shortMonthNames[dateRange.initialDate.getMonth()]
                        : ''}
                    {dateRange.initialDate.getFullYear() ===
                    dateRange.finalDate.getFullYear()
                        ? ''
                        : `, ${dateRange.initialDate.getFullYear()}`}{' '}
                    - {dateRange.finalDate.getDate()}{' '}
                    {shortMonthNames[dateRange.finalDate.getMonth()]}
                    {`, ${dateRange.finalDate.getFullYear()}`}
                </span>
                <button
                    className="bg-button-primary ml-4 px-2 py-1 rounded-md border border-primary-border hover:border-white text-sm"
                    onClick={() => {
                        setSelectedDate(new Date())
                    }}
                >
                    Today
                </button>
            </div>

            <div></div>
        </div>
    )
}
