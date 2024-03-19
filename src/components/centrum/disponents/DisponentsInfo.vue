<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { LoadingIcon, LocationEnterIcon, LocationExitIcon } from 'mdi-vue3';
import { useCentrumStore } from '~/store/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';
import DisponentsModal from '~/components/centrum/disponents/DisponentsModal.vue';

const { $grpc } = useNuxtApp();

const centrumStore = useCentrumStore();
const { getCurrentMode, disponents, isDisponent } = storeToRefs(centrumStore);

async function takeControl(signon: boolean): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().takeControl({
            signon,
        });
        await call;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (e: boolean) => {
    canSubmit.value = false;
    await takeControl(e).finally(() => useTimeoutFn(() => (canSubmit.value = true), 850));
}, 1000);

const disponentsNames = computed(() => disponents.value.map((u) => `${u.firstname} ${u.lastname}`));

const open = ref(false);
</script>

<template>
    <div class="h-full overflow-y-auto px-4 sm:px-6 lg:px-8">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h2 class="text-base font-semibold leading-6 text-gray-100">{{ $t('common.disponents', 2) }}</h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <div class="grid grid-cols-1 items-center justify-items-center sm:grid-cols-3">
                        <div class="inline-flex flex-1">
                            <button
                                v-if="!isDisponent"
                                type="button"
                                class="flex items-center justify-center rounded-full bg-success-500 text-neutral hover:bg-success-400"
                                @click="onSubmitThrottle(true)"
                            >
                                <LocationEnterIcon v-if="canSubmit" class="h-7 w-7" aria-hidden="true" />
                                <template v-else>
                                    <LoadingIcon class="mr-2 h-7 w-7 animate-spin" aria-hidden="true" />
                                </template>
                                <span class="px-1">{{ $t('common.join') }}</span>
                            </button>
                            <button
                                v-else
                                type="button"
                                class="flex items-center justify-center rounded-full bg-primary-500 text-neutral hover:bg-primary-400"
                                @click="onSubmitThrottle(false)"
                            >
                                <LocationExitIcon v-if="canSubmit" class="h-7 w-7" aria-hidden="true" />
                                <template v-else>
                                    <LoadingIcon class="mr-2 h-7 w-7 animate-spin" aria-hidden="true" />
                                </template>
                                <span class="px-1">{{ $t('common.leave') }}</span>
                            </button>
                        </div>
                        <div class="flex-1">
                            <DisponentsModal :open="open" @close="open = false" />

                            <p class="text-sm text-neutral">
                                <button
                                    class="inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset"
                                    :class="
                                        disponents.length === 0
                                            ? 'bg-warn-400/10 text-warn-500 ring-warn-400/20'
                                            : 'bg-success-500/10 text-success-400 ring-success-500/20'
                                    "
                                    :title="disponentsNames.join(', ')"
                                    @click="open = true"
                                >
                                    {{ $t('common.disponent', disponents.length) }}
                                </button>
                            </p>
                        </div>
                        <div class="flex-1">
                            <p class="text-sm text-neutral">
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
