export interface NormalizedError {
    message: string;
    statusCode: string;
    stack: string;
}

export function normalizeError(error: unknown): NormalizedError {
    if (error instanceof Error) {
        return { message: error.message || error.name || 'Unknown error', statusCode: 'N/A', stack: error.stack ?? 'N/A' };
    }

    if (typeof error === 'object' && error !== null) {
        const value = error as { message?: unknown; statusCode?: unknown; statusMessage?: unknown; stack?: unknown };
        let fallback = 'Unknown error';
        try {
            fallback = JSON.stringify(error) || fallback;
        } catch {
            // Keep the fallback message when the error cannot be serialized.
        }

        return {
            message:
                typeof value.statusMessage === 'string'
                    ? value.statusMessage
                    : typeof value.message === 'string'
                      ? value.message
                      : fallback,
            statusCode:
                typeof value.statusCode === 'number' || typeof value.statusCode === 'string' ? String(value.statusCode) : 'N/A',
            stack: typeof value.stack === 'string' ? value.stack : 'N/A',
        };
    }

    return { message: typeof error === 'string' ? error : 'Unknown error', statusCode: 'N/A', stack: 'N/A' };
}
