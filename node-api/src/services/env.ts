import fs from 'fs'
import { SecretManagerServiceClient } from '@google-cloud/secret-manager';
import { panic } from './error'
import { crc32 } from 'zlib';
import { fail } from 'assert';
import { createPrivateKey, createPublicKey, KeyObject } from "crypto";

/**
* Use for non sensitive env variables
 */
export interface NonSecureEnv {
	googleClientId: string
	googleClientCallbackUrl: string
	apiPort: string
	webDomain: string
	webBaseUrl: string
	postgresUrl: string
}

/**
* Use for secure variables that should not be exposed in production
 */
export interface SecureEnv {
	googleClientSecret: string
	jwtKeyPrivate: KeyObject
	jwtKeyPublic: KeyObject
	postgresDbName: string
	postgresUser: string
	postgresPassword: string
}

export class Config {
	constructor(public env: NonSecureEnv, public secEnv: SecureEnv){}
}

async function loadFromFile(fileName: string): Promise<string> {
    return new Promise((resolve, reject) => {
        fs.readFile(fileName, (err, value) => {
            if (err) {
                reject(err)
            } else {
                resolve(value.toString('utf8'))
            }
        })
    }) 
}

export async function environment(): Promise<Config> {
    console.log("Getting environment");
    const env: string = process.env.ENV ?? panic("Environment is undefined");
    
    const googleClientId: string = process.env.GOOGLE_CLIENT_ID ?? fail("GOOGLE_CLIENT_ID is undefined");
    const googleClientCallbackUrl: string = process.env.GOOGLE_CLIENT_CALLBACK_URL ?? fail("GOOGLE_CLIENT_CALLBACK_URL is undefined");
    const apiPort: string = process.env.NODE_API_PORT ?? fail("NODE_API_PORT is undefined");
    const webDomain: string = process.env.COOKIE_DOMAIN ?? fail("COOKIE_DOMAIN is undefined");
    const webBaseUrl: string = process.env.WEB_BASE_URL ?? fail("WEB_BASE_URL is undefined");
    let postgresUrl: string;

    switch (env) {
        case "dev":
            postgresUrl = process.env.POSTGRES_URL ?? fail("POSTGRES_URL is undefined");
            const googleClientSecret = process.env.GOOGLE_CLIENT_SECRET ?? fail("GOOGLE_CLIENT_SECRET is undefined");
            const jwtKeyPrivate = process.env.FILE_KEY_PRIVATE ?? fail("FILE_KEY_PRIVATE is undefined");
            const jwtKeyPublic = process.env.FILE_KEY_PUBLIC ?? fail("FILE_KEY_PUBLIC is undefined");
            const postgresDbNameEnv = process.env.POSTGRES_DBNAME ?? fail("POSTGRES_DBNAME is undefined");
            const postgresUser = process.env.POSTGRES_USER ?? fail("POSTGRES_USER is undefined");
            const postgresPassword = process.env.POSTGRES_PASSWORD ?? fail("POSTGRES_PASSWORD is undefined");

            const jwtKeyPrivatePromise = loadFromFile(jwtKeyPrivate ?? fail("FILE_KEY_PRIVATE is undefined"));
            const jwtKeyPublicPromise = loadFromFile(jwtKeyPublic ?? fail("FILE_KEY_PUBLIC is undefined"));
            
            return new Config({
                googleClientId,
                googleClientCallbackUrl,
                apiPort,
                webDomain,
                webBaseUrl,
                postgresUrl,
            }, {
                googleClientSecret,
                jwtKeyPrivate: createPrivateKey(await jwtKeyPrivatePromise),
                jwtKeyPublic: createPublicKey(await jwtKeyPublicPromise),
                postgresDbName: postgresDbNameEnv,
                postgresUser,
                postgresPassword,
            });
            break;
        case "staging": 
            const googleProjectId: string = process.env.GOOGLE_PROJECT_ID ?? fail("GOOGLE_PROJECT_ID is undefined");
            postgresUrl = "" 
            const secretmanagerClient = new SecretManagerServiceClient();
            
            async function accessSecret(key: string) {
                const [version] = await secretmanagerClient.accessSecretVersion({
                    name: `projects/${googleProjectId}/secrets/${key}/versions/latest`,
                });
                const crc32c = crc32(version.payload?.data ?? "")
                if (crc32c != version.payload?.dataCrc32c) {
                    fail(`Data corruption detected retrieving secret ${key}`)
                }
                return version.payload?.data?.toString();
            }
            const privateKeyPromise = accessSecret("FILE_KEY_PRIVATE");
            const publicKeyPromise = accessSecret("FILE_KEY_PUBLIC");
            const postgresDbNamePromise = accessSecret("POSTGRES_DBNAME");
            const postgresDbUserPromise = accessSecret("POSTGRES_USER");
            const postgresDbPwdPromise = accessSecret("POSTGRES_PASSWORD");
            const [privateKey, publicKey, postgresDbName, postgresDbUser, postgresDbPwd] = 
                await Promise.all([privateKeyPromise, publicKeyPromise, postgresDbNamePromise, postgresDbUserPromise, postgresDbPwdPromise]);

            const secEnv = {
                googleClientSecret: "",
                jwtKeyPrivate: createPrivateKey(privateKey ?? fail("FILE_KEY_PRIVATE is undefined")),
                jwtKeyPublic: createPublicKey(publicKey ?? fail("FILE_KEY_PUBLIC is undefined")),
                postgresDbName: postgresDbName ?? fail("POSTGRES_DBNAME is undefined"),
                postgresUser: postgresDbUser ?? fail("POSTGRES_USER is undefined"),
                postgresPassword: postgresDbPwd ?? fail("POSTGRES_PASSWORD is undefined"),
            }

            return new Config({
                googleClientId,
                googleClientCallbackUrl,
                apiPort,
                webDomain,
                webBaseUrl,
                postgresUrl,
            }, secEnv)
            break;
    
        default:
            panic("Environment is incorrectly defined");
    }
}