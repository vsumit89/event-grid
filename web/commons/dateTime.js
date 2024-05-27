import { dayNames, monthNames, shortMonthNames } from '@/constants/date'

export function areSameDate(firstDate, secondDate) {
    // Extract year, month, and day from the first date
    let firstYear = firstDate.getFullYear()
    let firstMonth = firstDate.getMonth()
    let firstDay = firstDate.getDate()

    // Extract year, month, and day from the second date
    let secondYear = secondDate.getFullYear()
    let secondMonth = secondDate.getMonth()
    let secondDay = secondDate.getDate()

    // Compare the dates
    return (
        firstYear === secondYear &&
        firstMonth === secondMonth &&
        firstDay === secondDay
    )
}

export function addDaysToDate(initialDate, daysToAdd) {
    var result = new Date(initialDate)
    result.setDate(result.getDate() + daysToAdd)
    return result
}

export function getDateList(startDate, endDate) {
    const dateList = []
    let currentDate = startDate
    while (currentDate <= endDate) {
        dateList.push(currentDate)
        currentDate = addDaysToDate(currentDate, 1)
    }
    return dateList
}

export function generateTimeParts(startTime, numParts) {
    let timeParts = []
    const interval = Math.floor(60 / numParts)

    for (let i = 0; i < numParts; i++) {
        const minutes = i * interval
        const hours = Math.floor(startTime + minutes / 60)
        const formattedMinutes = String(minutes % 60).padStart(2, '0')
        const formattedHours = String(hours % 24).padStart(2, '0')

        timeParts.push(`${formattedHours}:${formattedMinutes}`)
    }

    timeParts.push(`${String((startTime + 1) % 24).padStart(2, '0')}:00`)

    return timeParts
}

export function splitMeetings(meetings) {
    const daysWithMeetings = {}

    meetings.forEach((meeting) => {
        const startDate = new Date(meeting.start)
        const endDate = new Date(meeting.end)
        const startDay = startDate.toDateString()
        const endDay = endDate.toDateString()

        if (startDay === endDay) {
            // Meeting is within a single day
            if (!daysWithMeetings[startDay]) {
                daysWithMeetings[startDay] = []
            }
            daysWithMeetings[startDay].push(meeting)
        } else {
            // Meeting spans across multiple days
            const midnightTime = new Date(
                startDate.getFullYear(),
                startDate.getMonth(),
                startDate.getDate() + 1
            )

            const firstPart = {
                ...meeting,
                end: getDateForDateInput(midnightTime, '00:00:00+05:30'),
            }

            const secondPart = {
                ...meeting,
                start: getDateForDateInput(midnightTime, '00:00:00+05:30'),
            }

            if (!daysWithMeetings[startDay]) {
                daysWithMeetings[startDay] = []
            }
            daysWithMeetings[startDay].push(firstPart)

            if (!daysWithMeetings[endDay]) {
                daysWithMeetings[endDay] = []
            }
            daysWithMeetings[endDay].push(secondPart)
        }
    })

    return daysWithMeetings
}

export function minuteDiffBetweenTime(startTimestamp, endTimestamp) {
    // Convert timestamps to Date objects
    const startTime = new Date(startTimestamp)
    const endTime = new Date(endTimestamp)

    // Calculate the difference in milliseconds
    const timeDifference = endTime - startTime

    // Convert time difference to minutes
    const timeDifferenceMinutes = timeDifference / (1000 * 60)

    return timeDifferenceMinutes
}

// daysBetweenDates calculates the number of days between two dates
export function daysBetweenDates(startDate, endDate) {
    const diffTime = Math.abs(endDate - startDate)
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    return diffDays
}

export function formatDateRange(dateRange) {
    const initialDate = dateRange.initialDate
    const finalDate = dateRange.finalDate

    let formattedRange = ''

    // Start date
    if (initialDate.getDate() !== finalDate.getDate()) {
        formattedRange += initialDate.getDate() + ' '
    }

    if (initialDate.getMonth() !== finalDate.getMonth()) {
        formattedRange += shortMonthNames[initialDate.getMonth()] + ' '
    }

    if (initialDate.getFullYear() !== finalDate.getFullYear()) {
        formattedRange += `${initialDate.getFullYear()}, `
    }

    // Separator
    if (initialDate.getDate() !== finalDate.getDate()) {
        formattedRange += '- '
    }

    // End date
    formattedRange += `${finalDate.getDate()} ${shortMonthNames[finalDate.getMonth()]} ${finalDate.getFullYear()}`

    return formattedRange
}

export function getDateString(date) {
    let month = date.getMonth() + 1
    let monthStr
    if (month < 10) {
        monthStr = '0' + month
    } else {
        monthStr = `${month}`
    }

    let day = date.getDate()

    let dateStr
    if (day < 10) {
        dateStr = '0' + day
    } else {
        dateStr = `${day}`
    }

    return `${date.getFullYear()}-${monthStr}-${dateStr}`
}

// getFromToTime return the time range of a meeting in the format "HH:MM - HH:MM"
export function getFromToTime(start, end) {
    const startTime = new Date(start)
    const endTime = new Date(end)

    let startHour = startTime.getHours()
    if (startHour < 10) {
        startHour = '0' + startHour
    } else {
        startHour = `${startHour}`
    }

    let startMinute = startTime.getMinutes()
    if (startMinute < 10) {
        startMinute = '0' + startMinute
    } else {
        startMinute = `${startMinute}`
    }

    let endHour = endTime.getHours()
    if (endHour < 10) {
        endHour = '0' + endHour
    } else {
        endHour = `${endHour}`
    }

    let endMinute = endTime.getMinutes()
    if (endMinute < 10) {
        endMinute = '0' + endMinute
    } else {
        endMinute = `${endMinute}`
    }

    return `${startHour}:${startMinute} - ${endHour}:${endMinute}`
}

export function formatTimeForDetails(start, end) {
    const startDate = new Date(start)
    const endDate = new Date(end)
    const isSameDay = startDate.toDateString() === endDate.toDateString()

    const startDayName = dayNames[startDate.getDay()]
    const startMonth = shortMonthNames[startDate.getMonth()]
    const startDay = startDate.getDate()
    const startYear = startDate.getFullYear()
    const startHours = startDate.getHours()
    const startMinutes = startDate.getMinutes()
    const startTimeString = `${startHours.toString().padStart(2, '0')}:${startMinutes.toString().padStart(2, '0')}`

    const endHours = endDate.getHours()
    const endMinutes = endDate.getMinutes()
    const endTimeString = `${endHours.toString().padStart(2, '0')}:${endMinutes.toString().padStart(2, '0')}`

    if (isSameDay) {
        return `${startDayName}, ${startDay} ${startMonth} ${startYear}, ${startTimeString} - ${endTimeString}`
    } else {
        const endDayName = dayNames[endDate.getDay()]
        const endMonth = monthNames[endDate.getMonth()]
        const endDay = endDate.getDate()
        const endYear = endDate.getFullYear()

        return `${startDayName}, ${startDay} ${startMonth} ${startYear}, ${startTimeString} - ${endDayName}, ${endDay} ${endMonth} ${endYear}, ${endTimeString}`
    }
}

export const nowForDateInput = () => {
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const day = String(now.getDate()).padStart(2, '0')
    const hours = String(now.getHours()).padStart(2, '0')
    const minutes = String(now.getMinutes()).padStart(2, '0')

    const formattedDateTime = `${year}-${month}-${day}T${hours}:${minutes}`

    return formattedDateTime
}

export const nowForDateInputWithDelay = (delay) => {
    let now = new Date()

    now.setMinutes(now.getMinutes() + delay)

    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const day = String(now.getDate()).padStart(2, '0')
    const hours = String(now.getHours()).padStart(2, '0')
    const minutes = String(now.getMinutes()).padStart(2, '0')

    const formattedDateTime = `${year}-${month}-${day}T${hours}:${minutes}`

    return formattedDateTime
}

export const getDateForDateInput = (date, time) => {
    if (date === null) {
        return nowForDateInput()
    }

    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')

    const formattedDate = `${year}-${month}-${day}T${time}`

    return formattedDate
}

export const getDateForUpdateInput = (timestamp) => {
    const date = new Date(timestamp)

    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')

    const formattedDate = `${year}-${month}-${day}T${hours}:${minutes}`

    return formattedDate
}

export function wasBeforeNow(date) {
    const dateObj = new Date(date)

    const now = new Date()

    return dateObj < now
}
