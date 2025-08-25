<script lang="ts" setup>
import Documents from '~/components/citizens/info/Documents.vue';
import type { User } from '~~/gen/ts/resources/users/users';

const user = defineModel<User>('user', { required: true });
</script>

<template>
    <Documents :user-id="user.userId" />
</template>
