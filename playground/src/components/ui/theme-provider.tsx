import { createContext, useContext, useEffect, useState } from 'react';

type Theme = 'dark' | 'light';

type ThemeProviderProps = {
	children: React.ReactNode;
};

type ThemeProviderState = {
	theme: Theme;
};

const initialState: ThemeProviderState = {
	theme: 'light',
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({ children, ...props }: ThemeProviderProps) {
	const systemTheme = useSystemTheme();

	useEffect(() => {
		const root = window.document.documentElement;

		root.classList.remove('light', 'dark');

		root.classList.add(systemTheme);
	}, [systemTheme]);

	const value = {
		theme: systemTheme,
	}

	return (
		<ThemeProviderContext.Provider {...props} value={value}>
			{children}
		</ThemeProviderContext.Provider>
	);
}

export const useSystemTheme = () => {
	const [systemTheme, setSystemTheme] = useState<'dark' | 'light'>(
		window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
	);

	useEffect(() => {
		const listener = (e: MediaQueryListEvent) => {
			const systemTheme = e.matches ? 'dark' : 'light';
			setSystemTheme(systemTheme);
		};
		const media = window.matchMedia('(prefers-color-scheme: dark)');
		media.addEventListener('change', listener);
		return () => {
			media.removeEventListener('change', listener);
		};
	}, []);

	return systemTheme;
};

export const useTheme = () => {
	const context = useContext(ThemeProviderContext);

	if (context === undefined) throw new Error('useTheme must be used within a ThemeProvider');

	return context;
};
