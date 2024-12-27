import React, { useState } from 'react';
import { googleLogout, useGoogleLogin } from '@react-oauth/google';

const apiBaseUrl = window.env.API_HOST

function App() {
    const [profile, setProfile] = useState(null);
    const [data, setData] = useState(null);

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
                    setProfile(data);
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

    const authenticatedCall = () => {
        fetch(apiBaseUrl + '/authenticated', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            mode: "cors",
            credentials: "include",
        })
            .then(response => response.json())
            .then(response => 
                setData(response)
            )
            .catch(error => {
                setData({error: "Unauthorised"})
                console.error('Error:', error);
            });
    }

    // log out function to log the user out of google and set the profile array to null
    const logOut = () => {
        googleLogout();
        setProfile(null);
    };

    return (
        <div>
            <h2>React Google Login</h2>
            <br />
            <br />
            {profile ? (
                <div>
                    <img src={profile.picture} alt="user image" />
                    <h3>User Logged in</h3>
                    <p>Name: {profile.name}</p>
                    <p>Email Address: {profile.email}</p>
                    <br />
                    <br />
                    <button onClick={logOut}>Log out</button>
                </div>
            ) : (
                <>
                    <button onClick={() => googleLogin()}>Sign in with Google 🚀 </button>
                    <button onClick={() => authenticatedCall()}>Make authenticated call</button>
                </>
            )}
            {data ? (
                <div> {JSON.stringify(data)} </div>
            ) : (
                <div>No data</div>
            )
            }
        </div>
    );
}
export default App;