/**
 * Shared Vue ↔ Yjs primitive bindings **without feedback-loops**.
 *
 * Each composable spins up a two-way sync between a Vue ref and a Yjs data
 * structure.  To prevent infinite loops we wrap every *local* change in a Yjs
 * transaction tagged with a unique `LOCAL_ORIGIN` symbol, and observers simply
 * ignore events whose `transaction.origin` matches that symbol.
 */

import DMP from 'diff-match-patch';
import type * as Y from 'yjs';
import type GrpcProvider from './yjs/yjs';

const dmp = new DMP.diff_match_patch();

/** Unique marker used to tag local Yjs transactions */
export const LOCAL_ORIGIN = Symbol('vue-yjs-local');

/**
 * Options common to all composables that synchronise a Vue ref with a Yjs data
 * structure.
 */
export interface YjsSyncOptions {
    /** A provider that emits `synced` (WebsocketProvider, HocuspocusProvider …) and has authoritative prop */
    provider?: GrpcProvider;
}

/**
 * Composable to bind a Vue `Ref<string>` (or a ref to a reactive property) to a
 * Yjs `Y.Text` instance. Applies incremental deltas to avoid flicker and syncs
 * local changes.
 *
 * The initial value is chosen as follows:
 *   1. Wait for the provider's `synced` event (if a provider is given).
 *   If the remote text already contains data → remote wins.
 *   3. Otherwise, if `authoritative === true` → seed the remote document with
 *      the local ref's current value.
 *
 * @param yText    The `Y.Text` instance to sync with.
 * @param textRef  Optional Vue `Ref<string>`. Pass e.g. `toRef(obj, 'title')`
 *                 to bind a property. If provided the returned ref ===
 *                 `textRef`, so you can continue using it. If not provided a
 *                 new ref is created.
 * @param opts     {@link YjsSyncOptions}.
 * @returns        A Vue `Ref<string>` that stays in sync with the `Y.Text`
 *                 content (=== `textRef` if supplied).
 */
export function useYText(yText: Y.Text, textRef?: Ref<string>, opts: YjsSyncOptions = {}): Ref<string> {
    const { provider } = opts;
    const text = (textRef ?? ref('')) as Ref<string>;
    let remoteApplying = false;

    // REMOTE → LOCAL
    const handleUpdate = (evt: Y.YTextEvent): void => {
        if (evt.transaction.origin === LOCAL_ORIGIN) return; // ignore our own writes
        remoteApplying = true;
        text.value = applyDelta(text.value, evt.delta as DeltaOp[]);
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const onSync = (s: boolean) => s && init();

    // INITIAL SYNC
    const init = (): void => {
        const remote = yText.toString();
        remoteApplying = true;

        if (provider && provider.isAuthoritative) {
            // remote empty → seed with local value
            yText.doc?.transact(() => {
                yText.delete(0, yText.length);
                yText.insert(0, text.value);
            }, LOCAL_ORIGIN);
        } else {
            text.value = remote; // remote wins
        }

        nextTick(() => {
            remoteApplying = false;
        });

        yText.observe(handleUpdate);
        if (getCurrentInstance())
            onUnmounted(() => {
                provider?.off('sync', onSync);
                yText.unobserve(handleUpdate);
            });

        // LOCAL → REMOTE
        watch(
            text,
            (val, oldVal) => {
                if (remoteApplying || val === oldVal) return;
                const diffs = dmp.diff_main(oldVal, val);
                let pos = 0;
                yText.doc?.transact(() => {
                    for (const [op, data] of diffs) {
                        if (op === 0)
                            pos += data.length; // equal
                        else if (op === -1)
                            yText.delete(pos, data.length); // delete
                        else if (op === 1) {
                            // insert
                            yText.insert(pos, data);
                            pos += data.length;
                        }
                    }
                }, LOCAL_ORIGIN);
            },
            { flush: 'sync' },
        );
    };

    if (provider) {
        provider.once('sync', onSync);
    } else {
        init();
    }

    return text;
}

// Helper for incremental text patches
interface DeltaOp {
    retain?: number;
    delete?: number;
    insert?: string;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    attributes?: Record<string, any>;
}
function applyDelta(str: string, delta: DeltaOp[]): string {
    let res = '';
    let idx = 0;
    for (const op of delta) {
        if (op.retain) {
            res += str.slice(idx, idx + op.retain);
            idx += op.retain;
        }
        if (op.delete) idx += op.delete;
        if (op.insert) res += op.insert;
    }
    return res + str.slice(idx);
}

/**
 * Composable to bind a Vue `Ref<boolean>` (or a ref to a reactive property) to
 * a Yjs `Y.Map` entry. Syncs boolean values without flicker and handles
 * remote/local updates.
 *
 * @param yMap     The `Y.Map` instance holding boolean values.
 * @param key      The key within the `Y.Map` to bind.
 * @param boolRef  Optional Vue `Ref<boolean>`. Pass `toRef(obj, 'flag')` to
 *                 bind a property. If provided the returned ref === `boolRef`.
 * @param opts     {@link YjsSyncOptions}.
 * @returns        A Vue `Ref<boolean>` that stays in sync with the `Y.Map`
 *                 entry (=== `boolRef` if supplied).
 */
export function useYBoolean(
    yMap: Y.Map<unknown>,
    key: string,
    boolRef?: Ref<boolean>,
    opts: YjsSyncOptions = {},
): Ref<boolean> {
    const { provider } = opts;
    const bool = (boolRef ?? ref(false)) as Ref<boolean>;
    let remoteApplying = false;

    const handleUpdate = (evt: Y.YMapEvent<unknown>): void => {
        if (evt.transaction.origin === LOCAL_ORIGIN) return;
        remoteApplying = true;
        const v = yMap.get(key);
        bool.value = typeof v === 'boolean' ? (v as boolean) : false;
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const onSync = (s: boolean) => s && init();

    // INITIAL SYNC
    const init = (): void => {
        const remote = yMap.get(key);
        const hasRemote = typeof remote === 'boolean';
        if (hasRemote) {
            bool.value = remote as boolean;
        } else if (provider && provider.isAuthoritative) {
            yMap.doc?.transact(() => yMap.set(key, bool.value), LOCAL_ORIGIN);
        }

        yMap.observe(handleUpdate);
        if (getCurrentInstance())
            onUnmounted(() => {
                provider?.off('sync', onSync);
                yMap.unobserve(handleUpdate);
            });

        watch(
            bool,
            (val, oldVal) => {
                if (remoteApplying || val === oldVal) return;
                yMap.doc?.transact(() => yMap.set(key, val), LOCAL_ORIGIN);
            },
            { flush: 'sync' },
        );
    };

    if (provider) {
        provider.once('sync', onSync);
    } else {
        init();
    }

    return bool;
}

/**
 * Composable to bind a Vue `Ref<number>` (or a ref to a reactive property) to
 * a Yjs `Y.Map` entry.
 *
 * @param yMap    The `Y.Map` instance holding number values.
 * @param key     The key within the `Y.Map` to bind.
 * @param numRef  Optional Vue `Ref<number>`. Pass `toRef(obj, 'count')` to bind
 *                a property. If provided the returned ref === `numRef`.
 * @param opts    {@link YjsSyncOptions}.
 * @returns       A Vue `Ref<number>` that stays in sync with the `Y.Map` entry
 *                (=== `numRef` if supplied).
 */
export function useYNumber(yMap: Y.Map<unknown>, key: string, numRef?: Ref<number>, opts: YjsSyncOptions = {}): Ref<number> {
    const { provider } = opts;
    const num = (numRef ?? ref(0)) as Ref<number>;
    let remoteApplying = false;

    const handleUpdate = (evt: Y.YMapEvent<unknown>): void => {
        if (evt.transaction.origin === LOCAL_ORIGIN) return;
        remoteApplying = true;
        const v = yMap.get(key);
        num.value = typeof v === 'number' ? (v as number) : 0;
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const onSync = (s: boolean) => s && init();

    // INITIAL SYNC
    const init = (): void => {
        const remote = yMap.get(key);
        const hasRemote = typeof remote === 'number';
        if (hasRemote) {
            num.value = remote as number;
        } else if (provider && provider.isAuthoritative) {
            yMap.doc?.transact(() => yMap.set(key, num.value), LOCAL_ORIGIN);
        }

        yMap.observe(handleUpdate);
        if (getCurrentInstance())
            onUnmounted(() => {
                provider?.off('sync', onSync);
                yMap.unobserve(handleUpdate);
            });

        watch(
            num,
            (val, oldVal) => {
                if (remoteApplying || val === oldVal) return;
                yMap.doc?.transact(() => yMap.set(key, val), LOCAL_ORIGIN);
            },
            { flush: 'sync' },
        );
    };

    if (provider) {
        provider.once('sync', onSync);
    } else {
        init();
    }

    return num;
}
