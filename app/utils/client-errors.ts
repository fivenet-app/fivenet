import { normalizeError } from '~/utils/normalize-error';

const STORAGE_KEY = 'fivenet.client-errors';
export const MAX_CLIENT_ERRORS = 25;

export interface ClientError {
    timestamp: string;
    source: string;
    route: string;
    message: string;
    statusCode: string;
    stack: string;
}

function canUseStorage(): boolean {
    if (!import.meta.client) return false;
    try {
        return window.sessionStorage !== undefined;
    } catch {
        return false;
    }
}

function isClientError(value: unknown): value is ClientError {
    if (typeof value !== 'object' || value === null) return false;
    const error = value as Partial<ClientError>;
    return (
        typeof error.timestamp === 'string' &&
        typeof error.source === 'string' &&
        typeof error.route === 'string' &&
        typeof error.message === 'string' &&
        typeof error.statusCode === 'string' &&
        typeof error.stack === 'string'
    );
}

function escapeInline(value: string): string {
    return value.replaceAll('\\', '\\\\').replaceAll('`', '\\`');
}

function formatStack(stack: string): string {
    return stack
        .split('\n')
        .map((line) => `    ${line}`)
        .join('\n');
}

export function getClientErrors(): ClientError[] {
    if (!canUseStorage()) return [];
    try {
        const value: unknown = JSON.parse(window.sessionStorage.getItem(STORAGE_KEY) ?? '[]');
        return Array.isArray(value) ? value.filter(isClientError).slice(-MAX_CLIENT_ERRORS) : [];
    } catch {
        return [];
    }
}

export function recordClientError(error: unknown, options: { source: string; route?: string; dedupeWindowMs?: number }): void {
    if (!canUseStorage()) return;
    const normalized = normalizeError(error);
    const entry: ClientError = {
        timestamp: new Date().toISOString(),
        source: options.source,
        route: options.route ?? `${window.location.pathname}${window.location.search}${window.location.hash}`,
        message: normalized.message,
        statusCode: normalized.statusCode,
        stack: normalized.stack,
    };
    try {
        const previous = getClientErrors();
        const last = previous.at(-1);
        const isDuplicate =
            last &&
            last.source === entry.source &&
            last.route === entry.route &&
            last.message === entry.message &&
            Date.now() - new Date(last.timestamp).getTime() < (options.dedupeWindowMs ?? 1000);
        const errors = (isDuplicate ? previous : [...previous, entry]).slice(-MAX_CLIENT_ERRORS);
        window.sessionStorage.setItem(STORAGE_KEY, JSON.stringify(errors));
    } catch {
        // Storage can be unavailable or full. Error reporting must never throw.
    }
}

export function clearClientErrors(): void {
    if (!canUseStorage()) return;
    try {
        window.sessionStorage.removeItem(STORAGE_KEY);
    } catch {
        // Storage can be unavailable even after the capability check.
    }
}

export function formatClientErrors(): string {
    const errors = getClientErrors();
    if (errors.length === 0) return '';
    return `
### Recent Client Errors (${errors.length})

${errors
    .map(
        (error, index) => `#### ${index + 1}. ${escapeInline(error.source)} — ${escapeInline(error.timestamp)}
- Page: \`${escapeInline(error.route)}\`
- Message: \`${escapeInline(error.message)}\`
- Status Code: \`${escapeInline(error.statusCode)}\`

${formatStack(error.stack)}`,
    )
    .join('\n\n')}`;
}
