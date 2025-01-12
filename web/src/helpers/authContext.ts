const keys = {
    auth: "auth_user"
}

interface IStateSetter {
    (newState: AuthState | null): void
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
    static newState(name: string): AuthState {
        return new AuthState(new Profile(name))
    }
}

class Profile {
    constructor(public name: string){}
}

export { AuthState }
export default AuthContext