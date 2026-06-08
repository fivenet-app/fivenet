import { describe, expect, it } from 'vitest';
import { isProxy, reactive, ref } from 'vue';
import { deepToRaw } from './deepToRaw';

describe('deepToRaw', () => {
    it('unwraps nested refs and proxies into a structured-clone-safe object', () => {
        const source = reactive({
            count: ref(1),
            nested: reactive({
                label: 'hello',
                items: [reactive({ id: 2 })],
            }),
        });

        const raw = deepToRaw(source);

        expect(isProxy(raw)).toBe(false);
        expect(isProxy(raw.nested)).toBe(false);
        expect(isProxy(raw.nested.items[0])).toBe(false);
        expect(raw.count).toBe(1);
        expect(raw.nested.label).toBe('hello');
        expect(raw.nested.items[0].id).toBe(2);
        expect(() => structuredClone(raw)).not.toThrow();
    });
});
