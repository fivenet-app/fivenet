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
