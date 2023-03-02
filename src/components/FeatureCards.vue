<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    data() {
        return {
            features: {
                citizens: {
                    name: 'Citizen Search',
                    image: '/images/features/citizens_database.jpg',
                    description: 'Search and find information about Citizens, including basic info, their licenses and related documents.',
                    to: '/citizens',
                },
                documents: {
                    name: 'Documents',
                    image: '/images/features/documents.jpg',
                    description: 'Create and share documents for and about citizens easily with your colleagues and other jobs.',
                    to: '/documents',
                },
                livemap: {
                    name: 'Livemap',
                    image: '/images/features/livemap.jpg',
                    description: 'Live position of dispatches and your colleagues.',
                    to: '/livemap',
                },
            },
        };
    },
});
</script>

<template>
    <div class="grid grid-cols-3">
        <div v-for="feature in features" class="card card-compact w-96 bg-base-100 shadow-xl image-full">
            <figure><img :src="feature.image" :alt="feature.name" /></figure>
            <div class="card-body">
                <h2 class="card-title">{{ feature.name }}</h2>
                <p>{{ feature.description }}</p>
                <div class="card-actions justify-end">
                    <router-link :to="feature.to" class="btn btn-primary">Goto</router-link>
                </div>
            </div>
        </div>
    </div>
</template>
