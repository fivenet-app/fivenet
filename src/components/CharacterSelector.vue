<script lang="ts">
import { defineComponent } from 'vue';
import CharacterCard from './CharacterCard.vue';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import { GetCharactersRequest } from '@arpanet/gen/auth/auth_pb';
import { Character } from '@arpanet/gen/common/character_pb';
import authInterceptor from '../grpcauth';
import { RpcError } from 'grpc-web';

export default defineComponent({
    components: {
        CharacterCard,
    },
    data: function () {
        return {
            'chars': [] as Array<Character>,
        };
    },
    methods: {
        fetchCharacters() {
            const req = new GetCharactersRequest();
            client.
                getCharacters(req, null).
                then((resp) => {
                    this.chars = resp.getCharsList();
                }).catch((err: RpcError) => {
                    authInterceptor.handleError(err, this.$route);
                });
        },
    },
    mounted() {
        // Fetch user's characters
        this.fetchCharacters();
    },
});
const client = new AccountServiceClient('https://localhost:8181', null, {
    unaryInterceptors: [authInterceptor],
    streamInterceptors: [authInterceptor],
});

</script>

<template>
    <div class="grid place-items-center">
        <div class="flex w-full">
            <CharacterCard v-for="char in chars" :char="char" :identifier="char.getIdentifier()"
                :key="char.getIdentifier()" />
        </div>
    </div>
</template>
