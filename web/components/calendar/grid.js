'use client'

import { useState } from 'react'

export const Grid = () => {
    const [selectedTime, setSelectedTime] = useState(null)
    const [isFormOpen, setIsFormOpen] = useState(false)

    const handleTimeClick = (day, hour, minutes) => {
        setSelectedTime({ day, hour, minutes })
        setIsFormOpen(true)
    }

    const handleCloseForm = () => {
        setIsFormOpen(false)
    }

    return (
        <div className="h-[90vh] overflow-auto w-full">
            {/* Render your week grid */}
            <div className="flex">
                {/* Loop through each day */}
                {[...Array(7)].map((_, index) => (
                    <div key={index} className="flex-1">
                        {/* Render each hour */}
                        {[...Array(24)].map((_, hour) => (
                            <div key={hour} className="flex flex-col">
                                {/* Loop through each 15-minute interval */}
                                {[0, 15, 30, 45].map((minutes) => (
                                    <div
                                        key={hour * 60 + minutes}
                                        className="h-12 border border-gray-200 flex items-center justify-center cursor-pointer hover:bg-gray-100 text-white"
                                        onClick={() =>
                                            handleTimeClick(
                                                index,
                                                hour,
                                                minutes
                                            )
                                        }
                                    >
                                        {/* Display time */}
                                        {hour < 10 ? '0' + hour : hour}:
                                        {minutes === 0 ? '00' : minutes}
                                    </div>
                                ))}
                            </div>
                        ))}
                    </div>
                ))}
            </div>

            {/* Render form */}
            {isFormOpen && (
                <div className="fixed inset-0 flex items-center justify-center bg-white bg-opacity-50">
                    <div className="bg-white p-4 rounded shadow-lg">
                        {/* Your form components go here */}
                        <h2 className="text-lg font-semibold mb-4">
                            Appointment Form
                        </h2>
                        <p className="mb-2">
                            Selected Time: Day {selectedTime.day},{' '}
                            {selectedTime.hour < 10
                                ? '0' + selectedTime.hour
                                : selectedTime.hour}
                            :
                            {selectedTime.minutes === 0
                                ? '00'
                                : selectedTime.minutes}
                        </p>
                        {/* Add form inputs, buttons, etc. */}
                        <button
                            className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded"
                            onClick={handleCloseForm}
                        >
                            Close
                        </button>
                    </div>
                </div>
            )}
        </div>
    )
}

export default CalendarGrid
