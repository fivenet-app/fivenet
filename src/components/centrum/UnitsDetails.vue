<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiClose } from '@mdi/js';
import { Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = defineProps<{
    open: boolean;
    unit: Unit;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const {
    data: activities,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`centrum-unit-${props.unit.id.toString()}-activity-${offset.value}`, () => listUnitActivity());

async function listUnitActivity(): Promise<Array<UnitStatus>> {
    return new Promise(async (res, rej) => {
        try {
            const req = {
                pagination: {
                    offset: offset.value,
                },
                id: props.unit.id,
            };

            const call = $grpc.getCentrumClient().listUnitActivity(req);
            const { response } = await call;

            pagination.value = response.pagination;
            return res(response.activity);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function addUserToUnit(): Promise<void> {
    return new Promise(async (res, rej) => {
        $grpc.getCentrumClient().assignUnit({
            unitId: props.unit.id,
            toAdd: [],
            toRemove: [],
        });
    });
}
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
                                <div class="flex h-full flex-col overflow-y-scroll bg-base-850 py-6 shadow-xl">
                                    <div class="px-4 sm:px-6">
                                        <div class="flex items-start justify-between">
                                            <DialogTitle class="text-base font-semibold leading-6 text-gray-100">
                                                {{ unit.id.toString() }} - {{ unit.name }} ({{ $t('common.initials') }}:
                                                {{ unit.initials }})
                                            </DialogTitle>
                                            <div class="ml-3 flex h-7 items-center">
                                                <button
                                                    type="button"
                                                    class="rounded-md bg-white text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                                    @click="$emit('close')"
                                                >
                                                    <span class="sr-only">Close panel</span>
                                                    <SvgIcon type="mdi" :path="mdiClose" class="h-6 w-6" aria-hidden="true" />
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="relative mt-6 flex-1 px-4 sm:px-6 text-gray-100">
                                        <ul>
                                            <li v-for="user in unit.users">
                                                {{ user.user?.firstname }} {{ user.user?.lastname }}
                                            </li>
                                        </ul>

                                        <button type="submit" @click="addUserToUnit">Add User To Unit</button>

                                        <hr />
                                        <ul>
                                            <li v-for="activity in activities">
                                                {{ activity.id }}
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
