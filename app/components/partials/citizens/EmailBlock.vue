<script lang="ts" setup>
withDefaults(
    defineProps<{
        email?: string;
        hideIcon?: boolean;
        hideEmail?: boolean;
        hideNaText?: boolean;
        padded?: boolean;
    }>(),
    {
        email: undefined,
        hideIcon: false,
        hideEmail: false,
        hideNaText: false,
    },
);
</script>

<template>
    <UButton
        v-if="email"
        :icon="!hideIcon ? 'i-mdi-email' : undefined"
        :padded="padded"
        variant="link"
        :label="!hideEmail ? email : $t('common.mail', 1)"
        :to="email ? `/mail?to=${email}` : undefined"
        class="break-all"
    />
    <span v-else-if="!hideNaText" class="inline-flex items-center gap-1">
        <UIcon v-if="!hideIcon" name="i-mdi-email" class="h-5 w-5" />

        {{ $t('common.na') }}
    </span>
</template>
