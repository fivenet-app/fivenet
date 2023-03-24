<script lang="ts" setup>
import Footer from '../components/partials/Footer.vue';
import ContentWrapper from '../components/partials/ContentWrapper.vue';
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
    <ContentWrapper>
        <CharacterSelector />
    </ContentWrapper>
    <Footer />
</template>
