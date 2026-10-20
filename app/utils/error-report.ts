import { formatClientErrors } from '~/utils/client-errors';
import { normalizeError } from '~/utils/normalize-error';

function escapeInline(value: string): string {
    return value.replaceAll('\\', '\\\\').replaceAll('`', '\\`');
}

function formatStack(stack: string): string {
    return stack
        .split('\n')
        .map((line) => `    ${line}`)
        .join('\n');
}

export function formatErrorReport(error: unknown, options: { context: string; source: string; statusCode?: string }): string {
    const normalized = normalizeError(error);

    return `## FiveNet Error Report

### Technical Information

- Source: \`${escapeInline(options.source)}\`
- Timestamp: \`${new Date().toISOString()}\`
- Error Message: \`${escapeInline(normalized.message)}\`
- Status Code: \`${escapeInline(options.statusCode ?? normalized.statusCode)}\`
${options.context}
${formatClientErrors()}

### Stack Trace

${formatStack(normalized.stack)}`;
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
