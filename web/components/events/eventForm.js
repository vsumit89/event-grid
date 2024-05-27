import { CircleNotch, X } from '@phosphor-icons/react'
import { useEffect, useState } from 'react'
import { Input } from '../form/input'
import CreatableSelect from 'react-select/creatable'
import {
    getDateForDateInput,
    getDateForUpdateInput,
    nowForDateInput,
    nowForDateInputWithDelay,
} from '@/commons/dateTime'
import { searchUsers } from '@/services/user'
import { isEmail } from '@/commons/validation'

export function EventForm({
    options,
    handleSubmit,
    onClose,
    updateEvent,
    ownerID,
    formError,
}) {
    const [eventDetails, setEventDetails] = useState({
        title: {
            value: '',
            error: '',
        },
        description: {
            value: '',
            error: '',
        },
        start: {
            value: options?.start
                ? getDateForDateInput(options.date, options.start)
                : nowForDateInput(),
            error: '',
        },
        end: {
            value: options?.end
                ? getDateForDateInput(options.date, options.end)
                : nowForDateInputWithDelay(30),
            error: '',
        },
        attendees: {
            value: [],
            error: '',
        },
        meetingUrl: {
            value: '',
            error: '',
        },
    })

    useEffect(() => {
        if (updateEvent) {
            setEventDetails((prev) => {
                return {
                    ...prev,
                    title: {
                        ...prev.title,
                        value: updateEvent.title,
                    },
                    description: {
                        ...prev.description,
                        value: updateEvent.description,
                    },
                    start: {
                        ...prev.start,
                        value: getDateForUpdateInput(updateEvent.start),
                    },
                    end: {
                        ...prev.end,
                        value: getDateForUpdateInput(updateEvent.end),
                    },
                    attendees: {
                        ...prev.attendees,
                        value: updateEvent.attendees
                            .filter((attendee) => attendee.id !== ownerID)
                            .map((attendee) => {
                                return attendee?.email
                            }),
                    },
                    meetingUrl: {
                        ...prev.meetingUrl,
                        value: updateEvent?.meeting_url || '',
                    },
                }
            })
        }
    }, [])

    const [loading, setLoading] = useState(false)

    const validateForm = () => {
        let isValid = true

        if (eventDetails.title.value === '') {
            setEventDetails((prev) => {
                return {
                    ...prev,
                    title: {
                        ...prev.title,
                        error: 'Title is required',
                    },
                }
            })
            isValid = false
        }

        if (
            eventDetails.description.value === '' ||
            eventDetails.description.value.length < 10
        ) {
            setEventDetails((prev) => {
                return {
                    ...prev,
                    description: {
                        ...prev.description,
                        error: 'Description is required',
                    },
                }
            })
            isValid = false
        }

        if (eventDetails.start.value === '') {
            setEventDetails((prev) => {
                return {
                    ...prev,
                    start: {
                        ...prev.start,
                        error: 'Start time is required',
                    },
                }
            })
            isValid = false
        }

        if (eventDetails.end.value === '') {
            setEventDetails((prev) => {
                return {
                    ...prev,
                    end: {
                        ...prev.end,
                        error: 'End time is required',
                    },
                }
            })
            isValid = false
        } else {
            if (eventDetails.end.value <= eventDetails.start.value) {
                setEventDetails((prev) => {
                    return {
                        ...prev,
                        end: {
                            ...prev.end,
                            error: 'End time cannot be before or equal to start time',
                        },
                    }
                })
                isValid = false
            }
        }

        if (eventDetails.attendees.value.length <= 0) {
            setEventDetails((prev) => {
                return {
                    ...prev,
                    attendees: {
                        ...prev.attendees,
                        error: 'Attendees are required',
                    },
                }
            })
            isValid = false
        } else {
            eventDetails.attendees.value.map((attendee) => {
                if (!isEmail(attendee)) {
                    setEventDetails((prev) => {
                        return {
                            ...prev,
                            attendees: {
                                ...prev.attendees,
                                error: 'please check the list of emails entered, it is not valid',
                            },
                        }
                    })
                    isValid = false
                }
            })
        }

        return isValid
    }

    const [userOptions, setUserOptions] = useState([])

    const [searchTerm, setSearchTerm] = useState('')

    return (
        <div
            className={`bg-modal-background p-4 w-2/5 rounded-md text-white flex flex-col gap-2 overflow-y-scroll`}
            onClick={(e) => {
                e.stopPropagation()
            }}
        >
            <div className="text-xl flex items-center justify-between">
                <span className="opacity-90">
                    {updateEvent ? 'Update' : 'Create'} an event
                </span>
                <span>
                    <X size={16} cursor={'pointer'} onClick={onClose} />
                </span>
            </div>
            <hr className="opacity-30 mb-2" />
            {<span className="text-xs text-error-text">{formError}</span>}
            <Input
                label="Title"
                placeholder="Title"
                value={eventDetails.title.value}
                onChange={(e) => {
                    setEventDetails((prev) => {
                        return {
                            ...prev,
                            title: {
                                value: e.target.value,
                                error: '',
                            },
                        }
                    })
                }}
                error={eventDetails.title.error}
            />
            <Input
                label="Description"
                placeholder="Description"
                value={eventDetails.description.value}
                onChange={(e) => {
                    setEventDetails((prev) => {
                        return {
                            ...prev,
                            description: {
                                value: e.target.value,
                                error: '',
                            },
                        }
                    })
                }}
                error={eventDetails.description.error}
            />
            <Input
                label="Meeting Link"
                placeholder="Meeting Link"
                value={eventDetails.meetingUrl.value}
                onChange={(e) => {
                    setEventDetails((prev) => {
                        return {
                            ...prev,
                            meetingUrl: {
                                value: e.target.value,
                                error: '',
                            },
                        }
                    })
                }}
                error={eventDetails.meetingUrl.error}
            />
            <div className="flex flex-col gap-1">
                <label className="text-primary-label text-sm">Attendees</label>
                <CreatableSelect
                    isMulti
                    isSearchable
                    inputValue={searchTerm}
                    value={eventDetails.attendees.value.map((attendee) => {
                        return {
                            value: attendee,
                            label: attendee,
                        }
                    })}
                    options={userOptions}
                    className="bg-secondary-background text-sm"
                    styles={{
                        control: (provided) => ({
                            ...provided,
                            border: 'none',
                            boxShadow: 'none',
                            '&:hover': {
                                border: 'none',
                                boxShadow: 'none',
                            },
                            backgroundColor: '#1f1f1f',
                            color: 'white',
                        }),
                        valueContainer: (provided) => ({
                            ...provided,
                            color: '#a6a6a6',
                            backgroundColor: '#1f1f1f',
                        }),
                        option: (provided) => ({
                            ...provided,
                            color: '#a6a6a6',
                            backgroundColor: '#1f1f1f',
                            '&:hover': {
                                color: 'white',
                                backgroundColor: '#a6a6a6',
                            },
                        }),
                        dropdownIndicator: (provided) => ({
                            ...provided,
                            color: '#a6a6a6',
                            '&:hover': {
                                color: '#a6a6a6',
                            },
                        }),
                        menu: (provided) => ({
                            ...provided,
                            backgroundColor: '#1f1f1f',
                        }),
                        input: (provided) => ({
                            ...provided,
                            color: '#FCFCFD',
                        }),
                    }}
                    onChange={(e) => {
                        setEventDetails((prev) => {
                            return {
                                ...prev,
                                attendees: {
                                    value: e.map((e) => e.value),
                                },
                            }
                        })
                    }}
                    onInputChange={async (e) => {
                        setSearchTerm(e)

                        setUserOptions([])
                        try {
                            const users = await searchUsers(e)

                            const newUsers = users.map((user) => {
                                return {
                                    value: user.email,
                                    label: user.email,
                                }
                            })

                            setUserOptions(newUsers)
                        } catch (error) {
                            console.log(error)
                        }
                    }}
                />
                {eventDetails.attendees.error && (
                    <span className="text-xs text-error-text">
                        {eventDetails.attendees.error}
                    </span>
                )}
            </div>
            <Input
                label="Start"
                type="datetime-local"
                value={eventDetails.start.value}
                onChange={(e) => {
                    setEventDetails((prev) => {
                        return {
                            ...prev,
                            start: {
                                value: e.target.value,
                                error: '',
                            },
                        }
                    })
                }}
                error={eventDetails.start.error}
            />
            <Input
                label="End"
                type="datetime-local"
                value={eventDetails.end.value}
                onChange={(e) => {
                    setEventDetails((prev) => {
                        return {
                            ...prev,
                            end: {
                                value: e.target.value,
                                error: '',
                            },
                        }
                    })
                }}
                error={eventDetails.end.error}
            />
            <button
                className="mt-2 bg-white text-black px-2 py-1 rounded-md text-sm w-fit ring-white"
                onClick={async () => {
                    if (loading) return

                    const isValid = validateForm()

                    if (!isValid) return

                    setLoading(true)

                    try {
                        handleSubmit({
                            title: eventDetails.title.value,
                            description: eventDetails.description.value,
                            meeting_url: eventDetails.meetingUrl.value,
                            attendees: eventDetails.attendees.value,
                            start_time: `${eventDetails.start.value}:00+05:30`,
                            end_time: `${eventDetails.end.value}:00+05:30`,
                        })
                    } catch (error) {
                        setFormError(error.message)
                    }

                    setLoading(false)
                }}
            >
                {loading ? (
                    <CircleNotch size={20} className="animate-spin" />
                ) : (
                    <>{updateEvent ? 'Update Event' : 'Create Event'}</>
                )}
            </button>
        </div>
    )
}
