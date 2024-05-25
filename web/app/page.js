'use client'

import { addDaysToDate } from '@/commons/dateTime'
import { Avatar } from '@/components/avatar'
import Calendar from '@/components/calendar'
import { CalendarGrid } from '@/components/calendar/calendarGrid'
import { CalendarSettings } from '@/components/calendar/settings'
import { Input } from '@/components/form/input'
import { Backdrop } from '@/components/modal/backdrop'
import { useState } from 'react'

export default function Home() {
    const [selectedDate, setSelectedDate] = useState(new Date())

    const [selectedDuration, setSelectedDuration] = useState(15)

    const [dateRange, setDateRange] = useState({
        initialDate: selectedDate,
        finalDate: addDaysToDate(selectedDate, 6),
    })

    const [showCreateForm, setShowCreateForm] = useState(false)

    const handleCreateFormVisibility = () => {
        setShowCreateForm(!showCreateForm)
    }

    return (
        <main className="flex w-screen h-screen bg-primary-background relative">
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
                        toggleForm={handleCreateFormVisibility}
                    />
                </div>
            </div>
            {showCreateForm ? (
                <Backdrop 
                  onClose={() => setShowCreateForm(false)}
                >
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
