<script lang="ts" setup>
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
    <div class="w-full">
        <CitizensList />
    </div>
</template>
