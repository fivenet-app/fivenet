<script lang="ts" setup>
import MarkerList from './MarkerList.vue';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const isDismissible = ref<boolean>(true);
</script>

<template>
    <UDrawer
        :title="$t('common.marker', 2)"
        handle-only
        :modal="false"
        :dismissible="isDismissible"
        :overlay="false"
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{
            container: 'flex-1 overflow-y-hidden',
            content: 'flex h-[35dvh] min-h-[30dvh] flex-col',
            title: 'flex flex-row gap-2',
            body: 'flex min-h-0 flex-1 overflow-y-hidden',
        }"
    >
        <template #title>
            <span class="flex-1">{{ $t('common.marker', 2) }}</span>

            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto flex h-full min-h-0 w-full max-w-(--breakpoint-xl) flex-col overflow-hidden">
                <MarkerList class="h-full px-2" @editing="isDismissible = !isDismissible" />
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </UDrawer>
</template>
