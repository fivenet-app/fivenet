<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';
import DisponentsModal from '~/components/centrum/disponents/DisponentsModal.vue';

withDefaults(
    defineProps<{
        hideJoin?: boolean;
    }>(),
    {
        hideJoin: false,
    },
);

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

const modal = useModal();
</script>

<template>
    <div class="flex w-full items-center justify-items-center gap-2">
        <UButton
            :icon="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN ? 'i-mdi-monitor' : 'i-mdi-robot'"
            :color="getCurrentMode === CentrumMode.AUTO_ROUND_ROBIN ? 'primary' : disponents.length === 0 ? 'amber' : 'green'"
            truncate
            :title="disponentsNames"
            @click="modal.open(DisponentsModal, {})"
        >
            <template v-if="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN">
                {{ $t('common.disponent', disponents.length) }}
            </template>
            <template v-else>
                {{ $t('enums.centrum.CentrumMode.AUTO_ROUND_ROBIN') }}
            </template>
        </UButton>

        <template v-if="!hideJoin">
            <UButton
                v-if="!isDisponent"
                class="inline-flex items-center justify-center rounded-full"
                :disabled="!canSubmit"
                :loading="!canSubmit"
                icon="i-mdi-location-enter"
                @click="onSubmitThrottle(true)"
            >
                <span class="px-1">{{ $t('common.join') }}</span>
            </UButton>
            <UButton
                v-else
                class="inline-flex items-center justify-center rounded-full"
                :disabled="!canSubmit"
                :loading="!canSubmit"
                icon="i-mdi-location-exit"
                @click="onSubmitThrottle(false)"
            >
                <span class="px-1">{{ $t('common.leave') }}</span>
            </UButton>
        </template>
    </div>
</template>
