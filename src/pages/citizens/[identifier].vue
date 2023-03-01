<script lang="ts">
import { Character } from '@arpanet/gen/common/character_pb';
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import { GetUserRequest } from '@arpanet/gen/users/users_pb';
import { RpcError } from 'grpc-web';
import { defineComponent } from 'vue';
import { ref, watch } from 'vue';
import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';
import UserInfo from '../../components/UserInfo.vue';
import authInterceptor from '../../grpcauth';

export default defineComponent({
    data: function () {
        return {
            "user": null as null | Character,
        };
    },
    methods: {
        async fetchCitizen(identifier: string) {
            const req = new GetUserRequest();
            req.setIdentifier(identifier);
            client.
                getUser(req, null).
                then((resp) => {
                this.user = resp.getUser() as Character;
            }).catch((err: RpcError) => {
                authInterceptor.handleError(err, this.$route);
            });
        },
    },
    mounted() {
        const userData = ref();
        // Fetch the user information when params change
        //@ts-ignore
        watch(() => this.$route.params.identifier, async (newIdentifier) => {
            userData.value = await this.fetchCitizen(newIdentifier);
        });
    },
    components: {
        Navbar,
        Footer,
        UserInfo,
    }
});

const client = new UsersServiceClient('https://localhost:8181', null, {
    unaryInterceptors: [authInterceptor],
    streamInterceptors: [authInterceptor],
});
</script>

<route lang="json">
{
    "name": "citizens-byid",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <Navbar />
    <UserInfo />
    <Footer />
</template>
