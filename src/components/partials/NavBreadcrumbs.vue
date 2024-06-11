<script lang="ts" setup>
const { t } = useI18n();

const breadcrumbs = useBreadcrumbs();

const links = computed(() => breadcrumbs.value.map((b) => ({ label: t(b.title), to: b.to, icon: b?.icon })));
</script>

<template>
    <UBreadcrumb class="hidden md:flex" :links="links" />
</template>
