<script lang="ts" setup>
import ContentCenterWrapper from '../components/partials/ContentCenterWrapper.vue';
import CharacterSelector from '../components/login/CharacterSelector.vue';
</script>

<route lang="json">
{
    "name": "Character Selector",
    "meta": {
        "requiresAuth": true,
        "authOnlyToken": true
    }
}
</route>

<template>
    <ContentCenterWrapper>
        <CharacterSelector />
    </ContentCenterWrapper>
</template>
