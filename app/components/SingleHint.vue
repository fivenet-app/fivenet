<script lang="ts" setup generic="T extends RoutesNamesList, P extends string, E extends boolean = false">
import type { NuxtRoute, RoutesNamesList } from '@typed-router';

const props = defineProps<{
    hintId: string;
    showKey?: boolean;
    to?: NuxtRoute<T, P, E>;
    external?: E;
    linkTarget?: '_blank';
}>();
</script>

<template>
    <UAlert
        v-bind="$attrs"
        :ui="{
            body: { padding: 'px-2 py-3 sm:p-3' },
            header: { padding: 'px-2 py-3 sm:p-3' },
            footer: { padding: 'px-2 py-2 sm:p-3' },
        }"
    >
        <template #title>
            <div class="inline-flex items-center gap-1">
                <UIcon name="i-mdi-information-slab-circle" class="size-6" />
                <span class="shrink-0 font-semibold">{{ $t('components.hints.start_text') }}</span>
            </div>
        </template>

        <template #description>
            <div class="mx-auto mb-2 flex items-center gap-1 text-base">
                <span class="grow">{{ $t(`components.hints.${hintId}.content`) }} </span>

                <div v-if="showKey || to" class="flex-initial">
                    <UKbd v-if="showKey" :value="$t(`components.hints.${hintId}.keyboard`)" />

                    <UButton v-else-if="to" variant="soft" :to="to" :external="props.external" :target="linkTarget">
                        {{ $t('components.hints.click_me') }}
                    </UButton>
                </div>
            </div>
        </template>
    </UAlert>
</template>
