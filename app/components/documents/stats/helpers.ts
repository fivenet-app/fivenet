import { addDays, addWeeks, differenceInDays, startOfDay, startOfWeek } from 'date-fns';
import { toDate, toTimestamp } from '~/utils/time';
import { type DailyValue, type PeriodSeriesValue, StatsPeriod } from '~~/gen/ts/resources/stats/stats';

export interface Range {
    start: Date;
    end: Date;
}

export type DataRecord = {
    date: Date;
    amount: number;
    fine: number;
    detention: number;
    points: number;
    absent: number;
    vacation: number;
};

export type ChartStats = {
    periodValues: DailyValue[];
    periodSeriesValues: PeriodSeriesValue[];
    totalValue: number;
    averageValue: number;
};

type BuildChartDataOptions = {
    stats: ChartStats;
    isPenalties: boolean;
    period: StatsPeriod;
    range: Range;
};

export function getSelectedPeriod(period: StatsPeriod, range: Range): StatsPeriod {
    if (period !== StatsPeriod.UNSPECIFIED) {
        return period;
    }

    const days = Math.max(differenceInDays(range.end, range.start), 1);
    if (days > 365) {
        return StatsPeriod.MONTHLY;
    }

    return days > 60 ? StatsPeriod.WEEKLY : StatsPeriod.DAILY;
}

function getBucketStart(date: Date, period: StatsPeriod): Date {
    if (period === StatsPeriod.MONTHLY) {
        return startOfDay(new Date(date.getFullYear(), date.getMonth(), 1));
    }

    return period === StatsPeriod.WEEKLY ? startOfWeek(date, { weekStartsOn: 1 }) : startOfDay(date);
}

function getNextBucket(date: Date, period: StatsPeriod): Date {
    if (period === StatsPeriod.MONTHLY) {
        return startOfDay(new Date(date.getFullYear(), date.getMonth() + 1, 1));
    }

    return period === StatsPeriod.WEEKLY ? addWeeks(date, 1) : addDays(date, 1);
}

function createDataRecord(time: number): DataRecord {
    return {
        date: new Date(time),
        amount: 0,
        fine: 0,
        detention: 0,
        points: 0,
        absent: 0,
        vacation: 0,
    };
}

export function buildChartData(options: BuildChartDataOptions): DataRecord[] {
    const period = getSelectedPeriod(options.period, options.range);
    const valueByBucket = new Map<number, DataRecord>();

    const ensureBucket = (time: number): DataRecord => {
        const existing = valueByBucket.get(time);
        if (existing) {
            return existing;
        }

        const created = createDataRecord(time);
        valueByBucket.set(time, created);
        return created;
    };

    if (options.isPenalties) {
        const values = options.stats.periodSeriesValues ?? [];
        for (const item of values) {
            if (!item.day) {
                continue;
            }

            const bucket = getBucketStart(toDate(item.day), period).getTime();
            const target = ensureBucket(bucket);

            switch (item.key) {
                case 'fine_total':
                    target.fine += item.value;
                    break;
                case 'detention_time_total':
                    target.detention += item.value;
                    break;
                case 'stvo_points_total':
                    target.points += item.value;
                    break;
            }

            target.amount = target.fine + target.detention + target.points;
        }
    } else {
        const seriesValues = options.stats.periodSeriesValues ?? [];
        for (const item of seriesValues) {
            if (!item.day) {
                continue;
            }

            const bucket = getBucketStart(toDate(item.day), period).getTime();
            const target = ensureBucket(bucket);

            switch (item.key) {
                case 'on_vacation_count':
                case 'vacation_count':
                    target.vacation += item.value;
                    break;
                case 'absence_count':
                case 'absent_count':
                case 'on_absence_count':
                    target.absent += item.value;
                    break;
            }
        }

        const values = options.stats.periodValues ?? [];
        for (const item of values) {
            if (!item.day) {
                continue;
            }

            const bucket = getBucketStart(toDate(item.day), period).getTime();
            const target = ensureBucket(bucket);
            target.amount += item.value;
        }
    }

    const data: DataRecord[] = [];
    let cursor = getBucketStart(options.range.start, period);
    const end = getBucketStart(options.range.end, period).getTime();

    while (cursor.getTime() <= end) {
        const key = cursor.getTime();
        const item = valueByBucket.get(key);
        data.push({
            date: new Date(key),
            amount: item?.amount ?? 0,
            fine: item?.fine ?? 0,
            detention: item?.detention ?? 0,
            points: item?.points ?? 0,
            absent: item?.absent ?? 0,
            vacation: item?.vacation ?? 0,
        });
        cursor = getNextBucket(cursor, period);
    }

    return data;
}

type FillEmployeeSeriesOptions = {
    stats: ChartStats;
    period: StatsPeriod;
    range: Range;
};

export function fillEmployeeSeriesData(options: FillEmployeeSeriesOptions): ChartStats {
    const data = buildChartData({
        stats: options.stats,
        isPenalties: false,
        period: options.period,
        range: options.range,
    });

    return {
        ...options.stats,
        periodValues: data.map((item) => ({
            day: toTimestamp(item.date),
            value: item.amount,
        })),
        periodSeriesValues: data.flatMap((item) => [
            {
                day: toTimestamp(item.date),
                key: 'employee_count',
                label: 'employee_count',
                value: item.amount,
            },
            {
                day: toTimestamp(item.date),
                key: 'on_vacation_count',
                label: 'on_vacation_count',
                value: item.vacation,
            },
        ]),
    };
}
