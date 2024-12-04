<script lang="ts" setup>
import type { RoutesNamedLocations } from '@typed-router';

const props = defineProps<{
    to?: RoutesNamedLocations | string;
    fallbackTo?: RoutesNamedLocations | string;
}>();

const router = useRouter();

function goBack(): void {
    if (props.to) {
        return;
    }

    if (history.length === 0 && props.fallbackTo) {
        navigateTo(props.fallbackTo);
        return;
    }

    router.back();
}
</script>

<template>
    <UButton color="black" icon="i-mdi-arrow-back" :to="to" :label="$t('common.back')" @click="goBack" />
</template>
