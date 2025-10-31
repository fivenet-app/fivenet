export function randomNumber(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1) + min);
}

export function randomUUID(): string {
    if (window.isSecureContext) {
        return self.crypto.randomUUID();
    }
    return randomNumber(10000000, 99999999).toString();
}
