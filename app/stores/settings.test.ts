import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it } from 'vitest';
import { useSettingsStore } from './settings';

describe('useSettingsStore quick access helpers', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('pins routes idempotently', () => {
        const settingsStore = useSettingsStore();

        settingsStore.pinOverviewQuickAccess('/mail');
        settingsStore.pinOverviewQuickAccess('/mail');
        settingsStore.pinOverviewQuickAccess('/calendar');

        expect(settingsStore.overviewQuickAccess).toEqual(['/mail', '/calendar']);
        expect(settingsStore.isOverviewQuickAccess('/mail')).toBe(true);
        expect(settingsStore.isOverviewQuickAccess('/wiki')).toBe(false);
    });

    it('unpins routes and toggles them back on', () => {
        const settingsStore = useSettingsStore();

        settingsStore.pinOverviewQuickAccess('/mail');
        settingsStore.pinOverviewQuickAccess('/calendar');
        settingsStore.unpinOverviewQuickAccess('/mail');

        expect(settingsStore.overviewQuickAccess).toEqual(['/calendar']);

        settingsStore.toggleOverviewQuickAccess('/calendar');
        expect(settingsStore.overviewQuickAccess).toEqual([]);

        settingsStore.toggleOverviewQuickAccess('/wiki');
        expect(settingsStore.overviewQuickAccess).toEqual(['/wiki']);
    });

    it('normalizes duplicate entries while preserving order', () => {
        const settingsStore = useSettingsStore();

        settingsStore.overviewQuickAccess = ['/mail', '/calendar', '/mail', '/wiki', '/calendar'];
        settingsStore.normalizeOverviewQuickAccess();

        expect(settingsStore.overviewQuickAccess).toEqual(['/mail', '/calendar', '/wiki']);
    });

    it('reorders only visible quick access entries and keeps hidden ones in place', () => {
        const settingsStore = useSettingsStore();

        settingsStore.overviewQuickAccess = ['/mail', '/hidden-1', '/calendar', '/hidden-2', '/wiki'];

        settingsStore.reorderOverviewQuickAccess(['/wiki', '/mail', '/calendar']);

        expect(settingsStore.overviewQuickAccess).toEqual(['/wiki', '/hidden-1', '/mail', '/hidden-2', '/calendar']);
    });
});
