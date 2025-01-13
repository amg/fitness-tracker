import React, { useContext } from 'react'

const keys = {
    auth: "auth_user"
}

interface IStateSetter {
    (newState: AuthState | null): void
}

export const GlobalAuthContext = React.createContext<AuthContext | null>(null);
// The `| null` will be removed via the check in the Hook.
export const useGlobalAuthContext = () => {
    const object = useContext(GlobalAuthContext);
    if (!object) { throw new Error("useGlobalAuthContext must be used within a Provider") }
    return object;
  }


class AuthContext {
    constructor(public state: AuthState | null, private setState: IStateSetter){}
    
    static getAuthState(): AuthState | null {
        try {
            const json = localStorage.getItem(keys.auth)
            return json ? JSON.parse(json) : null
        } catch (e) {
            console.log((e as Error).message)
            return null
        }
    }

    setAuthState(value: AuthState | null) {
        this.setState(value) // update state
        localStorage.setItem(keys.auth, JSON.stringify(value));
    }
}

class AuthState {
    constructor(public profile: Profile) {}
    static newState(name: string, url: string | null): AuthState {
        return new AuthState(new Profile(name, url))
    }
}

class Profile {
    constructor(public name: string, public url: string | null){}
}

export { AuthState }
export default AuthContext