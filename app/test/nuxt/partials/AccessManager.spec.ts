import { mountSuspended } from '@nuxt/test-utils/runtime';
import { flushPromises } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent, nextTick } from 'vue';
import AccessEntry from '~/components/partials/access/AccessEntry.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import type { AccessLevelEnum, JobAccessEntry, MixedAccessEntry } from '~/components/partials/access/helpers';

const mocks = vi.hoisted(() => ({
    completorStore: {
        completeCitizens: vi.fn(),
        listJobs: vi.fn().mockResolvedValue([]),
    },
    getQualificationsQualificationsClient: vi.fn().mockResolvedValue({
        getQualification: vi.fn(),
        listQualifications: vi.fn(),
    }),
    activeChar: { value: { job: 'police' } },
}));

vi.mock('~/stores/completor', () => ({
    useCompletorStore: () => mocks.completorStore,
}));

vi.mock('~/composables/useAuth', () => ({
    useAuth: () => ({ activeChar: mocks.activeChar }),
}));

vi.mock('~~/gen/ts/clients', () => ({
    getQualificationsQualificationsClient: mocks.getQualificationsQualificationsClient,
}));

const UButtonStub = defineComponent({
    name: 'UButton',
    props: {
        disabled: {
            type: Boolean,
            default: false,
        },
        label: {
            type: String,
            default: '',
        },
    },
    emits: ['click'],
    template:
        '<button class="u-button-stub" :disabled="disabled" :data-label="label" @click="$emit(\'click\')">{{ label }}</button>',
});

const UFormFieldStub = defineComponent({
    name: 'UFormField',
    template: '<div class="u-form-field-stub"><slot /></div>',
});

const UTooltipStub = defineComponent({
    name: 'UTooltip',
    template: '<div class="u-tooltip-stub"><slot /></div>',
});

const UBadgeStub = defineComponent({
    name: 'UBadge',
    props: {
        label: {
            type: String,
            default: '',
        },
    },
    template: '<span class="u-badge-stub">{{ label }}</span>',
});

const UCheckboxStub = defineComponent({
    name: 'UCheckbox',
    props: {
        modelValue: {
            type: Boolean,
            default: false,
        },
        disabled: {
            type: Boolean,
            default: false,
        },
    },
    emits: ['update:modelValue'],
    template: '<input class="u-checkbox-stub" type="checkbox" :checked="modelValue" :disabled="disabled" />',
});

const UInputStub = defineComponent({
    name: 'UInput',
    template: '<input class="u-input-stub" />',
});

const USelectMenuStub = defineComponent({
    name: 'USelectMenu',
    props: {
        disabled: {
            type: Boolean,
            default: false,
        },
        items: {
            type: Array,
            default: () => [],
        },
    },
    template: '<div class="u-select-menu-stub" :data-disabled="String(disabled)"><slot /><slot name="empty" /></div>',
});

const ClientOnlyStub = defineComponent({
    name: 'ClientOnly',
    template: '<div class="client-only-stub"><slot /></div>',
});

const SelectMenuStub = defineComponent({
    name: 'SelectMenu',
    template: '<div class="select-menu-stub"><slot /><slot name="empty" /></div>',
});

const AccessEntryStub = defineComponent({
    name: 'AccessEntry',
    props: {
        modelValue: {
            type: Object,
            required: true,
        },
    },
    emits: ['delete', 'update:modelValue'],
    template:
        '<div class="access-entry-stub"><span class="entry-type">{{ modelValue.type }}</span><button class="entry-delete" @click="$emit(\'delete\')">delete</button></div>',
});

const accessRoles: AccessLevelEnum[] = [
    { label: 'View', value: 2 },
    { label: 'Edit', value: 3 },
];

function managerStubs() {
    return {
        AccessEntry: AccessEntryStub,
        UButton: UButtonStub,
        UTooltip: UTooltipStub,
    };
}

function entryStubs() {
    return {
        UButton: UButtonStub,
        UFormField: UFormFieldStub,
        UTooltip: UTooltipStub,
        UBadge: UBadgeStub,
        UCheckbox: UCheckboxStub,
        UInput: UInputStub,
        USelectMenu: USelectMenuStub,
        ClientOnly: ClientOnlyStub,
        SelectMenu: SelectMenuStub,
    };
}

describe('AccessManager', () => {
    it('renders existing entries and adds a new entry to the model', async () => {
        const jobs: JobAccessEntry[] = [
            {
                id: 10,
                targetId: 7,
                job: 'police',
                minimumGrade: 0,
                access: 2,
            },
        ];
        const wrapper = await mountSuspended(AccessManager, {
            props: {
                jobs,
                targetId: 7,
                accessRoles,
                accessTypes: [{ label: 'Jobs', value: 'job' }],
                name: 'access',
            },
            global: { stubs: managerStubs(), mocks: { $t: (key: string) => key } },
        });

        expect(wrapper.findAllComponents(AccessEntryStub)).toHaveLength(1);

        await wrapper.find('button[data-label="components.access.add_entry"]').trigger('click');
        await flushPromises();

        expect(wrapper.findAllComponents(AccessEntryStub)).toHaveLength(2);
        expect(jobs).toHaveLength(2);
        expect(jobs[1]).toEqual(
            expect.objectContaining({
                id: expect.any(Number),
                targetId: 7,
                access: 2,
            }),
        );
    });

    it('removes a non-required entry', async () => {
        const jobs: JobAccessEntry[] = [{ id: 10, targetId: 7, job: 'police', minimumGrade: 0, access: 2 }];
        const wrapper = await mountSuspended(AccessManager, {
            props: {
                jobs,
                targetId: 7,
                accessRoles,
                accessTypes: [{ label: 'Jobs', value: 'job' }],
            },
            global: { stubs: managerStubs(), mocks: { $t: (key: string) => key } },
        });

        await wrapper.find('.entry-delete').trigger('click');
        await flushPromises();

        expect(wrapper.findAllComponents(AccessEntryStub)).toHaveLength(0);
        expect(jobs).toEqual([]);
    });

    it('disables adding entries at the configured limit', async () => {
        const wrapper = await mountSuspended(AccessManager, {
            props: {
                jobs: [{ id: 10, targetId: 7, job: 'police', minimumGrade: 0, access: 2 }],
                targetId: 7,
                totalLimit: 1,
                accessRoles,
                accessTypes: [{ label: 'Jobs', value: 'job' }],
            },
            global: { stubs: managerStubs(), mocks: { $t: (key: string) => key } },
        });

        expect(wrapper.find('button[data-label="components.access.add_entry"]').attributes('disabled')).toBeDefined();
    });

    it('updates entry target IDs when the target changes', async () => {
        const jobs: JobAccessEntry[] = [{ id: 10, targetId: 7, job: 'police', minimumGrade: 0, access: 2 }];
        const wrapper = await mountSuspended(AccessManager, {
            props: {
                jobs,
                targetId: 7,
                accessRoles,
                accessTypes: [{ label: 'Jobs', value: 'job' }],
            },
            global: { stubs: managerStubs(), mocks: { $t: (key: string) => key } },
        });

        await wrapper.setProps({ targetId: 8 });
        await nextTick();

        expect(jobs[0]?.targetId).toBe(8);
    });
});

describe('AccessEntry', () => {
    it('does not emit delete when attempting to remove a required entry', async () => {
        const entry: MixedAccessEntry = {
            id: 1,
            type: 'job',
            job: 'police',
            minimumGrade: 0,
            access: 2,
            required: true,
        };
        const wrapper = await mountSuspended(AccessEntry, {
            props: {
                modelValue: entry,
                accessTypes: [{ label: 'Jobs', value: 'job' }],
                accessRoles,
                jobs: [{ name: 'police', label: 'Police', grades: [{ grade: 0, name: 'Cadet', label: 'Cadet' }] }],
            },
            global: { stubs: entryStubs(), mocks: { $t: (key: string) => key } },
        });

        const removeButton = wrapper.find('button[data-label="components.access.remove_entry"]');
        expect(removeButton.attributes('disabled')).toBeDefined();

        await removeButton.trigger('click');
        expect(wrapper.emitted('delete')).toBeUndefined();
    });

    it('raises access to the required access floor', async () => {
        const entry: MixedAccessEntry = {
            id: 1,
            type: 'job',
            job: 'police',
            minimumGrade: 0,
            access: 2,
            required: true,
            requiredAccess: 3,
        };
        const wrapper = await mountSuspended(AccessEntry, {
            props: {
                modelValue: entry,
                accessTypes: [{ label: 'Jobs', value: 'job' }],
                accessRoles,
                jobs: [{ name: 'police', label: 'Police', grades: [{ grade: 0, name: 'Cadet', label: 'Cadet' }] }],
            },
            global: { stubs: entryStubs(), mocks: { $t: (key: string) => key } },
        });

        await flushPromises();

        expect(wrapper.emitted('update:modelValue')?.at(-1)?.[0]).toEqual(
            expect.objectContaining({ access: 3, requiredAccess: 3 }),
        );
    });

    it('clears type-specific data when switching entry types', async () => {
        const entry: MixedAccessEntry = {
            id: 1,
            type: 'user',
            userId: 42,
            user: {
                userId: 42,
                job: 'police',
                jobGrade: 0,
                firstname: 'Alex',
                lastname: 'Example',
                dateofbirth: '2000-01-01',
            },
            access: 2,
        };
        const wrapper = await mountSuspended(AccessEntry, {
            props: {
                modelValue: entry,
                accessTypes: [
                    { label: 'Users', value: 'user' },
                    { label: 'Jobs', value: 'job' },
                ],
                accessRoles,
            },
            global: { stubs: entryStubs(), mocks: { $t: (key: string) => key } },
        });

        await wrapper.setProps({ modelValue: { ...entry, type: 'job' } });
        await nextTick();

        const updated = wrapper.emitted('update:modelValue')?.at(-1)?.[0] as MixedAccessEntry | undefined;
        expect(updated).toEqual(
            expect.objectContaining({
                type: 'job',
                userId: undefined,
                user: undefined,
                job: undefined,
                qualificationId: undefined,
            }),
        );
    });

    it('disables already-used job and grade combinations', async () => {
        const entry: MixedAccessEntry = {
            id: 2,
            type: 'job',
            minimumGrade: 0,
            access: 2,
        };
        const wrapper = await mountSuspended(AccessEntry, {
            props: {
                modelValue: entry,
                existingEntries: [{ id: 1, type: 'job', job: 'police', minimumGrade: 0, access: 2 }],
                accessTypes: [{ label: 'Jobs', value: 'job' }],
                accessRoles,
                jobs: [
                    {
                        name: 'police',
                        label: 'Police',
                        grades: [{ grade: 0, name: 'Cadet', label: 'Cadet' }],
                    },
                    {
                        name: 'ems',
                        label: 'EMS',
                        grades: [{ grade: 0, name: 'Trainee', label: 'Trainee' }],
                    },
                ],
            },
            global: { stubs: entryStubs(), mocks: { $t: (key: string) => key } },
        });

        const menus = wrapper.findAllComponents(USelectMenuStub);
        const jobsMenu = menus.find((menu) => menu.props('items').some((item: { name?: string }) => item.name === 'police'));
        expect(jobsMenu?.props('items')).toEqual(
            expect.arrayContaining([expect.objectContaining({ name: 'police', disabled: true })]),
        );
    });

    it('shows a fallback label when a user lookup fails', async () => {
        mocks.completorStore.completeCitizens.mockRejectedValueOnce(new Error('lookup failed'));
        const entry: MixedAccessEntry = {
            id: 1,
            type: 'user',
            userId: 42,
            access: 2,
        };
        const wrapper = await mountSuspended(AccessEntry, {
            props: {
                modelValue: entry,
                accessTypes: [{ label: 'Users', value: 'user' }],
                accessRoles,
            },
            global: { stubs: entryStubs(), mocks: { $t: (key: string) => key } },
        });

        await flushPromises();

        expect(wrapper.text()).toContain('#UserID 42');
        expect(entry.userId).toBe(42);
    });

    it('does not initialize the qualifications client for a user entry', async () => {
        mocks.getQualificationsQualificationsClient.mockClear();
        const entry: MixedAccessEntry = {
            id: 1,
            type: 'user',
            access: 2,
        };

        await mountSuspended(AccessEntry, {
            props: {
                modelValue: entry,
                accessTypes: [{ label: 'Users', value: 'user' }],
                accessRoles,
            },
            global: { stubs: entryStubs(), mocks: { $t: (key: string) => key } },
        });

        expect(mocks.getQualificationsQualificationsClient).not.toHaveBeenCalled();
    });
});
