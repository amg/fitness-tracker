import { Request, Response } from 'express';

export default async function handlerUnsupported(req: Request, res: Response): Promise<void> {
    res.send(`This is unsupported path: ${req.path}`);
}