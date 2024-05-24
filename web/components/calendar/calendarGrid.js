import { getDateList } from "@/commons/date"
import { dayNames } from "@/constants/date"

export function CalendarGrid({ startDate, endDate }) {
    const dateList = getDateList(startDate, endDate)
    return (
        <div className="flex h-[90vh] overflow-y-scroll w-full">
            <div className="p-2">
                {[...Array(24)].map((_, hour) => (
                    <div key={hour} className="flex flex-col">
                        {
                            <div
                                key={hour}
                                className="h-12 flex items-center cursor-pointer text-white text-xs"
                            >
                                {hour}:{'00'} AM
                            </div>
                        }
                    </div>
                ))}
            </div>
            <div className="grid grid-cols-7 flex-1">
                {dateList.map((date) => (
                    <div
                        key={date}
                        className="h-12 flex items-center justify-center cursor-pointer text-white text-xs"
                    >
                        {dayNames[date.getDay()]} {date.getDate()} 
                    </div>
                ))}
            </div>
        </div>
    )
}
