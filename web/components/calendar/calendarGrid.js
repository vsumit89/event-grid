import {
    minuteDiffBetweenTime,
    getDateList,
    splitMeetings,
} from '@/commons/dateTime'
import { dayNames } from '@/constants/date'
import { DayTile } from './dayTile'

export function CalendarGrid({ startDate, endDate, toggleForm }) {
    const dateList = getDateList(startDate, endDate)

    const today = new Date()

    const topToday =
        today.getHours() * 20 * 4 + (today.getMinutes() * 20 * 4) / 60

    const meetings = [
        {
            startTime: '2024-05-25T09:00:00',
            endTime: '2024-05-25T10:30:00',
            title: 'Meeting 1',
        },
        {
            startTime: '2024-05-28T11:00:00',
            endTime: '2024-05-28T12:30:00',
            title: 'Meeting 2',
        },
        {
            startTime: '2024-05-27T14:00:00',
            endTime: '2024-05-27T10:30:00',
            title: 'Meeting 3',
        },
        {
            startTime: '2024-05-26T23:45:00',
            endTime: '2024-05-27T00:15:00',
            title: 'Meeting 4',
        },
    ]

    const daysWithMeetings = splitMeetings(meetings)

    return (
        <div className="flex h-[90vh] overflow-y-scroll w-full">
            <div className="mt-16 pl-2">
                {[...Array(24)].map((_, hour) => (
                    <div
                        key={hour}
                        className="h-20 flex text-white text-xs relative"
                    >
                        <span className="relative top-[-8px]">
                            {hour}:{'00'}
                        </span>
                    </div>
                ))}
            </div>
            <div className="grid grid-cols-7 flex-1 gap-y-0">
                {dateList.map((date) => {
                    let top = 20 * 4 * 3

                    const topD = `${top}px`

                    // const numMinute = 36

                    const stringFormatDate = date.toDateString()
                    return (
                        <div key={date} className="flex flex-col relative">
                            <hr
                                className="border-white border absolute w-full mt-16"
                                style={{
                                    top: topToday,
                                }}
                            ></hr>
                            {daysWithMeetings[stringFormatDate] &&
                                daysWithMeetings[stringFormatDate].map(
                                    (meeting) => {
                                        const meetingDate = new Date(
                                            meeting.startTime
                                        )

                                        const topPos =
                                            meetingDate.getHours() * 20 * 4 +
                                            (meetingDate.getMinutes() *
                                                20 *
                                                4) /
                                                60

                                        const numMinute = minuteDiffBetweenTime(
                                            meeting.startTime,
                                            meeting.endTime
                                        )

                                        return (
                                            <div
                                                key={meeting.startTime}
                                                className={`absolute bg-white text-sm mt-16 w-full text-center rounded-md`}
                                                style={{
                                                    top: topPos,
                                                    height: numMinute * (4 / 3),
                                                }}
                                            >
                                                Test render
                                            </div>
                                        )
                                    }
                                )}
                            <div className="h-16 flex items-center justify-center cursor-pointer text-white text-xs">
                                {dayNames[date.getDay()]} {date.getDate()}
                            </div>
                            {[...Array(24)].map((_, hour) => (
                                <div
                                    key={`${date}-${hour}`}
                                    className={`h-20 border-t border-primary-border cursor-pointer p-1`}
                                >
                                    <DayTile
                                        startTime={hour}
                                        toggleForm={toggleForm}
                                    ></DayTile>
                                </div>
                            ))}
                        </div>
                    )
                })}
            </div>
        </div>
    )
}
