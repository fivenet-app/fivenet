<script lang="ts" setup>
import DispatchDetailsByIDSlideover from '~/components/centrum/dispatches/DispatchDetailsByIDSlideover.vue';
import type { DispatchReferences } from '~~/gen/ts/resources/centrum/dispatches';
import { DispatchReferenceType } from '~~/gen/ts/resources/centrum/dispatches';

defineProps<{
    references?: DispatchReferences;
}>();

const modal = useModal();

const selectedDispatch = ref<number | undefined>();
</script>

<template>
    <template v-if="references !== undefined && references?.references.length > 0">
        <div class="grid grid-cols-2 gap-1">
            <span
                v-for="reference in references?.references"
                :key="reference.targetDispatchId"
                class="flex items-center rounded-md bg-info-400/10 px-2 py-1 text-xs font-medium text-info-400 ring-1 ring-inset ring-info-400/20"
            >
                <span class="flex-1">
                    {{ $t(`enums.centrum.DispatchReferenceType.${DispatchReferenceType[reference.referenceType]}`) }} DSP-{{
                        reference.targetDispatchId
                    }}
                </span>

                <UButton
                    @click="
                        selectedDispatch = reference.targetDispatchId;
                        modal.open(DispatchDetailsByIDSlideover, {
                            dispatchId: reference.targetDispatchId,
                            onClose: () => (selectedDispatch = undefined),
                        });
                    "
                >
                    <span class="sr-only">{{ $t('common.open') }}</span>
                    <UIcon class="ml-1 size-5" name="i-mdi-link-variant" />
                </UButton>
            </span>
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
