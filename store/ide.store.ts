import { config } from "../lib/config.ts";
import type { IDE } from "../types/store.ts";
import { v4 as uuid } from "npm:uuid";

/** Find All IDE*/
export const findAllIDE = () => {
    return config.get("ides") || [];
};

/** Find Many IDE*/
export const findManyIDE = (condition?: (item: IDE) => boolean) => {
    const allIDEs = findAllIDE();
    if (condition && typeof condition === "function") {
        return allIDEs.filter(condition);
    } else {
        return allIDEs;
    }
};

/** Find One IDE */
export function findOneIDE(idOrName: string) {
    const projects = findManyIDE((item) => {
        return item.id === idOrName || item.name === idOrName;
    });
    return projects.length >= 1 ? projects[0] : undefined;
}

/** Save IDE */
export function saveIDE(ides: IDE[]) {
    config.set("ides", ides);
}

export function addIde(ide: Omit<IDE, "id">) {
    const ides = findAllIDE();
    ides.push(
        {
            ...ide,
            id: uuid(),
        },
    );
    saveIDE(ides);
}
