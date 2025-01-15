import { Request, Response } from 'express';

export default function handlerUnsupported(req: Request, res: Response) {
    res.send(`This is unsupported path: ${req.path}`);
}