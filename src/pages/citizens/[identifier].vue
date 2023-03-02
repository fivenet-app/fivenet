<script lang="ts">
import { Character } from '@arpanet/gen/common/character_pb';
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import { GetUserRequest } from '@arpanet/gen/users/users_pb';
import { RpcError } from 'grpc-web';
import { defineComponent } from 'vue';
import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';
import CitizenInfo from '../../components/CitizenInfo.vue';
import authInterceptor from '../../grpcauth';
import ContentWrapper from '../../components/ContentWrapper.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        CitizenInfo,
        ContentWrapper,
    },
    data() {
        return {
            client: new UsersServiceClient('https://localhost:8181', null, {
                unaryInterceptors: [authInterceptor],
                streamInterceptors: [authInterceptor],
            }),
            char: null as null | Character,
        };
    },
    methods: {
        async fetchCitizen(identifier: string) {
            const req = new GetUserRequest();
            req.setIdentifier(identifier);
            this.client.
                getUser(req, null).
                then((resp) => {
                    this.char = resp.getUser() as Character;
                }).catch((err: RpcError) => {
                    authInterceptor.handleError(err, this.$route);
                });
        },
    },
    mounted() {
        // Fetch the user information when params change
        //@ts-ignore
        this.fetchCitizen(this.$route.params.identifier);
    },
});
</script>

<route lang="json">
{
    "name": "citizens-byid",
    "meta": {
        "requiresAuth": true
    }
}
</route>

<template>
    <Navbar />
    <ContentWrapper>
        <CitizenInfo v-if="char" :char="char" />
    </ContentWrapper>
    <Footer />
</template>
