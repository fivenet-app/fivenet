<script lang="ts" setup>
import List from '~/components/vehicles/List.vue';
import type { User } from '~~/gen/ts/resources/users/users';

const user = defineModel<User>('user', { required: true });
</script>

<template>
    <List :user-id="user.userId" hide-owner hide-citizen-link search-key="citizen_vehicles" />
</template>
