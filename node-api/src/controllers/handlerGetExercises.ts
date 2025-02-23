import { Request, Response } from 'express';
import { listDefaultExercises } from '../db/generated/exercises_query_sql';
import { Client } from 'pg';

export default async function handlerGetExercises(req: Request, res: Response, client: Client): Promise<void> {
    await client.connect();
    await listDefaultExercises(client).then((exercises) => {
        res.status(200).send(exercises);
    }).catch((error) => {
        res.status(500).send(error);
    });
    client.end();
}