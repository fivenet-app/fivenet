<script lang="ts" setup>
import type { RoutesNamedLocations } from '@typed-router';

const props = defineProps<{
    to?: RoutesNamedLocations | string;
    fallbackTo?: RoutesNamedLocations | string;
}>();

const router = useRouter();

async function goBack(): Promise<void> {
    if (props.to) return;

    if (history.length === 0) {
        if (props.fallbackTo) {
            // @ts-expect-error string can be valid for route paths
            await navigateTo(props.fallbackTo);
            return;
        }

        return;
    }

    router.back();
}
</script>

<template>
    <UTooltip :text="$t('common.back')">
        <UButton color="neutral" variant="subtle" icon="i-mdi-arrow-back" :to="to" @click="goBack">
            <span class="hidden truncate sm:block">
                {{ $t('common.back') }}
            </span>
        </UButton>
    </UTooltip>
</template>
