<script lang="ts">
import { Character } from '@arpanet/gen/common/character_pb';
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import { GetUserRequest } from '@arpanet/gen/users/users_pb';
import { RpcError } from 'grpc-web';
import { defineComponent } from 'vue';
import { ref, watch } from 'vue';
import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';
import CitizenInfo from '../../components/CitizenInfo.vue';
import authInterceptor from '../../grpcauth';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        CitizenInfo,
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
    <div class="container mx-auto py-8">
        <div class="text-sm breadcrumbs">
            <ul>
                <li>
                    <a>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            class="w-4 h-4 mr-2 stroke-current">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
                        </svg>
                        Overview
                    </a>
                </li>
                <li>
                    <a>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            class="w-4 h-4 mr-2 stroke-current">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
                        </svg>
                        Citizens
                    </a>
                </li>
                <li>
                    <a>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            class="w-4 h-4 mr-2 stroke-current">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
                        </svg>
                        <span v-if="!char">Citizen</span>
                        <span v-else>{{ char.getFirstname() }}, {{ char.getLastname() }}</span>
                    </a>
                </li>
            </ul>
        </div>
        <CitizenInfo v-if="char" :char="char" />
    </div>
    <Footer />
</template>
