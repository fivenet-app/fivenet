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
const { getCurrentMode, getOwnDisponents, disponents, isDisponent } = storeToRefs(centrumStore);

async function takeControl(signon: boolean): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.takeControl({
            signon: signon,
        });
        await call;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const ownDisponents = computed(() => getOwnDisponents.value);

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
        <UButtonGroup>
            <UTooltip :text="usersToLabel(ownDisponents)">
                <UButton
                    :icon="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN ? 'i-mdi-monitor' : 'i-mdi-robot'"
                    :color="
                        getCurrentMode === CentrumMode.AUTO_ROUND_ROBIN
                            ? 'gray'
                            : ownDisponents.length === 0
                              ? 'amber'
                              : 'success'
                    "
                    truncate
                    @click="modal.open(DisponentsModal, {})"
                >
                    <template v-if="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN">
                        {{ $t('common.disponent', ownDisponents.length) }}
                    </template>
                    <template v-else>
                        {{ $t('enums.centrum.CentrumMode.AUTO_ROUND_ROBIN') }}
                    </template>
                </UButton>
            </UTooltip>

            <UButton v-if="disponents.length > 1" icon="i-mdi-multicast" @click="modal.open(DisponentsModal, {})" />
        </UButtonGroup>

        <template v-if="!hideJoin">
            <UTooltip :text="`${$t('common.join')}/ ${$t('common.leave')}`" :shortcuts="['C', 'Q']">
                <UButton
                    v-if="!isDisponent"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    icon="i-mdi-location-enter"
                    @click="onSubmitThrottle(true)"
                >
                    {{ $t('common.join') }}
                </UButton>
                <UButton
                    v-else
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    color="amber"
                    icon="i-mdi-location-exit"
                    @click="onSubmitThrottle(false)"
                >
                    {{ $t('common.leave') }}
                </UButton>
            </UTooltip>
        </template>
    </div>
</template>
