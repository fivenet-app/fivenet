<script lang="ts" setup>
import { checkUnitAccess, unitStatusToBGColor, unitStatusToIcon } from '~/components/centrum//helpers';
import UnitAttributes from '~/components/centrum/partials/UnitAttributes.vue';
import UnitAssignUsersModal from '~/components/centrum/units/UnitAssignUsersModal.vue';
import UnitFeed from '~/components/centrum/units/UnitFeed.vue';
import UnitStatusUpdateModal from '~/components/centrum/units/UnitStatusUpdateModal.vue';
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { type Unit, StatusUnit } from '~~/gen/ts/resources/centrum/units';
import { UnitAccessLevel } from '~~/gen/ts/resources/centrum/units_access';

const props = defineProps<{
    unit: Unit;
    statusSelected?: StatusUnit;
}>();

defineEmits<{
    close: [boolean];
}>();

const { canDo } = useCentrumStore();
const { goto } = useLivemapStore();

const overlay = useOverlay();

const unitStatusUpdateModal = overlay.create(UnitStatusUpdateModal);
const unitAssignUsersModal = overlay.create(UnitAssignUsersModal);

const unitStatusColors = computed(() => unitStatusToBGColor(props.unit.status?.status));
</script>

<template>
    <USlideover :title="`${$t('common.unit')} ${unit.initials} - ${unit.name}`" :overlay="false">
        <template #body>
            <dl class="divide-neutral/10 divide-y">
                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.description') }}
                    </dt>
                    <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                        <p class="max-h-14 overflow-y-scroll break-words">
                            {{ unit.description ?? $t('common.na') }}
                        </p>
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ `${$t('common.department')} ${$t('common.postal')}` }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        {{ unit.homePostal ?? $t('common.na') }}
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.last_update') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <GenericTime :value="unit.status?.createdAt" />
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.status') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <UButton
                            class="rounded-sm px-2 py-1 text-sm font-semibold shadow-xs"
                            :class="unitStatusColors"
                            :disabled="!checkUnitAccess(unit.access, UnitAccessLevel.JOIN)"
                            :icon="unitStatusToIcon(props.unit.status?.status)"
                            @click="
                                unitStatusUpdateModal.open({
                                    unit: unit,
                                    status: statusSelected,
                                })
                            "
                        >
                            {{ $t(`enums.centrum.StatusUnit.${StatusUnit[unit.status?.status ?? 0]}`) }}
                        </UButton>
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.code') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        {{ unit.status?.code ?? $t('common.na') }}
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.reason') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        {{ unit.status?.reason ?? $t('common.na') }}
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.location') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <div class="sm:inline-flex sm:flex-row sm:gap-2">
                            <span class="block">
                                {{ $t('common.postal') }}:
                                {{ unit.status?.postal ?? $t('common.na') }}
                            </span>

                            <UButton
                                v-if="unit.status?.x !== undefined && unit.status?.y !== undefined"
                                size="xs"
                                variant="link"
                                icon="i-mdi-map-marker"
                                @click="goto({ x: unit.status?.x, y: unit.status?.y })"
                            >
                                {{ $t('common.go_to_location') }}
                            </UButton>
                            <span v-else>{{ $t('common.no_location') }}</span>
                        </div>
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.attributes', 2) }}
                    </dt>
                    <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                        <UnitAttributes :attributes="unit.attributes" />
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.access') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <AccessBadges
                            :access-level="UnitAccessLevel"
                            :jobs="unit.access?.jobs"
                            :qualifications="unit.access?.qualifications"
                            i18n-key="enums.centrum"
                            i18n-access-level-key="UnitAccessLevel"
                        />
                    </dd>
                </div>

                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.members') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <span v-if="unit.users.length === 0" class="block">
                            {{ $t('common.member', 0) }}
                        </span>
                        <div v-else class="dark:bg-base-900 rounded-md bg-neutral-100">
                            <ul class="divide-y divide-gray-100 text-sm font-medium dark:divide-gray-800" role="list">
                                <li
                                    v-for="user in unit.users"
                                    :key="user.userId"
                                    class="flex items-center justify-between py-3 pr-4 pl-3"
                                >
                                    <div class="flex flex-1 items-center">
                                        <CitizenInfoPopover
                                            class="flex items-center justify-center"
                                            :user="user.user"
                                            show-avatar-in-name
                                        />
                                    </div>
                                </li>
                            </ul>
                        </div>

                        <span class="isolate mt-2 inline-flex rounded-md shadow-xs">
                            <UButton
                                v-if="canDo('TakeControl')"
                                icon="i-mdi-pencil"
                                truncate
                                @click="
                                    unitAssignUsersModal.open({
                                        unit: unit,
                                    })
                                "
                            >
                                {{ $t('common.assign') }}
                            </UButton>
                        </span>
                    </dd>
                </div>

                <div>
                    <UnitFeed :unit-id="unit.id" />
                </div>
            </dl>
        </template>

        <template #footer>
            <UButton class="flex-1" color="neutral" block @click="$emit('close', false)">
                {{ $t('common.close', 1) }}
            </UButton>
        </template>
    </USlideover>
</template>
