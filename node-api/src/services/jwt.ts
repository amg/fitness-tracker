import { Request } from 'express';
import { environment } from './../services/env'
import jwt from 'jsonwebtoken';

const HEADER_JWT = "jwt_token"

/**
 * Checks that JWT token is present in cookies and returns true if valid
 * @param req Request
 * @returns 
 */
export async function validateToken(req: Request): Promise<boolean> {
    const allCookies = req.headers.cookie?.split("; "); // [name=cookieValue]
    if (!allCookies || allCookies.length < 1) {
        return false;
    }
    let cookies = new Map<string, string>();
    allCookies.forEach((cookie) => {
        const [name, value] = cookie.split("=");
        cookies.set(name, value);
    });
    const jwtCookie = cookies.get(HEADER_JWT);
	if (!jwtCookie) {
		return false;
	}

	const cryptoPublicKey = (await environment()).secEnv.jwtKeyPublic;
    try { 
        jwt.verify(jwtCookie, cryptoPublicKey);
        return true
    } catch (e) { 
        return false
    }
}