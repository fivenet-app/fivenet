<script lang="ts">
import { defineComponent } from 'vue';
import { mapState } from 'vuex';
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Login from '../components/Login.vue';
import ContentWrapper from '../components/ContentWrapper.vue';
import NavPageHeader from '../components/NavPageHeader.vue';
import CharacterSelector from '../components/CharacterSelector.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        Login,
        NavPageHeader,
        ContentWrapper,
        CharacterSelector,
    },
    computed: {
        ...mapState({
            accessToken: 'accessToken',
        }),
    },
});
</script>

<route lang="json">
{
    "name": "login",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <Navbar />
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
