<script lang="ts" setup generic="T extends RoutesNamesList, P extends string, E extends boolean = false">
import type { NuxtRoute, RoutesNamesList } from '@typed-router';

defineProps<{
    hintId: string;
    keyboard?: boolean;
    to?: NuxtRoute<T, P, E>;
    external?: E;
    linkTarget?: '_blank' | null;
}>();
</script>

<template>
    <UCard
        :ui="{
            body: { padding: 'px-2 py-3 sm:p-3' },
            header: { padding: 'px-2 py-3 sm:p-3' },
            footer: { padding: 'px-2 py-2 sm:p-3' },
        }"
    >
        <template #header>
            <div class="inline-flex items-center">
                <UIcon name="i-mdi-information-slab-circle" class="size-6" />
                <strong class="ml-1 shrink-0 font-semibold">{{ $t('components.hints.start_text') }}</strong>
            </div>
        </template>

        <div class="mx-auto mb-2 flex items-center gap-1 text-base">
            <span class="grow">{{ $t(`components.hints.${hintId}.content`) }} </span>

            <div v-if="keyboard || to" class="flex-initial">
                <UKbd v-if="keyboard" :value="$t(`components.hints.${hintId}.keyboard`)" />

                <UButton v-else-if="to" variant="soft" :to="to" :external="external" :target="linkTarget ?? null">
                    {{ $t('components.hints.click_me') }}
                </UButton>
            </div>
        </div>
    </UCard>
</template>
