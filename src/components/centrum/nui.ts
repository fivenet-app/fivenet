import { useConfigStore } from '~/store/config';

// Checking for `GetParentResourceName` existance doesn't work (anymore) in FiveM NUI iframes
export function checkForNUI(): boolean {
    //return typeof window.GetParentResourceName !== 'undefined';
    return useConfigStore().clientConfig.NUIEnabled;
}

function getParentResourceName(): string {
    //return (window as any).GetParentResourceName();
    return useConfigStore().clientConfig.NUIResourceName ?? 'fivenet';
}

export async function setWaypoint(x: number, y: number): Promise<void> {
    if (!checkForNUI()) return;

    return await fetchNui('setWaypoint', { x: x, y: y });
}

export async function fetchNui<T = any, V = any>(event: string, data: T): Promise<V> {
    // @ts-ignore FiveM NUI functions
    const resp = await fetch(`https://${getParentResourceName()}/${event}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json; charset=UTF-8',
        },
        body: JSON.stringify(data),
    });

    const parsed = resp.json();
    return parsed as V;
}
