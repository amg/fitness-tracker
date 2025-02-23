import { QueryArrayConfig, QueryArrayResult } from "pg";

interface Client {
    query: (config: QueryArrayConfig) => Promise<QueryArrayResult>;
}

export const getExerciseQuery = `-- name: GetExercise :one
SELECT id, user_id, name, description FROM exercise
WHERE id = $1 LIMIT 1`;

export interface GetExerciseArgs {
    id: string;
}

export interface GetExerciseRow {
    id: string;
    userId: string | null;
    name: string;
    description: string;
}

export async function getExercise(client: Client, args: GetExerciseArgs): Promise<GetExerciseRow | null> {
    const result = await client.query({
        text: getExerciseQuery,
        values: [args.id],
        rowMode: "array"
    });
    if (result.rows.length !== 1) {
        return null;
    }
    const row = result.rows[0];
    return {
        id: row[0],
        userId: row[1],
        name: row[2],
        description: row[3]
    };
}

export const listDefaultExercisesQuery = `-- name: ListDefaultExercises :many
SELECT id, user_id, name, description FROM exercise
WHERE user_id IS NULL
ORDER BY id`;

export interface ListDefaultExercisesRow {
    id: string;
    userId: string | null;
    name: string;
    description: string;
}

export async function listDefaultExercises(client: Client): Promise<ListDefaultExercisesRow[]> {
    const result = await client.query({
        text: listDefaultExercisesQuery,
        values: [],
        rowMode: "array"
    });
    return result.rows.map(row => {
        return {
            id: row[0],
            userId: row[1],
            name: row[2],
            description: row[3]
        };
    });
}

export const listExercisesForUserQuery = `-- name: ListExercisesForUser :many
SELECT id, user_id, name, description FROM exercise
WHERE user_id = $1
ORDER BY id`;

export interface ListExercisesForUserArgs {
    userId: string | null;
}

export interface ListExercisesForUserRow {
    id: string;
    userId: string | null;
    name: string;
    description: string;
}

export async function listExercisesForUser(client: Client, args: ListExercisesForUserArgs): Promise<ListExercisesForUserRow[]> {
    const result = await client.query({
        text: listExercisesForUserQuery,
        values: [args.userId],
        rowMode: "array"
    });
    return result.rows.map(row => {
        return {
            id: row[0],
            userId: row[1],
            name: row[2],
            description: row[3]
        };
    });
}

