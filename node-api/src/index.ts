import express, { RequestHandler } from 'express';
import cors from 'cors';
import handlerRoot from './controllers/handlerRoot'
import handlerGetExercises from './controllers/handlerGetExercises'
import handlerUnsupported from './controllers/handlerUnsupported'
import { panic } from './services/error'
import { Request, Response } from 'express';
import sql from './db/db';
import { validateToken } from './services/jwt';
import { Client } from 'pg';
import dns from 'node:dns';

type Handler = (req: Request, res: Response) => Promise<void>;
type SQLHandler = (req: Request, res: Response, client: Client) => Promise<void>;

dns.setDefaultResultOrder('ipv4first');
console.log(`web: ${process.env.WEB_BASE_URL}`);
const corsOptions: cors.CorsOptions = {
  origin: `${process.env.WEB_BASE_URL ?? panic("WEB_BASE_URL is not specified")}`,
  allowedHeaders: ['Content-Type', 'Origin', 'Accept', 'googleTokens'],
  credentials: true,
  optionsSuccessStatus: 200
};
const port = process.env.NODE_API_PORT ?? panic("NODE_API_PORT is not specified");

const app = express();
app.use(cors(corsOptions));
app.use(express.json())

app.options('*',cors(corsOptions));


app.get('/node/exercises', authenticated(dbHandler(handlerGetExercises)));

app.get('/node/*', authenticated(handlerRoot));
app.get('/*', authenticated(handlerUnsupported));
  
app.listen(port, () => {
    console.log(`Server running at port:${port}`);
});

// Middleware functions
function authenticated(next: Handler): Handler {
    return async (req: Request, res: Response) => { 
      if (!await validateToken(req)) {
        res.status(401).send(`Token validation failed; Not authenticated`); 
        return;
      }
      return next(req, res);
    };
}

function dbHandler(sqlHandler: SQLHandler): Handler {
  return async (req: Request, res: Response) => { 
    let sqlInstance = await sql();
    return sqlHandler(req, res, sqlInstance);
  };
}