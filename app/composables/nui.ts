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

export const isNUIEnabled = createSharedComposable(_isNUIEnabled);

function getParentResourceName(): string {
    return useSettingsStore().nuiResourceName ?? 'fivenet';
}

export const focusNUITargets = ['input', 'textarea', 'select'] as const;

let textInputFocused = false;
let focusUpdateQueued = false;

export function isNUITextTarget(target: EventTarget | null): target is HTMLElement {
    if (!(target instanceof HTMLElement)) return false;

    return (
        focusNUITargets.includes(target.tagName.toLowerCase() as (typeof focusNUITargets)[number]) || target.isContentEditable
    );
}

function updateNUIInputFocus(): void {
    if (focusUpdateQueued) return;

    focusUpdateQueued = true;

    // focusout and focusin are emitted as separate events when focus moves
    // between controls. Read activeElement after both events have settled so
    // we only send the final state to FiveM.
    queueMicrotask(() => {
        focusUpdateQueued = false;

        const nextState = isNUITextTarget(document.activeElement);
        if (nextState === textInputFocused) return;

        textInputFocused = nextState;
        logger.debug('text input focus changed:', nextState);
        void focusTablet(nextState);
    });
}

/**
 *
 * @param event FocusEvent `focusin`/`focusout` event
 * @returns void
 */
export function onFocusHandler(event: FocusEvent): void {
    if (!isNUITextTarget(event.target) && !isNUITextTarget(document.activeElement)) return;

    logger.debug('focus handler event:', event.type);
    updateNUIInputFocus();
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
