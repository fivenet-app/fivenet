import { flushPromises, mount } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent, nextTick, ref } from 'vue';
import ExclusionsPanel from '~/components/jobs/groups/details/ExclusionsPanel.vue';
import LeadersPanel from '~/components/jobs/groups/details/LeadersPanel.vue';
import ManualMembersPanel from '~/components/jobs/groups/details/ManualMembersPanel.vue';
import RulesPanel from '~/components/jobs/groups/details/RulesPanel.vue';
import { GroupType } from '~~/gen/ts/resources/jobs/groups/group';

const mocks = vi.hoisted(() => ({
    jobsGroupsClient: {
        listGroupManualMembers: vi.fn().mockResolvedValue({ response: { manualMembers: [], pagination: { totalCount: 0 } } }),
        listGroupRules: vi.fn().mockResolvedValue({ response: { rules: [], pagination: { totalCount: 0 } } }),
        listGroupMemberExclusions: vi.fn().mockResolvedValue({ response: { exclusions: [], pagination: { totalCount: 0 } } }),
        listGroupLeaders: vi.fn().mockResolvedValue({ response: { leaders: [], pagination: { totalCount: 0 } } }),
    },
    qualificationsQualificationsClient: {
        getQualification: vi.fn().mockResolvedValue({ response: { qualification: undefined } }),
    },
    completorStore: {
        completeColleagues: vi.fn(),
        listJobs: vi.fn().mockResolvedValue([]),
    },
}));

vi.mock('#imports', async () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const actual = await vi.importActual<any>('#imports');

    return {
        ...actual,
        useOverlay: () => ({
            create: () => ({
                open: vi.fn(),
            }),
        }),
        useLazyAsyncData: () => ({
            data: ref(undefined),
            status: ref('success'),
            error: ref(undefined),
            refresh: vi.fn(),
        }),
    };
});

vi.mock('vue-i18n', () => ({
    createI18n: () => ({
        global: {
            t: (key: string) => key,
            tm: (key: string) => key,
            rt: (value: string) => value,
            locale: ref('en'),
        },
        install: vi.fn(),
    }),
    useI18n: () => ({
        t: (key: string) => key,
        tm: (key: string) => key,
        rt: (value: string) => value,
        locale: ref('en'),
    }),
}));

vi.mock('~~/gen/ts/clients', () => ({
    getJobsGroupsClient: vi.fn(async () => mocks.jobsGroupsClient),
    getQualificationsQualificationsClient: vi.fn(async () => mocks.qualificationsQualificationsClient),
}));

vi.mock('~/stores/completor', () => ({
    useCompletorStore: () => mocks.completorStore,
}));

const UAlertStub = defineComponent({
    name: 'UAlert',
    props: {
        title: {
            type: String,
            default: '',
        },
        description: {
            type: String,
            default: '',
        },
    },
    template: '<div class="u-alert-stub">{{ title }} {{ description }}</div>',
});

const DataNoDataBlockStub = defineComponent({
    name: 'DataNoDataBlock',
    props: {
        message: {
            type: String,
            default: '',
        },
    },
    template: '<div class="data-no-data-stub">{{ message }}</div>',
});

const PolicyPanelWrapper = defineComponent({
    name: 'PolicyPanelWrapper',
    components: {
        ExclusionsPanel,
        LeadersPanel,
        ManualMembersPanel,
        RulesPanel,
    },
    props: {
        panel: {
            type: String,
            required: true,
        },
        panelProps: {
            type: Object,
            required: true,
        },
    },
    template: `
        <Suspense>
            <component :is="panel" v-bind="panelProps" />
        </Suspense>
    `,
});

function mountManualMembersPanel(groupType = GroupType.SMART) {
    return mount(PolicyPanelWrapper, {
        props: {
            panel: 'ManualMembersPanel',
            panelProps: {
                groupId: 1,
                groupType,
                canView: true,
                canManage: false,
            },
        },
        global: {
            mocks: {
                $t: (key: string) => key,
            },
            stubs: {
                UAlert: UAlertStub,
                DataNoDataBlock: DataNoDataBlockStub,
            },
        },
    });
}

function mountRulesPanel(groupType = GroupType.MANUAL) {
    return mount(PolicyPanelWrapper, {
        props: {
            panel: 'RulesPanel',
            panelProps: {
                groupId: 1,
                groupType,
                canView: true,
                canManage: false,
            },
        },
        global: {
            mocks: {
                $t: (key: string) => key,
            },
            stubs: {
                UAlert: UAlertStub,
                DataNoDataBlock: DataNoDataBlockStub,
            },
        },
    });
}

function mountExclusionsPanel(groupType = GroupType.MANUAL) {
    return mount(PolicyPanelWrapper, {
        props: {
            panel: 'ExclusionsPanel',
            panelProps: {
                groupId: 1,
                groupType,
                canView: true,
                canManage: false,
            },
        },
        global: {
            mocks: {
                $t: (key: string) => key,
            },
            stubs: {
                UAlert: UAlertStub,
                DataNoDataBlock: DataNoDataBlockStub,
            },
        },
    });
}

function mountLeadersPanel() {
    return mount(PolicyPanelWrapper, {
        props: {
            panel: 'LeadersPanel',
            panelProps: {
                groupId: 1,
                canView: true,
                canManage: false,
            },
        },
        global: {
            mocks: {
                $t: (key: string) => key,
            },
            stubs: {
                DataNoDataBlock: DataNoDataBlockStub,
            },
        },
    });
}

describe('job group policy panels', () => {
    it('hides manual member controls for smart groups', async () => {
        const wrapper = mountManualMembersPanel();

        await flushPromises();
        await nextTick();

        expect(wrapper.text()).toContain('components.jobs.groups.policy.manual_members_disabled');
        expect(wrapper.text()).not.toContain('components.jobs.groups.details.add_manual_member');
    });

    it('hides rule controls for manual groups', async () => {
        const wrapper = mountRulesPanel();

        await flushPromises();
        await nextTick();

        expect(wrapper.text()).toContain('components.jobs.groups.policy.rules_disabled');
        expect(wrapper.text()).not.toContain('common.add');
    });

    it('hides exclusion controls for manual groups', async () => {
        const wrapper = mountExclusionsPanel();

        await flushPromises();
        await nextTick();

        expect(wrapper.text()).toContain('components.jobs.groups.policy.exclusions_disabled');
        expect(wrapper.text()).not.toContain('components.jobs.groups.details.add_exclusion');
    });

    it('hides rule controls when a supported panel is read-only', async () => {
        const wrapper = mountRulesPanel(GroupType.MIXED);

        await flushPromises();
        await nextTick();

        expect(wrapper.text()).not.toContain('components.jobs.groups.policy.rules_disabled');
        expect(wrapper.text()).not.toContain('common.add');
    });

    it('hides manual member controls when a supported panel is read-only', async () => {
        const wrapper = mountManualMembersPanel(GroupType.MIXED);

        await flushPromises();
        await nextTick();

        expect(wrapper.text()).not.toContain('components.jobs.groups.policy.manual_members_disabled');
        expect(wrapper.text()).not.toContain('components.jobs.groups.details.add_manual_member');
    });

    it('hides exclusion controls when a supported panel is read-only', async () => {
        const wrapper = mountExclusionsPanel(GroupType.MIXED);

        await flushPromises();
        await nextTick();

        expect(wrapper.text()).not.toContain('components.jobs.groups.policy.exclusions_disabled');
        expect(wrapper.text()).not.toContain('components.jobs.groups.details.add_exclusion');
    });

    it('hides leader controls when the panel is read-only', async () => {
        const wrapper = mountLeadersPanel();

        await flushPromises();
        await nextTick();

        expect(wrapper.text()).not.toContain('common.add');
    });
});
