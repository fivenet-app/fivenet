<script lang="ts" setup>
import type { FormSubmitEvent, TabsItem } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { CentrumAccessLevel, type CentrumJobAccess } from '~~/gen/ts/resources/centrum/access';
import { CentrumMode, CentrumType, type Settings } from '~~/gen/ts/resources/centrum/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const { activeChar, attr } = useAuth();

const notifications = useNotificationsStore();

const { maxAccessEntries } = useAppConfig();

const centrumCentrumClient = await getCentrumCentrumClient();

const { data: settings, status, refresh, error } = useLazyAsyncData('settings-centrum-settings', () => getCentrumSettings());

async function getCentrumSettings(): Promise<Settings> {
    try {
        const call = centrumCentrumClient.getSettings({});
        const { response } = await call;

        return response.settings!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const modes = ref<{ mode: CentrumMode; selected?: boolean }[]>([
    { mode: CentrumMode.MANUAL },
    { mode: CentrumMode.SIMPLIFIED },
    { mode: CentrumMode.CENTRAL_COMMAND },
    { mode: CentrumMode.AUTO_ROUND_ROBIN },
]);

const schema = z.object({
    enabled: z.coerce.boolean().default(false),
    type: z.nativeEnum(CentrumType).default(CentrumType.DISPATCH),
    public: z.coerce.boolean(),
    mode: z.nativeEnum(CentrumMode).default(CentrumMode.MANUAL),
    fallbackMode: z.nativeEnum(CentrumMode).default(CentrumMode.AUTO_ROUND_ROBIN),
    predefinedStatus: z.object({
        unitStatus: z.string().array().max(20).default([]),
        dispatchStatus: z.string().array().max(20).default([]),
    }),
    timings: z.object({
        dispatchMaxWait: z.coerce.number().min(30).max(6000).default(900),
        requireUnit: z.coerce.boolean(),
        requireUnitReminderSeconds: z.coerce.number().min(60).max(6000).default(180),
    }),
    access: z
        .object({
            jobs: z.custom<CentrumJobAccess>().array().max(maxAccessEntries).default([]),
        })
        .default({ jobs: [] }),
    offeredAccess: z
        .object({
            jobs: z.custom<CentrumJobAccess>().array().max(maxAccessEntries).default([]),
        })
        .default({ jobs: [] }),
    configuration: z.object({
        deduplicationEnabled: z.coerce.boolean().default(true),
        deduplicationRadius: z.coerce.number().min(5).max(1000000).default(45),
        deduplicationDuration: z.coerce.number().max(1000000).positive().default(180),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    enabled: false,
    type: CentrumType.DISPATCH,
    public: false,
    mode: CentrumMode.MANUAL,
    fallbackMode: CentrumMode.AUTO_ROUND_ROBIN,
    predefinedStatus: {
        unitStatus: [],
        dispatchStatus: [],
    },
    timings: {
        dispatchMaxWait: 900,
        requireUnit: false,
        requireUnitReminderSeconds: 180,
    },
    access: {
        jobs: [],
    },
    offeredAccess: {
        jobs: [],
    },
    configuration: {
        deduplicationEnabled: true,
        deduplicationRadius: 45,
        deduplicationDuration: 180,
    },
});

async function updateSettings(values: Schema): Promise<void> {
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));
    values.offeredAccess.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    try {
        const call = centrumCentrumClient.updateSettings({
            settings: {
                job: '',
                enabled: values.enabled,
                type: values.type,
                public: values.public,
                mode: values.mode,
                fallbackMode: values.fallbackMode,
                predefinedStatus: values.predefinedStatus,
                timings: values.timings,
                access: values.access,
                offeredAccess: values.offeredAccess,
                configuration: {
                    deduplicationEnabled: values.configuration?.deduplicationEnabled ?? true,
                    deduplicationRadius: values.configuration?.deduplicationRadius ?? 45,
                    deduplicationDuration: {
                        seconds: values.configuration?.deduplicationDuration ?? 180,
                        nanos: 0,
                    },
                },
            },
        });
        await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await refresh();

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setSettingsValues(): void {
    if (!settings.value) {
        return;
    }

    state.enabled = settings.value.enabled;
    state.mode = settings.value.mode;
    state.fallbackMode = settings.value.fallbackMode;
    state.predefinedStatus = settings.value.predefinedStatus ?? {
        unitStatus: [],
        dispatchStatus: [],
    };
    state.timings = settings.value.timings ?? {
        dispatchMaxWait: 900,
        requireUnit: false,
        requireUnitReminderSeconds: 180,
    };
    state.access = settings.value.access ?? { jobs: [] };
    state.offeredAccess = settings.value.offeredAccess ?? { jobs: [] };
    state.configuration = {
        deduplicationEnabled: settings.value.configuration?.deduplicationEnabled ?? true,
        deduplicationRadius: settings.value.configuration?.deduplicationRadius ?? 45,
        deduplicationDuration: settings.value.configuration?.deduplicationDuration?.seconds ?? 180,
    };
}

watch(settings, () => setSettingsValues());
setSettingsValues();

const items: TabsItem[] = [
    {
        slot: 'settings' as const,
        label: t('common.settings'),
        icon: 'i-mdi-settings',
        value: 'settings',
    },
    {
        slot: 'predefined' as const,
        label: `${t('common.predefined', 2)} ${t('common.status', 2)}`,
        icon: 'i-mdi-selection',
        value: 'predefined',
    },
    {
        slot: 'timings' as const,
        label: t('common.timings'),
        icon: 'i-mdi-access-time',
        value: 'timings',
    },
    {
        slot: 'access' as const,
        label: t('common.access'),
        icon: 'i-mdi-lock',
        value: 'access',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'settings';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.submitter?.getAttribute('role') === 'tab') {
        return;
    }

    canSubmit.value = false;
    await updateSettings(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('components.centrum.settings.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/centrum" />

                    <UButton
                        v-if="!!settings"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="() => formRef?.submit()"
                    >
                        {{ $t('common.save', 1) }}
                    </UButton>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.settings')])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.settings')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!settings" icon="i-mdi-tune" :type="$t('common.settings')" />

            <UForm
                v-else
                ref="formRef"
                class="flex w-full max-w-full flex-1 flex-col overflow-y-auto"
                :schema="schema"
                :state="state"
                @submit="onSubmitThrottle"
            >
                <UTabs
                    v-model="selectedTab"
                    class="flex flex-1 flex-col"
                    :items="items"
                    variant="link"
                    :ui="{ content: 'flex flex-col p-4 gap-4 sm:gap-4' }"
                    :unmount-on-hide="false"
                >
                    <template #settings>
                        <UPageCard
                            :title="$t('components.centrum.settings.title')"
                            :description="$t('components.centrum.settings.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="enabled"
                                :label="$t('common.enabled')"
                                :ui="{ container: '' }"
                            >
                                <USwitch v-model="state.enabled" name="enabled" :disabled="!canSubmit" />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="public"
                                :label="$t('common.public')"
                                :description="$t('components.centrum.settings.public.description')"
                                :ui="{ container: '' }"
                            >
                                <USwitch
                                    v-model="state.public"
                                    :disabled="
                                        !canSubmit || !attr('centrum.CentrumService/UpdateSettings', 'Access', 'Public').value
                                    "
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="mode"
                                :label="$t('common.mode')"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.mode"
                                        :items="modes"
                                        value-key="mode"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        :disabled="!canSubmit"
                                    >
                                        <template #default>
                                            {{ $t(`enums.centrum.CentrumMode.${CentrumMode[state.mode ?? 0]}`) }}
                                        </template>

                                        <template #item="{ item }">
                                            {{ $t(`enums.centrum.CentrumMode.${CentrumMode[item.mode ?? 0]}`) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="fallbackMode"
                                :label="$t('common.fallback_mode')"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.fallbackMode"
                                        :items="modes"
                                        value-key="mode"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        :disabled="!canSubmit"
                                    >
                                        <template #default>
                                            {{ $t(`enums.centrum.CentrumMode.${CentrumMode[state.fallbackMode ?? 0]}`) }}
                                        </template>

                                        <template #item="{ item }">
                                            {{ $t(`enums.centrum.CentrumMode.${CentrumMode[item.mode ?? 0]}`) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.centrum.settings.deduplication.title')"
                            :description="$t('components.centrum.settings.deduplication.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="configuration.deduplicationEnabled"
                                :label="$t('common.enabled')"
                                :ui="{ container: '' }"
                            >
                                <USwitch v-model="state.configuration.deduplicationEnabled" :disabled="!canSubmit" />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="configuration.deduplicationDuration"
                                :label="$t('components.centrum.settings.deduplication.deduplication_duration')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.configuration.deduplicationDuration"
                                    type="number"
                                    :min="30"
                                    :placeholder="$t('common.time_ago.second', 2)"
                                    trailing-icon="i-mdi-access-time"
                                    :disabled="!canSubmit"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="configuration.deduplicationRadius"
                                :label="$t('components.centrum.settings.deduplication.deduplication_radius')"
                                :description="$t('components.centrum.settings.deduplication.deduplication_radius_description')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.configuration.deduplicationRadius"
                                    type="number"
                                    :min="5"
                                    :placeholder="$t('common.meters', 2)"
                                    :disabled="!canSubmit"
                                    :ui="{ base: 'pr-16!' }"
                                >
                                    <template #trailing>
                                        <span class="text-xs text-muted">
                                            {{ $t('common.meters', 2) }}
                                        </span>
                                    </template>
                                </UInput>
                            </UFormField>
                        </UPageCard>
                    </template>

                    <template #predefined>
                        <UPageCard
                            :title="$t('components.centrum.settings.predefined.title')"
                            :description="$t('components.centrum.settings.predefined.description')"
                        >
                            <!-- Predefined Unit Status Reason -->
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="unitStatus"
                                :label="`${$t('common.unit')} ${$t('common.status')}`"
                                :ui="{ container: '' }"
                            >
                                <div class="flex flex-col gap-1">
                                    <div
                                        v-for="(_, idx) in state.predefinedStatus.unitStatus"
                                        :key="idx"
                                        class="flex items-center gap-1"
                                    >
                                        <UFormField class="flex-1" :name="`unitStatus.${idx}`" :ui="{ container: '' }">
                                            <UInput
                                                v-model="state.predefinedStatus.unitStatus[idx]"
                                                class="w-full flex-1"
                                                type="text"
                                                :placeholder="$t('common.reason')"
                                                :disabled="!canSubmit"
                                            />
                                        </UFormField>

                                        <UTooltip :text="$t('common.delete')">
                                            <UButton
                                                icon="i-mdi-close"
                                                :disabled="!canSubmit"
                                                @click="state.predefinedStatus.unitStatus.splice(idx, 1)"
                                            />
                                        </UTooltip>
                                    </div>
                                </div>

                                <UTooltip :text="$t('common.add')">
                                    <UButton
                                        :class="state.predefinedStatus.unitStatus.length ? 'mt-2' : ''"
                                        icon="i-mdi-plus"
                                        :disabled="!canSubmit || state.predefinedStatus.unitStatus.length >= 8"
                                        @click="state.predefinedStatus.unitStatus.push('')"
                                    />
                                </UTooltip>
                            </UFormField>

                            <!-- Predefined Dispatch Status Reason -->
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="dispatchStatus"
                                :label="`${$t('common.dispatches')} ${$t('common.status')}`"
                                :ui="{ container: '' }"
                            >
                                <div class="flex flex-col gap-1">
                                    <div
                                        v-for="(_, idx) in state.predefinedStatus.dispatchStatus"
                                        :key="idx"
                                        class="flex items-center gap-1"
                                    >
                                        <UFormField class="flex-1" :name="`dispatchStatus.${idx}`" :ui="{ container: '' }">
                                            <UInput
                                                v-model="state.predefinedStatus.dispatchStatus[idx]"
                                                class="w-full flex-1"
                                                type="text"
                                                :placeholder="$t('common.reason')"
                                                :disabled="!canSubmit"
                                            />
                                        </UFormField>

                                        <UTooltip :text="$t('common.delete')">
                                            <UButton
                                                icon="i-mdi-close"
                                                :disabled="!canSubmit"
                                                @click="state.predefinedStatus.dispatchStatus.splice(idx, 1)"
                                            />
                                        </UTooltip>
                                    </div>
                                </div>

                                <UTooltip :text="$t('common.add')">
                                    <UButton
                                        :class="state.predefinedStatus.dispatchStatus.length ? 'mt-2' : ''"
                                        icon="i-mdi-plus"
                                        :disabled="!canSubmit || state.predefinedStatus.dispatchStatus.length >= 8"
                                        @click="state.predefinedStatus.dispatchStatus.push('')"
                                    />
                                </UTooltip>
                            </UFormField>
                        </UPageCard>
                    </template>

                    <template #timings>
                        <UPageCard
                            :title="$t('components.centrum.settings.timings.title')"
                            :description="$t('components.centrum.settings.timings.description')"
                        >
                            <!-- Timings -->
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="timings.dispatchMaxWait"
                                :label="$t('components.centrum.settings.timings.dispatch_max_wait')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.timings.dispatchMaxWait"
                                    type="number"
                                    :min="30"
                                    :placeholder="$t('common.time_ago.second', 2)"
                                    trailing-icon="i-mdi-access-time"
                                    :disabled="!canSubmit"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="timings.requireUnit"
                                :label="$t('components.centrum.settings.timings.require_unit')"
                                :ui="{ container: '' }"
                            >
                                <USwitch v-model="state.timings.requireUnit" :disabled="!canSubmit" />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="timings.requireUnitReminderSeconds"
                                :label="$t('components.centrum.settings.timings.require_unit_reminder_seconds')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.timings.requireUnitReminderSeconds"
                                    type="number"
                                    :min="60"
                                    :placeholder="$t('common.time_ago.second', 2)"
                                    trailing-icon="i-mdi-access-time"
                                    :disabled="!canSubmit"
                                />
                            </UFormField>
                        </UPageCard>
                    </template>

                    <template #access>
                        <UPageCard
                            :title="$t('components.centrum.settings.access.title')"
                            :description="$t('components.centrum.settings.access.description')"
                        >
                            <!-- Access -->
                            <UFormField name="access" :label="$t('common.access')">
                                <AccessManager
                                    v-model:jobs="state.access.jobs"
                                    :target-id="0"
                                    :access-roles="
                                        enumToAccessLevelEnums(
                                            CentrumAccessLevel,
                                            'enums.centrum.CentrumAccessLevel',
                                            (val) => val > CentrumAccessLevel.BLOCKED,
                                        )
                                    "
                                    :access-types="[{ label: $t('common.job', 2), value: 'job' }]"
                                    hide-grade
                                    :hide-jobs="[activeChar!.job]"
                                    :disabled="!attr('centrum.CentrumService/UpdateSettings', 'Access', 'Shared').value"
                                    name="access"
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.centrum.settings.offered_access.title')"
                            :description="$t('components.centrum.settings.offered_access.description')"
                        >
                            <UFormField name="offeredAccess" :label="$t('common.access')">
                                <UPageGrid class="mt-2">
                                    <UCard v-for="(ja, idx) in state.offeredAccess.jobs" :key="ja.id">
                                        <div class="flex items-center gap-2">
                                            <USwitch
                                                :model-value="state.offeredAccess.jobs[idx]!.acceptedAt !== undefined"
                                                @update:model-value="
                                                    $event
                                                        ? (state.offeredAccess.jobs[idx]!.acceptedAt = toTimestamp(new Date()))
                                                        : (state.offeredAccess.jobs[idx]!.acceptedAt = undefined)
                                                "
                                            />
                                            <h3 class="flex-1 text-lg">{{ ja.jobLabel ?? ja.job }}</h3>

                                            <UBadge
                                                :label="`${$t('common.access')}: ${$t('enums.centrum.CentrumAccessLevel.' + CentrumAccessLevel[ja.access])}`"
                                            />
                                        </div>
                                    </UCard>
                                </UPageGrid>
                            </UFormField>
                        </UPageCard>
                    </template>
                </UTabs>
            </UForm>
        </template>
    </UDashboardPanel>
</template>
