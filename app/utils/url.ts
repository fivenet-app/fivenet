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

export function useGenerateImageURL(filePath: string | File | undefined | Ref<string | File | undefined>) {
    const imageURL = ref('');

    const cleanupURL = (path: string | File | undefined) => {
        if (path === undefined) {
            imageURL.value = '/images/broken_link.png';
            return;
        }

        const resolvedPath = typeof path === 'object' ? path.filePath : path;

        if (
            !resolvedPath.startsWith('http') &&
            !resolvedPath.startsWith('/images') &&
            !resolvedPath.startsWith('/api/filestore')
        ) {
            imageURL.value = `/api/filestore/${resolvedPath.replace(/^\//, '')}`;
        } else {
            imageURL.value = resolvedPath;
        }
    };

    watch(
        () => (typeof filePath === 'object' && 'value' in filePath ? filePath.value : filePath),
        (newPath) => {
            cleanupURL(newPath);
        },
        { immediate: true },
    );

    return imageURL;
}
