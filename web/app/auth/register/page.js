'use client'

import { isEmail, isValidPassword } from '@/commons/validation'
import { Input } from '@/components/form/input'
import { useDebouncedCallback } from '@/hooks/useDebounce'
import { register } from '@/services/auth'
import { CircleNotch } from '@phosphor-icons/react/dist/ssr'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'

export default function RegistrationPage() {
    const router = useRouter()

    const [formError, setFormError] = useState('')

    const [loading, setLoading] = useState(false)

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
    }, 500)

    const validateForm = () => {
        var isValid = true
        if (!registrationData.email) {
            setRegistrationData((prev) => {
                return {
                    ...prev,
                    email: {
                        ...prev.email,
                        error: 'Email is required',
                    },
                }
            })
            isValid = false
        } else if (!isEmail(registrationData.email.value)) {
            setRegistrationData((prev) => {
                return {
                    ...prev,
                    email: {
                        ...prev.email,
                        error: 'Please enter a valid email',
                    },
                }
            })
            isValid = false
        }

        if (!registrationData.password) {
            setRegistrationData((prev) => {
                return {
                    ...prev,
                    password: {
                        ...prev.password,
                        error: 'Password is required',
                    },
                }
            })
            isValid = false
        } else if (!isValidPassword(registrationData.password.value)) {
            setRegistrationData((prev) => {
                return {
                    ...prev,
                    password: {
                        ...prev.password,
                        error: 'Password must contain at least one uppercase letter, one lowercase letter, one number, one special character, and be at least 10 characters long.',
                    },
                }
            })
            isValid = false
        }

        if (registrationData.name.value === '') {
            setRegistrationData((prev) => {
                return {
                    ...prev,
                    name: {
                        ...prev.name,
                        error: 'Name is required',
                    },
                }
            })
            isValid = false
        }

        return isValid
    }

    const handleSubmit = async () => {
        const isValid = validateForm()

        console.log(isValid)
        if (!isValid) {
            return
        }

        setLoading(true)
        try {
            await register({
                email: registrationData.email.value,
                password: registrationData.password.value,
                name: registrationData.name.value,
            })
            router.push('/auth/login')
        } catch (error) {
            setFormError(error.message)
        } finally {
            setLoading(false)
        }
    }

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
                <button
                    className="bg-primary-text rounded-lg px-4 py-2 w-full text-sm hover:bg-white hover:text-black transition transform duration-300 flex justify-center"
                    onClick={() => handleSubmit()}
                >
                    {loading ? (
                        <CircleNotch
                            weight="bold"
                            className="w-4 h-4 animate-spin font-bold"
                        />
                    ) : (
                        'Register'
                    )}
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
            {formError && (
                <div className="text-error-text text-sm">{formError}</div>
            )}
        </div>
    )
}
