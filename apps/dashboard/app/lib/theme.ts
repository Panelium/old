export interface ThemeColors {
    primary: ColorScale;
    secondary: ColorScale;
    success: ColorScale;
    warning: ColorScale;
    danger: ColorScale;
    background: BackgroundColors;
    dark: DarkModeColors;
}

interface ColorScale {
    DEFAULT: string;
    foreground: string;
    50: string;
    100: string;
    200: string;
    300: string;
    400: string;
    500: string;
    600: string;
    700: string;
    800: string;
    900: string;
    950?: string;
}

interface BackgroundColors {
    DEFAULT: string;
    card: string;
    panel: string;
}

interface DarkModeColors {
    background: BackgroundColors;
    border: string;
    input: string;
    ring: string;
}

interface ThemeOptions {
    colors: Partial<ThemeColors>;
    shadows?: Record<string, string>;
    radius?: Record<string, string>;
    animation?: Record<string, string>;
}

export function createTheme(options: ThemeOptions) {
    return {
        colors: {
            // Primary brand colors - modern vibrant blues
            primary: {
                DEFAULT: "hsl(222, 85%, 45%)",
                foreground: "hsl(0, 0%, 100%)",
                50: "hsl(222, 95%, 97%)",
                100: "hsl(222, 90%, 93%)",
                200: "hsl(222, 85%, 85%)",
                300: "hsl(222, 80%, 75%)",
                400: "hsl(222, 75%, 65%)",
                500: "hsl(222, 85%, 45%)",
                600: "hsl(222, 80%, 40%)",
                700: "hsl(222, 75%, 35%)",
                800: "hsl(222, 70%, 30%)",
                900: "hsl(222, 65%, 25%)",
                950: "hsl(222, 60%, 15%)",
                ...options.colors?.primary,
            },

            // Secondary accent colors - modern purples for accents
            secondary: {
                DEFAULT: "hsl(260, 65%, 55%)",
                foreground: "hsl(0, 0%, 100%)",
                50: "hsl(260, 95%, 97%)",
                100: "hsl(260, 90%, 93%)",
                200: "hsl(260, 85%, 85%)",
                300: "hsl(260, 80%, 75%)",
                400: "hsl(260, 75%, 65%)",
                500: "hsl(260, 65%, 55%)",
                600: "hsl(260, 60%, 45%)",
                700: "hsl(260, 55%, 40%)",
                800: "hsl(260, 50%, 35%)",
                900: "hsl(260, 45%, 25%)",
                950: "hsl(260, 40%, 15%)",
                ...options.colors?.secondary,
            },

            // Status colors
            success: {
                DEFAULT: "hsl(142, 70%, 45%)",
                foreground: "hsl(0, 0%, 100%)",
                50: "hsl(142, 95%, 97%)",
                100: "hsl(142, 90%, 93%)",
                200: "hsl(142, 85%, 85%)",
                300: "hsl(142, 80%, 75%)",
                400: "hsl(142, 75%, 65%)",
                500: "hsl(142, 70%, 45%)",
                600: "hsl(142, 65%, 40%)",
                700: "hsl(142, 60%, 35%)",
                800: "hsl(142, 55%, 30%)",
                900: "hsl(142, 50%, 25%)",
                ...options.colors?.success,
            },

            warning: {
                DEFAULT: "hsl(38, 95%, 50%)",
                foreground: "hsl(0, 0%, 100%)",
                50: "hsl(38, 95%, 97%)",
                100: "hsl(38, 90%, 93%)",
                200: "hsl(38, 85%, 85%)",
                300: "hsl(38, 80%, 75%)",
                400: "hsl(38, 75%, 65%)",
                500: "hsl(38, 95%, 50%)",
                600: "hsl(38, 90%, 45%)",
                700: "hsl(38, 85%, 40%)",
                800: "hsl(38, 80%, 35%)",
                900: "hsl(38, 75%, 30%)",
                ...options.colors?.warning,
            },

            danger: {
                DEFAULT: "hsl(358, 85%, 55%)",
                foreground: "hsl(0, 0%, 100%)",
                50: "hsl(358, 95%, 97%)",
                100: "hsl(358, 90%, 93%)",
                200: "hsl(358, 85%, 85%)",
                300: "hsl(358, 80%, 75%)",
                400: "hsl(358, 75%, 65%)",
                500: "hsl(358, 85%, 55%)",
                600: "hsl(358, 80%, 50%)",
                700: "hsl(358, 75%, 45%)",
                800: "hsl(358, 70%, 40%)",
                900: "hsl(358, 65%, 35%)",
                ...options.colors?.danger,
            },

            // Enhanced background colors
            background: {
                DEFAULT: "hsl(220, 15%, 98%)",
                card: "hsl(0, 0%, 100%)",
                panel: "hsl(220, 25%, 97%)",
                ...options.colors?.background,
            },

            // Dark mode colors
            dark: {
                background: {
                    DEFAULT: "hsl(230, 15%, 12%)",
                    card: "hsl(230, 15%, 15%)",
                    panel: "hsl(230, 15%, 18%)",
                    ...options.colors?.dark?.background,
                },
                border: "hsl(230, 10%, 22%)",
                input: "hsl(230, 15%, 20%)",
                ring: "hsl(230, 15%, 28%)",
                ...options.colors?.dark,
            },
        },

        // Enhanced shadows for depth and elevation
        shadows: {
            sm: "0 1px 2px 0 rgba(0, 0, 0, 0.05)",
            DEFAULT: "0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)",
            md: "0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)",
            lg: "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)",
            xl: "0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)",
            "2xl": "0 25px 50px -12px rgba(0, 0, 0, 0.25)",
            inner: "inset 0 2px 4px 0 rgba(0, 0, 0, 0.05)",
            soft: "0 4px 20px rgba(0, 0, 0, 0.08)",
            glow: "0 0 15px rgba(66, 153, 225, 0.5)",
            ...options.shadows,
        },

        // Refined border radius
        radius: {
            none: "0",
            sm: "0.125rem",
            DEFAULT: "0.25rem",
            md: "0.375rem",
            lg: "0.5rem",
            xl: "0.75rem",
            "2xl": "1rem",
            "3xl": "1.5rem",
            full: "9999px",
            ...options.radius,
        },

        // Animation properties
        animation: {
            DEFAULT: "200ms cubic-bezier(0.4, 0, 0.2, 1)",
            fast: "150ms cubic-bezier(0.4, 0, 0.2, 1)",
            slow: "300ms cubic-bezier(0.4, 0, 0.2, 1)",
            elastic: "300ms cubic-bezier(0.68, -0.55, 0.27, 1.55)",
            bounce: "300ms cubic-bezier(0.34, 1.56, 0.64, 1)",
            ...options.animation,
        },
    };
}

// Create default theme
const theme = createTheme({
    colors: {
        // Override with your custom colors if needed
    }
});

export default theme;
