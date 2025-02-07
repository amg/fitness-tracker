interface Error {
    message: string
}
export default function GenericError(error: Error) {
    return (
        <>
            <div>Whoops, something went wrong...</div>
            <div>{error.message}</div>
        </>
    );
}