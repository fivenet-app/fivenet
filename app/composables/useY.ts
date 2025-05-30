// Extended Vue ↔ Yjs array & object helpers – **complete file**
// ---------------------------------------------------------------------------
//  This module complements `useYPrimitives.ts` by adding higher‑level helpers
//  that bind Vue arrays / nested objects to Yjs containers.  All helpers use
//  the unified fallback‑to‑local initialisation controlled through
//  `YjsSyncOptions` (authoritative client + provider with `synced` event).
// ---------------------------------------------------------------------------

import { getCurrentInstance, onUnmounted, ref, toRef, watch, type Ref } from 'vue';
import * as Y from 'yjs';
import { LOCAL_ORIGIN, useYBoolean, useYNumber, useYText, type YjsSyncOptions } from './useYPrimitives';

// ---------------------------------------------------------------------------
//  Primitive helpers & types
// ---------------------------------------------------------------------------

/** Only string, boolean or number are supported as “primitive” Y.Map values. */
export type Primitive = string | boolean | number;

/** A flat record or nested structure that can be mirrored into Yjs. */
export interface YStateMap {
    [key: string]: Primitive | YStateMap | Array<Primitive | YStateMap>;
}

/** Union for arrays that may contain primitives or nested Y.Map objects. */
export type MaybeMap = Primitive | Y.Map<unknown>;

// ---------------------------------------------------------------------------
//  useYArray  ⇄  Y.Array ⇄ Vue Ref<T[]>
// ---------------------------------------------------------------------------

/**
 * Bind a Vue `Ref<T[]>` to a Yjs `Y.Array` that contains primitives or
 * objects.  Deep objects are handled recursively via `useYSyncStructure`.
 *
 * Fallback logic:
 *   • If remote already contains items → remote wins.
 *   • Otherwise, if `opts.authoritative === true` → local array is serialised
 *     and pushed to the Yjs array.
 */
export function useYArray<T extends object | Primitive>(
    yArray: Y.Array<MaybeMap>,
    itemsRef?: Ref<T[]>,
    opts: YjsSyncOptions = {},
): Ref<T[]> {
    const { provider } = opts;
    const items = (itemsRef ?? ref([])) as Ref<T[]>;
    let remoteApplying = false;

    const isPrimitiveMode = () => {
        const first = yArray.length ? yArray.get(0) : undefined;
        return typeof first === 'string' || typeof first === 'number' || typeof first === 'boolean';
    };

    const syncFromY = () => {
        remoteApplying = true;
        if (isPrimitiveMode()) {
            items.value = yArray.toArray() as T[];
        } else {
            while (items.value.length > yArray.length) items.value.pop();
            while (items.value.length < yArray.length) items.value.push({} as T);
            yArray.forEach((m, i) => {
                if (m instanceof Y.Map) {
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    useYSyncStructure(m, items.value[i]! as any);
                }
            });
        }
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const writeLocalToY = () => {
        yArray.doc?.transact(() => {
            yArray.delete(0, yArray.length);
            if (isPrimitiveMode()) {
                yArray.insert(0, items.value as unknown as MaybeMap[]);
            } else {
                const maps = (items.value as unknown as object[]).map((o) => {
                    const m = new Y.Map(Object.entries(o));
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    useYSyncStructure(m, o as any);
                    return m;
                });
                yArray.insert(0, maps as MaybeMap[]);
            }
        }, LOCAL_ORIGIN);
    };

    const init = () => {
        remoteApplying = true;

        if (yArray.length > 0) {
            syncFromY();
        } else if (provider && provider.isAuthoritative) {
            writeLocalToY();
        }

        nextTick(() => {
            remoteApplying = false;
        });

        yArray.observe(handleYArray);
        if (getCurrentInstance()) onUnmounted(() => yArray.unobserve(handleYArray));

        watch(
            items,
            () => {
                if (remoteApplying) return;
                writeLocalToY();
            },
            { deep: true, flush: 'post' },
        );
    };

    if (provider) {
        provider.once('sync', (s) => s && init());
    } else {
        init();
    }

    const handleYArray = (_evt: Y.YArrayEvent<MaybeMap>, tr: Y.Transaction) => {
        if (tr.origin === LOCAL_ORIGIN) return;
        syncFromY();
    };

    return items;
}

// ---------------------------------------------------------------------------
//  useYArrayFiltered  ⇄  object array with key filtering
// ---------------------------------------------------------------------------

export type OptsKeyFilter = { omit?: string[]; only?: string[] };

export function useYArrayFiltered<T extends object>(
    yArr: Y.Array<unknown>,
    items: Ref<T[]>,
    filter: OptsKeyFilter,
    opts: YjsSyncOptions = {},
) {
    const { provider } = opts;
    let remoteApplying = false;

    const serialize = (o: T): Record<string, unknown> => {
        if (filter.only) return Object.fromEntries(Object.entries(o).filter(([k]) => filter.only!.includes(k)));
        if (filter.omit) return Object.fromEntries(Object.entries(o).filter(([k]) => !filter.omit!.includes(k)));
        return o as Record<string, unknown>;
    };

    const hydrate = (target: T, map: Y.Map<unknown>) => {
        Object.entries(map.toJSON() as Record<string, unknown>).forEach(([k, v]) => {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            (target as any)[k] = v;
        });
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        useYSyncStructure(map, target as any);
    };

    const syncFromY = () => {
        remoteApplying = true;
        while (items.value.length > yArr.length) items.value.pop();
        while (items.value.length < yArr.length) items.value.push({} as T);
        (yArr as Y.Array<Y.Map<unknown>>).forEach((raw, i) => {
            if (raw instanceof Y.Map) hydrate(items.value[i]!, raw);
        });
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const writeLocalToY = () => {
        yArr.doc?.transact(() => {
            yArr.delete(0, yArr.length);
            const maps = items.value.map((o) => new Y.Map(Object.entries(serialize(o))));
            yArr.insert(0, maps);
        }, LOCAL_ORIGIN);
    };

    const handleYArr = (_evt: Y.YArrayEvent<Y.Map<unknown>>, tr: Y.Transaction) => {
        if (tr.origin === LOCAL_ORIGIN) return;
        syncFromY();
    };

    const init = () => {
        remoteApplying = true;

        console.log('useYArrayFiltered init', yArr.length, items.value.length);
        if (provider && provider.isAuthoritative) {
            writeLocalToY();
        } else {
            syncFromY();
        }

        nextTick(() => {
            remoteApplying = false;
        });

        (yArr as Y.Array<Y.Map<unknown>>).observe(handleYArr);
        if (getCurrentInstance()) onUnmounted(() => (yArr as Y.Array<Y.Map<unknown>>).unobserve(handleYArr));

        watch(
            items,
            () => {
                if (remoteApplying) return;
                writeLocalToY();
            },
            { deep: true, flush: 'post' },
        );
    };

    if (provider) {
        provider.once('sync', (s) => s && init());
    } else {
        init();
    }
}

// ---------------------------------------------------------------------------
//  useYMap  ⇄  flat object ⇄ Y.Map
// ---------------------------------------------------------------------------

export function useYMap<T extends Record<string, Primitive>>(
    yMap: Y.Map<Primitive>,
    state: T,
    fields?: Array<keyof T>,
    opts: YjsSyncOptions = {},
): T {
    const { provider } = opts;
    const keys: Array<keyof T> = fields ?? (Object.keys(state) as Array<keyof T>);
    const stops: Array<() => void> = [];
    let remoteApplying = false;

    const pullRemote = () => {
        remoteApplying = true;
        keys.forEach((k) => {
            const v = yMap.get(k as string);
            if (v !== undefined) state[k] = v as T[typeof k];
        });
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const init = () => {
        remoteApplying = true;

        const hasRemote = keys.some((k) => yMap.get(k as string) !== undefined);
        if (hasRemote) {
            pullRemote();
        } else if (provider && provider.isAuthoritative) {
            yMap.doc?.transact(() => {
                keys.forEach((k) => yMap.set(k as string, state[k]));
            }, LOCAL_ORIGIN);
        }

        nextTick(() => {
            remoteApplying = false;
        });

        // Observe remote map
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        yMap.observe(pullRemote as any);

        // Watch local fields
        keys.forEach((k) => {
            const stop = watch(
                () => state[k],
                (v) => {
                    if (remoteApplying) return;
                    yMap.set(k as string, v);
                },
            );
            stops.push(stop);
        });

        if (getCurrentInstance())
            onUnmounted(() => {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                yMap.unobserve(pullRemote as any);
                stops.forEach((s) => s());
            });
    };

    if (provider) {
        provider.once('sync', (s) => s && init());
    } else {
        init();
    }

    return state;
}

// ---------------------------------------------------------------------------
//  Internal utilities: container creators & recursive binder
// ---------------------------------------------------------------------------

const isDoc = (c: Y.Doc | Y.Map<unknown>): c is Y.Doc => (c as Y.Doc).clientID !== undefined;

const ensureText = (c: Y.Doc | Y.Map<unknown>, key: string): Y.Text => {
    // Always return a Y.Text. If a primitive was stored, upgrade in‑place.
    const replace = (txt: Y.Text) => {
        if (isDoc(c))
            c.getMap('root')?.set?.(key, txt); // doc root edge‑case
        else c.set(key, txt);
        return txt;
    };
    if (isDoc(c)) return c.getText(key);
    const cur = c.get(key);
    if (cur instanceof Y.Text) return cur;
    const t = new Y.Text();
    if (typeof cur === 'string' && cur.length) t.insert(0, cur);
    return replace(t);
};

const ensureArray = (c: Y.Doc | Y.Map<unknown>, key: string): Y.Array<unknown> => {
    if (isDoc(c)) return c.getArray(key);
    const cur = c.get(key);
    if (cur instanceof Y.Array) return cur;
    const arr = new Y.Array();
    c.set(key, arr);
    return arr;
};

const ensureMap = (c: Y.Doc | Y.Map<unknown>, key: string): Y.Map<unknown> => {
    if (isDoc(c)) return c.getMap(key);
    const cur = c.get(key);
    if (cur instanceof Y.Map) return cur;
    const m = new Y.Map();
    c.set(key, m);
    return m;
};

function useYSyncStructure<T extends YStateMap>(
    ycontainer: Y.Doc | Y.Map<unknown>,
    state: T,
    fields?: Array<keyof T>,
    opts: YjsSyncOptions = {},
): void {
    const keys = (fields ?? Object.keys(state)) as Array<keyof T>;
    keys.forEach((k) => {
        const name = k as string;
        const val = state[k];

        if (typeof val === 'string') {
            useYText(ensureText(ycontainer, name), toRef(state, name) as Ref<string>, opts);
            return;
        }
        if (typeof val === 'boolean') {
            const target = isDoc(ycontainer) ? ensureMap(ycontainer, 'primitives') : (ycontainer as Y.Map<unknown>);
            useYBoolean(target, name, toRef(state, name) as Ref<boolean>, opts);
            return;
        }
        if (typeof val === 'number') {
            const target = isDoc(ycontainer) ? ensureMap(ycontainer, 'primitives') : (ycontainer as Y.Map<unknown>);
            useYNumber(target, name, toRef(state, name) as Ref<number>, opts);
            return;
        }

        if (Array.isArray(val)) {
            const yarr = ensureArray(ycontainer, name);
            const isObjs = val.every((v) => v && typeof v === 'object' && !Array.isArray(v));
            if (isObjs) {
                useYArrayOfObjects(yarr as Y.Array<Y.Map<unknown>>, toRef(state, name) as Ref<YStateMap[]>, opts);
            } else {
                useYArray<Primitive>(yarr as Y.Array<MaybeMap>, toRef(state, name) as Ref<Primitive[]>, opts);
            }
            return;
        }

        if (val && typeof val === 'object') {
            const nested = ensureMap(ycontainer, name);
            const isFlat = Object.values(val).every((v) => ['string', 'number', 'boolean'].includes(typeof v));
            if (isFlat) {
                useYMap<Record<string, Primitive>>(
                    nested as Y.Map<Primitive>,
                    val as Record<string, Primitive>,
                    undefined,
                    opts,
                );
            } else {
                useYSyncStructure(nested, val as YStateMap);
            }
        }
    });
}

function useYArrayOfObjects(yarr: Y.Array<Y.Map<unknown>>, list: Ref<YStateMap[]>, opts: YjsSyncOptions = {}) {
    const { provider } = opts;
    let remoteApplying = false;

    const syncFromY = () => {
        remoteApplying = true;
        while (list.value.length > yarr.length) list.value.pop();
        while (list.value.length < yarr.length) list.value.push({});
        list.value.forEach((obj, i) => useYSyncStructure(yarr.get(i)!, obj, undefined, opts));
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const writeLocal = () => {
        yarr.doc?.transact(() => {
            yarr.delete(0, yarr.length);
            const maps = list.value.map((o) => {
                const m = new Y.Map(Object.entries(o));
                useYSyncStructure(m, o, undefined, opts);
                return m;
            });
            yarr.insert(0, maps);
        }, LOCAL_ORIGIN);
    };

    const init = () => {
        remoteApplying = true;

        syncFromY();
        watch(
            list,
            () => {
                if (remoteApplying) return;
                writeLocal();
            },
            { deep: true, flush: 'post' },
        );

        nextTick(() => {
            remoteApplying = false;
        });

        yarr.observe(syncFromY);
        if (getCurrentInstance()) onUnmounted(() => yarr.unobserve(syncFromY));
    };

    if (provider) {
        provider.once('sync', (s) => s && init());
    } else {
        init();
    }
}
