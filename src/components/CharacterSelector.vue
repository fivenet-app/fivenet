<script lang="ts">
import { defineComponent } from 'vue';
import CharacterSelectorCard from './CharacterSelectorCard.vue';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import { GetCharactersRequest } from '@arpanet/gen/auth/auth_pb';
import { Character } from '@arpanet/gen/common/character_pb';
import authInterceptor from '../grpcauth';
import { RpcError } from 'grpc-web';
import { mapActions } from 'vuex';

export default defineComponent({
    components: {
        CharacterSelectorCard,
    },
    data: function () {
        return {
            'chars': [] as Array<Character>,
        };
    },
    methods: {
        ...mapActions([
            'updateActiveChar',
            'updateActiveCharIdentifier',
        ]),
        fetchCharacters() {
            client.
                getCharacters(new GetCharactersRequest(), null).
                then((resp) => {
                    this.chars = resp.getCharsList();
                }).catch((err: RpcError) => {
                    authInterceptor.handleError(err, this.$route);
                });
        },
    },
    beforeMount() {
        this.updateActiveChar(null);
        this.updateActiveCharIdentifier(null);
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
            <CharacterSelectorCard v-for="char in chars" :char="char" :identifier="char.getIdentifier()"
                :key="char.getIdentifier()" />
        </div>
    </div>
</template>
