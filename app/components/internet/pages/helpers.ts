import type { Component } from 'vue';
import AdsManagement from './AdsManagement.vue';
import HomePage from './HomePage.vue';
import NicRegistrar from './NicRegistrar.vue';

export const localPages: Record<string, Component> = {
    'internet.search': HomePage,
    'nic.ls': NicRegistrar,
    'ads.ls': AdsManagement,
};

export const localPagesDomains = Object.keys(localPages);
