<script lang="ts" setup>
import type { Label } from '~~/gen/ts/resources/citizens/labels/labels';

defineProps<{
    id?: number;
    label?: Label;
}>();

// TODO resolve label if only id is given
</script>

<template>
    <UBadge
        v-if="label"
        :class="[isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!']"
        :style="{ backgroundColor: label.color }"
        :icon="label.icon && label.icon !== '' ? convertComponentIconNameToDynamic(label.icon) : undefined"
    >
        <div class="inline-flex flex-col gap-1">
            <span>{{ label.name }}</span>

            <div v-if="label.expiresAt">({{ $t('common.expires_at') }} {{ $d(toDate(label.expiresAt), 'short') }})</div>
        </div>
    </UBadge>
</template>
