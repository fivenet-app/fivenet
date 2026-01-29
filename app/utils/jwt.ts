export async function parseJWTPayload<T>(token: string): Promise<T | undefined> {
    const parts = token.split('.');
    if (parts.length !== 3) {
        throw new Error('Invalid JWT token');
    }

    try {
        const parsed = JSON.parse(atob(parts[1]!)) as T;
        return parsed;
    } catch {
        throw new Error('Invalid JWT payload encoding');
    }
}
