'use client'

import { addDaysToDate } from '@/commons/dateTime'

import { Avatar } from '@/components/avatar'
import Calendar from '@/components/calendar'
import { CalendarGrid } from '@/components/calendar/calendarGrid'
import { CalendarSettings } from '@/components/calendar/settings'
import { EventForm } from '@/components/events/eventForm'
import { EventDetailTile } from '@/components/events/eventTile'
import { Backdrop } from '@/components/modal/backdrop'

import { useEvents } from '@/hooks/useEvents'
import useProfile from '@/hooks/useProfile'

import {
    createEvent,
    deleteEvent,
    getEvent,
    updateEvent,
} from '@/services/event'
// import { cookies } from 'next/headers'

import { useEffect, useState } from 'react'

export default function Home() {
    const [selectedDate, setSelectedDate] = useState(new Date())

    const [numDays, setNumDays] = useState(6)

    const [dateRange, setDateRange] = useState({
        initialDate: selectedDate,
        finalDate: addDaysToDate(selectedDate, numDays),
    })

    const [showCreateForm, setShowCreateForm] = useState(false)

    const [showUpdateForm, setShowUpdateForm] = useState(false)

    const [timeFormat, setTimeFormat] = useState(24)

    const [createFormOptions, setCreateFormOptions] = useState({
        date: null,
        start: null,
        end: null,
    })

    const handleCreateFormVisibility = (date, start, end) => {
        setFormError('')

        setCreateFormOptions({ date, start, end })

        setShowCreateForm(!showCreateForm)
    }

    useEffect(() => {
        setDateRange({
            initialDate: selectedDate,
            finalDate: addDaysToDate(selectedDate, numDays),
        })
    }, [numDays])

    const { profile, error } = useProfile()

    const {
        events,
        error: eventsError,
        refreshEvents,
    } = useEvents(dateRange.initialDate, dateRange.finalDate)

    const [eventInFocus, setEventInFocus] = useState(null)

    const [formError, setFormError] = useState('')

    return (
        <main className="flex w-screen h-screen bg-primary-background relative">
            <div className="md:w-[424px] h-full border-r border-primary-border flex flex-col justify-between">
                <div className="flex flex-col gap-2 p-6">
                    <div className="flex items-center justify-between">
                        <Avatar name={profile?.name ? profile.name : '--'} />
                    </div>
                    <span className="text-white opacity-70 mb-4">
                        {profile?.name ? profile.name : '--'}
                    </span>
                    <button
                        className="bg-primary-text rounded-lg px-4 py-2 w-full text-sm hover:bg-white hover:text-black transition transform duration-300"
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
                        showEventDetails={async (eventDetails) => {
                            try {
                                let eventData = await getEvent(eventDetails.id)
                                setEventInFocus(eventData)
                            } catch (error) {
                                setEventInFocus(null)
                            }
                        }}
                    />
                </div>
            </div>
            {showCreateForm ? (
                <Backdrop onClose={() => setShowCreateForm(false)}>
                    <EventForm
                        options={createFormOptions}
                        onClose={() => setShowCreateForm(false)}
                        handleSubmit={async (data) => {
                            try {
                                await createEvent(data)
                                refreshEvents()
                                setShowCreateForm(false)
                            } catch (error) {
                                setFormError(error.message)
                            }
                        }}
                        ownerID={profile?.id}
                        formError={formError}
                    />
                </Backdrop>
            ) : null}
            {eventInFocus && !showUpdateForm && (
                <Backdrop onClose={() => setEventInFocus(null)}>
                    <EventDetailTile
                        event={eventInFocus}
                        userId={profile?.id}
                        onClose={() => setEventInFocus(null)}
                        onDelete={async () => {
                            try {
                                await deleteEvent(eventInFocus.id)
                                setEventInFocus(null)
                            } catch (error) {
                                console.log(error)
                            } finally {
                                refreshEvents()
                            }
                        }}
                        onEdit={() => {
                            setFormError('')
                            setShowUpdateForm(true)
                        }}
                    />
                </Backdrop>
            )}
            {showUpdateForm && eventInFocus && (
                <Backdrop onClose={() => setShowUpdateForm(false)}>
                    <EventForm
                        onClose={() => setShowUpdateForm(false)}
                        handleSubmit={async (data) => {
                            try {
                                await updateEvent(eventInFocus.id, data)
                                refreshEvents()
                                setShowUpdateForm(false)
                                setEventInFocus(null)
                            } catch (error) {
                                setFormError(error.message)
                            }
                        }}
                        updateEvent={eventInFocus}
                        ownerID={profile?.id}
                        formError={formError}
                    />
                </Backdrop>
            )}
        </main>
    )
}
