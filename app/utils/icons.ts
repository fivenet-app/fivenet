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

export function convertComponentIconNameToDynamic(search: string): string {
    if (search.startsWith('i-mdi-')) return search;

    return (
        'i-mdi-' +
        search.replace(/Icon$/, '').replace(/[A-Z]+(?![a-z])|[A-Z1-9]/g, ($, ofs) => (ofs ? '-' : '') + $.toLowerCase())
    );
}
