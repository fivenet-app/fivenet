<script lang="ts" setup>
import { LinkVariantIcon } from 'mdi-vue3';
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
                        $emit('selectedDispatch', reference.targetDispatchId);
                    "
                >
                    <span class="sr-only">{{ $t('common.open') }}</span>
                    <LinkVariantIcon class="ml-1 size-5" aria-hidden="true" />
                </UButton>
            </span>
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
