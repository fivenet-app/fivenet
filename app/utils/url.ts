export function generateDerefedURL(target: string): string {
    return (
        '/dereferer?' +
        new URLSearchParams({
            target: target,
            source: window.location.href,
        })
    );
}
