import React, { useState, useCallback } from 'react';
import { createTheme, ThemeProvider } from '@mui/material';
import { AuthComponent, AuthStore } from './components/AuthComponent.js'
import Button from '@mui/material/Button';

const apiBaseUrl = window.env.API_BASE_URL
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
    const [profile, setProfile] = useState(null);
    const [data, setData] = useState(null);

    const childSetProfile = useCallback((val) => {
        setProfile(val);
    }, [setProfile]);

    const clearData = useCallback((val) => {
        setData({status: "Unauthenticated"});
    }, [setData])

    const authenticatedCall = () => {
        fetch(apiBaseUrl + '/authenticated', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            mode: "cors",
            credentials: "include",
        })
            .then(response => {
                if (response.status == 401) {
                    setProfile(null);
                    AuthStore.name = null;
                    throw new Error(`Response status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setData(data);
            })
            .catch(error => {
                setData({message: error.message});
                console.info('Caught error:', error);
            });
    }

    return (
        <ThemeProvider theme={theme}>
            <div>
                <h2>React Google Login</h2>
                <br />
                <br />
                <AuthComponent profile={profile} setProfile={childSetProfile} onUnauthenticated={clearData} />
                <br />
                <Button variant="contained" size="small" onClick={() => authenticatedCall()}>Make authenticated call</Button>
                <br />
                <br />
                {data ? (
                    <div> {JSON.stringify(data)} </div>
                ) : (
                    <div>No data</div>
                )
                }
            </div>
        </ThemeProvider>
    );
}
export default App;