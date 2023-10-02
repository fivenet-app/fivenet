<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { GroupIcon, LocationEnterIcon, LocationExitIcon } from 'mdi-vue3';
import { useCentrumStore } from '~/store/centrum';
import { CentrumMode } from '~~/gen/ts/resources/dispatch/settings';

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

const disponentsNames = computed(() => {
    const names: string[] = [];

    disponents.value.forEach((u) => {
        names.push(`${u.firstname} ${u.lastname}`);
    });

    return names.join(', ');
});
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
                            class="absolute inset-0 flex flex-col justify-center items-center z-20 bg-gray-600/70"
                        >
                            <button
                                @click="takeControl(true)"
                                type="button"
                                class="relative block w-full p-12 text-center border-2 border-dotted rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                            >
                                <LocationEnterIcon class="w-12 h-12 mx-auto text-neutral" />
                                <span class="block mt-2 text-sm font-semibold text-gray-300">
                                    {{ $t('components.centrum.dispatch_center.join_center') }}
                                </span>
                            </button>
                            <div class="flex flex-row">
                                <NuxtLink
                                    :to="{ name: 'centrum-units' }"
                                    class="mt-4 px-2 py-1 flex items-center justify-center rounded-full bg-primary-500 text-neutral hover:bg-primary-400"
                                >
                                    <GroupIcon class="w-8 h-8" />
                                    <span class="px-1">{{ $t('common.units') }}</span>
                                </NuxtLink>
                            </div>
                        </div>

                        <div class="flex-1 inline-flex">
                            <button
                                v-if="!isDisponent"
                                type="button"
                                @click="takeControl(true)"
                                class="flex items-center justify-center rounded-full bg-success-500 text-neutral hover:bg-success-400"
                            >
                                <LocationEnterIcon class="w-8 h-8" />
                                <span class="px-1">{{ $t('common.join') }}</span>
                            </button>
                            <button
                                v-else
                                type="button"
                                @click="takeControl(false)"
                                class="flex items-center justify-center rounded-full bg-primary-500 text-neutral hover:bg-primary-400"
                            >
                                <LocationExitIcon class="w-8 h-8" />
                                <span class="px-1">{{ $t('common.leave') }}</span>
                            </button>
                        </div>
                        <div class="flex-1">
                            <p class="text-neutral text-sm">
                                <span
                                    class="inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset"
                                    :class="
                                        disponents.length === 0
                                            ? 'bg-warn-400/10 text-warn-500 ring-warn-400/20'
                                            : 'bg-success-500/10 text-success-400 ring-success-500/20'
                                    "
                                    :title="disponentsNames"
                                >
                                    {{ $t('common.disponent', disponents.length) }}
                                </span>
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
