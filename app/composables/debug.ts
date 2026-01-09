import { useGRPCWebsocketTransport } from './grpc/grpcws';

export function collectDebugInfo(): string {
    const authStore = useAuthStore();
    const { activeChar, sessionExpiration, attributes, permissions } = storeToRefs(authStore);

    const settingsStore = useSettingsStore();

    const { webSocket } = useGRPCWebsocketTransport();

    const { name: browserName, platform: browserPlatform } = getBrowserNameAndPlatform();

    return `## Debug Info
- Version: ${APP_VERSION} / ${settingsStore.version}
- Access Token Expiration: ${sessionExpiration.value ? sessionExpiration.value.toISOString() : 'N/A'}
- NUI: ${settingsStore.nuiEnabled ? 'Enabled' : 'Disabled'} (${settingsStore.nuiResourceName ?? 'N/A'})
- WebSocket Status: ${webSocket.status.value}
- Active Char ID: ${activeChar.value ? activeChar.value.userId : 'N/A'} (Identifier: ${activeChar.value ? activeChar.value.identifier : 'N/A'})
- Active Char Job: ${activeChar.value ? `${activeChar.value.job} (Rank: ${activeChar.value.jobGrade})` : 'N/A'}
- Permissions: ${permissions.value.length} (Attributes: ${attributes.value.length})

### Browser Info

- Browser: ${browserName}
- Platform: ${browserPlatform}
- Resolution: ${window.screen.width}x${window.screen.height} (Device Pixel Ratio: ${window.devicePixelRatio})
- Language: ${navigator.language}
- Timezone: ${Intl.DateTimeFormat().resolvedOptions().timeZone}
- Online: ${navigator.onLine ? 'Yes' : 'No'}
- Cookies Enabled: ${navigator.cookieEnabled ? 'Yes' : 'No'}
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
