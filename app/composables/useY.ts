import { onMounted, onUnmounted, ref, type Ref, watch } from 'vue';
import type { YTextEvent } from 'yjs';
import * as Y from 'yjs';

/**
 * Only string, boolean or number are supported as “primitive” Y-Map values.
 */
type Primitive = string | boolean | number;

/**
 * Recursive state type: can be a primitive, a nested map, or an array of nested maps
 */
type YStateType = Primitive | YStateMap | YStateMap[];

/**
 * A flat record whose values are YStateType (for nested structures)
 */
interface YStateMap {
    [key: string]: YStateType;
}

/**
 * Represents a Yjs delta operation for Text updates.
 */
type DeltaOp = {
    retain?: number;
    delete?: number;
    insert?: string;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    attributes?: Record<string, any>;
};

/**
 * Composable to bind a Vue Ref<string> (or a ref to a reactive property) to a Yjs Y.Text instance.
 * Applies incremental deltas to avoid flicker and syncs local changes.
 *
 * @param yText   - The Y.Text instance to sync with.
 * @param textRef - Optional Vue Ref<string>. Pass e.g. `toRef(reactiveObj, 'title')` to bind a property.
 *                  If provided, this ref will be used directly (the returned ref === textRef), so you can continue using it.
 *                  If not provided, a new ref is created.
 * @returns A Vue Ref<string> that stays in sync with the Y.Text content (=== textRef if supplied).
 */
export function useYText(yText: Y.Text, textRef: Ref<string>): Ref<string> {
    const text = (textRef ?? ref('')) as Ref<string>;

    let remoteApplying = false;
    let initialSync = true;

    /**
     * Applies a Yjs delta to a string, returning the updated text.
     */
    function applyDelta(str: string, delta: DeltaOp[]): string {
        let result = '';
        let index = 0;

        for (const op of delta) {
            if (op.retain) {
                result += str.slice(index, index + op.retain);
                index += op.retain;
            }
            if (op.delete) {
                index += op.delete;
            }
            if (op.insert) {
                result += op.insert;
            }
        }
        // append any trailing text
        result += str.slice(index);
        return result;
    }

    /**
     * Yjs update handler: applies remote deltas to our ref value.
     */
    function handleUpdate(event: YTextEvent): void {
        remoteApplying = true;
        text.value = applyDelta(text.value, event.delta as DeltaOp[]);
        remoteApplying = false;
    }

    onMounted(() => {
        // initialize from current Y.Text state
        text.value = yText.toString();
        initialSync = false;

        // subscribe to Yjs updates
        yText.observe(handleUpdate);
    });

    // watch for local changes and push them into Yjs
    watch(text, (newVal) => {
        if (remoteApplying || initialSync) return;
        // full replace strategy; for large docs consider diffing
        yText.delete(0, yText.length);
        yText.insert(0, newVal);
    });

    onUnmounted(() => {
        yText.unobserve(handleUpdate);
    });

    return text;
}

/**
 * Composable to bind a Vue Ref<boolean> (or a ref to a reactive property) to a Yjs Y.Map entry.
 * Syncs boolean values without flicker and handles remote/local updates.
 *
 * @param yMap     - The Y.Map instance holding boolean values.
 * @param key      - The key within the Y.Map to bind.
 * @param boolRef  - Optional Vue Ref<boolean>. Pass `toRef(reactiveObj, 'flag')` to bind a property.
 *                   If provided, the returned ref === boolRef, so you can continue using it.
 *                   If not provided, a new ref is created.
 * @returns A Vue Ref<boolean> that stays in sync with the Y.Map entry (=== boolRef if supplied).
 */
export function useYBoolean(yMap: Y.Map<unknown>, key: string, boolRef?: Ref<boolean>): Ref<boolean> {
    const bool = (boolRef ?? ref(false)) as Ref<boolean>;
    let remoteApplying = false;
    let initialSync = true;

    function handleMapUpdate(_: Y.YMapEvent<unknown>): void {
        // Only update on remote changes or initial sync
        remoteApplying = true;
        const value = yMap.get(key);
        bool.value = typeof value === 'boolean' ? value : false;
        remoteApplying = false;
    }

    onMounted(() => {
        const initial = yMap.get(key);
        bool.value = typeof initial === 'boolean' ? initial : false;
        initialSync = false;
        yMap.observe(handleMapUpdate);
    });

    watch(bool, (newVal) => {
        if (remoteApplying || initialSync) return;
        yMap.set(key, newVal);
    });

    onUnmounted(() => {
        yMap.unobserve(handleMapUpdate);
    });

    return bool;
}

/**
 * Composable to bind a Vue Ref<number> (or a ref to a reactive property) to a Yjs Y.Map entry.
 */
export function useYNumber(yMap: Y.Map<unknown>, key: string, numRef?: Ref<number>): Ref<number> {
    const num = (numRef ?? ref(0)) as Ref<number>;
    let remoteApplying = false;
    let initialSync = true;

    function handleMapUpdate(_: Y.YMapEvent<unknown>): void {
        remoteApplying = true;
        const value = yMap.get(key);
        num.value = typeof value === 'number' ? value : 0;
        remoteApplying = false;
    }

    onMounted(() => {
        const initial = yMap.get(key);
        num.value = typeof initial === 'number' ? initial : 0;
        initialSync = false;
        yMap.observe(handleMapUpdate);
    });

    watch(num, (newVal) => {
        if (remoteApplying || initialSync) return;
        yMap.set(key, newVal);
    });

    onUnmounted(() => {
        yMap.unobserve(handleMapUpdate);
    });

    return num;
}

/**
 * Composable to bind a Vue Ref<T[]> (or a ref to a reactive property) to a Yjs Y.Array of Y.Map entries.
 * @param yArray   - The Y.Array instance holding Y.Map entries.
 * @param itemsRef - Optional Vue Ref<T[]> for existing reactive array. Returned ref === itemsRef if supplied.
 * @returns A Vue Ref<T[]> that stays in sync with the Y.Array content.
 */
export function useYArray<T extends Record<string, unknown>>(yArray: Y.Array<Y.Map<unknown>>, itemsRef?: Ref<T[]>): Ref<T[]> {
    const items = (itemsRef ?? ref([])) as Ref<T[]>;
    let remoteApplying = false;
    let initialSync = true;

    function snapshot(): T[] {
        return yArray.map((map) => map.toJSON() as T);
    }

    function handleArrayUpdate(_: Y.YArrayEvent<Y.Map<unknown>>): void {
        remoteApplying = true;
        items.value = snapshot();
        remoteApplying = false;
    }

    onMounted(() => {
        items.value = snapshot();
        initialSync = false;
        yArray.observe(handleArrayUpdate);
    });

    watch(
        items,
        (newArr) => {
            if (remoteApplying || initialSync) return;
            // overwrite entire array: for small lists this is simple
            yArray.delete(0, yArray.length);
            const maps = newArr.map((obj) => new Y.Map(Object.entries(obj)));
            yArray.insert(0, maps);
        },
        { deep: true },
    );

    onUnmounted(() => {
        yArray.unobserve(handleArrayUpdate);
    });

    return items;
}

/**
 * Generic composable to bind a reactive object of primitive fields to a Yjs Y.Map.
 * Syncs string, boolean, and number fields.
 * @param yMap   - The Y.Map instance holding fields.
 * @param state  - A reactive object with primitive values.
 * @param fields - Optional list of keys to sync; defaults to all keys of state.
 * @returns The same state object, kept in sync with yMap.
 */
export function useYMap<T extends Record<string, string | boolean | number>>(
    yMap: Y.Map<string | boolean | number>,
    state: T,
    fields?: Array<keyof T>,
): T {
    // Determine which keys to sync
    const keys: Array<keyof T> = fields ?? (Object.keys(state) as Array<keyof T>);
    const stops: Array<() => void> = [];
    let remoteApplying = false;
    let initialSync = true;

    function handleMapUpdate(_: Y.YMapEvent<string | boolean | number>): void {
        remoteApplying = true;
        // Apply remote changes for all watched keys
        for (const key of keys) {
            const val = yMap.get(key as string);
            if (val !== undefined) {
                state[key] = val as T[typeof key];
            }
        }
        remoteApplying = false;
    }

    onMounted(() => {
        // Initial sync from Yjs to state
        for (const key of keys) {
            const val = yMap.get(key as string);
            if (val !== undefined) {
                state[key] = val as T[typeof key];
            }
        }
        initialSync = false;
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        yMap.observe(handleMapUpdate as any);

        // Watch local state changes and propagate to Yjs
        for (const key of keys) {
            const stop = watch(
                () => state[key],
                (newVal) => {
                    if (remoteApplying || initialSync) return;
                    yMap.set(key as string, newVal);
                },
            );
            stops.push(stop);
        }
    });

    onUnmounted(() => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        yMap.unobserve(handleMapUpdate as any);
        stops.forEach((stop) => stop());
    });

    return state;
}

/**
 * Automatically sync a reactive state tree with Yjs types, recursing into nested objects.
 *
 * @param ydoc   - The Y.Doc or nested Y.Map (cast to Y.Doc for recursion)
 * @param state  - A reactive object of primitives, nested objects, or arrays of objects
 * @param fields - Optional subset of keys to sync (defaults to all keys)
 */
export function useYSyncStructure<T extends YStateMap>(ydoc: Y.Doc, state: T, fields?: Array<keyof T>): void {
    const keys = (fields ?? Object.keys(state)) as Array<keyof T>;

    for (const key of keys) {
        const name = key as string;
        const val = state[key];

        // Primitive string
        if (typeof val === 'string') {
            useYText(ydoc.getText(name), toRef(state, name) as Ref<string>);

            // Primitive boolean
        } else if (typeof val === 'boolean') {
            useYBoolean(ydoc.getMap('primitives'), name, toRef(state, name) as Ref<boolean>);

            // Primitive number
        } else if (typeof val === 'number') {
            useYNumber(ydoc.getMap('primitives'), name, toRef(state, name) as Ref<number>);

            // Array of nested maps
        } else if (Array.isArray(val)) {
            useYArray<YStateMap>(ydoc.getArray(name), toRef(state, name) as Ref<YStateMap[]>);

            // Nested object
        } else if (val !== null && typeof val === 'object') {
            const nestedMap = ydoc.getMap(name);
            const values = Object.values(val);

            // Check if all fields are primitives
            const isFlat = values.every((v) => typeof v === 'string' || typeof v === 'number' || typeof v === 'boolean');

            if (isFlat) {
                // Flat map: sync via useYMap
                useYMap<Record<string, Primitive>>(nestedMap as Y.Map<Primitive>, val as Record<string, Primitive>);
            } else {
                // Deeper nesting: recurse with the nested Y.Map cast as Y.Doc
                useYSyncStructure(nestedMap as unknown as Y.Doc, val as YStateMap);
            }
        }
    }
}
