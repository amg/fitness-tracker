import express from 'express';
import cors from 'cors';
import handlerRoot from './controllers/handlerRoot'
import handlerAuthenticated from './controllers/handlerAuthenticated'
import handlerUnsupported from './controllers/handlerUnsupported'
import { panic } from './services/error'

const corsOptions: cors.CorsOptions = {
  origin: `${process.env.WEB_BASE_URL ?? panic("WEB_BASE_URL is not specified")}`
};

const app = express();
app.use(cors(corsOptions));
app.use(express.json())

const port = process.env.NODE_API_PORT ?? panic("NODE_API_PORT is not specified");

app.get('/node/authenticated', handlerAuthenticated);
app.get('/node/*', handlerRoot);
app.get('/*', handlerUnsupported);
  
app.listen(port, () => {
    console.log(`Server running at port:${port}`);
});