<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';
import NavPageHeader from '../../components/NavPageHeader.vue';
import ContentWrapper from '../../components/ContentWrapper.vue';

export default defineComponent({
    components: {
    Navbar,
    Footer,
    ContentWrapper,
    NavPageHeader
},
});
</script>

<route lang="json">
{
    "name": "dispatches",
    "meta": {
        "requiresAuth": true
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Dispatches" />
    <ContentWrapper>
        <h1>DISPATCHES SYSTEM</h1>
    </ContentWrapper>
    <Footer />
</template>
