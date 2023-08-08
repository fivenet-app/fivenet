<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccount, mdiClose } from '@mdi/js';
import Time from '~/components/partials/elements/Time.vue';
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';
import UnitFeed from './Feed.vue';
import StatusUpdateModal from './StatusUpdateModal.vue';

defineProps<{
    open: boolean;
    unit: Unit;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

const statusOpen = ref(false);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-500 sm:duration-700"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-500 sm:duration-700"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-3xl">
                                <form class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl">
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-white">
                                                    {{ $t('common.unit') }}: {{ unit.initials }} -
                                                    {{ unit.name }}
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-white"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">Close panel</span>
                                                        <SvgIcon
                                                            type="mdi"
                                                            :path="mdiClose"
                                                            class="h-6 w-6"
                                                            aria-hidden="true"
                                                        />
                                                    </button>
                                                </div>
                                            </div>
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">
                                                    Description: {{ unit.description ?? 'N/A' }}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="border-b border-white/10 divide-y divide-white/10">
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">Status</dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <StatusUpdateModal
                                                                    :open="statusOpen"
                                                                    :unit="unit"
                                                                    @close="statusOpen = false"
                                                                />
                                                                <button
                                                                    type="button"
                                                                    @click="statusOpen = true"
                                                                    class="rounded bg-white/10 px-2 py-1 text-xs font-semibold text-white shadow-sm hover:bg-white/20"
                                                                >
                                                                    {{ UNIT_STATUS[unit.status?.status ?? 0] }}
                                                                </button>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                Last Unit Update
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <Time :value="unit.status?.createdAt" />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                Status Reason
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ unit.status?.reason ?? 'N/A' }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                Status Code
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ unit.status?.code ?? 'N/A' }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">Location</dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <button
                                                                    v-if="unit.status?.x && unit.status?.y"
                                                                    type="button"
                                                                    class="text-primary-400 hover:text-primary-600"
                                                                >
                                                                    Go to Location
                                                                </button>
                                                                <span v-else>No Location</span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">Members</dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span v-if="unit.users.length === 0">No members </span>
                                                                <ul
                                                                    v-else
                                                                    role="list"
                                                                    class="border divide-y rounded-md divide-base-200 border-base-200"
                                                                >
                                                                    <li
                                                                        v-for="user in unit.users"
                                                                        class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                                                                    >
                                                                        <div class="flex items-center flex-1">
                                                                            <SvgIcon
                                                                                class="flex-shrink-0 w-5 h-5 text-base-400"
                                                                                aria-hidden="true"
                                                                                type="mdi"
                                                                                :path="mdiAccount"
                                                                            />
                                                                            <span class="flex-1 ml-2 truncate">
                                                                                {{ user.user?.firstname }}
                                                                                {{ user.user?.lastname }}
                                                                                ({{ user.user?.dateofbirth }})
                                                                            </span>
                                                                        </div>
                                                                    </li>
                                                                </ul>
                                                            </dd>
                                                        </div>
                                                    </dl>
                                                </div>

                                                <UnitFeed :unit-id="unit.id" />
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <button
                                            type="button"
                                            class="rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                                            @click="$emit('close')"
                                        >
                                            Close
                                        </button>
                                        <button
                                            type="submit"
                                            class="ml-4 inline-flex justify-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
                                        >
                                            Save
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
