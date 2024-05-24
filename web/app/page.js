'use client'

import { addDaysToDate } from '@/commons/date'
import { Avatar } from '@/components/avatar'
import Calendar from '@/components/calendar'
import { CalendarGrid } from '@/components/calendar/calendarGrid'
import { CalendarSettings } from '@/components/calendar/settings'
import { Input } from '@/components/form/input'
import { useState } from 'react'

export default function Home() {
    const [selectedDate, setSelectedDate] = useState(new Date())

    const [selectedDuration, setSelectedDuration] = useState(15)

    const [dateRange, setDateRange] = useState({
        initialDate: selectedDate,
        finalDate: addDaysToDate(selectedDate, 7),
    })

    return (
        <main className="flex w-screen h-screen bg-primary-background">
            <div className="md:w-[424px] h-full border-r border-primary-border flex flex-col justify-between">
                <div className="flex flex-col gap-2 p-6">
                    <Avatar name={'Steve Jobs'} />
                    <span className="text-white opacity-70 mb-4">
                        Sumit Vishwakarma
                    </span>
                    <Input
                        name={'duration'}
                        label={'Duration (in mins)'}
                        type={'number'}
                        placeholder={'enter the duration in minutes'}
                        value={selectedDuration}
                        onChange={(e) => {
                            setSelectedDuration(e.target.value)
                        }}
                        step={15}
                        min={0}
                    />
                </div>
                <div>
                    <Calendar
                        selectedDate={selectedDate}
                        handleSelectDate={(selectedDate) => {
                            setSelectedDate(selectedDate)
                        }}
                    />
                </div>
            </div>
            <div className="flex-1 h-full flex flex-col w-full">
                <CalendarSettings
                    selectedDate={selectedDate}
                    setSelectedDate={setSelectedDate}
                    dateRange={dateRange}
                    setDateRange={setDateRange}
                />
                <div className="flex flex-1 border-t border-primary-border">
                    <CalendarGrid
                        startDate={dateRange.initialDate}
                        endDate={dateRange.finalDate}
                    />
                </div>
            </div>
        </main>
    )
}
