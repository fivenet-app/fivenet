<script lang="ts" setup>
import { checkUnitAccess, defaultUnitIcon, unitStatusToBGColor, unitStatusToIcon } from '~/components/dispatch/helpers';
import UnitAttributes from '~/components/dispatch/partials/UnitAttributes.vue';
import UnitAssignUsersModal from '~/components/dispatch/units/UnitAssignUsersModal.vue';
import UnitFeed from '~/components/dispatch/units/UnitFeed.vue';
import UnitStatusUpdateModal from '~/components/dispatch/units/UnitStatusUpdateModal.vue';
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useLivemapStore } from '~/stores/livemap';
import { UnitAccessLevel } from '~~/gen/ts/resources/centrum/units/access/access';
import { type Unit, StatusUnit } from '~~/gen/ts/resources/centrum/units/units';

const props = defineProps<{
    unit: Unit;
    statusSelected?: StatusUnit;
}>();

const open = defineModel<boolean>('open', { default: true });

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { canDo } = useCentrumStore();
const { gotoCoords } = useLivemapStore();

const overlay = useOverlay();

const unitStatusUpdateModal = overlay.create(UnitStatusUpdateModal);
const unitAssignUsersModal = overlay.create(UnitAssignUsersModal);
const unitFeed = useTemplateRef<InstanceType<typeof UnitFeed>>('unitFeed');
const activityPage = ref(1);

const unitStatusColor = computed(() => unitStatusToBGColor(props.unit.status?.status));
const hasAccess = computed(
    () => (props.unit.access?.jobs?.length ?? 0) > 0 || (props.unit.access?.qualifications?.length ?? 0) > 0,
);

function refreshActivity(): Promise<unknown> {
    return unitFeed.value?.refresh() ?? Promise.resolve();
}
</script>

<template>
    <USlideover v-model:open="open" :overlay="false" :ui="{ content: 'sm:max-w-lg', body: 'p-3 sm:p-3' }">
        <template #title>
            <div class="flex min-w-0 items-center gap-2">
                <span class="shrink-0">{{ $t('common.unit') }}</span>

                <UIcon
                    v-if="unit.icon && unit.icon !== defaultUnitIcon"
                    class="size-5 shrink-0"
                    :name="convertComponentIconNameToDynamic(unit.icon)"
                    :style="{ color: unit.color ?? 'currentColor' }"
                />

                <span class="shrink-0 font-semibold text-highlighted">{{ unit.initials }}</span>
                <span class="truncate">{{ unit.name }}</span>
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
                                :class="unitStatusColor"
                                :disabled="!checkUnitAccess(unit.access, UnitAccessLevel.JOIN)"
                                :icon="unitStatusToIcon(props.unit.status?.status)"
                                :label="$t(`enums.centrum.StatusUnit.${StatusUnit[unit.status?.status ?? 0]}`)"
                                @click="
                                    unitStatusUpdateModal.open({
                                        unit: unit,
                                        status: statusSelected,
                                    })
                                "
                            />
                        </div>

                        <div class="space-y-1">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.last_update') }}</p>
                            <GenericTime :value="unit.status?.createdAt" />
                        </div>

                        <div class="space-y-1">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.code') }}</p>
                            <p class="text-sm text-highlighted">{{ unit.status?.code ?? $t('common.na') }}</p>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.reason') }}</p>
                            <p class="text-sm break-words text-highlighted">{{ unit.status?.reason ?? $t('common.na') }}</p>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">{{ $t('common.location') }}</p>
                            <div class="flex flex-wrap items-center gap-x-2 gap-y-1 text-sm">
                                <span>
                                    {{ $t('common.postal') }}:
                                    {{ unit.status?.postal ?? $t('common.na') }}
                                </span>

                                <UButton
                                    v-if="unit.status?.x !== undefined && unit.status?.y !== undefined"
                                    size="xs"
                                    variant="link"
                                    icon="i-mdi-map-marker"
                                    :label="$t('common.go_to_location')"
                                    @click="gotoCoords({ x: unit.status?.x, y: unit.status?.y })"
                                />
                                <span v-else class="text-muted">{{ $t('common.no_location') }}</span>
                            </div>
                        </div>
                    </div>
                </UCard>

                <UCard variant="subtle" :ui="{ header: 'p-3 sm:p-3', body: 'p-3 sm:p-3' }">
                    <template #header>
                        <div class="flex items-center justify-between gap-2">
                            <div class="flex items-center gap-2">
                                <UIcon name="i-mdi-account-group-outline" class="size-5 text-primary" />
                                <h2 class="font-semibold text-highlighted">{{ $t('common.members') }}</h2>
                                <UBadge color="neutral" variant="soft" :label="`${unit.users.length}`" />
                            </div>

                            <UButton
                                v-if="canDo('TakeControl')"
                                size="sm"
                                variant="soft"
                                icon="i-mdi-pencil"
                                :label="$t('common.assign')"
                                @click="unitAssignUsersModal.open({ unit: unit })"
                            />
                        </div>
                    </template>

                    <p v-if="unit.users.length === 0" class="text-sm text-muted">{{ $t('common.member', 0) }}</p>
                    <div v-else class="overflow-hidden rounded-md border border-default bg-elevated">
                        <ul class="divide-y divide-default text-sm font-medium" role="list">
                            <li v-for="user in unit.users" :key="user.userId" class="flex items-center py-3 pr-4 pl-3">
                                <CitizenInfoPopover
                                    class="flex items-center justify-center"
                                    :user="user.user"
                                    show-avatar-in-name
                                />
                            </li>
                        </ul>
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
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('common.description') }}
                            </dt>
                            <dd class="text-sm break-words whitespace-pre-wrap text-highlighted">
                                {{ unit.description ?? $t('common.na') }}
                            </dd>
                        </div>

                        <div class="space-y-1">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ `${$t('common.department')} ${$t('common.postal')}` }}
                            </dt>
                            <dd class="text-sm text-highlighted">{{ unit.homePostal ?? $t('common.na') }}</dd>
                        </div>

                        <div class="space-y-1 sm:col-span-2">
                            <dt class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('common.attributes', 2) }}
                            </dt>
                            <dd class="text-sm text-highlighted"><UnitAttributes :attributes="unit.attributes" /></dd>
                        </div>
                    </dl>
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
                        <div class="flex items-center justify-between gap-2">
                            <div class="flex items-center gap-2">
                                <UIcon name="i-mdi-pulse" class="size-5 text-primary" />
                                <h2 class="font-semibold text-highlighted">{{ $t('common.feed') }}</h2>
                            </div>
                        </div>
                    </template>

                    <UnitFeed ref="unitFeed" v-model:page="activityPage" :unit-id="unit.id" />

                    <template #footer>
                        <Pagination
                            v-model="activityPage"
                            :pagination="unitFeed?.pagination"
                            disable-border
                            hide-text
                            :status="unitFeed?.status"
                            :refresh="refreshActivity"
                        />
                    </template>
                </UCard>

                <UCard variant="subtle" :ui="{ header: 'p-3 sm:p-3', body: 'p-3 sm:p-3' }">
                    <template #header>
                        <div class="flex items-center gap-2">
                            <UIcon name="i-mdi-shield-account-outline" class="size-5 text-primary" />
                            <h2 class="font-semibold text-highlighted">{{ $t('common.access') }}</h2>
                        </div>
                    </template>

                    <AccessBadges
                        v-if="hasAccess"
                        :access-level="UnitAccessLevel"
                        :jobs="unit.access?.jobs"
                        :qualifications="unit.access?.qualifications"
                        i18n-key="enums.centrum"
                        i18n-access-level-key="UnitAccessLevel"
                    />
                    <p v-else class="text-sm text-muted">{{ $t('common.no_access') }}</p>
                </UCard>
            </div>
        </template>

        <template #footer>
            <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="emit('close', false)" />
        </template>
    </USlideover>
</template>
