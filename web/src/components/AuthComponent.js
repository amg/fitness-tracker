import React, { useEffect } from 'react';
import { googleLogout, useGoogleLogin } from '@react-oauth/google';
import Button from '@mui/material/Button';

const apiBaseUrl = window.env.API_BASE_URL

class AuthStore {
    constructor() { }

    static get name() {
        const json = JSON.parse(localStorage.getItem('auth_user'))
        return json ? json.name : null;
    }
    static set name(value) {
        localStorage.setItem('auth_user', JSON.stringify({ name: value }));
    }
}


function AuthComponent({profile, setProfile, onUnauthenticated}) {
    // Retrieving and parsing the object from LocalStorage
    useEffect(() => {
        setProfile({ name: AuthStore.name });
    }, [setProfile, AuthStore.name]);

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
                    AuthStore.name = data.name
                    setProfile({ name: AuthStore.name });
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
            .then(data => {
                googleLogout();
                setProfile(null);
                AuthStore.name = null
                onUnauthenticated()
            })
            .catch(error => {
                console.error('Error:', error);
            });
    };

    return (
        <div className="popup">
            {profile?.name ? (
                <div>
                    <h3>User Logged in</h3>
                    <p>Name: {profile.name}</p>
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
        </div>
    );
}

export { AuthComponent, AuthStore};