import React from "react";
import { NavLink } from "react-router-dom";

import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import './Header.css';
import { useGlobalAuthContext } from "../../helpers/authContext";

export default function Header() {
    const authContext = useGlobalAuthContext()
    const profile = authContext.state?.profile
    return (
        <Box sx={{ flexGrow: 1 }}>
        <AppBar position="static">
            <Toolbar>
            <nav>
                <NavLink to="/">Home</NavLink>
                <NavLink to="/record">Record</NavLink>
                <NavLink to="/about">About</NavLink>
                <span className="spacer" />
                <NavLink to="/login">
                    {profile ? (
                        <div className="account">
                            { profile.url ? (<img src={profile.url!!} alt="icon"></img>) : "" }
                            Account
                        </div>
                    ) : (
                        "Login"
                    )}
                </NavLink>
            </nav>
            </Toolbar>
        </AppBar>
        </Box>
    );
}