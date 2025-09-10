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
        class="text-highlighted"
        :class="dispatchStatusAnimate(status) ? 'animate-pulse' : ''"
        :color="dispatchStatusToBadgeColor(status)"
        size="xs"
        :label="$t(`enums.centrum.StatusDispatch.${StatusDispatch[status ?? 0]}`)"
    />
</template>
