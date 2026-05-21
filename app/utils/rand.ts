export function randomNumber(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1) + min);
}

export function randomUUID(): string {
    if (window.isSecureContext) {
        return self.crypto.randomUUID();
    }
    return randomNumber(10000000, 99999999).toString();
}

export function shuffleArray<T>(arr: T[]): T[] {
    for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [arr[i], arr[j]] = [arr[j]!, arr[i]!];
    }

    return arr;
}
