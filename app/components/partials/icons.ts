import * as mdiVue3 from 'mdi-vue3';
import type { DefineComponent } from 'vue';

export const fallbackIcon = mdiVue3.MapMarkerQuestionIcon;

export const availableIcons = Object.values(mdiVue3) as readonly DefineComponent[];
