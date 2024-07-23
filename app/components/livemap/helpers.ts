import * as mdiVue3 from 'mdi-vue3';
import type { DefineComponent } from 'vue';

export const markerFallbackIcon = mdiVue3.MapMarkerQuestionIcon;

export const markerIcons = Object.values(mdiVue3) as readonly DefineComponent[];

export function convertDynamicIconNameToComponent(search: string): string {
    return search.includes('i-mdi-')
        ? uppercaseFirstLetter(
              search
                  .replace('i-mdi-', '')
                  .toLowerCase()
                  .replace(/([-_][a-z])/g, (group) => group.toUpperCase().replace('-', '').replace('_', '')),
          ) + 'Icon'
        : search;
}
