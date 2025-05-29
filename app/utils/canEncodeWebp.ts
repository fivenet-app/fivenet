let cached: boolean | null = null;

export async function canEncodeWebp(): Promise<boolean> {
    if (cached !== null) return cached;

    const canvas = document.createElement('canvas');
    // quickest check: dataURL sniff
    cached = canvas.toDataURL('image/webp').startsWith('data:image/webp');
    return cached;
}
