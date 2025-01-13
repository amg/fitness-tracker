import { useState, useEffect } from 'react';
import { Button } from "@mui/material"
import {useGlobalAuthContext} from "../../helpers/authContext"

const apiBaseUrl = window.env.API_BASE_URL

function HomePage() {
    const authContext = useGlobalAuthContext();
       
    const [data, setData] = useState<string | null>("");

    useEffect(() => {
        if (!authContext.state?.profile) {
            setData("Logged out");
        }
    }, [authContext.state])

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
    
    return <div>
                <h2>Welcome Home</h2>
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
}

export default HomePage