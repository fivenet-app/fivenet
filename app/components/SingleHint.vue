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
    <UPageCard
        icon="i-mdi-information-outline"
        :title="
            $te(`components.hints.${hintId}.title`) ? $t(`components.hints.${hintId}.title`) : $t('components.hints.start_text')
        "
        :ui="{
            wrapper: 'flex-row gap-2',
            leadingIcon: 'size-6',
            container: 'p-2 sm:p-2',
        }"
        v-bind="$attrs"
    >
        <template #description>
            <div class="mx-auto mb-2 flex items-center gap-2 text-sm">
                <span class="grow">{{ $t(`components.hints.${hintId}.content`) }} </span>

                <div v-if="showKey || to" class="flex-initial shrink-0">
                    <UKbd v-if="showKey" size="md" :value="$t(`components.hints.${hintId}.keyboard`)" />
                    <UButton
                        v-else-if="to"
                        size="sm"
                        :to="to.toString()"
                        :external="props.external"
                        :target="linkTarget"
                        :label="$t('components.hints.click_me')"
                    />
                </div>
            </div>
        </template>
    </UPageCard>
</template>
