import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { Playground } from './components/Playground'
import { Toaster } from './components/ui/sonner'
import { ThemeProvider } from './components/ui/theme-provider'

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<ThemeProvider>
			<Toaster visibleToasts={5}/>
			<Playground />
		</ThemeProvider>	
	</StrictMode>,
)
