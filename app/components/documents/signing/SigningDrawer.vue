<script lang="ts" setup>
import SignaturePad from '~/components/partials/SignaturePad.vue';

defineProps<{
    documentId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

// TODO

const signature = ref('');
</script>

<template>
    <UDrawer
        :title="$t('common.sign')"
        :overlay="false"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ content: 'min-h-[50%]', title: 'flex flex-row gap-2' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('common.sign') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                Signing Drawer

                <SignaturePad v-model="signature" />
            </div>
        </template>
    </UDrawer>
</template>
