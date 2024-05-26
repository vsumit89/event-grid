'use client'

import { addDaysToDate } from '@/commons/dateTime'

import { Avatar } from '@/components/avatar'
import Calendar from '@/components/calendar'
import { CalendarGrid } from '@/components/calendar/calendarGrid'
import { CalendarSettings } from '@/components/calendar/settings'
import { Backdrop } from '@/components/modal/backdrop'

import { useEvents } from '@/hooks/useEvents'
import useProfile from '@/hooks/useProfile'


import { useEffect, useState } from 'react'

export default function Home() {
    const [selectedDate, setSelectedDate] = useState(new Date())

    const [numDays, setNumDays] = useState(6)

    const [dateRange, setDateRange] = useState({
        initialDate: selectedDate,
        finalDate: addDaysToDate(selectedDate, numDays),
    })

    const [showCreateForm, setShowCreateForm] = useState(false)

    const [timeFormat, setTimeFormat] = useState(24)

    const handleCreateFormVisibility = () => {
        setShowCreateForm(!showCreateForm)
    }

    useEffect(() => {
        setDateRange({
            initialDate: selectedDate,
            finalDate: addDaysToDate(selectedDate, numDays),
        })
    }, [numDays])


    const {
        profile,
        error
    } = useProfile()


    const {
        events,
        loading,
        error: eventsError
    } = useEvents(dateRange.initialDate, dateRange.finalDate)


    return (
        <main className="flex w-screen h-screen bg-primary-background relative">
            <div className="md:w-[424px] h-full border-r border-primary-border flex flex-col justify-between">
                <div className="flex flex-col gap-2 p-6">
                    <Avatar name={profile?.name ? profile.name : '--'} />
                    <span className="text-white opacity-70 mb-4">
                        {profile?.name ? profile.name : '--'}
                    </span>
                    <button 
                        className='bg-primary-text rounded-lg px-4 py-2 w-full text-sm hover:bg-white hover:text-black transition transform duration-300'
                        onClick={() => setShowCreateForm(true)}
                    >
                        Create Event 
                    </button>
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
                    numDays={numDays}
                    setNumDays={setNumDays}
                    timeFormat={timeFormat}
                    setTimeFormat={setTimeFormat}
                />
                <div className="flex flex-1 border-t border-primary-border">
                    <CalendarGrid
                        startDate={dateRange.initialDate}
                        endDate={dateRange.finalDate}
                        toggleForm={handleCreateFormVisibility}
                        timeFormat={timeFormat}
                        openDate={(date) => {
                            setSelectedDate(date)
                            setDateRange({
                                initialDate: date,
                                finalDate: addDaysToDate(date, 0),
                            })

                            setNumDays(0)
                        }}
                        eventData={events}
                    />
                </div>
            </div>
            {showCreateForm ? (
                <Backdrop onClose={() => setShowCreateForm(false)}>
                    <div
                        className="p-6 bg-white text-black w-2/5 rounded-md"
                        onClick={(e) => {
                            e.stopPropagation()
                        }}
                    >
                        <div className="">Sumit</div>
                    </div>
                </Backdrop>
            ) : null}
        </main>
    )
}
