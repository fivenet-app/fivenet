<script lang="ts" setup>
import ActivityFeed from '~/components/citizens/info/ActivityFeed.vue';
import type { User } from '~~/gen/ts/resources/users/users';

const user = defineModel<User>('user', { required: true });
</script>

<template>
    <ActivityFeed :user-id="user.userId" />
</template>
