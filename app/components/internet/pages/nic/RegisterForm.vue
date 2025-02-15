<script lang="ts" setup>
import { z } from 'zod';
import type { CheckDomainAvailabilityResponse } from '~~/gen/ts/services/internet/domain';

defineProps<{
    status?: CheckDomainAvailabilityResponse | undefined;
}>();

defineEmits<{
    (e: 'cancel'): void;
}>();

const schema = z.object({
    transferCode: z.string().length(6).optional(),
});

const state = reactive({
    transferCode: undefined,
});

// TODO
</script>

<template>
    <UForm :schema="schema" :state="state">
        <UFormGroup v-if="status?.transferable" :label="$t('components.internet.transfer_code')">
            <UInput v-model="state.transferCode" type="text" />
        </UFormGroup>

        <!-- Even without a transfer code "are you sure" confirmation -->
        <UFormGroup :label="$t('components.internet.pages.nic_registrar.register_form.submit')">
            <div class="flex gap-2">
                <UButton :label="$t('common.yes')" color="green" />

                <UButton :label="$t('common.cancel')" color="red" class="flex-1" @click="$emit('cancel')" />
            </div>
        </UFormGroup>
    </UForm>
</template>
