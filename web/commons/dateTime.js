import { shortMonthNames } from '@/constants/date'

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
        console.log(meeting)
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
                endTime: midnightTime.toLocaleString(),
            }

            const secondPart = {
                ...meeting,
                startTime: midnightTime.toLocaleString(),
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
    let monthStr;
    if (month < 10) {
        monthStr = '0' + month
    } else {
        monthStr = `${month}`
    }

    let day = date.getDate()

    let dateStr;
    if (day < 10) {
        dateStr = '0' + day
    } else {
        dateStr = `${day}`
    }

    console.log({
        date,
        monthStr,
        dateStr,
        year: date.getFullYear()
    })
    return `${date.getFullYear()}-${monthStr}-${dateStr}`
}