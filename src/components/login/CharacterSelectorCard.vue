<script lang="ts">
import { getAccountClient, handleGRPCError } from '../../grpc';
import { ChooseCharacterRequest } from '@arpanet/gen/services/auth/auth_pb';
import { defineComponent } from 'vue';
import { mapActions, mapState } from 'vuex';
import { RpcError } from 'grpc-web';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { parseQuery } from 'vue-router/auto';
import CharSexBadge from '../misc/CharSexBadge.vue';
import { getSecondsFormattedAsDuration } from '../../utils/time';

export default defineComponent({
    computed: {
        ...mapState(["lastCharID"]),
    },
    props: {
        char: {
            required: true,
            type: User,
        },
    },
    methods: {
        ...mapActions(["updateAccessToken", "updateActiveChar", "updatePermissions"]),
        chooseCharacter() {
            const req = new ChooseCharacterRequest();
            req.setUserid(this.char.getUserid());

            getAccountClient()
                .chooseCharacter(req, null)
                .then((resp) => {
                    this.updateAccessToken(resp.getToken());
                    this.updateActiveChar(this.char);
                    this.updatePermissions(resp.getPermissionsList());
                    console.log(resp.getPermissionsList());
                    const path = this.$route.query.redirect?.toString() || "/overview";
                    const url = new URL("https://example.com" + path);
                    this.$router.push({ path: url.pathname, query: parseQuery(url.search), hash: url.hash });
                }).catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                });
        },
        getSecondsFormattedAsDuration,
    },
    components: {
        CharSexBadge,
    }
});
</script>

<template>
    <li :key="char.getUserid()"
        class="col-span-2 flex flex-col divide-y divide-white rounded-lg bg-gray-800 text-center shadow">
        <div class="flex flex-1 flex-col p-8">
            <h2 class="mt-6 text-2xl font-medium text-white">
                {{ char.getFirstname() }}, {{ char.getLastname() }}
                <CharSexBadge :sex="char.getSex()" />
            </h2>
            <dl class="mt-1 flex flex-grow flex-col justify-between">
                <dd>
                    <span v-if="lastCharID == char.getUserid()"
                        class="inline-flex items-center rounded-full bg-green-100 px-3 py-0.5 text-sm font-medium text-green-800">
                        Last Used
                    </span>
                    <br v-else />
                </dd>
                <dd class="mt-3">
                    <span
                        class="inline-flex items-center rounded-md bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800">{{
                            char.getJob() }} (Rank: {{ char.getJobgrade() }})</span>
                </dd>
                <dt class="text-sm text-white">Date of Birth</dt>
                <dd class="text-sm text-gray-300">{{ char.getDateofbirth() }}</dd>
                <dt class="text-sm text-white">Height</dt>
                <dd class="text-sm text-gray-300">{{ char.getHeight() }}cm</dd>
                <dt class="text-sm text-white">Visum</dt>
                <dd class="text-sm text-gray-300">{{ char.getVisum() }}</dd>
                <dt class="text-sm text-white">Playtime</dt>
                <dd class="text-sm text-gray-300">{{ getSecondsFormattedAsDuration(char.getPlaytime()) }}</dd>
            </dl>
        </div>
        <div>
            <div class="-mt-px flex divide-x divide-white">
                <div class="flex w-0 flex-1">
                    <button @click="chooseCharacter()"
                        class="relative -mr-px inline-flex w-0 flex-1 items-center justify-center gap-x-3 rounded-bl-lg border border-transparent py-4 text-sm font-semibold text-white bg-gray-600">
                        Choose
                    </button>
                </div>
            </div>
        </div>
    </li>
</template>
