import { Request, Response } from 'express';
import { validateToken } from '../services/jwt';

export default async function handlerAuthenticated(req: Request, res: Response) {
    if (await validateToken(req)) {
        res.send(`Token validation successful`);
    } else {
        res.send(`Token validation failed; Not authenticated`);
    }
}