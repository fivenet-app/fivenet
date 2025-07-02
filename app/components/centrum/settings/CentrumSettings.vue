<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { CentrumAccessLevel, type CentrumJobAccess } from '~~/gen/ts/resources/centrum/access';
import { CentrumMode, CentrumType, type Settings } from '~~/gen/ts/resources/centrum/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { isOpen } = useModal();

const { activeChar } = useAuth();

const notifications = useNotificationsStore();

const { maxAccessEntries } = useAppConfig();

const {
    data: settings,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData('settings-centrum-settings', () => getCentrumSettings());

async function getCentrumSettings(): Promise<Settings> {
    try {
        const call = $grpc.centrum.centrum.getSettings({});
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
    access: z.object({
        jobs: z.custom<CentrumJobAccess>().array().max(maxAccessEntries).default([]),
    }),
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
    configuration: {
        deduplicationEnabled: true,
        deduplicationRadius: 45,
        deduplicationDuration: 180,
    },
});

async function updateSettings(values: Schema): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.updateSettings({
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

        isOpen.value = false;
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
    state.access = settings.value.access ?? { jobs: [] };
    state.configuration = {
        deduplicationEnabled: settings.value.configuration?.deduplicationEnabled ?? true,
        deduplicationRadius: settings.value.configuration?.deduplicationRadius ?? 45,
        deduplicationDuration: settings.value.configuration?.deduplicationDuration?.seconds ?? 180,
    };
}

watch(settings, () => setSettingsValues());
setSettingsValues();

const items = [
    {
        slot: 'settings',
        label: t('common.settings'),
        icon: 'i-mdi-settings',
    },
    {
        slot: 'predefined',
        label: `${t('common.predefined', 2)} ${t('common.status', 2)}`,
        icon: 'i-mdi-selection',
    },
    {
        slot: 'timings',
        label: t('common.timings'),
        icon: 'i-mdi-access-time',
    },
    {
        slot: 'access',
        label: t('common.access'),
        icon: 'i-mdi-lock',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
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
</script>

<template>
    <UForm
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
        :schema="schema"
        :state="state"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('components.centrum.settings.title')">
            <template #right>
                <PartialsBackButton fallback-to="/centrum" />

                <UButton
                    v-if="!!settings"
                    type="submit"
                    trailing-icon="i-mdi-content-save"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                >
                    {{ $t('common.save', 1) }}
                </UButton>
            </template>
        </UDashboardNavbar>

        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.settings')])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.settings')])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="!settings" icon="i-mdi-tune" :type="$t('common.settings')" />

        <template v-else>
            <UDashboardPanelContent class="p-0 sm:pb-0">
                <UTabs
                    v-model="selectedTab"
                    class="flex flex-1 flex-col"
                    :items="items"
                    :ui="{
                        wrapper: 'space-y-0 overflow-y-hidden',
                        container: 'flex flex-1 flex-col overflow-y-hidden',
                        base: 'flex flex-1 flex-col overflow-y-hidden',
                        list: { rounded: '' },
                    }"
                >
                    <template #settings>
                        <UDashboardPanelContent>
                            <UDashboardSection
                                :title="$t('components.centrum.settings.title')"
                                :description="$t('components.centrum.settings.description')"
                            >
                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="enabled"
                                    :label="$t('common.enabled')"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.enabled" name="enabled" :disabled="!canSubmit" />
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="mode"
                                    :label="$t('common.mode')"
                                    :ui="{ container: '' }"
                                >
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.mode"
                                            :options="modes"
                                            value-attribute="mode"
                                            :searchable-placeholder="$t('common.search_field')"
                                            :disabled="!canSubmit"
                                        >
                                            <template #label>
                                                <span class="truncate">{{
                                                    $t(`enums.centrum.CentrumMode.${CentrumMode[state.mode ?? 0]}`)
                                                }}</span>
                                            </template>

                                            <template #option="{ option }">
                                                <span class="truncate">{{
                                                    $t(`enums.centrum.CentrumMode.${CentrumMode[option.mode ?? 0]}`)
                                                }}</span>
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="fallbackMode"
                                    :label="$t('common.fallback_mode')"
                                    :ui="{ container: '' }"
                                >
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.fallbackMode"
                                            :options="modes"
                                            value-attribute="mode"
                                            :searchable-placeholder="$t('common.search_field')"
                                            :disabled="!canSubmit"
                                        >
                                            <template #label>
                                                <span class="truncate">{{
                                                    $t(`enums.centrum.CentrumMode.${CentrumMode[state.fallbackMode ?? 0]}`)
                                                }}</span>
                                            </template>

                                            <template #option="{ option }">
                                                <span class="truncate">{{
                                                    $t(`enums.centrum.CentrumMode.${CentrumMode[option.mode ?? 0]}`)
                                                }}</span>
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormGroup>
                            </UDashboardSection>

                            <UDivider class="mb-4" />

                            <UDashboardSection
                                :title="$t('components.centrum.settings.deduplication.title')"
                                :description="$t('components.centrum.settings.deduplication.description')"
                            >
                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="configuration.deduplicationEnabled"
                                    :label="$t('common.enabled')"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.configuration.deduplicationEnabled" :disabled="!canSubmit" />
                                </UFormGroup>

                                <UFormGroup
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
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="configuration.deduplicationRadius"
                                    :label="$t('components.centrum.settings.deduplication.deduplication_radius')"
                                    :description="
                                        $t('components.centrum.settings.deduplication.deduplication_radius_description')
                                    "
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.configuration.deduplicationRadius"
                                        type="number"
                                        :min="5"
                                        :placeholder="$t('common.meters', 2)"
                                        :disabled="!canSubmit"
                                        :ui="{ base: '!pr-16' }"
                                    >
                                        <template #trailing>
                                            <span class="text-xs text-gray-500 dark:text-gray-400">
                                                {{ $t('common.meters', 2) }}
                                            </span>
                                        </template>
                                    </UInput>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #predefined>
                        <UDashboardPanelContent>
                            <UDashboardSection
                                :title="$t('components.centrum.settings.predefined.title')"
                                :description="$t('components.centrum.settings.predefined.description')"
                            >
                                <!-- Predefined Unit Status Reason -->
                                <UFormGroup
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
                                            <UFormGroup class="flex-1" :name="`unitStatus.${idx}`" :ui="{ container: '' }">
                                                <UInput
                                                    v-model="state.predefinedStatus.unitStatus[idx]"
                                                    class="w-full flex-1"
                                                    type="text"
                                                    :placeholder="$t('common.reason')"
                                                    :disabled="!canSubmit"
                                                />
                                            </UFormGroup>

                                            <UTooltip :text="$t('common.delete')">
                                                <UButton
                                                    :ui="{ rounded: 'rounded-full' }"
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
                                            :ui="{ rounded: 'rounded-full' }"
                                            icon="i-mdi-plus"
                                            :disabled="!canSubmit || state.predefinedStatus.unitStatus.length >= 8"
                                            @click="state.predefinedStatus.unitStatus.push('')"
                                        />
                                    </UTooltip>
                                </UFormGroup>

                                <!-- Predefined Dispatch Status Reason -->
                                <UFormGroup
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
                                            <UFormGroup class="flex-1" :name="`dispatchStatus.${idx}`" :ui="{ container: '' }">
                                                <UInput
                                                    v-model="state.predefinedStatus.dispatchStatus[idx]"
                                                    class="w-full flex-1"
                                                    type="text"
                                                    :placeholder="$t('common.reason')"
                                                    :disabled="!canSubmit"
                                                />
                                            </UFormGroup>

                                            <UTooltip :text="$t('common.delete')">
                                                <UButton
                                                    :ui="{ rounded: 'rounded-full' }"
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
                                            :ui="{ rounded: 'rounded-full' }"
                                            icon="i-mdi-plus"
                                            :disabled="!canSubmit || state.predefinedStatus.dispatchStatus.length >= 8"
                                            @click="state.predefinedStatus.dispatchStatus.push('')"
                                        />
                                    </UTooltip>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #timings>
                        <UDashboardPanelContent>
                            <UDashboardSection
                                :title="$t('components.centrum.settings.timings.title')"
                                :description="$t('components.centrum.settings.timings.description')"
                            >
                                <!-- Timings -->
                                <UFormGroup
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
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="timings.requireUnit"
                                    :label="$t('components.centrum.settings.timings.require_unit')"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.timings.requireUnit" :disabled="!canSubmit" />
                                </UFormGroup>

                                <UFormGroup
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
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #access>
                        <UDashboardPanelContent>
                            <UDashboardSection
                                :title="$t('components.centrum.settings.access.title')"
                                :description="$t('components.centrum.settings.access.description')"
                            >
                                <!-- Access -->
                                <UFormGroup name="access" :label="$t('common.access')">
                                    <AccessManager
                                        v-model:jobs="state.access.jobs"
                                        :target-id="0"
                                        :access-roles="
                                            enumToAccessLevelEnums(CentrumAccessLevel, 'enums.centrum.CentrumAccessLevel')
                                        "
                                        :access-types="[{ type: 'job', name: $t('common.job', 2) }]"
                                        hide-grade
                                        :hide-jobs="[activeChar!.job]"
                                    />
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>
                </UTabs>
            </UDashboardPanelContent>
        </template>
    </UForm>
</template>
