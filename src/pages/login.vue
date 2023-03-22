<script lang="ts" setup>
import { useStore } from '../store/store';
import { computed } from 'vue';
import Footer from '../components/partials/Footer.vue';
import Login from '../components/Login.vue';
import ContentWrapper from '../components/partials/ContentWrapper.vue';
import NavPageHeader from '../components/partials/NavPageHeader.vue';
import CharacterSelector from '../components/login/CharacterSelector.vue';

const store = useStore();
const accessToken = computed(() => store.state.auth?.accessToken);
</script>

<route lang="json">
{
    "name": "Login",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <NavPageHeader v-if="!accessToken" title="Login" />
    <NavPageHeader v-else title="Character Selector" />
    <ContentWrapper>
        <transition name="fade" mode="out-in">
            <Login v-if="!accessToken" />
            <CharacterSelector v-else />
        </transition>
    </ContentWrapper>
    <Footer />
</template>
