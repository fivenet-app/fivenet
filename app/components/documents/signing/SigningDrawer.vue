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
        :ui="{ container: 'flex-1', content: 'min-h-[50%]', title: 'flex flex-row gap-2', body: 'h-full' }"
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

        <template #footer>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <UButtonGroup class="w-full flex-1">
                    <UButton
                        class="flex-1"
                        color="neutral"
                        block
                        :label="$t('common.close', 1)"
                        @click="$emit('close', false)"
                    />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
