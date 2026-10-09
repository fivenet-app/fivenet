export interface StreamErrorDetails {
    timestamp: string;
    code: string;
    message: string;
    cause: string;
    stack: string;
}

function formatErrorValue(value: unknown): string {
    if (value instanceof Error) return value.stack ?? `${value.name}: ${value.message}`;
    if (typeof value === 'string') return value;
    if (value === undefined) return 'N/A';

    try {
        return JSON.stringify(value) ?? String(value);
    } catch {
        return String(value);
    }
}

export function getStreamErrorDetails(error: unknown): StreamErrorDetails {
    const value = error as { code?: unknown; message?: unknown; cause?: unknown; stack?: unknown };

    return {
        timestamp: new Date().toISOString(),
        code: typeof value?.code === 'string' ? value.code : 'N/A',
        message: typeof value?.message === 'string' ? value.message : formatErrorValue(error),
        cause: formatErrorValue(value?.cause),
        stack: typeof value?.stack === 'string' ? value.stack : 'N/A',
    };
}
