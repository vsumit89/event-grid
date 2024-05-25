'use client'

import { areSameDate } from '@/commons/dateTime'
import { dayNames, monthNames } from '@/constants/date'
import { CaretLeft, CaretRight } from '@phosphor-icons/react'
import { memo, useEffect, useState } from 'react'

export default memo(({ selectedDate, handleSelectDate }) => {
    const todaysDate = new Date()

    const [selectedMonthAndYear, setSelectedMonthAndYear] = useState({
        month: selectedDate ? selectedDate.getMonth() : todaysDate.getMonth(),
        year: selectedDate
            ? selectedDate.getFullYear()
            : todaysDate.getFullYear(),
    })

    const [activeDate, setActiveDate] = useState(todaysDate)

    const getDaysInMonth = (month, year) => {
        return new Date(year, month + 1, 0).getDate()
    }

    const getFirstDayOfMonth = (month, year) => {
        return new Date(year, month, 1).getDay()
    }

    const generateCalendarDates = (month, year) => {
        const daysInMonth = getDaysInMonth(month, year)
        const firstDayOfMonth = getFirstDayOfMonth(month, year)
        const dates = []

        // Add empty blocks for days before the first day of the month
        for (let i = 0; i < firstDayOfMonth; i++) {
            dates.push(null)
        }

        // Add dates for the current month
        for (let i = 1; i <= daysInMonth; i++) {
            dates.push(new Date(year, month, i))
        }

        // Add empty blocks for days after the last day of the month
        const remaining = 42 - dates.length
        for (let i = 0; i < remaining; i++) {
            dates.push(null)
        }

        return dates
    }

    const [calendarDates, setCalendarDates] = useState(
        generateCalendarDates(
            selectedMonthAndYear.month,
            selectedMonthAndYear.year
        )
    )

    const updateMonth = (increment) => {
        setSelectedMonthAndYear((prev) => {
            const newMonth = prev.month + increment

            if (newMonth > 11) {
                return { month: 0, year: prev.year + 1 }
            } else if (newMonth < 0) {
                return { month: 11, year: prev.year - 1 }
            } else {
                return { month: newMonth, year: prev.year }
            }
        })
    }

    // Example usage for incrementing and decrementing month
    const incrementMonth = () => updateMonth(1)
    const decrementMonth = () => updateMonth(-1)

    useEffect(() => {
        setCalendarDates(
            generateCalendarDates(
                selectedMonthAndYear.month,
                selectedMonthAndYear.year
            )
        )
    }, [selectedMonthAndYear])

    useEffect(() => {
        setSelectedMonthAndYear({
            month: selectedDate.getMonth(),
            year: selectedDate.getFullYear(),
        })

        setActiveDate(selectedDate)
    }, [selectedDate])

    return (
        <div className="px-5 py-3 flex flex-col gap-4">
            <div className="flex items-center justify-between">
                <span className="text-white">
                    {monthNames[selectedMonthAndYear.month]}{' '}
                    {selectedMonthAndYear.year}
                </span>
                <div className="flex gap-2">
                    <div
                        className="p-2 hover:bg-secondary-background rounded-md cursor-pointer"
                        onClick={() => decrementMonth()}
                    >
                        <CaretLeft className="text-white font-bold w-4 h-4" />
                    </div>
                    <div
                        className="p-2 hover:bg-secondary-background rounded-md cursor-pointer"
                        onClick={() => incrementMonth()}
                    >
                        <CaretRight className="text-white font-bold w-4 h-4" />
                    </div>
                </div>
            </div>
            <div className="grid grid-cols-7 gap-2 rounded-sm">
                {dayNames.map((day) => (
                    <span key={day} className="text-white text-xs mb-4 pl-1">
                        {day}
                    </span>
                ))}
                {calendarDates.map((date, index) => (
                    <button
                        key={index}
                        className={`p-1 text-center rounded-md cursor-pointer ${date ? (areSameDate(date, activeDate) ? 'bg-white opacity-100 text-' : 'bg-secondary-background opacity-60 text-white hover:bg-secondary-background hover:opacity-100') : 'bg-transparent'} h-12`}
                        onClick={() => {
                            handleSelectDate(date)
                            setActiveDate(date)
                        }}
                    >
                        {date ? date.getDate() : ''}
                    </button>
                ))}
            </div>
        </div>
    )
})
