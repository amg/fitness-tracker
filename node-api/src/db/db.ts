import { Client } from 'pg';
import { environment } from './../services/env'

export default async function client(): Promise<Client> {
    let config = await environment();
    return new Client({
        user: config.secEnv.postgresUser,
        database: config.secEnv.postgresDbName,
        password: config.secEnv.postgresPassword,
        port: 5432,
        host: `/cloudsql/${config.env.postgresUrl!}`,
        ssl: false
    });
}
