import React, { useState, useCallback, useMemo, useContext } from 'react';
import { createTheme, ThemeProvider } from '@mui/material';
import AuthComponent from './components/AuthComponent'
import Button from '@mui/material/Button';
import AuthContext from './helpers/authContext'

declare global {
    interface Window {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        env:any;
    }
}

export const GlobalAuthContext = React.createContext<AuthContext | null>(null);
// The `| null` will be removed via the check in the Hook.
export const useGlobalAuthContext = () => {
    const object = useContext(GlobalAuthContext);
    if (!object) { throw new Error("useGlobalAuthContext must be used within a Provider") }
    return object;
  }

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
    const [state, setState] = useState(AuthContext.getAuthState())
    const authContext = useMemo<AuthContext>(() => { return new AuthContext(state, setState) }, [state, setState]);
       
    const [data, setData] = useState<string | null>("");

    const setDataWrapper = useCallback((newData: string |  null) => {
        setData(newData);
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
                if (response.status === 401) {
                    authContext.setAuthState(null)
                    throw new Error(`Response status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setData(JSON.stringify(data));
            })
            .catch(error => {
                setData(`message: ${error.message}`);
                console.info('Caught error:', error);
            });
    }

    return (
        <ThemeProvider theme={theme}>
            <GlobalAuthContext.Provider value={authContext}>
            <div>
                <h2>React Google Login</h2>
                <br />
                <br />
                <AuthComponent onUnauthenticated={setDataWrapper} />
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
            </GlobalAuthContext.Provider>
        </ThemeProvider>
    );
}
export default App;