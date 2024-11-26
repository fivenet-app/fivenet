export function toTitleCase(s: string): string {
    return s.replace(/\w\S*/g, (w) => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase());
}

export function camelCaseToTitleCase(s: string): string {
    return (
        s
            .replace(/([A-Z0-9])/g, ' $1')
            // uppercase the first character
            .replace(/^./, function (str) {
                return str.toUpperCase();
            })
    );
}

export function lowercaseFirstLetter(s: string): string {
    return s.charAt(0).toLowerCase() + s.slice(1);
}

export function uppercaseFirstLetter(s: string): string {
    return s.charAt(0).toUpperCase() + s.slice(1);
}

const initialsCleanerRegex = /(Prof\.|Dr\.|Sr(\.| ))[ ]*/gm;

export function getInitials(input: string): string {
    input = input.replaceAll(initialsCleanerRegex, '');
    const names = input.split(' ');
    // Indicates a "broken" name if there are not at least "two parts"
    if (!names[0] || names.length < 2) {
        return input;
    }

    let initials = names[0].substring(0, 1).toUpperCase();
    if (names.length > 1) {
        initials += names[names.length - 1]?.substring(0, 1).toUpperCase();
    }

    return initials;
}

// Taken from https://stackoverflow.com/a/18650828
export function formatBytes(bytes: number, decimals = 2) {
    if (!+bytes) return '0 Bytes';

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

export interface UserLike {
    userId: number;
    firstname: string;
    lastname: string;
    dateofbirth?: string;
    job: string;
    jobGrade: number;
}

export function usersToLabel(users: UserLike[]): string {
    return users.map((c) => `${c?.firstname} ${c?.lastname} (${c?.dateofbirth})`).join(', ');
}

export function userToLabel(user: UserLike): string {
    return `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})`;
}
