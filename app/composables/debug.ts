import { useGRPCWebsocketTransport } from './grpcws';

export interface DebugContext {
    version: string;
    url: string;
    character: string;
    medium: string;
    resolution: string;
    languageTimezone: string;
    connection: string;
    nui: string;
}

export function getDebugContext(): DebugContext {
    const authStore = useAuthStore();
    const { activeChar } = storeToRefs(authStore);
    const settingsStore = useSettingsStore();
    const { webSocket } = useGRPCWebsocketTransport();
    const { name: browserName, platform: browserPlatform } = getBrowserNameAndPlatform();

    return {
        version: `v${APP_VERSION} / v${settingsStore.version}`,
        url: `${window.location.origin}${window.location.pathname}`,
        character: activeChar.value
            ? `${activeChar.value.userId} (${activeChar.value.job} - ${activeChar.value.jobGrade})`
            : 'N/A',
        medium: `${browserName} on ${browserPlatform}`,
        resolution: `${window.screen.width}x${window.screen.height} / ${window.devicePixelRatio}`,
        languageTimezone: `${navigator.language} / ${Intl.DateTimeFormat().resolvedOptions().timeZone}`,
        connection: `${navigator.onLine ? 'Online' : 'Offline'} / WebSocket ${webSocket.status.value}`,
        nui: `${settingsStore.nuiEnabled ? 'Enabled' : 'Disabled'} (${settingsStore.nuiResourceName ?? 'N/A'})`,
    };
}

export function formatDebugContext(context: DebugContext): string {
    return `- FiveNet/Server Version: \`${context.version}\`
- URL: \`${context.url}\`
- Character: \`${context.character}\`
- Browser/Platform: \`${context.medium}\`
- Resolution/DPR: \`${context.resolution}\`
- Language/Timezone: \`${context.languageTimezone}\`
- Connection: \`${context.connection}\`
- NUI: \`${context.nui}\``;
}

export function collectDebugInfo(): string {
    const context = getDebugContext();

    return `## Debug Info
${formatDebugContext(context)}
`;
}

export function getBrowserNameAndPlatform(): { name: string; platform: string } {
    // Retrieve browser name
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const browserName = (navigator as any).userAgentData
        ? (['Edge', 'Brave', 'Opera', 'Chrome', 'Chromium'].find((n) =>
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              (navigator as any).userAgentData.brands?.some((b: { brand: string; version: string }) => b.brand.includes(n)),
          ) ?? 'Chromium')
        : /Firefox/i.test(navigator.userAgent)
          ? 'Firefox'
          : /Edg/i.test(navigator.userAgent)
            ? 'Edge'
            : /Chrome/i.test(navigator.userAgent)
              ? 'Chrome'
              : /Safari/i.test(navigator.userAgent)
                ? 'Safari'
                : 'Unknown';

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const browserPlatform = (navigator as any).userAgentData?.platform ?? navigator.platform;

    return { name: browserName, platform: browserPlatform };
}
