<script lang="ts" setup generic="T extends RoutesNamesList, P extends string, E extends boolean = false">
import type { NuxtRoute, RoutesNamesList } from '@typed-router';

const props = withDefaults(
    defineProps<{
        hintId: string;
        showKey?: boolean;
        to?: NuxtRoute<T, P, E>;
        external?: E;
        linkTarget?: '_blank';
        hideTitle?: boolean;
    }>(),
    {
        showKey: false,
        to: undefined,
        external: undefined,
        linkTarget: undefined,
        hideTitle: false,
    },
);
</script>

<template>
    <UAlert
        :ui="{
            icon: 'size-6',
        }"
        icon="i-mdi-information-outline"
        v-bind="$attrs"
    >
        <template v-if="!hideTitle" #title>
            <div class="shrink-0 font-semibold">
                {{
                    $te(`components.hints.${hintId}.title`)
                        ? $t(`components.hints.${hintId}.title`)
                        : $t('components.hints.start_text')
                }}
            </div>
        </template>

        <template #description>
            <div class="mx-auto mb-2 flex items-center gap-2 text-base">
                <span class="grow">{{ $t(`components.hints.${hintId}.content`) }} </span>

                <div v-if="showKey || to" class="flex-initial">
                    <UKbd v-if="showKey" :value="$t(`components.hints.${hintId}.keyboard`)" />

                    <UButton v-else-if="to" :to="to.toString()" :external="props.external" :target="linkTarget">
                        {{ $t('components.hints.click_me') }}
                    </UButton>
                </div>
            </div>
        </template>
    </UAlert>
</template>
