export async function remoteImageURLToBase64Data(url: string): Promise<string | undefined> {
    const resp = await fetch(url).then((r) => r.blob());

    const dataUrl = await blobToBase64(resp);
    if (!dataUrl) {
        return;
    }

    return dataUrl.toString();
}
