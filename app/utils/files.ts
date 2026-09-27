export function getDisplayNameFromFilePath(filePath?: string): string {
    const name = filePath?.split('/').pop() ?? '';
    return name.replace(/^[0-9]+-/, '');
}

export function getJobAssetDisplayName(asset: { displayName?: string; file?: { filePath?: string } }): string {
    return asset.displayName || getDisplayNameFromFilePath(asset.file?.filePath);
}
