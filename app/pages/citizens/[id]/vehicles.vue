<script lang="ts" setup>
import Vehicles from '~/components/citizens/info/Vehicles.vue';
import type { User } from '~~/gen/ts/resources/users/users';

const user = defineModel<User>('user', { required: true });
</script>

<template>
    <Vehicles :user-id="user.userId" />
</template>
