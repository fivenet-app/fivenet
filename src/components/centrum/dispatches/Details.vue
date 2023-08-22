<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { AccountIcon, CloseIcon, PencilIcon, PlusIcon } from 'mdi-vue3';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import Time from '~/components/partials/elements/Time.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { TAKE_DISPATCH_RESP } from '~~/gen/ts/services/centrum/centrum';
import AssignDispatchModal from './AssignDispatchModal.vue';
import Feed from './Feed.vue';
import StatusUpdateModal from './StatusUpdateModal.vue';

defineProps<{
    open: boolean;
    dispatch: Dispatch;
}>();

defineEmits<{
    (e: 'close'): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

async function selfAssign(dispatch: Dispatch): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().takeDispatch({
                dispatchIds: [dispatch.id],
                resp: TAKE_DISPATCH_RESP.ACCEPTED,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const openAssign = ref(false);
const openStatus = ref(false);
const openSelfAssign = ref(false);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
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
                                                <DialogTitle class="inline-flex text-base font-semibold leading-6 text-white">
                                                    {{ $t('common.dispatch') }}:
                                                    <IDCopyBadge class="ml-2 mr-2" :id="dispatch.id" prefix="DSP" />
                                                    {{ dispatch.message }}
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-white"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">
                                                    {{ $t('common.description') }}: {{ dispatch.description ?? 'N/A' }}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="border-b border-white/10 divide-y divide-white/10">
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.last_update') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <Time :value="dispatch.status?.createdAt" />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.status') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <StatusUpdateModal
                                                                    v-if="openStatus"
                                                                    :open="openStatus"
                                                                    :dispatch="dispatch"
                                                                    @close="openStatus = false"
                                                                />
                                                                <button
                                                                    type="button"
                                                                    @click="openStatus = true"
                                                                    class="rounded bg-white/10 px-2 py-1 text-xs font-semibold text-white shadow-sm hover:bg-white/20"
                                                                >
                                                                    {{
                                                                        $t(
                                                                            `enums.centrum.DISPATCH_STATUS.${
                                                                                DISPATCH_STATUS[
                                                                                    dispatch.status?.status ?? (0 as number)
                                                                                ]
                                                                            }`,
                                                                        )
                                                                    }}
                                                                    <span v-if="dispatch.status?.code">
                                                                        ({{ $t('common.code') }}: '{{ dispatch.status.code }}')
                                                                    </span>
                                                                </button>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.code') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ dispatch.status?.code ?? 'N/A' }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.reason') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ dispatch.status?.reason ?? 'N/A' }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.location') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <button
                                                                    v-if="dispatch.x && dispatch.y"
                                                                    type="button"
                                                                    class="text-primary-400 hover:text-primary-600"
                                                                    @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                                                                >
                                                                    {{ $t('common.go_to_location') }}
                                                                </button>
                                                                <span v-else>{{ $t('common.no_location') }}</span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.sent_by') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <NuxtLink
                                                                    :to="{
                                                                        name: 'citizens-id',
                                                                        params: { id: dispatch.status?.user?.userId ?? 0 },
                                                                    }"
                                                                    class="underline hover:text-neutral hover:transition-all"
                                                                >
                                                                    {{ dispatch.status?.user?.firstname }}
                                                                    {{ dispatch.status?.user?.lastname }}
                                                                </NuxtLink>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.units') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span v-if="dispatch.units.length === 0">
                                                                    {{ $t('common.unit', 0) }}
                                                                </span>
                                                                <ul
                                                                    v-else
                                                                    role="list"
                                                                    class="border divide-y rounded-md divide-base-200 border-base-200"
                                                                >
                                                                    <li
                                                                        v-for="unit in dispatch.units"
                                                                        class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                                                                    >
                                                                        <div class="flex items-center flex-1">
                                                                            <AccountIcon
                                                                                class="flex-shrink-0 w-5 h-5 text-base-400"
                                                                                aria-hidden="true"
                                                                            />
                                                                            <span class="flex-1 ml-2 truncate">
                                                                                {{ unit.unit?.name }}
                                                                                ({{ unit.unit?.initials }})
                                                                            </span>
                                                                        </div>
                                                                    </li>
                                                                </ul>

                                                                <AssignDispatchModal
                                                                    v-if="openAssign"
                                                                    :open="openAssign"
                                                                    :dispatch="dispatch"
                                                                    @close="openAssign = false"
                                                                />

                                                                <button
                                                                    v-if="can('CentrumService.TakeControl')"
                                                                    type="button"
                                                                    @click="openAssign = true"
                                                                    class="ml-2 rounded bg-white/10 px-2 py-1 text-xs font-semibold text-white shadow-sm hover:bg-white/20"
                                                                >
                                                                    <PencilIcon class="h-6 w-6" />
                                                                </button>
                                                                <button
                                                                    v-if="can('CentrumService.TakeDispatch')"
                                                                    type="button"
                                                                    @click="openSelfAssign = true"
                                                                    class="ml-2 rounded bg-white/10 px-2 py-1 text-xs font-semibold text-white shadow-sm hover:bg-white/20"
                                                                >
                                                                    <PlusIcon class="h-6 w-6" />
                                                                </button>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                {{ $t('common.attributes', 2) }}
                                                            </dt>
                                                            <dd class="mt-2 text-sm text-gray-400 sm:col-span-2 sm:mt-0">
                                                                <template v-if="dispatch.attributes?.list.length === 0">
                                                                    <span
                                                                        v-for="attribute in dispatch.attributes?.list"
                                                                        class="inline-flex items-center rounded-md bg-red-400/10 px-2 py-1 text-xs font-medium text-red-400 ring-1 ring-inset ring-red-400/20"
                                                                    >
                                                                        {{ attribute }}
                                                                    </span>
                                                                </template>
                                                                <span v-else>
                                                                    {{
                                                                        $t('common.none_selected', [$t('common.attributes', 2)])
                                                                    }}
                                                                </span>
                                                            </dd>
                                                        </div>
                                                    </dl>
                                                </div>

                                                <Feed :dispatch-id="dispatch.id" />
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <button
                                            type="button"
                                            class="w-full rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
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
