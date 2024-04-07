<script lang="ts" setup>
import { unitStatusToBGColor } from '~/components/centrum//helpers';
import UnitAssignUsersSlideover from '~/components/centrum/units/UnitAssignUsersSlideover.vue';
import UnitFeed from '~/components/centrum/units/UnitFeed.vue';
import UnitStatusUpdateModal from '~/components/centrum/units/UnitStatusUpdateModal.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { StatusUnit, Unit } from '~~/gen/ts/resources/centrum/units';
import UnitAttributes from '../partials/UnitAttributes.vue';

const props = defineProps<{
    unit: Unit;
    statusSelected?: StatusUnit;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const slideover = useSlideover();
const { isOpen } = useSlideover();

const modal = useModal();

const unitStatusColors = computed(() => unitStatusToBGColor(props.unit.status?.status));
</script>

<template>
    <USlideover>
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 max-h-[calc(100vh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.units') }}: {{ unit.initials }} - {{ unit.name }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <div class="flex flex-1 flex-col justify-between">
                    <div class="divide-y divide-gray-200 px-2 sm:px-6">
                        <div class="mt-1">
                            <dl class="divide-neutral/10 border-neutral/10 divide-y border-b">
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.description') }}
                                    </dt>
                                    <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                                        <p class="max-h-14 overflow-y-scroll break-words">
                                            {{ unit.description ?? $t('common.na') }}
                                        </p>
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ `${$t('common.department')} ${$t('common.postal')}` }}
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        {{ unit.homePostal ?? $t('common.na') }}
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.last_update') }}
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        <GenericTime :value="unit.status?.createdAt" />
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.status') }}
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        <UButton
                                            class="hover:bg-neutral/20 rounded px-2 py-1 text-sm font-semibold shadow-sm"
                                            :class="unitStatusColors"
                                            @click="
                                                modal.open(UnitStatusUpdateModal, {
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
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.code') }}
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        {{ unit.status?.code ?? $t('common.na') }}
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.reason') }}
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        {{ unit.status?.reason ?? $t('common.na') }}
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.location') }}
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        <span class="block">
                                            {{ $t('common.postal') }}:
                                            {{ unit.status?.postal ?? $t('common.na') }}
                                        </span>
                                        <UButton
                                            v-if="unit.status?.x && unit.status?.y"
                                            size="xs"
                                            variant="link"
                                            icon="i-mdi-map-marker"
                                            @click="$emit('goto', { x: unit.status?.x, y: unit.status?.y })"
                                        >
                                            {{ $t('common.go_to_location') }}
                                        </UButton>
                                        <span v-else>{{ $t('common.no_location') }}</span>
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.attributes', 2) }}
                                    </dt>
                                    <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                                        <UnitAttributes :attributes="unit.attributes" />
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        {{ $t('common.members') }}
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        <span v-if="unit.users.length === 0" class="block">
                                            {{ $t('common.member', unit.users.length) }}
                                        </span>
                                        <div v-else class="rounded-md bg-base-800">
                                            <ul role="list" class="divide-y divide-gray-200 text-sm font-medium">
                                                <li
                                                    v-for="user in unit.users"
                                                    :key="user.userId"
                                                    class="flex items-center justify-between py-3 pl-3 pr-4"
                                                >
                                                    <div class="flex flex-1 items-center">
                                                        <CitizenInfoPopover
                                                            :user="user.user"
                                                            class="flex items-center justify-center"
                                                            text-class="text-gray-300"
                                                        >
                                                            <template #before>
                                                                <UIcon
                                                                    name="i-mdi-account"
                                                                    class="mr-1 size-5 shrink-0 text-base-300"
                                                                />
                                                            </template>
                                                        </CitizenInfoPopover>
                                                    </div>
                                                </li>
                                            </ul>
                                        </div>

                                        <span class="isolate mt-2 inline-flex rounded-md shadow-sm">
                                            <UButton
                                                v-if="can('CentrumService.TakeControl')"
                                                icon="i-mdi-pencil"
                                                truncate
                                                @click="
                                                    slideover.open(UnitAssignUsersSlideover, {
                                                        unit: unit,
                                                    })
                                                "
                                            >
                                                {{ $t('common.assign') }}
                                            </UButton>
                                        </span>
                                    </dd>
                                </div>
                            </dl>
                        </div>

                        <UnitFeed :unit-id="unit.id" @goto="$emit('goto', $event)" />
                    </div>
                </div>
            </div>

            <template #footer>
                <UButton color="black" block class="flex-1" @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </USlideover>
</template>
