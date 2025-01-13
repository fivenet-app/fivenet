import type { TypedRouteFromName } from '@typed-router';
import { useSettingsStore } from '~/store/settings';

export const logger = useLogger('🎮 NUI');

// Checking for `GetParentResourceName` existance doesn't work (anymore) in FiveM NUI iframes
export function isNUIAvailable(): boolean {
    return useSettingsStore().isNUIAvailable;
}

function getParentResourceName(): string {
    return useSettingsStore().nuiResourceName ?? 'fivenet';
}

const focusNUITargets = ['input', 'textarea'];

/**
 *
 * @param event FocusEvent `focusin`/`focusout` event
 * @returns void promise
 */
export async function onFocusHandler(event: FocusEvent): Promise<void> {
    if (event.target === window) {
        return;
    }

    const element = event.target as HTMLElement;
    if (!focusNUITargets.includes(element.tagName.toLowerCase())) {
        return;
    }
    event.stopPropagation();
    logger.debug('focus handler event:', event.type, element.tagName.toLowerCase());

    focusTablet(event.type === 'focusin');
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
          type: undefined;
      };

export async function onNUIMessage(event: MessageEvent<NUIMessage>): Promise<void> {
    if (event.data.type === 'navigateTo') {
        await navigateTo(event.data.route);
    } else {
        logger.error('Message - Unknown message type received', event.data);
    }
}

// NUI Callbacks

export async function toggleTablet(state: boolean): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return await fetchNUI(state ? 'openTablet' : 'closeTablet', { ok: true });
}

export async function focusTablet(state: boolean): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return await fetchNUI('focusTablet', { state: state });
}

export async function openTokenMgmt(): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return await fetchNUI('openTokenMgmt', { ok: true });
}

export async function setWaypoint(x: number, y: number): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNUI('setWaypoint', { x: x, y: y });
}

export async function phoneCallNumber(phoneNumber: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNUI('phoneCallNumber', { phoneNumber: phoneNumber });
}

export async function copyToClipboard(text: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNUI('copyToClipboard', { text: text });
}

export async function setRadioFrequency(frequency: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNUI('setRadioFrequency', { frequency: frequency });
}

export async function setWaypointPLZ(plz: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNUI('setWaypointPLZ', { plz: plz });
}

export async function openURLInWindow(url: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNUI('openURLInWindow', { url: url });
}
