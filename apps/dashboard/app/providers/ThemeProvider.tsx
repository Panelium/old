"use client"

import React, {createContext, useContext, useEffect, useState} from "react"

type Theme = 'light' | 'dark' | 'system'

interface ThemeProviderProps {
    children: React.ReactNode
    defaultTheme?: Theme
    storageKey?: string
}

interface ThemeProviderState {
    theme: Theme
    setTheme: (theme: Theme) => void
}

const initialState: ThemeProviderState = {
    theme: 'system',
    setTheme: () => null,
}

const ThemeProviderContext = createContext<ThemeProviderState>(initialState)

export default function ThemeProvider({
                                          children,
                                          defaultTheme = 'system',
                                          storageKey = 'theme',
                                          ...props
                                      }: ThemeProviderProps) {
    const [theme, setTheme] = useState<Theme>(
        () => {
            if (typeof window !== 'undefined') {
                return (localStorage.getItem(storageKey) as Theme) || defaultTheme
            }
            return defaultTheme
        }
    )

    useEffect(() => {
        const root = window.document.documentElement

        root.classList.remove('light', 'dark')

        if (theme === 'system') {
            const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light'

            root.classList.add(systemTheme)
            return
        }

        root.classList.add(theme)
    }, [theme])

    const value = {
        theme,
        setTheme: (theme: Theme) => {
            if (typeof window !== 'undefined') {
                localStorage.setItem(storageKey, theme)
            }
            setTheme(theme)
        },
    }

    return (
        <ThemeProviderContext.Provider {...props} value={value}>
            {children}
        </ThemeProviderContext.Provider>
    )
}

export const useTheme = () => useContext(ThemeProviderContext)
