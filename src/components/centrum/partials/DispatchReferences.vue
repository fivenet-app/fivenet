<script lang="ts" setup>
import { type DispatchReferences, DispatchReferenceType } from '~~/gen/ts/resources/centrum/dispatches';
import DispatchDetailsByID from '~/components/centrum/dispatches/DispatchDetailsByID.vue';

defineProps<{
    references?: DispatchReferences;
}>();

defineEmits<{
    (e: 'selectedDispatch', dspId: string): void;
}>();

const selectedDispatch = ref<string | undefined>();
</script>

<template>
    <template v-if="references !== undefined && references?.references.length > 0">
        <DispatchDetailsByID
            v-if="selectedDispatch"
            :open="true"
            :dispatch-id="selectedDispatch"
            @close="selectedDispatch = undefined"
        />

        <div class="inline-flex gap-1 line-clamp-2">
            <span
                v-for="reference in references?.references"
                :key="reference.targetDispatchId"
                class="inline-flex items-center rounded-md bg-info-400/10 px-2 py-1 text-xs font-medium text-info-400 ring-1 ring-inset ring-info-400/20"
                @click="
                    selectedDispatch = reference.targetDispatchId;
                    $emit('selectedDispatch', reference.targetDispatchId);
                "
            >
                {{ $t(`enums.centrum.DispatchReferenceType.${DispatchReferenceType[reference.referenceType]}`) }} - DSP-{{
                    reference.targetDispatchId
                }}
            </span>
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
