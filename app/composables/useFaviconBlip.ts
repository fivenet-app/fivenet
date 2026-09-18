import { isNUIEnabled } from '~/composables/nui';

const faviconSource = '/favicon.webp';
const faviconBlipSource = '/favicon-blip.webp';
const faviconLinkId = 'fivenet-favicon';
const blipDuration = 900;

let faviconLink: HTMLLinkElement | undefined;
let resetTimer: ReturnType<typeof useTimeoutFn> | undefined;

function getFaviconLink(): HTMLLinkElement {
    if (faviconLink) return faviconLink;

    let link = document.getElementById(faviconLinkId) as HTMLLinkElement | null;
    if (link) {
        faviconLink = link;
        return link;
    }

    link = document.createElement('link');
    link.id = faviconLinkId;
    link.rel = 'icon';
    link.type = 'image/webp';
    link.href = faviconSource;
    document.head.appendChild(link);
    faviconLink = link;
    return link;
}

/** Briefly marks the browser tab for an incoming stream notification. */
export function useFaviconBlip(): void {
    if (import.meta.server || isNUIEnabled().value) return;

    const link = getFaviconLink();
    link.href = faviconBlipSource;

    resetTimer ??= useTimeoutFn(
        () => {
            if (faviconLink) faviconLink.href = faviconSource;
        },
        blipDuration,
        { immediate: false },
    );
    resetTimer.stop();
    resetTimer.start();
}
