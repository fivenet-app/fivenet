export function shuffleArray<T>(arr: T[]): T[] {
    for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [arr[i], arr[j]] = [arr[j]!, arr[i]!];
    }

    return arr;
}

export function range(size: number, startAt?: number): number[] {
    return [...Array(size).keys()].map((i) => i + (startAt ?? 0));
}

export function writeUInt32BE(arr: Uint8Array, value: number, offset: number) {
    value = +value;
    offset = offset | 0;
    arr[offset] = value >>> 24;
    arr[offset + 1] = value >>> 16;
    arr[offset + 2] = value >>> 8;
    arr[offset + 3] = value & 0xff;
    return offset + 4;
}
