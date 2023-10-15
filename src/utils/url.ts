export function marshalObjectToHash(object: Object): string {
    return '#' + new URLSearchParams(object as any).toString();
}

export function unmarshalHashToObject<T>(hash: string): T {
    const params = new URLSearchParams(hash.replace(/^#/, ''));
    const entries = params.entries();
    return Object.fromEntries(entries) as T;
}
