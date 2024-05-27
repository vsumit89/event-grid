'use client'

import { formatTimeForDetails, getFromToTime } from '@/commons/dateTime'
import { CircleNotch, PencilSimple, Trash, X } from '@phosphor-icons/react'
import { useState } from 'react'

export function EventTile({ event, top, height, onClick }) {
    if (height === 0) return null
    return (
        <div
            className={`absolute bg-white ${height > 40 ? 'text-sm' : 'text-xs'} mt-16 w-full rounded-md flex flex-col items-center justify-center cursor-pointer`}
            style={{
                top: top,
                height: height,
            }}
            onClick={onClick}
        >
            <span className="text-sm text-center">{event.title}</span>
            {height > 40 && (
                <span className="text-xs">
                    {getFromToTime(event.start, event.end)}
                </span>
            )}
        </div>
    )
}

export function EventDetailTile({ event, onClose, userId, onDelete, onEdit }) {
    let owner = event.attendees.find((attendee) => attendee.id === userId)

    const [loading, setLoading] = useState(false)

    return (
        <div
            className="bg-modal-background p-4 w-2/5 rounded-md text-white flex flex-col gap-2 max-h-1/2 overflow-y-scroll"
            onClick={(e) => {
                e.stopPropagation()
            }}
        >
            <div className="text-xl flex items-center justify-between">
                <span className="opacity-90">{event.title}</span>
                <span>
                    <X size={16} cursor={'pointer'} onClick={onClose} />
                </span>
            </div>
            <div className="text-sm opacity-70">
                {formatTimeForDetails(event.start, event.end)}
            </div>
            <div className="text-sm opacity-70">{event.description}</div>
            <span className="text-sm opacity-70">
                Created by : {owner?.name}
            </span>
            <div>
                <span className="text-sm opacity-70">Attendees</span>
                <ul className="h-16 overflow-y-scroll">
                    {event.attendees
                        .filter((attendee) => attendee.id !== userId)
                        .map((attendee) => (
                            <li key={attendee} className="text-sm opacity-70">
                                -{' '}
                                {attendee?.name
                                    ? attendee.name
                                    : attendee.email}
                            </li>
                        ))}
                </ul>
            </div>
            <div className="flex flex-row-reverse justify-between items-center opacity-70">
                <div className="flex gap-4 items-center">
                    <PencilSimple
                        size={20}
                        cursor={'pointer'}
                        className="hover:text-primary-text"
                        onClick={() => onEdit(event?.id)}
                    />
                    <div>
                        {loading ? (
                            <CircleNotch size={20} className="animate-spin" />
                        ) : (
                            <Trash
                                size={20}
                                cursor={'pointer'}
                                className="hover:text-error-text"
                                onClick={() => onDelete(event?.id)}
                            />
                        )}
                    </div>
                </div>
                {event?.meeting_url && (
                    <a
                        className="text-sm bg-white text-black px-4 py-1 rounded-md w-fit"
                        href={event?.meeting_url}
                        target="_blank"
                        rel="noreferrer"
                    >
                        Join Event
                    </a>
                )}
            </div>
        </div>
    )
}
