<script lang="ts" setup>
import DispatcherModal from '~/components/centrum/dispatchers/DispatcherModal.vue';
import { useCentrumStore } from '~/stores/centrum';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

const props = withDefaults(
    defineProps<{
        hideJoin?: boolean;
    }>(),
    {
        hideJoin: false,
    },
);

const modal = useOverlay();

const centrumStore = useCentrumStore();
const { getCurrentMode, getJobDispatchers, isDispatcher } = storeToRefs(centrumStore);

const centrumCentrumClient = await getCentrumCentrumClient();

async function takeControl(signon: boolean): Promise<void> {
    try {
        const call = centrumCentrumClient.takeControl({
            signon,
        });
        await call;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const dispatchers = computed(() => getJobDispatchers.value ?? { dispatchers: [] });

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (e: boolean) => {
    canSubmit.value = false;
    await takeControl(e).finally(() => useTimeoutFn(() => (canSubmit.value = true), 850));
}, 1000);

if (!props.hideJoin) {
    defineShortcuts({
        'c-q': () => onSubmitThrottle(!isDispatcher.value),
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
                    :icon="!isDispatcher ? 'i-mdi-location-enter' : 'i-mdi-location-exit'"
                    :color="!isDispatcher ? 'primary' : 'warning'"
                    :label="!isDispatcher ? $t('common.join', 1) : $t('common.leave', 1)"
                    @click="onSubmitThrottle(!isDispatcher)"
                />
            </UTooltip>
        </template>

        <UTooltip :text="usersToLabel(dispatchers.dispatchers)">
            <UButton
                :icon="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN ? 'i-mdi-monitor' : 'i-mdi-robot'"
                :color="
                    getCurrentMode === CentrumMode.AUTO_ROUND_ROBIN
                        ? 'neutral'
                        : dispatchers.dispatchers.length === 0
                          ? 'warning'
                          : 'success'
                "
                truncate
                @click="modal.open(DispatcherModal, {})"
            >
                <template v-if="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN">
                    {{ $t('common.dispatcher', dispatchers.dispatchers.length) }}
                </template>
                <template v-else>
                    {{ $t('enums.centrum.CentrumMode.AUTO_ROUND_ROBIN') }}
                </template>
            </UButton>
        </UTooltip>
    </div>
</template>
