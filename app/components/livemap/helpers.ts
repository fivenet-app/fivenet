import * as mdiVue3 from 'mdi-vue3';
import type { DefineComponent } from 'vue';

export const markerFallbackIcon = mdiVue3.MapMarkerQuestionIcon;

export const markerIcons = Object.values(mdiVue3) as readonly DefineComponent[];
