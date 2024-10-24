/* eslint-disable  @typescript-eslint/no-explicit-any */
import { LogLevels, consola, type ConsolaInstance, type LogLevel } from 'consola';

export interface ILogger {
    log(message?: any, ...optionalParams: any[]): void;
    debug(message?: any, ...optionalParams: any[]): void;
    info(message?: any, ...optionalParams: any[]): void;
    warn(message?: any, ...optionalParams: any[]): void;
    error(message?: any, ...optionalParams: any[]): void;
}

export class Logger {
    readonly consola: ConsolaInstance;
    readonly prefix?: string;

    constructor(prefix?: string, level?: LogLevel) {
        this.prefix = prefix?.trim();

        if (this.prefix !== undefined) {
            this.consola = consola.withTag(this.prefix);
        } else {
            this.consola = consola.create({});
        }

        if (level === undefined) {
            // Debug level for dev env if unset
            if (import.meta.dev) {
                level = LogLevels.debug;
            } else {
                level = LogLevels.warn;
            }
            this.consola.level = level;
        }
    }

    log(message?: any, ...optionalParams: any[]): void {
        this.consola.log(message, ...optionalParams);
    }

    debug(message?: any, ...optionalParams: any[]): void {
        this.consola.debug(message, ...optionalParams);
    }

    info(message?: any, ...optionalParams: any[]): void {
        this.consola.info(message, ...optionalParams);
    }

    warn(message?: any, ...optionalParams: any[]): void {
        this.consola.warn(message, ...optionalParams);
    }

    error(message?: any, ...optionalParams: any[]): void {
        this.consola.error(message, ...optionalParams);
    }
}

export function setDefaultLogLevel(level: LogLevel): void {
    consola.level = level;
}

export function getDefaultLogLevel(): LogLevel {
    return consola.level;
}
