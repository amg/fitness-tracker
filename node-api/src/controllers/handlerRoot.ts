import { Request, Response } from 'express';

export default function handlerRoot(req: Request, res: Response) {
    res.send(`Catch all for anything under /node/*: ${req.path}`);
}