<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { AccountIcon, CloseIcon, MapMarkerIcon, PencilIcon } from 'mdi-vue3';
import { unitStatusToBGColor } from '~/components/centrum//helpers';
import UnitAssignUsersModal from '~/components/centrum/units/UnitAssignUsersModal.vue';
import UnitFeed from '~/components/centrum/units/UnitFeed.vue';
import UnitStatusUpdateModal from '~/components/centrum/units/UnitStatusUpdateModal.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { StatusUnit, Unit } from '~~/gen/ts/resources/centrum/units';

const props = defineProps<{
    open: boolean;
    unit: Unit;
    statusSelected?: StatusUnit;
}>();

defineEmits<{
    (e: 'close'): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const unitStatusColors = computed(() => unitStatusToBGColor(props.unit.status?.status));

const openAssign = ref(false);
const openStatus = ref(false);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-2xl pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-150 sm:duration-300"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-150 sm:duration-300"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-2xl">
                                <form class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl">
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-neutral">
                                                    {{ $t('common.unit') }}: {{ unit.initials }} -
                                                    {{ unit.name }}
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-neutral"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="divide-y divide-neutral/10 border-b border-neutral/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.description') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-2 max-h-24 text-sm text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <p class="whitespace-pre break-words">
                                                                    {{ unit.description ?? $t('common.na') }}
                                                                </p>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ `${$t('common.department')} ${$t('common.postal')}` }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ unit.homePostal ?? $t('common.na') }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.last_update') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <GenericTime :value="unit.status?.createdAt" />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.status') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <UnitStatusUpdateModal
                                                                    :open="openStatus"
                                                                    :unit="unit"
                                                                    :status="statusSelected"
                                                                    @close="openStatus = false"
                                                                />

                                                                <button
                                                                    type="button"
                                                                    class="rounded px-2 py-1 text-sm font-semibold text-neutral shadow-sm hover:bg-neutral/20"
                                                                    :class="unitStatusColors"
                                                                    @click="openStatus = true"
                                                                >
                                                                    {{
                                                                        $t(
                                                                            `enums.centrum.StatusUnit.${
                                                                                StatusUnit[unit.status?.status ?? 0]
                                                                            }`,
                                                                        )
                                                                    }}
                                                                </button>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.code') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ unit.status?.code ?? $t('common.na') }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.reason') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ unit.status?.reason ?? $t('common.na') }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.location') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span class="block">
                                                                    {{ $t('common.postal') }}:
                                                                    {{ unit.status?.postal ?? $t('common.na') }}
                                                                </span>
                                                                <button
                                                                    v-if="unit.status?.x && unit.status?.y"
                                                                    type="button"
                                                                    class="inline-flex items-center text-primary-400 hover:text-primary-600"
                                                                    @click="
                                                                        $emit('goto', { x: unit.status?.x, y: unit.status?.y })
                                                                    "
                                                                >
                                                                    <MapMarkerIcon class="mr-1 h-5 w-5" aria-hidden="true" />
                                                                    {{ $t('common.go_to_location') }}
                                                                </button>
                                                                <span v-else>{{ $t('common.no_location') }}</span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.attributes', 2) }}
                                                            </dt>
                                                            <dd class="mt-2 text-sm text-gray-300 sm:col-span-2 sm:mt-0">
                                                                <template
                                                                    v-if="
                                                                        unit.attributes !== undefined &&
                                                                        unit.attributes?.list.length > 0
                                                                    "
                                                                >
                                                                    <span
                                                                        v-for="attribute in unit.attributes?.list"
                                                                        :key="attribute"
                                                                        class="inline-flex items-center rounded-md bg-warn-400/10 px-2 py-1 text-xs font-medium text-warn-400 ring-1 ring-inset ring-warn-400/20"
                                                                    >
                                                                        {{
                                                                            $t(
                                                                                `components.centrum.units.attributes.${attribute}`,
                                                                            )
                                                                        }}
                                                                    </span>
                                                                </template>
                                                                <span v-else>
                                                                    {{ $t('common.none', [$t('common.attributes', 2)]) }}
                                                                </span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.members') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span v-if="unit.users.length === 0" class="block">
                                                                    {{ $t('common.member', unit.users.length) }}
                                                                </span>
                                                                <div v-else class="rounded-md bg-base-800">
                                                                    <ul
                                                                        role="list"
                                                                        class="divide-y divide-gray-200 text-sm font-medium"
                                                                    >
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
                                                                                        <AccountIcon
                                                                                            class="mr-1 h-5 w-5 flex-shrink-0 text-base-400"
                                                                                            aria-hidden="true"
                                                                                        />
                                                                                    </template>
                                                                                </CitizenInfoPopover>
                                                                            </div>
                                                                        </li>
                                                                    </ul>
                                                                </div>

                                                                <UnitAssignUsersModal
                                                                    :open="openAssign"
                                                                    :unit="unit"
                                                                    @close="openAssign = false"
                                                                />

                                                                <span class="isolate mt-2 inline-flex rounded-md shadow-sm">
                                                                    <button
                                                                        v-if="can('CentrumService.TakeControl')"
                                                                        type="button"
                                                                        class="flex flex-row items-center rounded bg-neutral/10 px-2 py-1 text-xs font-semibold text-neutral shadow-sm hover:bg-neutral/20"
                                                                        @click="openAssign = true"
                                                                    >
                                                                        <PencilIcon class="h-5 w-5" />
                                                                        <span class="ml-0.5 truncate">
                                                                            {{ $t('common.assign') }}
                                                                        </span>
                                                                    </button>
                                                                </span>
                                                            </dd>
                                                        </div>
                                                    </dl>
                                                </div>

                                                <UnitFeed :unit-id="unit.id" @goto="$emit('goto', $event)" />
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <button
                                            type="button"
                                            class="w-full rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                                            @click="$emit('close')"
                                        >
                                            {{ $t('common.close') }}
                                        </button>
                                    </div>
                                </form>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
