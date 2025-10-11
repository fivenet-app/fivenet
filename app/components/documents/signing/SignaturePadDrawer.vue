<script lang="ts" setup>
import SignaturePad from '~/components/partials/SignaturePad.vue';

const props = defineProps<{
    policyId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const signature = defineModel<string>({ required: true });

const signatureRef = useTemplateRef('signatureRef');

function handleSave(): void {
    const sig = signatureRef.value?.signature?.saveSignature('image/svg+xml') ?? '';
    // atob? Yes, because supporting FiveM's NUI CEF version 103 is fun..
    signature.value = atob(sig.replace(/^data:image\/svg\+xml;base64,/, ''));
}

// TODO
</script>

<template>
    <UDrawer
        :title="$t('common.approve')"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ container: 'flex-1', title: 'flex flex-row gap-2', body: 'h-full' }"
    >
        <slot />

        <template #title>
            <span class="flex-1">{{ $t('common.sign') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <SignaturePad ref="signatureRef" />
            </div>
        </template>

        <template #footer>
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col gap-4">
                <UButtonGroup class="w-full flex-1">
                    <UButton class="flex-1" color="neutral" block :label="$t('common.cancel')" @click="$emit('close', false)" />

                    <UButton class="flex-1" color="success" block :label="$t('common.save')" @click="() => handleSave()" />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
