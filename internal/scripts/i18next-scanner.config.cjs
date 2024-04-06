module.exports = {
    input: [
        './src/**/*.{htm,html,js,jsx,vue,ts}',
        // Use ! to filter out files or directories
        '!**/node_modules/**',
    ],
    output: './',
    options: {
        debug: false,
        func: {
            list: ['t'],
            extensions: ['.js', '.jsx', '.vue', '.html', '.htm', '.ts'],
        },
        trans: {
            component: 'Trans',
            extensions: [],
        },
        lngs: ['en'],
        ns: ['lang'],
        defaultLng: 'en',
        defaultNs: 'lang',
        defaultValue: '__STRING_NOT_TRANSLATED__',
        resource: {
            loadPath: './src/lang/{{lng}}.json',
            savePath: 'lang/{{lng}}.json',
            jsonIndent: 4,
            lineEnding: '\n',
        },
        nsSeparator: false,
        keySeparator: '.',
        interpolation: {
            prefix: '{{',
            suffix: '}}',
        },
        allowDynamicKeys: true,
    },
};
