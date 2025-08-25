<script lang="ts" setup>
import DispatchDetailsByIDSlideover from '~/components/centrum/dispatches/DispatchDetailsByIDSlideover.vue';
import { type DispatchReferences, DispatchReferenceType } from '~~/gen/ts/resources/centrum/dispatches';

defineProps<{
    references?: DispatchReferences;
}>();

const overlay = useOverlay();
const dispatchDetailsByIDSlideover = overlay.create(DispatchDetailsByIDSlideover);

const selectedDispatch = ref<number | undefined>();
</script>

<template>
    <template v-if="references !== undefined && references?.references.length > 0">
        <div class="grid grid-cols-2 gap-1">
            <span
                v-for="reference in references?.references"
                :key="reference.targetDispatchId"
                class="flex items-center rounded-md bg-info-400/10 px-2 py-1 text-xs font-medium text-info-400 ring-1 ring-info-400/20 ring-inset"
            >
                <span class="flex-1">
                    {{ $t(`enums.centrum.DispatchReferenceType.${DispatchReferenceType[reference.referenceType]}`) }} DSP-{{
                        reference.targetDispatchId
                    }}
                </span>

                <UButton
                    icon="i-mdi-link-variant"
                    @click="
                        selectedDispatch = reference.targetDispatchId;
                        dispatchDetailsByIDSlideover.open({
                            dispatchId: reference.targetDispatchId,
                            onClose: () => (selectedDispatch = undefined),
                        });
                    "
                />
            </span>
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
