import { mountSuspended } from '@nuxt/test-utils/runtime';
import { flushPromises } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent, nextTick, ref } from 'vue';
import SystemStatus from '~/components/settings/SystemStatus.vue';
import SystemStatusDBSyncDrawer from '~/components/settings/SystemStatusDBSyncDrawer.vue';
import { toTimestamp } from '~/utils/time';

const mockStatus = {
    database: {
        version: '8.0.36',
        connected: true,
        migrationVersion: 12,
        migrationDirty: false,
        dbCharset: 'utf8mb4',
        dbCollation: 'utf8mb4_unicode_ci',
        tablesOk: true,
    },
    nats: {
        version: '2.10.17',
        connected: true,
    },
    dbsync: {
        enabled: true,
        streamConnected: true,
        lastSyncedData: toTimestamp(new Date('2026-08-25T12:00:00Z')),
        lastSyncedActivity: toTimestamp(new Date('2026-08-25T12:05:00Z')),
        lastDbsyncVersion: '1.2.3',
        syncState: {
            tables: [
                {
                    table: 'jobs',
                    checkpoint: {
                        lastCheck: toTimestamp(new Date('2026-08-25T11:00:00Z')),
                        lastId: '42',
                    },
                    lastSyncedAt: toTimestamp(new Date('2026-08-25T12:10:00Z')),
                    lastAttemptAt: toTimestamp(new Date('2026-08-25T12:11:00Z')),
                },
                {
                    table: 'vehicles_resync',
                    lastSyncedAt: toTimestamp(new Date('2026-08-25T12:12:00Z')),
                    lastAttemptAt: toTimestamp(new Date('2026-08-25T12:13:00Z')),
                    lastError: 'sync failed',
                },
            ],
        },
    },
    version: {
        current: '2026.8.25',
    },
};

const getStatusMock = vi.fn().mockImplementation(async () => ({
    response: {
        status: structuredClone(mockStatus),
    },
}));

let resolveInitialStatusLoad: (() => void) | undefined;

const translations: Record<string, string> = {
    'common.copy': 'Copy',
    'common.loading': 'Loading',
    'common.na': 'N/A',
    'common.no': 'No',
    'common.not_found': 'Not found',
    'common.status': 'status',
    'common.version': 'Version',
    'components.settings.system_status.title': 'System Status',
    'components.settings.system_status.database.title': 'Database',
    'components.settings.system_status.db_sync.title': 'DB Sync',
    'components.settings.system_status.db_sync.drawer_title': 'DB Sync Details',
    'components.settings.system_status.db_sync.summary': 'Summary',
    'components.settings.system_status.db_sync.per_table': 'Per Table',
    'components.settings.system_status.db_sync.last_synced_at': 'Last Synced At',
    'components.settings.system_status.db_sync.last_attempt_at': 'Last Attempt At',
    'components.settings.system_status.db_sync.checkpoint': 'Checkpoint',
    'components.settings.system_status.db_sync.last_error': 'Last Error',
    'components.settings.system_status.db_sync.state.healthy': 'Healthy',
    'components.settings.system_status.db_sync.state.error': 'Error',
    'components.settings.system_status.db_sync.state.idle': 'Idle',
    'components.settings.system_status.db_sync.no_tables': 'No table sync state is available',
    'components.settings.system_status.db_sync.last_data_received': 'Last Data Received',
    'components.settings.system_status.db_sync.last_activity_received': 'Last Activity Received',
    'components.settings.system_status.db_sync.last_dbsync_version': 'Last seen DB Sync Version',
    'components.settings.system_status.db_sync.tables.jobs': 'Jobs',
    'components.settings.system_status.db_sync.tables.users_resync': 'Users Resync',
    'components.settings.system_status.db_sync.tables.vehicles_resync': 'Vehicles Resync',
    'components.settings.system_status.nats.title': 'NATS',
};

vi.mock('#imports', async () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const actual = await vi.importActual<any>('#imports');

    return {
        ...actual,
        useLazyAsyncData: (_key: string, handler: () => Promise<unknown>) => {
            const data = ref<unknown>();
            const error = ref<unknown>();
            const status = ref<'pending' | 'success' | 'error'>('pending');
            const refresh = vi.fn(async () => {
                status.value = 'pending';
                try {
                    data.value = await handler();
                    status.value = 'success';
                } catch (err) {
                    error.value = err;
                    status.value = 'error';
                }
            });

            void new Promise<void>((resolve) => {
                resolveInitialStatusLoad = resolve;
            }).then(refresh);

            return { data, error, status, refresh };
        },
    };
});

vi.mock('~~/gen/ts/clients', () => ({
    getSettingsSystemClient: vi.fn(async () => ({
        getStatus: getStatusMock,
    })),
}));

const UButtonStub = defineComponent({
    name: 'UButton',
    props: {
        label: {
            type: String,
            default: '',
        },
    },
    emits: ['click'],
    template: '<button class="u-button-stub" :data-label="label" @click="$emit(\'click\')"><slot />{{ label }}</button>',
});

const UDrawerStub = defineComponent({
    name: 'UDrawer',
    props: {
        open: {
            type: Boolean,
            default: false,
        },
    },
    emits: ['update:open'],
    setup(props, { emit }) {
        const open = ref(props.open);

        watch(
            () => props.open,
            (value) => {
                open.value = value;
            },
        );

        function openDrawer() {
            open.value = true;
            emit('update:open', true);
        }

        return { open, openDrawer };
    },
    template:
        '<div class="u-drawer-stub"><div class="u-drawer-trigger" @click="openDrawer"><slot /></div><div v-if="open" class="u-drawer-panel"><div class="u-drawer-title"><slot name="title" /></div><div class="u-drawer-body"><slot name="body" /></div><div class="u-drawer-footer"><slot name="footer" /></div></div></div>',
});

const UCardStub = defineComponent({
    name: 'UCard',
    props: {
        title: {
            type: String,
            default: '',
        },
    },
    template:
        '<section class="u-card-stub"><header><slot name="title">{{ title }}</slot></header><div><slot /></div></section>',
});

const UPopoverStub = defineComponent({
    name: 'UPopover',
    template: '<div class="u-popover-stub"><slot /><slot name="content" /></div>',
});

const UTooltipStub = defineComponent({
    name: 'UTooltip',
    template: '<div class="u-tooltip-stub"><slot /></div>',
});

const UBadgeStub = defineComponent({
    name: 'UBadge',
    props: {
        label: {
            type: [String, Number],
            default: '',
        },
        color: {
            type: String,
            default: '',
        },
    },
    template: '<span class="u-badge-stub" :data-color="color">{{ label }}<slot /></span>',
});

const UChipStub = defineComponent({
    name: 'UChip',
    props: {
        color: {
            type: String,
            default: '',
        },
        show: {
            type: Boolean,
            default: false,
        },
    },
    template: '<div class="u-chip-stub" :data-color="color" :data-show="show"><slot /></div>',
});

const UIconStub = defineComponent({
    name: 'UIcon',
    template: '<i class="u-icon-stub" />',
});

const USkeletonStub = defineComponent({
    name: 'USkeleton',
    template: '<div class="u-skeleton-stub" />',
});

const DataPendingBlockStub = defineComponent({
    name: 'DataPendingBlock',
    props: {
        message: {
            type: String,
            default: '',
        },
    },
    template: '<div class="data-pending-stub">{{ message }}</div>',
});

const DataErrorBlockStub = defineComponent({
    name: 'DataErrorBlock',
    props: {
        title: {
            type: String,
            default: '',
        },
    },
    template: '<div class="data-error-stub">{{ title }}</div>',
});

const DataNoDataBlockStub = defineComponent({
    name: 'DataNoDataBlock',
    props: {
        title: {
            type: String,
            default: '',
        },
        message: {
            type: String,
            default: '',
        },
    },
    template: '<div class="data-nodata-stub">{{ title }} {{ message }}</div>',
});

const GenericTimeStub = defineComponent({
    name: 'GenericTime',
    props: {
        value: {
            type: Object,
            default: undefined,
        },
    },
    template: '<time class="generic-time-stub">{{ value ? "time" : "" }}</time>',
});

async function mountSystemStatus() {
    return mountSuspended(SystemStatus, {
        global: {
            mocks: {
                $t: (key: string) => translations[key] ?? key,
                $d: (value: Date) => value.toISOString(),
            },
            stubs: {
                UButton: UButtonStub,
                UDrawer: UDrawerStub,
                UCard: UCardStub,
                UPopover: UPopoverStub,
                UTooltip: UTooltipStub,
                UBadge: UBadgeStub,
                UChip: UChipStub,
                UIcon: UIconStub,
                USkeleton: USkeletonStub,
                DataPendingBlock: DataPendingBlockStub,
                DataErrorBlock: DataErrorBlockStub,
                DataNoDataBlock: DataNoDataBlockStub,
                GenericTime: GenericTimeStub,
            },
        },
    });
}

describe('SystemStatus', () => {
    it('keeps the db sync drawer open on refresh', async () => {
        const wrapper = await mountSystemStatus();

        resolveInitialStatusLoad?.();
        await flushPromises();
        await nextTick();

        expect(wrapper.text()).toContain('Database');
        expect(wrapper.text()).toContain('NATS');

        const dbSyncButton = wrapper.find('button[data-label="DB Sync"]');
        expect(dbSyncButton.exists()).toBe(true);
        expect(wrapper.find('button[data-label="Database"]').exists()).toBe(true);
        expect(wrapper.find('button[data-label="NATS"]').exists()).toBe(true);

        await dbSyncButton.trigger('click');
        await nextTick();

        expect(wrapper.text()).toContain('DB Sync Details');
        expect(wrapper.text()).toContain('Summary');
        expect(wrapper.text()).toContain('Per Table');
        expect(wrapper.text()).toContain('Jobs');
        expect(wrapper.text()).toContain('Vehicles Resync');
        expect(wrapper.text()).toContain('Last Synced At');
        expect(wrapper.text()).toContain('Last Attempt At');
        expect(wrapper.text()).toContain('Checkpoint');
        expect(wrapper.text()).toContain('Last Error');
        expect(wrapper.text()).toContain('Healthy');
        expect(wrapper.text()).toContain('Error');
        expect(wrapper.text()).toContain('sync failed');
        expect(wrapper.find('.u-badge-stub').attributes('data-color')).toBe('success');
        expect(wrapper.find('.u-chip-stub').attributes('data-color')).toBe('success');

        const dbSyncDrawer = wrapper.findComponent(SystemStatusDBSyncDrawer);
        expect(dbSyncDrawer.exists()).toBe(true);

        await wrapper.find('button[data-label="DB Sync"]').trigger('click');
        await nextTick();

        expect(wrapper.text()).toContain('DB Sync Details');

        mockStatus.dbsync.streamConnected = false;
        dbSyncDrawer.vm.$emit('refresh');
        await flushPromises();
        await nextTick();

        expect(wrapper.text()).toContain('DB Sync Details');
        expect(dbSyncDrawer.find('.u-badge-stub').attributes('data-color')).toBe('warning');
        expect(dbSyncDrawer.find('.u-chip-stub').attributes('data-color')).toBe('warning');
    });
});
