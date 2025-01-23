import { isAbsolute, resolve } from "jsr:@std/path";
import { cwd } from "node:process";

export function parsePath(path: string): string {
    if (isAbsolute(path)) {
        return path;
    }
    return resolve(cwd(), path);
}
