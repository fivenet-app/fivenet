<script lang="ts">
import { mapState, mapActions } from 'vuex';

import { LoginRequest } from '@arpanet/gen/auth/auth_pb';

export default {
    data() {
        return {
            username: '',
            password: '',
        };
    },
    computed: {
        ...mapState([
            'loggingIn',
            'loginError',
            'accessToken',
        ]),
    },
    methods: {
        ...mapActions([
            'doLogin',
        ]),
        loginSubmit() {
            const req = new LoginRequest();
            req.setUsername(this.username);
            req.setPassword(this.password);
            this.doLogin(req);
        },
    },
};
</script>

<template>
    <div class="hero min-h-screen bg-base-200">
        <div class="hero-content flex-col lg:flex-row-reverse">
            <div class="text-center lg:text-left">
                <h1 class="text-5xl font-bold">Login now!</h1>
                <p class="py-6">Welcome to aRPaNet! Custom made to be better integrated with your FiveM servers.</p>
            </div>
            <form @submit.prevent="loginSubmit">
                <div class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
                    <div class="card-body">
                        <p v-if="loginError">{{ loginError }}</p>
                        <p v-if="accessToken">Login Successful</p>
                        <div class="form-control">
                            <label class="label">
                                <span class="label-text">Username</span>
                            </label>
                            <input type="text" placeholder="username" class="input input-bordered" v-model="username" />
                        </div>
                        <div class="form-control">
                            <label class="label">
                                <span class="label-text">Password</span>
                            </label>
                            <input type="password" placeholder="password" class="input input-bordered" v-model="password" />
                            <label class="label">
                                <a href="#" class="label-text-alt link link-hover">Forgot password?</a>
                            </label>
                        </div>
                        <div class="form-control mt-6">
                            <button class="btn btn-primary">Login</button>
                        </div>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>
