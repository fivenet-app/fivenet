<script lang="ts" setup>
import DispatchAssignModal from '~/components/dispatch/dispatches/DispatchAssignModal.vue';
import DispatchFeed from '~/components/dispatch/dispatches/DispatchFeed.vue';
import DispatchStatusUpdateModal from '~/components/dispatch/dispatches/DispatchStatusUpdateModal.vue';
import { checkDispatchAccess, dispatchStatusToButtonColor, dispatchStatusToIcon } from '~/components/dispatch/helpers';
import DispatchAttributes from '~/components/dispatch/partials/DispatchAttributes.vue';
import DispatchReferences from '~/components/dispatch/partials/DispatchReferences.vue';
import UnitInfoPopover from '~/components/dispatch/units/UnitInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { getCentrumDispatchesClient } from '~~/gen/ts/clients';
import { CentrumAccessLevel } from '~~/gen/ts/resources/centrum/access/access';
import { type Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    dispatchId: number;
    dispatch?: Dispatch;
}>();

const open = defineModel<boolean>('open', { default: true });

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { can } = useAuth();
const overlay = useOverlay();
const { gotoCoords } = useLivemapStore();

const centrumStore = useCentrumStore();
const { dispatches, timeCorrection } = storeToRefs(centrumStore);
const { canDo, selfAssign } = centrumStore;
const notifications = useNotificationsStore();
const now = useSecondClock();
const formatTimeAgo = useLocaleTimeAgoFormatter();

const centrumDispatchesClient = await getCentrumDispatchesClient();

const dispatch = computed(() => (props.dispatch ? props.dispatch : dispatches.value.get(props.dispatchId)));
const dispatchFeed = useTemplateRef<InstanceType<typeof DispatchFeed>>('dispatchFeed');
const activityPage = ref(1);

async function deleteDispatch(id: number): Promise<void> {
    try {
        const call = centrumDispatchesClient.deleteDispatch({ id });
        await call;

        notifications.add({
            title: { key: 'notifications.centrum.sidebar.dispatch_deleted.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.dispatch_deleted.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const dispatchStatusColor = computed(() => dispatchStatusToButtonColor(dispatch.value?.status?.status));

watch(dispatch, () => {
    if (dispatch.value === undefined) {
        emit('close', false);
    }
});

const canAccessDispatch = computed(() => ({
    participate: checkDispatchAccess(dispatch.value?.jobs, CentrumAccessLevel.PARTICIPATE),
    dispatch: checkDispatchAccess(dispatch.value?.jobs, CentrumAccessLevel.DISPATCH),
}));

function refreshActivity(): Promise<unknown> {
    return dispatchFeed.value?.refresh() ?? Promise.resolve();
}

const confirmModal = overlay.create(ConfirmModal);
const dispatchAssignModal = overlay.create(DispatchAssignModal);
const dispatchStatusUpdateModal = overlay.create(DispatchStatusUpdateModal);
</script>

<template>
    <USlideover v-model:open="open" :overlay="false" :ui="{ content: 'sm:max-w-lg', body: 'p-3 sm:p-3' }">
        <template #title>
            <div class="flex min-w-0 items-center gap-2">
                <IDCopyBadge :id="dispatch?.id ?? 0" prefix="DSP" />
                <p class="truncate text-highlighted">{{ dispatch?.message ?? $t('common.na') }}</p>
            </div>
        </template>

        <template #body>
            <div class="flex flex-col gap-4">
                <UCard variant="subtle" :ui="{ header: 'p-3 sm:p-3', body: 'p-3 sm:p-3' }">
                    <template #header>
                        <div class="flex items-center gap-2">
                            <UIcon name="i-mdi-information-outline" class="size-5 text-primary" />
                            <h2 class="font-semibold text-highlighted">{{ $t('common.status') }}</h2>
                        </div>
                    </template>

                    <div class="grid gap-4 sm:grid-cols-2">
                        <div class="space-y-1 sm:col-span-2">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.status') }}</p>
                            <UButton
                                class="font-semibold shadow-xs"
                                :color="dispatchStatusColor"
                                :icon="dispatchStatusToIcon(dispatch?.status?.status)"
                                :disabled="!canAccessDispatch.participate"
                                :label="$t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch?.status?.status ?? 0]}`)"
                                @click="dispatch && dispatchStatusUpdateModal.open({ dispatchId: dispatch.id })"
                            />
                        </div>

                        <div class="space-y-1">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.last_update') }}</p>
                            <GenericTime v-if="dispatch?.status?.createdAt" :value="dispatch.status.createdAt" />
                            <span v-else class="text-sm text-muted">{{ $t('common.na') }}</span>
                        </div>

                        <div class="space-y-1">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.code') }}</p>
                            <p class="text-sm text-highlighted">{{ dispatch?.status?.code ?? $t('common.na') }}</p>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.reason') }}</p>
                            <p class="text-sm break-words text-highlighted">
                                {{ dispatch?.status?.reason ?? $t('common.na') }}
                            </p>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.location') }}</p>
                            <div class="flex flex-wrap items-center gap-x-2 gap-y-1 text-sm">
                                <span>
                                    {{ $t('common.postal') }}:
                                    {{ dispatch?.status?.postal ?? $t('common.na') }}
                                </span>

                                <UButton
                                    v-if="dispatch?.status?.x !== undefined && dispatch?.status?.y !== undefined"
                                    size="xs"
                                    variant="link"
                                    icon="i-mdi-map-marker"
                                    :label="$t('common.go_to_location')"
                                    @click="gotoCoords({ x: dispatch.status.x, y: dispatch.status.y })"
                                />
                                <span v-else class="text-muted">{{ $t('common.no_location') }}</span>
                            </div>
                        </div>
                    </div>
                </UCard>

                <UCard variant="subtle" :ui="{ header: 'p-3 sm:p-3', body: 'p-3 sm:p-3' }">
                    <template #header>
                        <div class="flex items-center gap-2">
                            <UIcon name="i-mdi-card-account-details-outline" class="size-5 text-primary" />
                            <h2 class="font-semibold text-highlighted">{{ $t('common.info') }}</h2>
                        </div>
                    </template>

                    <dl class="grid gap-4 sm:grid-cols-2">
                        <div class="space-y-1 sm:col-span-2">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.job') }}</dt>
                            <dd class="text-sm break-words text-highlighted">
                                {{ dispatch?.jobs?.jobs?.map((job) => job.label ?? job.name).join(', ') || $t('common.na') }}
                            </dd>
                        </div>

                        <div class="space-y-1">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.sent_at') }}</dt>
                            <dd class="text-sm text-highlighted">
                                <GenericTime v-if="dispatch?.createdAt" :value="dispatch.createdAt" />
                                <span v-else class="text-muted">{{ $t('common.na') }}</span>
                            </dd>
                        </div>

                        <div class="space-y-1">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.sent_by') }}</dt>
                            <dd class="text-sm text-highlighted">
                                <span v-if="dispatch?.anon">{{ $t('common.anon') }}</span>
                                <CitizenInfoPopover v-else-if="dispatch?.creator" :user="dispatch.creator" />
                                <span v-else class="text-muted">{{ $t('common.unknown') }}</span>
                            </dd>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.location') }}</dt>
                            <dd class="flex flex-wrap items-center gap-x-2 gap-y-1 text-sm text-highlighted">
                                <span>
                                    {{ $t('common.postal') }}:
                                    {{ dispatch?.postal ?? $t('common.na') }}
                                </span>
                                <UButton
                                    v-if="dispatch"
                                    size="xs"
                                    variant="link"
                                    icon="i-mdi-map-marker"
                                    :label="$t('common.go_to_location')"
                                    @click="gotoCoords({ x: dispatch.x, y: dispatch.y })"
                                />
                            </dd>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('common.description') }}
                            </dt>
                            <dd class="text-sm break-words whitespace-pre-wrap text-highlighted">
                                {{ dispatch?.description ?? $t('common.na') }}
                            </dd>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('common.attributes', 2) }}
                            </dt>
                            <dd class="text-sm text-highlighted"><DispatchAttributes :attributes="dispatch?.attributes" /></dd>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('common.reference', 2) }}
                            </dt>
                            <dd class="text-sm text-highlighted"><DispatchReferences :references="dispatch?.references" /></dd>
                        </div>
                    </dl>
                </UCard>

                <UCard variant="subtle" :ui="{ header: 'p-3 sm:p-3', body: 'p-3 sm:p-3' }">
                    <template #header>
                        <div class="flex items-center justify-between gap-2">
                            <div class="flex items-center gap-2">
                                <UIcon name="i-mdi-account-group-outline" class="size-5 text-primary" />
                                <h2 class="font-semibold text-highlighted">{{ $t('common.unit', 2) }}</h2>
                                <UBadge color="neutral" variant="soft" :label="`${dispatch?.units?.length ?? 0}`" />
                            </div>

                            <UButton
                                v-if="canDo('TakeControl') && canAccessDispatch.dispatch"
                                size="sm"
                                variant="soft"
                                icon="i-mdi-pencil"
                                :label="$t('common.assign')"
                                @click="dispatchAssignModal.open({ dispatchId: dispatchId })"
                            />
                        </div>
                    </template>

                    <p v-if="!dispatch?.units?.length" class="text-sm text-muted">{{ $t('common.units', 0) }}</p>
                    <div v-else class="overflow-hidden rounded-md border border-default bg-elevated">
                        <ul class="divide-y divide-default text-sm font-medium" role="list">
                            <li
                                v-for="unit in dispatch.units"
                                :key="unit.unitId"
                                class="flex items-center justify-between py-3 pr-4 pl-3"
                            >
                                <div class="flex min-w-0 flex-1 items-center">
                                    <UnitInfoPopover
                                        class="flex items-center justify-center"
                                        :unit-id="unit.unitId"
                                        :unit="unit.unit"
                                        :assignment="unit"
                                        show-icon
                                        size="md"
                                    />

                                    <span
                                        v-if="unit.expiresAt"
                                        class="ml-2 inline-flex min-w-0 flex-1 items-center truncate text-muted"
                                    >
                                        {{ formatTimeAgo(toDate(unit.expiresAt, timeCorrection), { showSecond: true }, now) }}
                                    </span>
                                </div>
                            </li>
                        </ul>
                    </div>

                    <UFieldGroup v-if="canDo('TakeDispatch') && canAccessDispatch.participate" class="mt-3 flex w-full">
                        <UButton
                            class="flex-1"
                            icon="i-mdi-plus"
                            :label="$t('common.self_assign')"
                            size="sm"
                            @click="dispatch && selfAssign(dispatch.id)"
                        />
                    </UFieldGroup>
                </UCard>

                <UCard
                    variant="subtle"
                    :ui="{
                        header: 'p-3 sm:p-3',
                        body: 'max-h-[50rem] overflow-y-auto p-2 sm:p-2',
                        footer: 'px-2 py-1 sm:px-2 sm:py-1',
                    }"
                >
                    <template #header>
                        <div class="flex items-center gap-2">
                            <UIcon name="i-mdi-pulse" class="size-5 text-primary" />
                            <h2 class="font-semibold text-highlighted">{{ $t('common.feed') }}</h2>
                        </div>
                    </template>

                    <DispatchFeed ref="dispatchFeed" v-model:page="activityPage" :dispatch-id="dispatch?.id" />

                    <template #footer>
                        <Pagination
                            v-model="activityPage"
                            :pagination="dispatchFeed?.pagination"
                            disable-border
                            hide-text
                            :status="dispatchFeed?.status"
                            :refresh="refreshActivity"
                        />
                    </template>
                </UCard>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="emit('close', false)" />

                <UTooltip
                    v-if="can('centrum.DispatchesService/DeleteDispatch').value && canAccessDispatch.dispatch"
                    :text="$t('common.delete')"
                >
                    <UButton
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            () =>
                                confirmModal.open({
                                    confirm: async () => dispatch && deleteDispatch(dispatch.id),
                                })
                        "
                    />
                </UTooltip>
            </UFieldGroup>
        </template>
    </USlideover>
</template>
