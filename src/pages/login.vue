<script lang="ts" setup>
import Login from '../components/Login.vue';
import ContentCenterWrapper from '../components/partials/ContentCenterWrapper.vue';
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
    <ContentCenterWrapper>
        <Login />
    </ContentCenterWrapper>
</template>
