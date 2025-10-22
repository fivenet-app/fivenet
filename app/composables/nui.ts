import type { TypedRouteFromName } from '@typed-router';
import mitt from 'mitt';
import { useSettingsStore } from '~/stores/settings';

type NUIEvent<T = string> = {
    type: T;
};

type NUIEvents = {
    openTablet: NUIEvent<'openTablet'> & { state: boolean };
    closeTablet: NUIEvent<'closeTablet'> & { state: boolean };
};

export const nuiEvents = mitt<NUIEvents>();

const logger = useLogger('🎮 NUI');

// Use settings store to see if NUI is enabled
const _isNUIEnabled = (): Ref<boolean> => {
    const settingsStore = useSettingsStore();
    const { nuiEnabled } = storeToRefs(settingsStore);

    return nuiEnabled;
};

const isNUIEnabled = createSharedComposable(_isNUIEnabled);

function getParentResourceName(): string {
    return useSettingsStore().nuiResourceName ?? 'fivenet';
}

const focusNUITargets = ['input', 'textarea'] as const;

/**
 *
 * @param event FocusEvent `focusin`/`focusout` event
 * @returns void promise
 */
export async function onFocusHandler(event: FocusEvent): Promise<void> {
    if (event.target === window) return;

    const element = event.target as HTMLElement;
    if (!focusNUITargets.includes(element.tagName.toLowerCase())) return;
    event.stopPropagation();
    logger.debug('focus handler event:', event.type, element.tagName.toLowerCase());

    await focusTablet(event.type === 'focusin');
}

type NUIRequest = boolean | string | object;
type NUIResponse = boolean | string | object;

export async function fetchNUI<T = NUIRequest, V = NUIResponse>(method: string, data: T): Promise<V> {
    const body = JSON.stringify(data);
    logger.debug(`Fetch ${method}:`, body);
    const resp = await fetch(`https://${getParentResourceName()}/${method}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json; charset=UTF-8',
        },
        body,
    });

    const parsed = resp.json();
    return parsed as V;
}

type NUIMessage =
    | {
          type: 'navigateTo';
          /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
          route: TypedRouteFromName<any>;
      }
    | {
          type: 'openTablet';
          state: boolean;
      }
    | {
          type: undefined;
      };

export async function onNUIMessage(event: MessageEvent<NUIMessage>): Promise<void> {
    if (event.data.type === 'navigateTo') {
        await navigateTo(event.data.route);
    } else if (event.data.type === 'openTablet') {
        nuiEvents.emit('openTablet', { type: 'openTablet', state: event.data.state });
    } else {
        logger.error('Message - Unknown message type received', event.data);
    }
}

// NUI Callbacks

export async function toggleTablet(state: boolean): Promise<void> {
    if (!isNUIEnabled().value) return;

    return await fetchNUI(state ? 'openTablet' : 'closeTablet', { ok: true });
}

export async function focusTablet(state: boolean): Promise<void> {
    if (!isNUIEnabled().value) return;

    return await fetchNUI('focusTablet', { state: state });
}

export async function openTokenMgmt(): Promise<void> {
    if (!isNUIEnabled().value) return;

    return await fetchNUI('openTokenMgmt', { ok: true });
}

export async function setWaypoint(x: number, y: number): Promise<void> {
    if (!isNUIEnabled().value) return;

    return fetchNUI('setWaypoint', { x: x, y: y });
}

export async function phoneCallNumber(phoneNumber: string): Promise<void> {
    if (!isNUIEnabled().value) return;

    return fetchNUI('phoneCallNumber', { phoneNumber: phoneNumber });
}

export async function copyToClipboard(text: string): Promise<void> {
    if (!isNUIEnabled().value) return;

    return fetchNUI('copyToClipboard', { text: text });
}

export async function setRadioFrequency(frequency: string): Promise<void> {
    if (!isNUIEnabled().value) return;

    return fetchNUI('setRadioFrequency', { frequency: frequency });
}

export async function setWaypointPLZ(plz: string): Promise<void> {
    if (!isNUIEnabled().value) return;

    return fetchNUI('setWaypointPLZ', { plz: plz });
}

export async function openURLInWindow(url: string): Promise<void> {
    if (!isNUIEnabled().value) return;

    return fetchNUI('openURLInWindow', { url: url });
}

export async function setTabletColors(primary: string, gray: string): Promise<void> {
    if (!isNUIEnabled().value) return;

    return fetchNUI('setTabletColors', { primary: primary, gray: gray });
}
