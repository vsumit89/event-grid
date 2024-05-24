'use client'

export function Input({
    name,
    label,
    type = 'text',
    placeholder = '',
    value = '',
    onChange = () => {},
    labelColor = 'text-primary-label',
    step,
    min,
    max,
}) {
    return (
        <div className="flex flex-col gap-2">
            <label htmlFor={name} className={`${labelColor} text-sm`}>
                {label}
            </label>
            <input
                name={name}
                type={type}
                placeholder={placeholder}
                value={value}
                onChange={onChange}
                className="border border-primary-border rounded-lg px-4 py-2 focus:outline-none focus:border-secondary-background bg-primary-background text-sm text-primary-label"
                step={step}
                min={min}
                max={max}
            />
        </div>
    )
}
