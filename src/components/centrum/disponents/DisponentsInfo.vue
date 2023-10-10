<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useThrottleFn } from '@vueuse/core';
import { GroupIcon, LoadingIcon, LocationEnterIcon, LocationExitIcon, MonitorIcon } from 'mdi-vue3';
import { useCentrumStore } from '~/store/centrum';
import { CentrumMode } from '~~/gen/ts/resources/dispatch/settings';
import Modal from './Modal.vue';

const { $grpc } = useNuxtApp();

const centrumStore = useCentrumStore();
const { getCurrentMode, disponents, isDisponent } = storeToRefs(centrumStore);

async function takeControl(signon: boolean): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().takeControl({
                signon: signon,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (e: boolean) => {
    canSubmit.value = false;
    await takeControl(e).finally(() => setTimeout(() => (canSubmit.value = true), 850));
}, 1000);

const disponentsNames = computed(() => {
    const names: string[] = [];

    disponents.value.forEach((u) => {
        names.push(`${u.firstname} ${u.lastname}`);
    });

    return names.join(', ');
});

const open = ref(false);
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-auto">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h2 class="text-base font-semibold leading-6 text-gray-100">{{ $t('common.disponents', 2) }}</h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <div class="grid grid-cols-3 items-center justify-items-center">
                        <div
                            v-if="!isDisponent"
                            class="absolute z-20 inset-0 flex flex-col justify-center items-center bg-gray-600/70"
                        >
                            <button
                                @click="onSubmitThrottle(true)"
                                type="button"
                                class="relative block w-full p-12 text-center border-2 border-dotted rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                                :disabled="!canSubmit"
                            >
                                <LocationEnterIcon v-if="canSubmit" class="w-12 h-12 mx-auto text-neutral" />
                                <template v-else>
                                    <LoadingIcon class="animate-spin w-12 h-12 mx-auto text-neutral" />
                                </template>
                                <span class="block mt-2 text-sm font-semibold text-gray-300">
                                    {{ $t('components.centrum.dispatch_center.join_center') }}
                                </span>
                            </button>
                            <div class="flex flex-row gap-4">
                                <NuxtLink
                                    :to="{ name: 'centrum-units' }"
                                    class="mt-4 px-2 py-1 flex items-center justify-center rounded-full bg-primary-500 text-neutral hover:bg-primary-400"
                                >
                                    <GroupIcon class="w-8 h-8" />
                                    <span class="px-1">{{ $t('common.units') }}</span>
                                </NuxtLink>
                                <button
                                    type="button"
                                    class="mt-4 px-2 py-1 flex items-center justify-center rounded-full bg-primary-500 text-neutral hover:bg-primary-400"
                                    @click="open = true"
                                >
                                    <MonitorIcon class="w-8 h-8" />
                                    <span class="px-1">{{ $t('common.disponents', 2) }}</span>
                                </button>
                            </div>
                        </div>

                        <div class="flex-1 inline-flex">
                            <button
                                v-if="!isDisponent"
                                type="button"
                                @click="onSubmitThrottle(true)"
                                class="flex items-center justify-center rounded-full bg-success-500 text-neutral hover:bg-success-400"
                            >
                                <LocationEnterIcon v-if="canSubmit" class="w-7 h-7" />
                                <template v-else>
                                    <LoadingIcon class="animate-spin h-7 w-7 mr-2" />
                                </template>
                                <span class="px-1">{{ $t('common.join') }}</span>
                            </button>
                            <button
                                v-else
                                type="button"
                                @click="onSubmitThrottle(false)"
                                class="flex items-center justify-center rounded-full bg-primary-500 text-neutral hover:bg-primary-400"
                            >
                                <LocationExitIcon v-if="canSubmit" class="w-7 h-7" />
                                <template v-else>
                                    <LoadingIcon class="animate-spin h-7 w-7 mr-2" />
                                </template>
                                <span class="px-1">{{ $t('common.leave') }}</span>
                            </button>
                        </div>
                        <div class="flex-1">
                            <Modal :open="open" @close="open = false" />

                            <p class="text-neutral text-sm">
                                <button
                                    class="inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset"
                                    :class="
                                        disponents.length === 0
                                            ? 'bg-warn-400/10 text-warn-500 ring-warn-400/20'
                                            : 'bg-success-500/10 text-success-400 ring-success-500/20'
                                    "
                                    :title="disponentsNames"
                                    @click="open = true"
                                >
                                    {{ $t('common.disponent', disponents.length) }}
                                </button>
                            </p>
                        </div>
                        <div class="flex-1">
                            <p class="text-neutral text-sm">
                                <span
                                    class="inline-flex items-center rounded-md bg-gray-400/10 px-2 py-1 text-xs font-medium text-gray-400 ring-1 ring-inset ring-gray-400/20"
                                >
                                    {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
                                </span>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
