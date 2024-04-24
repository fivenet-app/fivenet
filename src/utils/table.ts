export function wrapRows(rows: { [key: string]: any }[], columns: { key: string; rowClass?: string }[]): any[] {
    return rows.map((row) => {
        columns.forEach((column) => {
            if (!column || !column.rowClass) {
                return;
            }

            row[column.key] = { class: column.rowClass, value: row[column.key] };
        });

        return row;
    });
}
