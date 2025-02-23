import { Request, Response } from 'express';

export default async function handlerRoot(req: Request, res: Response): Promise<void> {
    res.send(`Catch all for anything under /node/*: ${req.path}`);
}