<script lang="ts" setup>
import Footer from '../components/partials/Footer.vue';
import ContentWrapper from '../components/partials/ContentWrapper.vue';
import NavPageHeader from '../components/partials/NavPageHeader.vue';
import CharacterSelector from '../components/login/CharacterSelector.vue';
</script>

<route lang="json">
{
    "name": "Character Selector",
    "meta": {
        "requiresAuth": true
    }
}
</route>

<template>
    <NavPageHeader title="Character Selector" />
    <ContentWrapper>
        <CharacterSelector />
    </ContentWrapper>
    <Footer />
</template>
