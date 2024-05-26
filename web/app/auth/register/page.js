'use client'

import { isEmail, isValidPassword } from '@/commons/validation'
import { Input } from '@/components/form/input'
import { useDebouncedCallback } from '@/hooks/useDebounce'
import Link from 'next/link'
import { useState } from 'react'

export default function RegistrationPage() {
    const [registrationData, setRegistrationData] = useState({
        email: {
            value: '',
            error: '',
        },
        password: {
            value: '',
            error: '',
        },
        name: {
            value: '',
            error: '',
        },
    })

    const debouncedValidateEmail = useDebouncedCallback((value) => {
        if (isEmail(value)) {
            setRegistrationData({
                ...registrationData,
                email: {
                    value: value,
                    error: '',
                },
            })
        } else {
            setRegistrationData({
                ...registrationData,
                email: {
                    value: value,
                    error: 'please enter a valid email',
                },
            })
        }
    }, 500)

    const debouncedValidatePassword = useDebouncedCallback((value) => {
        if (isValidPassword(value)) {
            setRegistrationData({
                ...registrationData,
                password: {
                    value: value,
                    error: '',
                },
            })
        } else {
            setRegistrationData({
                ...registrationData,
                password: {
                    value: value,
                    error: 'Password must contain at least one uppercase letter, one lowercase letter, one number, one special character, and be at least 10 characters long.',
                },
            })
        }

        setRegistrationData({
            ...registrationData,
            password: {
                value: value,
                error: error,
            },
        })
    }, 500)

    return (
        <div className="flex flex-col gap-2">
            <Input
                label={'Name'}
                labelColor="text-primary-text"
                onChange={(e) => {
                    setRegistrationData({
                        ...registrationData,
                        name: {
                            value: e.target.value,
                            error: '',
                        },
                    })
                }}
                value={registrationData.name.value}
                error={registrationData.name.error}
            />
            <Input
                label={'Email'}
                labelColor="text-primary-text"
                onChange={(e) => {
                    const value = e.target.value
                    setRegistrationData({
                        ...registrationData,
                        email: {
                            value: value,
                            error: '',
                        },
                    })
                    debouncedValidateEmail(value)
                }}
                value={registrationData.email.value}
                error={registrationData.email.error}
            />
            <Input
                label={'Password'}
                labelColor="text-primary-text"
                type={'password'}
                value={registrationData.password.value}
                error={registrationData.password.error}
                onChange={(e) => {
                    const password = e.target.value

                    setRegistrationData({
                        ...registrationData,
                        password: {
                            value: e.target.value,
                            error: '',
                        },
                    })

                    debouncedValidatePassword(password)
                }}
            />
            <div className="mt-2 w-full flex flex-col gap-2">
                <button className="bg-primary-text rounded-lg px-4 py-2 w-full text-sm hover:bg-white hover:text-black transition transform duration-300">
                    Register
                </button>
                <div className="text-sm text-primary-text text-right">
                    Already have an account?{' '}
                    <Link
                        href="/auth/login"
                        className="text-primary-text hover:underline hover:text-white"
                    >
                        Login
                    </Link>
                </div>
            </div>
        </div>
    )
}
