/** Embedded image data is unsuitable for persisted or collaborative rich text. */
export function isEmbeddedImageSource(source: unknown): boolean {
    return typeof source === 'string' && source.trimStart().toLowerCase().startsWith('data:image/');
}
