import { Sql } from "postgres";

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

export async function getExercise(sql: Sql, args: GetExerciseArgs): Promise<GetExerciseRow | null> {
    const rows = await sql.unsafe(getExerciseQuery, [args.id]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0],
        userId: row[1],
        name: row[2],
        description: row[3]
    };
}

export const listDefaultExercisesQuery = `-- name: ListDefaultExercises :many
SELECT id, user_id, name, description FROM exercise
WHERE user_id = NULL
ORDER BY id`;

export interface ListDefaultExercisesRow {
    id: string;
    userId: string | null;
    name: string;
    description: string;
}

export async function listDefaultExercises(sql: Sql): Promise<ListDefaultExercisesRow[]> {
    return (await sql.unsafe(listDefaultExercisesQuery, []).values()).map(row => ({
        id: row[0],
        userId: row[1],
        name: row[2],
        description: row[3]
    }));
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

export async function listExercisesForUser(sql: Sql, args: ListExercisesForUserArgs): Promise<ListExercisesForUserRow[]> {
    return (await sql.unsafe(listExercisesForUserQuery, [args.userId]).values()).map(row => ({
        id: row[0],
        userId: row[1],
        name: row[2],
        description: row[3]
    }));
}

