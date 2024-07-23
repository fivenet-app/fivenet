<script lang="ts" setup>
defineProps<{
    sex: string;
}>();
</script>

<template>
    <UBadge :color="sex === 'f' ? 'pink' : 'blue'" size="lg">
        <UIcon :name="sex === 'f' ? 'i-mdi-gender-female' : 'i-mdi-gender-male'" />
    </UBadge>
</template>
