export function fail(message: string): never {
    console.log(message);
    throw new Error(message);
}

export function panic(message: string): never {
    console.log(message);
    console.log("Stopping the server");
    process.exit(1);
}