type JobItem = {
    job: string;
};

export function groupByJob<T extends JobItem>(items: Iterable<T>, predicate?: (item: T) => boolean): Map<string, T[]> {
    const grouped = new Map<string, T[]>();

    for (const item of items) {
        if (predicate && !predicate(item)) continue;

        const current = grouped.get(item.job);
        if (current) {
            current.push(item);
        } else {
            grouped.set(item.job, [item]);
        }
    }

    return grouped;
}
