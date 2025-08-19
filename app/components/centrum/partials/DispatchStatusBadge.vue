<script lang="ts" setup>
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { dispatchStatusAnimate, dispatchStatusToBadgeColor } from '../helpers';

withDefaults(
    defineProps<{
        status?: StatusDispatch | undefined;
    }>(),
    {
        status: StatusDispatch.UNSPECIFIED,
    },
);
</script>

<template>
    <UBadge
        class="text-gray-900 dark:text-white"
        :class="dispatchStatusAnimate(status) ? 'animate-pulse' : ''"
        :color="dispatchStatusToBadgeColor(status)"
        size="xs"
    >
        {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[status ?? 0]}`) }}
    </UBadge>
</template>
