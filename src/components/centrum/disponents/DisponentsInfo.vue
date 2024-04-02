<script lang="ts" setup>
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
    <div class="flex w-full items-center justify-items-center gap-2">
        <DisponentsModal :open="open" @close="open = false" />

        <p class="text-sm">
            <UButton
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
            </UButton>
        </p>

        <UBadge color="gray">
            {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
        </UBadge>

        <UButton
            v-if="!isDisponent"
            class="inline-flex items-center justify-center rounded-full"
            icon="i-mdi-location-enter"
            @click="onSubmitThrottle(true)"
        >
            <template v-if="!canSubmit">
                <LoadingIcon class="mr-2 size-5 animate-spin" />
            </template>
            <span class="px-1">{{ $t('common.join') }}</span>
        </UButton>
        <UButton
            v-else
            class="inline-flex items-center justify-center rounded-full"
            icon="i-mdi-location-exit"
            @click="onSubmitThrottle(false)"
        >
            <template v-if="!canSubmit">
                <LoadingIcon class="mr-2 size-5 animate-spin" />
            </template>
            <span class="px-1">{{ $t('common.leave') }}</span>
        </UButton>
    </div>
</template>
