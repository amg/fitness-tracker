import { Client } from 'pg';
import { environment } from './../services/env'

export default async function client(): Promise<Client> {
    let config = await environment();
    let connectionString = `postgres://${config.secEnv.postgresUser}:${config.secEnv.postgresPassword}@${config.env.postgresUrl}/${config.secEnv.postgresDbName}?sslmode=disable`
    return new Client({
        connectionString: connectionString
    });
}
