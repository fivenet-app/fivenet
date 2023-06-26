<script lang="ts" setup>
const statuses: { [key: string]: string } = {
    offline: 'animate-pulse text-gray-500 bg-gray-100/10',
    online: 'text-green-400 bg-green-400/10',
    error: 'text-rose-400 bg-rose-400/10',
};

const projects = [
    {
        id: 1,
        href: '#',
        name: 'ios-app',
        job: 'LSPD',
        status: 'offline',
        statusText: 'Initiated 1m 32s ago',
        description: 'Deploys from GitHub',
        environment: 'Preview',
    },
    {
        id: 2,
        href: '#',
        name: 'mobile-api',
        job: 'LSPD',
        status: 'online',
        statusText: 'Deployed 3m ago',
        description: 'Deploys from GitHub',
        environment: 'Production',
    },
    {
        id: 3,
        href: '#',
        name: 'tailwindcss.com',
        job: 'LSPD',
        status: 'offline',
        statusText: 'Deployed 3h ago',
        description: 'Deploys from GitHub',
        environment: 'Preview',
    },
];
</script>

<template>
    <div class="mx-2">TODO</div>
</template>
