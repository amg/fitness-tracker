import { useState, useMemo } from 'react';
import { createTheme, ThemeProvider } from '@mui/material';
import AuthContext, { GlobalAuthContext } from '../helpers/authContext'
import HomePage from './home/HomePage';
import AboutPage from './about/AboutPage';
import AccountPage from './account/AccountPage';
import Header from './common/Header';
import { Route, Routes } from 'react-router-dom';
import { Box } from '@mui/material';
import RecordExercisePage from './record/RecordExercisePage';

declare global {
    interface Window {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        env:any;
    }
}

const theme = createTheme({
    palette: {
        primary: {
            main: '#2196F3',
        },
        secondary: {
            main: '#f50057',
        },
    },
    typography: {
        fontFamily: 'Arial, sans-serif',
    },
});


function App() {
    const [state, setState] = useState(AuthContext.getAuthState())
    const authContext = useMemo<AuthContext>(() => { return new AuthContext(state, setState) }, [state, setState]);

    return (
        <ThemeProvider theme={theme}>
            <GlobalAuthContext.Provider value={authContext}>
                <Header />
                <Box m={2}>
                    <Routes>
                        <Route path="/" element={<HomePage />} />
                        <Route path="/record" element={<RecordExercisePage />} />
                        <Route path="/login" element={<AccountPage />} />
                        <Route path="/about" element={<AboutPage />} />
                    </Routes>
                </Box>

            </GlobalAuthContext.Provider>
        </ThemeProvider>
    );
}
export default App;