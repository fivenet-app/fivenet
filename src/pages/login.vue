<script lang="ts" setup>
import Login from '../components/Login.vue';
import HeroFull from '../components/partials/HeroFull.vue';
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
    <HeroFull>
        <Login />
    </HeroFull>
</template>
