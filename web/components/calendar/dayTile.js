'use client'

import { generateTimeParts } from '@/commons/dateTime'
import { useState } from 'react'

export function DayTile({ startTime, parts = 2, toggleForm, date }) {
    let sections = generateTimeParts(startTime, parts)

    const [activeIndex, setActiveIndex] = useState(null)

    return (
        <div className="flex flex-col justify-between items-center h-full">
            {sections.slice(0, sections.length - 1).map((_, index) => (
                <div
                    key={index}
                    className={`text-xs w-full h-full text-center justify-center items-center text-white`}
                    // ${activeIndex === index ? 'text-black bg-white rounded-md' : 'text-white'}
                    onMouseEnter={() => setActiveIndex(index)}
                    onMouseLeave={() => setActiveIndex(null)}
                    onClick={() => {
                        toggleForm(date, sections[index], sections[index + 1])
                    }}
                >
                    {activeIndex === index
                        ? `${sections[index]} - ${sections[index + 1]}`
                        : ''}
                    {/* {
                        `${sections[index]} - ${sections[index + 1]}`
                    } */}
                </div>
            ))}
        </div>
    )
}
