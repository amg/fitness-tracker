import { Request, Response } from 'express';

export default function handlerAuthenticated(req: Request, res: Response) {
    res.send(`This is authenticated call: ${req.path}`);
}