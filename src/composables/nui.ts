import { useConfigStore } from '~/store/config';

// Checking for `GetParentResourceName` existance doesn't work (anymore) in FiveM NUI iframes
export function isNUIAvailable(): boolean {
    return useConfigStore().isNUIAvailable;
}

function getParentResourceName(): string {
    return useConfigStore().nuiResourceName ?? 'fivenet';
}

export async function fetchNui<T = any, V = any>(event: string, data: T): Promise<V> {
    const body = jsonStringify(data);
    console.debug(`NUI: Fetch ${event}: ${body}`);
    // @ts-ignore FiveM NUI functions
    const resp = await fetch(`https://${getParentResourceName()}/${event}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json; charset=UTF-8',
        },
        body,
    });

    const parsed = resp.json();
    return parsed as V;
}

export async function toggleTablet(state: boolean): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return await fetchNui(state ? 'openTablet' : 'closeTablet', { ok: true });
}

export async function focusTablet(state: boolean): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return await fetchNui('focusTablet', { state });
}

export async function setWaypoint(x: number, y: number): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNui('setWaypoint', { x, y });
}

export async function phoneCallNumber(phoneNumber: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNui('phoneCallNumber', { phoneNumber });
}

export async function copyToClipboard(text: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNui('copyToClipboard', { text });
}

export async function setRadioFrequency(frequency: number): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNui('setRadioFrequency', { frequency });
}

export async function setWaypointPLZ(plz: string): Promise<void> {
    if (!isNUIAvailable()) {
        return;
    }

    return fetchNui('setWaypointPLZ', { plz });
}
