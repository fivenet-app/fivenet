<script lang="ts" setup>
import type { File } from '~~/gen/ts/resources/file/file';

defineProps<{
    files: File[];
}>();
</script>

<template>
    <div>
        <ul>
            <li v-for="file in files" :key="file.id">{{ file.filePath }} ({{ file.id }})</li>
        </ul>
    </div>
</template>
