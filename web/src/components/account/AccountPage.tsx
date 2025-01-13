import { googleLogout, useGoogleLogin } from '@react-oauth/google';
import Button from '@mui/material/Button';

import { AuthState } from '../../helpers/authContext'
import { useGlobalAuthContext } from '../../helpers/authContext'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const apiBaseUrl = window.env.API_BASE_URL

function AuthComponent() {
    const authContext = useGlobalAuthContext();
    
    const googleLogin = useGoogleLogin({
        onSuccess: (codeResponse) => {
            // Send the authorization code to the backend server

            // TODO: update this to allow easier customisation
            fetch(apiBaseUrl + '/api/auth/google', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                mode: "cors",
                credentials: "include",
                body: JSON.stringify({ code: codeResponse.code }),
            })
                .then(response => response.json())
                .then(data => {
                    authContext.setAuthState(AuthState.newState(data.name, data.picture))
                })
                .catch(error => {
                    console.error('Error:', error);
                });
        },
        onError: () => {
            // Handle login errors here
            console.error('Google login failed');
        },
        flow: 'auth-code',
    });

    // log out function to log the user out of google and set the profile array to null
    const logOut = () => {
        fetch(apiBaseUrl + '/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            mode: "cors",
            credentials: "include",
            body: JSON.stringify({}),
        })
            .then(response => response.json())
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            .then(data => {
                googleLogout();
                authContext.setAuthState(null)
            })
            .catch(error => {
                console.error('Error:', error);
            });
    };

    return (
        <>
        {authContext.state?.profile?.name ? (
            <div>
                <h3>User Logged in</h3>
                <p>Name: {authContext.state.profile.name}</p>
                <br />
                <br />
                <Button variant="contained" size="small" color="error" onClick={() => logOut()}>Logout</Button>
                <br />
            </div>
        ) : (
            <>
                <Button variant="contained" size="small" color="success" onClick={() => googleLogin()}>Sign in with Google 🚀</Button>
            </>
        )}
        </>
    );
}

export default AuthComponent;