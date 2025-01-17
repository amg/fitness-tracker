import { Request, Response } from 'express';
// import { Pool } from 'pg';
import { validateToken } from '../services/jwt';

export default async function handlerAuthenticated(req: Request, res: Response) {
    if (await validateToken(req)) {
        res.send(`Token validation successful`);
    } else {
        res.send(`Token validation failed; Not authenticated`);
    }
    // const pool = new Pool({
    //     host: env.POSTGRES_HOST,
    //     user: env.POSTGRES_USER,
    //     password: env.POSTGRES_PASSWORD,
    //     database: env.POSTGRES_DB,
    //     port: env.POSTGRES_PORT,
    //     idleTimeoutMillis: 30000,
    //   });
}