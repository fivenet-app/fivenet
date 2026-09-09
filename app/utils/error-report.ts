export interface NormalizedError {
    message: string;
    statusCode: string;
    stack: string;
}

export function normalizeError(error: unknown): NormalizedError {
    if (error instanceof Error) {
        return {
            message: error.message || error.name || 'Unknown error',
            statusCode: 'N/A',
            stack: error.stack ?? 'N/A',
        };
    }

    if (typeof error === 'object' && error !== null) {
        const value = error as {
            message?: unknown;
            statusCode?: unknown;
            statusMessage?: unknown;
            stack?: unknown;
        };

        let fallback = 'Unknown error';
        try {
            fallback = JSON.stringify(error) || fallback;
        } catch {
            // Keep the fallback message when the error cannot be serialized.
        }

        return {
            message:
                typeof value.statusMessage === 'string'
                    ? value.statusMessage
                    : typeof value.message === 'string'
                      ? value.message
                      : fallback,
            statusCode:
                typeof value.statusCode === 'number' || typeof value.statusCode === 'string' ? String(value.statusCode) : 'N/A',
            stack: typeof value.stack === 'string' ? value.stack : 'N/A',
        };
    }

    return {
        message: typeof error === 'string' ? error : 'Unknown error',
        statusCode: 'N/A',
        stack: 'N/A',
    };
}

export function formatErrorReport(error: unknown, options: { context: string; source: string; statusCode?: string }): string {
    const normalized = normalizeError(error);

    return `## FiveNet Error Report

### Technical Information

- Source: \`${options.source}\`
- Timestamp: \`${new Date().toISOString()}\`
- Error Message: \`${normalized.message}\`
- Status Code: \`${options.statusCode ?? normalized.statusCode}\`
${options.context}

### Stack Trace

\`\`\`text
${normalized.stack}
\`\`\``;
}

export function getSafeBrowserDebugContext(): string {
    const browser = /Firefox/i.test(navigator.userAgent)
        ? 'Firefox'
        : /Edg/i.test(navigator.userAgent)
          ? 'Edge'
          : /Chrome/i.test(navigator.userAgent)
            ? 'Chrome'
            : /Safari/i.test(navigator.userAgent)
              ? 'Safari'
              : 'Unknown';

    return `- FiveNet/Server Version: \`${APP_VERSION} / N/A\`
- URL: \`${window.location.origin}${window.location.pathname}\`
- Browser/Platform: \`${browser} / ${navigator.platform}\`
- Resolution/DPR: \`${window.screen.width}x${window.screen.height} / ${window.devicePixelRatio}\`
- Language/Timezone: \`${navigator.language} / ${Intl.DateTimeFormat().resolvedOptions().timeZone}\`
- Connection: \`${navigator.onLine ? 'Online' : 'Offline'}\``;
}
