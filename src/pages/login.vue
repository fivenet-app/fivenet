<script lang="ts">
import { defineComponent } from 'vue';
import { mapState } from 'vuex';
import Navbar from '../components/partials/Navbar.vue';
import Footer from '../components/partials/Footer.vue';
import Login from '../components/Login.vue';
import ContentWrapper from '../components/partials/ContentWrapper.vue';
import NavPageHeader from '../components/partials/NavPageHeader.vue';
import CharacterSelector from '../components/login/CharacterSelector.vue';

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
