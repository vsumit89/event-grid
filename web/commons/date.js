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