<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useNotificatorStore } from '~/stores/notificator';
import type { Settings } from '~~/gen/ts/resources/centrum/settings';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { isSuperuser } = useAuth();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const {
    data: settings,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData('rector-centrum-settings', () => getCentrumSettings());

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
    enabled: z.boolean(),
    mode: z.nativeEnum(CentrumMode),
    fallbackMode: z.nativeEnum(CentrumMode),
    predefinedStatus: z.object({
        unitStatus: z.string().array().max(10),
        dispatchStatus: z.string().array().max(10),
    }),
    timings: z.object({
        dispatchMaxWait: z.coerce.number().min(30).max(6000),
        requireUnit: z.boolean(),
        requireUnitReminderSeconds: z.coerce.number().min(60).max(6000),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    enabled: false,
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
});

async function updateSettings(values: Schema): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.updateSettings({
            settings: {
                job: '',
                enabled: values.enabled,
                mode: values.mode,
                fallbackMode: values.fallbackMode,
                predefinedStatus: values.predefinedStatus,
                timings: values.timings,
            },
        });
        await call;

        refresh();

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

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
    { slot: 'timings', label: t('common.timings'), icon: 'i-mdi-access-time' },
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
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
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
            <UTabs v-model="selectedTab" :items="items" class="w-full" :ui="{ list: { rounded: '' } }">
                <template #settings>
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.centrum.settings.title')"
                            :description="$t('components.centrum.settings.description')"
                        >
                            <UFormGroup
                                name="enabled"
                                :label="$t('common.enabled')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UToggle
                                    v-model="state.enabled"
                                    name="enabled"
                                    :disabled="!isSuperuser"
                                    :placeholder="$t('common.enabled')"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="mode"
                                :label="$t('common.mode')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.mode"
                                        :options="modes"
                                        value-attribute="mode"
                                        :searchable-placeholder="$t('common.search_field')"
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
                                name="fallbackMode"
                                :label="$t('common.fallback_mode')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.fallbackMode"
                                        :options="modes"
                                        value-attribute="mode"
                                        :searchable-placeholder="$t('common.search_field')"
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
                                name="unitStatus"
                                :label="`${$t('common.unit')} ${$t('common.status')}`"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <div class="flex flex-col gap-1">
                                    <div
                                        v-for="(_, idx) in state.predefinedStatus.unitStatus"
                                        :key="idx"
                                        class="flex items-center gap-1"
                                    >
                                        <UFormGroup :name="`unitStatus.${idx}`" :ui="{ container: '' }" class="flex-1">
                                            <UInput
                                                v-model="state.predefinedStatus.unitStatus[idx]"
                                                type="text"
                                                class="w-full flex-1"
                                                :placeholder="$t('common.reason')"
                                            />
                                        </UFormGroup>

                                        <UTooltip :text="$t('common.delete')">
                                            <UButton
                                                :ui="{ rounded: 'rounded-full' }"
                                                icon="i-mdi-close"
                                                @click="state.predefinedStatus.unitStatus.splice(idx, 1)"
                                            />
                                        </UTooltip>
                                    </div>
                                </div>

                                <UTooltip :text="$t('common.add')">
                                    <UButton
                                        :ui="{ rounded: 'rounded-full' }"
                                        icon="i-mdi-plus"
                                        :disabled="!canSubmit || state.predefinedStatus.unitStatus.length >= 8"
                                        :class="state.predefinedStatus.unitStatus.length ? 'mt-2' : ''"
                                        @click="state.predefinedStatus.unitStatus.push('')"
                                    />
                                </UTooltip>
                            </UFormGroup>

                            <!-- Predefined Dispatch Status Reason -->
                            <UFormGroup
                                name="dispatchStatus"
                                :label="`${$t('common.dispatches')} ${$t('common.status')}`"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <div class="flex flex-col gap-1">
                                    <div
                                        v-for="(_, idx) in state.predefinedStatus.dispatchStatus"
                                        :key="idx"
                                        class="flex items-center gap-1"
                                    >
                                        <UFormGroup :name="`dispatchStatus.${idx}`" :ui="{ container: '' }" class="flex-1">
                                            <UInput
                                                v-model="state.predefinedStatus.dispatchStatus[idx]"
                                                type="text"
                                                class="w-full flex-1"
                                                :placeholder="$t('common.reason')"
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
                                        :ui="{ rounded: 'rounded-full' }"
                                        icon="i-mdi-plus"
                                        :disabled="!canSubmit || state.predefinedStatus.dispatchStatus.length >= 8"
                                        :class="state.predefinedStatus.dispatchStatus.length ? 'mt-2' : ''"
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
                                name="timings.dispatchMaxWait"
                                :label="$t('components.centrum.settings.timings.dispatch_max_wait')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.timings.dispatchMaxWait"
                                    type="number"
                                    :placeholder="$t('common.time_ago.second', 2)"
                                    trailing-icon="i-mdi-access-time"
                                />
                            </UFormGroup>

                            <UFormGroup
                                name="timings.requireUnit"
                                :label="$t('components.centrum.settings.timings.require_unit')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UToggle v-model="state.timings.requireUnit" />
                            </UFormGroup>

                            <UFormGroup
                                name="timings.requireUnitReminderSeconds"
                                :label="$t('components.centrum.settings.timings.require_unit_reminder_seconds')"
                                class="grid grid-cols-2 items-center gap-2"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.timings.requireUnitReminderSeconds"
                                    type="number"
                                    :placeholder="$t('common.time_ago.second', 2)"
                                    trailing-icon="i-mdi-access-time"
                                />
                            </UFormGroup>
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </template>
            </UTabs>
        </template>
    </UForm>
</template>
