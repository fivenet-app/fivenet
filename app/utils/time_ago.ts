const secondsPerMinute = 60;
const secondsPerHour = secondsPerMinute * 60;
const secondsPerDay = secondsPerHour * 24;
const secondsPerWeek = secondsPerDay * 7;
const secondsPerYear = secondsPerWeek * 52;

export function fromSecondsToFormattedDuration(seconds: number, options?: { seconds?: boolean; emptyText?: string }): string {
    const { t } = useI18n();

    const years = Math.floor(seconds / secondsPerYear);
    seconds -= years * secondsPerYear;
    const weeks = Math.floor(seconds / secondsPerWeek);
    seconds -= weeks * secondsPerWeek;
    const days = Math.floor(seconds / secondsPerDay);
    seconds -= days * secondsPerDay;
    const hours = Math.floor(seconds / secondsPerHour);
    seconds -= hours * secondsPerHour;
    const minutes = Math.floor(seconds / secondsPerMinute);
    seconds -= minutes * secondsPerMinute;

    const parts: string[] = [];
    if (years > 0) {
        parts.push(`${years} ${t(`common.time_ago.year`, years)}`);
    }
    if (weeks > 0) {
        parts.push(`${weeks} ${t(`common.time_ago.week`, weeks)}`);
    }
    if (days > 0) {
        parts.push(`${days} ${t(`common.time_ago.day`, days)}`);
    }
    if (hours > 0) {
        parts.push(`${hours} ${t(`common.time_ago.hour`, hours)}`);
    }
    if (minutes > 0) {
        parts.push(`${minutes} ${t(`common.time_ago.minute`, minutes)}`);
    }
    if ((!options || options.seconds) && seconds > 0) {
        parts.push(`${seconds} ${t(`common.time_ago.second`, seconds)}`);
    }

    const text = parts.join(', ');
    return text.length > 0 ? text : t(options?.emptyText ?? 'common.unknown');
}
