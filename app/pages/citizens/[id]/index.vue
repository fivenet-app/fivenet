<script lang="ts" setup>
import Profile from '~/components/citizens/info/Profile.vue';
import type { User } from '~~/gen/ts/resources/users/users';

const user = defineModel<User>('user', { required: true });
</script>

<template>
    <Profile :user="user" />
</template>
