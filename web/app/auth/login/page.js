'use client'

import { isEmail, isValidPassword } from '@/commons/validation'
import { Input } from '@/components/form/input'
import { useDebouncedCallback } from '@/hooks/useDebounce'
import { login } from '@/services/auth'
import { CircleNotch } from '@phosphor-icons/react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'

export default function LoginPage() {
    const [loading, setLoading] = useState(false)

    const router = useRouter()

    const [formError, setFormError] = useState('')

    const [loginData, setLoginData] = useState({
        email: {
            value: '',
            error: '',
        },
        password: {
            value: '',
            error: '',
        },
    })

    const debouncedValidateEmail = useDebouncedCallback((value) => {
        if (isEmail(value)) {
            setLoginData({
                ...loginData,
                email: {
                    value: value,
                    error: '',
                },
            })
        } else {
            setLoginData({
                ...loginData,
                email: {
                    value: value,
                    error: 'please enter a valid email',
                },
            })
        }
    }, 500)

    const debouncedValidatePassword = useDebouncedCallback((value) => {
        if (isValidPassword(value)) {
            setLoginData({
                ...loginData,
                password: {
                    value: value,
                    error: '',
                },
            })
        } else {
            setLoginData({
                ...loginData,
                password: {
                    value: value,
                    error: 'Password must contain at least one uppercase letter, one lowercase letter, one number, one special character, and be at least 10 characters long.',
                },
            })
        }
    }, 500)

    const validateForm = () => {
        if (loginData.email.error || loginData.password.error) {
            return false
        }

        if (!loginData.email.value || !loginData.password.value) {
            if (!loginData.email.value) {
                setLoginData((prev) => {
                    return {
                        ...prev,
                        email: {
                            value: loginData.email.value,
                            error: 'please enter your email',
                        },
                    }
                })
            }

            if (!loginData.password.value) {
                setLoginData((prev) => {
                    return {
                        ...prev,
                        password: {
                            value: loginData.password.value,
                            error: 'please enter your password',
                        },
                    }
                })
            }
            return false
        }

        return true
    }

    const handleLoginSubmit = async () => {
        setLoading(true)

        const isValid = validateForm()

        if (!isValid) {
            setLoading(false)
            return
        }

        try {
            // handles login
            await login(loginData.email.value, loginData.password.value)

            // window.location.href = '/'
            router.push('/')
        } catch (error) {
            setFormError(error.message)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="flex flex-col gap-2">
            <Input
                label={'Email'}
                labelColor="text-primary-text"
                onChange={(e) => {
                    const value = e.target.value

                    setLoginData({
                        ...loginData,
                        email: {
                            value: value,
                            error: '',
                        },
                    })

                    debouncedValidateEmail(value)
                }}
                value={loginData.email.value}
                error={loginData.email.error}
            />
            <Input
                label={'Password'}
                labelColor="text-primary-text"
                type={'password'}
                value={loginData.password.value}
                onChange={(e) => {
                    const value = e.target.value

                    setLoginData({
                        ...loginData,
                        password: {
                            value: e.target.value,
                            error: '',
                        },
                    })

                    // debouncedValidatePassword(value)
                }}
                error={loginData.password.error}
            />
            <div className="mt-2 w-full flex flex-col gap-2">
                <button
                    className="bg-primary-text rounded-lg px-4 py-2 w-full text-sm hover:bg-white hover:text-black transition transform duration-300 flex justify-center"
                    onClick={() => handleLoginSubmit()}
                >
                    {loading ? (
                        <CircleNotch
                            weight="bold"
                            className="w-4 h-4 animate-spin font-bold"
                        />
                    ) : (
                        'Login'
                    )}
                </button>

                <div className="text-sm text-primary-text text-right">
                    Don't have an account?{' '}
                    <Link
                        href="/auth/register"
                        className="text-primary-text hover:underline hover:text-white"
                    >
                        Register
                    </Link>
                </div>
            </div>
            {formError && (
                <div className="text-error-text text-sm">{formError}</div>
            )}
        </div>
    )
}
