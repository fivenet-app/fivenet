import { ref, watch } from 'vue';
import type { File } from '~~/gen/ts/resources/file/file';

export function generateDerefURL(target: string, source = window.location.href): string {
    return (
        '/dereferer?' +
        new URLSearchParams({
            target: target,
            source: source,
        })
    );
}

export function generateDiscordConnectURL(provider: string, redirect?: string, params?: Record<string, string>): string {
    const url = new URL(`https://example.com/api/oauth2/login/${provider}`);

    if (params != undefined) Object.keys(params).forEach((key) => url.searchParams.set(key, params[key]!));

    url.searchParams.set('connect-only', 'true');
    if (redirect !== undefined) {
        url.searchParams.set('redirect', redirect);
    }

    return url.pathname + url.search;
}

export const safeImagePaths = ['/api/image_proxy', '/api/filestore'] as const;

export const brokenImageURL = '/images/broken_image.png' as const;

export function cleanupImageURL(path: string | File | undefined, fallback?: string | undefined): string | undefined {
    if (path === undefined) return fallback;

    const resolvedPath = typeof path === 'object' ? path.filePath : path;

    if (resolvedPath.startsWith('data:image') || resolvedPath.startsWith('/images')) {
        return resolvedPath;
    } else if (safeImagePaths.some((safePath) => resolvedPath.startsWith(safePath))) {
        const correctedPath = safeImagePaths.find((safePath) => resolvedPath.startsWith(safePath));
        if (correctedPath) {
            const remainingPath = resolvedPath.slice(correctedPath.length).replace(/^\//, '');
            return `${correctedPath}/${remainingPath}`;
        } else {
            return resolvedPath;
        }
    } else if (resolvedPath.startsWith('http')) {
        const url = new URL(resolvedPath);
        const isSameHost = url.host === window.location.host;
        const isServedPath = safeImagePaths.some((path) => url.pathname.startsWith(path + '/'));
        if (isSameHost && isServedPath) {
            return url.pathname.replace(/(?<!:)\/\//, '/');
        }

        return `/api/image_proxy/${encodeURIComponent(resolvedPath)}`;
    }

    return `/api/filestore/${resolvedPath.replace(/^\//, '')}`;
}

export function useImageURL(
    filePath: string | File | undefined | Ref<string | File | undefined>,
    fallback?: string | undefined,
): Ref<string | undefined> {
    const imageURL = ref<string | undefined>(undefined);

    watch(
        () => (typeof filePath === 'object' && 'value' in filePath ? filePath.value : filePath),
        (newPath) => {
            imageURL.value = cleanupImageURL(newPath, fallback);
        },
        { immediate: true },
    );

    return imageURL;
}
