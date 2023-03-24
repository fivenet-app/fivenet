<script lang="ts" setup>
import Login from '../components/Login.vue';
import ContentWrapper from '../components/partials/ContentWrapper.vue';
</script>

<route lang="json">
{
    "name": "Login",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <ContentWrapper>
        <Login />
    </ContentWrapper>
</template>
