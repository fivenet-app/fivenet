<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import { dispatchStatusAnimate, dispatchStatusToBadgeColor } from '../helpers';

interface Props extends /* @vue-ignore */ BadgeProps {
    status?: StatusDispatch | undefined;
}

const props = withDefaults(defineProps<Props>(), {
    status: StatusDispatch.UNSPECIFIED,
});
</script>

<template>
    <UBadge
        class="text-highlighted"
        :class="[dispatchStatusAnimate(props.status) ? 'animate-pulse' : '']"
        :color="dispatchStatusToBadgeColor(props.status)"
        size="xs"
        :label="$t(`enums.centrum.StatusDispatch.${StatusDispatch[status ?? 0]}`)"
        v-bind="$attrs"
    />
</template>
