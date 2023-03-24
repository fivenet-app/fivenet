<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import CitizensList from '../../components/citizens/CitizensList.vue';
</script>

<route lang="json">
{
    "name": "Citizens",
    "meta": {
        "requiresAuth": true,
        "permission": "CitizenStoreService.FindUsers"
    }
}
</route>

<template>
    <ContentWrapper>
        <CitizensList />
    </ContentWrapper>
</template>
