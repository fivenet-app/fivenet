export function wrapRows(rows: { [key: string]: any }[], columns: { key: string; rowClass?: string }[]): any[] {
    return rows.map((row) => {
        for (const key in row) {
            const column = columns.find((c) => c.key === key);
            if (!column || !column.rowClass) {
                continue;
            }

            row[key] = { class: column.rowClass, value: row[key] };
        }

        return row;
    });
}
