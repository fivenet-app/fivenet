<script lang="ts" setup>
import DisponentsModal from '~/components/centrum/disponents/DisponentsModal.vue';
import { useCentrumStore } from '~/stores/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

const props = withDefaults(
    defineProps<{
        hideJoin?: boolean;
    }>(),
    {
        hideJoin: false,
    },
);

const { $grpc } = useNuxtApp();

const modal = useModal();

const centrumStore = useCentrumStore();
const { getCurrentMode, disponents, isDisponent } = storeToRefs(centrumStore);

async function takeControl(signon: boolean): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.takeControl({
            signon,
        });
        await call;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (e: boolean) => {
    canSubmit.value = false;
    await takeControl(e).finally(() => useTimeoutFn(() => (canSubmit.value = true), 850));
}, 1000);

if (!props.hideJoin) {
    defineShortcuts({
        'c-q': () => onSubmitThrottle(!isDisponent.value),
    });
}
</script>

<template>
    <div class="flex w-full items-center justify-items-center gap-2">
        <template v-if="!hideJoin">
            <UTooltip :text="`${$t('common.join', 1)}/ ${$t('common.leave', 1)}`" :shortcuts="['C', 'Q']">
                <UButton
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :icon="!isDisponent ? 'i-mdi-location-enter' : 'i-mdi-location-exit'"
                    :color="!isDisponent ? 'primary' : 'amber'"
                    :label="!isDisponent ? $t('common.join', 1) : $t('common.leave', 1)"
                    @click="onSubmitThrottle(!isDisponent)"
                />
            </UTooltip>
        </template>

        <UTooltip :text="usersToLabel(disponents)">
            <UButton
                :icon="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN ? 'i-mdi-monitor' : 'i-mdi-robot'"
                :color="
                    getCurrentMode === CentrumMode.AUTO_ROUND_ROBIN ? 'gray' : disponents.length === 0 ? 'amber' : 'success'
                "
                truncate
                @click="modal.open(DisponentsModal, {})"
            >
                <template v-if="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN">
                    {{ $t('common.disponent', disponents.length) }}
                </template>
                <template v-else>
                    {{ $t('enums.centrum.CentrumMode.AUTO_ROUND_ROBIN') }}
                </template>
            </UButton>
        </UTooltip>
    </div>
</template>
