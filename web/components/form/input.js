'use client'

import { Eye, EyeSlash } from '@phosphor-icons/react'
import { useState } from 'react'

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
    error = '',
}) {
    const [showPassword, setShowPassword] = useState(false)

    return (
        <div className="flex flex-col gap-1">
            <label htmlFor={name} className={`${labelColor} text-sm`}>
                {label}
            </label>
            <div className="w-full">
                {type === 'password' ? (
                    <div className="flex items-center border border-primary-border rounded-lg px-4 py-2 focus:outline-none focus:border-white focus-within:border-white">
                        <input
                            name={name}
                            type={showPassword ? 'text' : 'password'}
                            placeholder={placeholder}
                            value={value}
                            onChange={onChange}
                            className="bg-primary-background text-sm text-primary-label w-full flex-1 focus:outline-none"
                            step={step}
                            min={min}
                            max={max}
                        />
                        <div>
                            {showPassword ? (
                                <Eye
                                    className="text-primary-text text-xl cursor-pointer"
                                    onClick={() =>
                                        setShowPassword(!showPassword)
                                    }
                                />
                            ) : (
                                <EyeSlash
                                    className="text-primary-text text-xl cursor-pointer"
                                    onClick={() =>
                                        setShowPassword(!showPassword)
                                    }
                                />
                            )}
                        </div>
                    </div>
                ) : (
                    <input
                        name={name}
                        type={type}
                        placeholder={placeholder}
                        value={value}
                        onChange={onChange}
                        className="border border-primary-border rounded-lg px-4 py-2 focus:outline-none focus:border-white bg-primary-background text-sm text-primary-label w-full"
                        step={step}
                        min={min}
                        max={max}
                    />
                )}
            </div>
            {error && <span className="text-error-text text-xs">{error}</span>}
        </div>
    )
}
