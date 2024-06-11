import * as mdiVue3 from 'mdi-vue3';

export const markerIcons = Object.values(mdiVue3);

export const markerFallbackIcon = mdiVue3.MapMarkerQuestionIcon;

export function convertDynamicIconNameToComponent(s: string): string {
    return s.includes('i-mdi-')
        ? uppercaseFirstLetter(
              s
                  .replace('i-mdi-', '')
                  .toLowerCase()
                  .replace(/([-_][a-z])/g, (group) => group.toUpperCase().replace('-', '').replace('_', '')),
          ) + 'Icon'
        : s;
}
