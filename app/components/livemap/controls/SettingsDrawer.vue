<script lang="ts" setup>
import { z } from 'zod';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { can } = useAuth();
const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const livemapStore = useLivemapStore();
const {
    markerDragEnabled,
    initiated,
    abort: livemapAbort,
    stopping: livemapStopping,
    reconnectBackoffTime: livemapReconnectBackoffTime,
    error: livemapError,
    jobsUsers,
    jobsMarkers,
} = storeToRefs(livemapStore);

const centrumStore = useCentrumStore();
const {
    abort: centrumAbort,
    stopping: centrumStopping,
    reconnectBackoffTime: centrumReconnectBackoffTime,
    error: centrumError,
    cleanupIntervalId: centrumCleanupIntervalId,
    acls: centrumAcls,
} = storeToRefs(centrumStore);

const livemapStreamStatus = computed(() => {
    if (livemapStopping.value) return 'stopping';
    if (livemapError.value) return 'error';
    if (!livemapAbort.value) return 'stopped';
    if (livemapAbort.value.signal.aborted) return 'aborted';
    return 'running';
});

const centrumStreamStatus = computed(() => {
    if (centrumStopping.value) return 'stopping';
    if (centrumError.value) return 'error';
    if (!centrumAbort.value) return 'stopped';
    if (centrumAbort.value.signal.aborted) return 'aborted';
    return 'running';
});

const livemapErrorText = computed(() =>
    livemapError.value ? `${livemapError.value.code}: ${livemapError.value.message}` : 'none',
);
const centrumErrorText = computed(() =>
    centrumError.value ? `${centrumError.value.code}: ${centrumError.value.message}` : 'none',
);
const centrumAclJobsCount = computed(() => centrumAcls.value?.dispatches?.jobs.length ?? 0);

const collectLivemapDebugInfo = (): string =>
    JSON.stringify(
        {
            livemap: {
                initiated: initiated.value,
                streamStatus: livemapStreamStatus.value,
                abortController: livemapAbort.value !== undefined,
                signalAborted: livemapAbort.value?.signal.aborted ?? false,
                stopping: livemapStopping.value,
                reconnectBackoffTime: livemapReconnectBackoffTime.value,
                jobsUsers: jobsUsers.value.length,
                jobsMarkers: jobsMarkers.value.length,
                error: livemapErrorText.value,
            },
            centrum: {
                streamStatus: centrumStreamStatus.value,
                abortController: centrumAbort.value !== undefined,
                signalAborted: centrumAbort.value?.signal.aborted ?? false,
                stopping: centrumStopping.value,
                reconnectBackoffTime: centrumReconnectBackoffTime.value,
                cleanupInterval: centrumCleanupIntervalId.value !== undefined,
                aclJobs: centrumAclJobsCount.value,
                error: centrumErrorText.value,
            },
        },
        undefined,
        2,
    );

function copyLivemapDebugInfo(): void {
    copyToClipboardWrapper(collectLivemapDebugInfo());

    notifications.add({
        title: { key: 'notifications.clipboard.copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const schema = z.object({
    markerSize: z.coerce.number().min(14).max(32),
    centerSelectedMarker: z.coerce.boolean(),
    showUnitNames: z.coerce.boolean(),
    showUnitStatus: z.coerce.boolean(),
    showAllDispatches: z.coerce.boolean(),
    showGrid: z.coerce.boolean(),
    useUnitColor: z.coerce.boolean(),
});
</script>

<template>
    <UDrawer
        :title="$t('common.setting', 2)"
        :overlay="false"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ title: 'flex flex-row gap-2' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('common.setting', 2) }}</span>
        </template>

        <template #body>
            <UContainer class="max-w-xl">
                <UForm :schema="schema" :state="livemap">
                    <UFormField name="centerSelectedMarker" :label="$t('components.livemap.center_selected_marker')">
                        <USwitch v-model="livemap.centerSelectedMarker" />
                    </UFormField>

                    <UFormField name="markerSize" :label="$t('components.livemap.settings.marker_size')">
                        <USlider
                            v-model="livemap.markerSize"
                            class="my-auto h-1.5 w-full cursor-grab rounded-full"
                            :min="14"
                            :max="32"
                            :step="2"
                        />
                        <span class="text-sm">{{ livemap.markerSize }}</span>
                    </UFormField>

                    <UFormField name="showUnitNames" :label="$t('components.livemap.show_unit_names')">
                        <USwitch v-model="livemap.showUnitNames" />
                    </UFormField>

                    <UFormField name="showUnitStatus" :label="$t('components.livemap.show_unit_status')">
                        <USwitch v-model="livemap.showUnitStatus" />
                    </UFormField>

                    <UFormField name="showAllDispatches" :label="$t('components.livemap.show_all_dispatches')">
                        <USwitch v-model="livemap.showAllDispatches" />
                    </UFormField>

                    <UFormField name="showGrid" :label="$t('components.livemap.show_grid')">
                        <USwitch v-model="livemap.showGrid" />
                    </UFormField>

                    <UFormField name="useUnitColor" :label="$t('components.livemap.use_unit_color')">
                        <USwitch v-model="livemap.useUnitColor" />
                    </UFormField>

                    <UFormField
                        v-if="can('livemap.LivemapService/CreateOrUpdateMarker').value"
                        name="markerDragEnabled"
                        :label="$t('components.livemap.enable_marker_dragging')"
                    >
                        <USwitch v-model="markerDragEnabled" />
                    </UFormField>
                </UForm>

                <USeparator class="my-4" />

                <UCollapsible>
                    <UButton
                        :label="$t('components.debug_info.title')"
                        color="neutral"
                        variant="subtle"
                        trailing-icon="i-mdi-chevron-down"
                        block
                    />

                    <template #content>
                        <div class="mt-2 space-y-4 rounded-md border border-default p-3 text-sm">
                            <UTooltip :text="$t('common.copy')">
                                <UButton
                                    color="neutral"
                                    icon="i-mdi-content-copy"
                                    size="xs"
                                    variant="outline"
                                    :label="$t('common.copy')"
                                    @click="copyLivemapDebugInfo"
                                />
                            </UTooltip>

                            <section class="space-y-2">
                                <h3 class="font-semibold">Livemap Stream</h3>
                                <dl class="grid grid-cols-2 gap-x-4 gap-y-1">
                                    <dt class="text-muted">Initiated</dt>
                                    <dd class="font-mono">{{ initiated }}</dd>
                                    <dt class="text-muted">Status</dt>
                                    <dd class="font-mono">{{ livemapStreamStatus }}</dd>
                                    <dt class="text-muted">Abort Controller</dt>
                                    <dd class="font-mono">{{ livemapAbort !== undefined }}</dd>
                                    <dt class="text-muted">Signal Aborted</dt>
                                    <dd class="font-mono">{{ livemapAbort?.signal.aborted ?? false }}</dd>
                                    <dt class="text-muted">Stopping</dt>
                                    <dd class="font-mono">{{ livemapStopping }}</dd>
                                    <dt class="text-muted">Reconnect Backoff</dt>
                                    <dd class="font-mono">{{ livemapReconnectBackoffTime }}s</dd>
                                    <dt class="text-muted">Jobs Users</dt>
                                    <dd class="font-mono">{{ jobsUsers.length }}</dd>
                                    <dt class="text-muted">Jobs Markers</dt>
                                    <dd class="font-mono">{{ jobsMarkers.length }}</dd>
                                    <dt class="text-muted">Error</dt>
                                    <dd class="font-mono break-all">{{ livemapErrorText }}</dd>
                                </dl>
                            </section>

                            <section class="space-y-2">
                                <h3 class="font-semibold">Centrum Stream</h3>
                                <dl class="grid grid-cols-2 gap-x-4 gap-y-1">
                                    <dt class="text-muted">Status</dt>
                                    <dd class="font-mono">{{ centrumStreamStatus }}</dd>
                                    <dt class="text-muted">Abort Controller</dt>
                                    <dd class="font-mono">{{ centrumAbort !== undefined }}</dd>
                                    <dt class="text-muted">Signal Aborted</dt>
                                    <dd class="font-mono">{{ centrumAbort?.signal.aborted ?? false }}</dd>
                                    <dt class="text-muted">Stopping</dt>
                                    <dd class="font-mono">{{ centrumStopping }}</dd>
                                    <dt class="text-muted">Reconnect Backoff</dt>
                                    <dd class="font-mono">{{ centrumReconnectBackoffTime }}s</dd>
                                    <dt class="text-muted">Cleanup Interval</dt>
                                    <dd class="font-mono">{{ centrumCleanupIntervalId !== undefined }}</dd>
                                    <dt class="text-muted">ACL Jobs</dt>
                                    <dd class="font-mono">{{ centrumAclJobsCount }}</dd>
                                    <dt class="text-muted">Error</dt>
                                    <dd class="font-mono break-all">{{ centrumErrorText }}</dd>
                                </dl>
                            </section>
                        </div>
                    </template>
                </UCollapsible>
            </UContainer>
        </template>

        <template #footer>
            <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />
        </template>
    </UDrawer>
</template>
