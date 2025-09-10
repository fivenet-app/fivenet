export async function blobToBase64(blob: Blob): Promise<string | undefined> {
    const reader = new FileReader();
    return new Promise((resolve) => {
        reader.onload = (ev) => {
            resolve(ev.target?.result ? ev.target.result.toString() : undefined);
        };
        reader.readAsDataURL(blob);
    });
}
