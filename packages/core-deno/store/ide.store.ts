import { config } from "../lib/config.ts";
import type { IDE } from "../types/store.ts";
import { v4 as uuid } from "npm:uuid";

/** 查找所有 IDE */
export const findAllIDE = () => {
    return config.get("ides") || [];
};

/** 查找多个 IDE */
export const findManyIDE = (condition?: (item: IDE) => boolean) => {
    const allIDEs = findAllIDE();
    if (condition && typeof condition === "function") {
        return allIDEs.filter(condition);
    } else {
        return allIDEs;
    }
};

/** 查找单个 IDE */
export function findOneIDE(idOrName: string) {
    const projects = findManyIDE((item) => {
        return item.id === idOrName || item.name === idOrName;
    });
    return projects.length >= 1 ? projects[0] : undefined;
}

/** 保存 IDE */
export function saveIDE(ides: IDE[]) {
    config.set("ides", ides);
}

/** 添加新的 IDE */
export function addIDE(ide: Omit<IDE, "id">) {
    const ides = findAllIDE();
    ides.push(
        {
            ...ide,
            id: uuid(),
        },
    );
    saveIDE(ides);
}

/** 删除 IDE */
export function deleteIDE(mark: string) {
    const projects = findAllIDE();
    saveIDE(projects.filter((e) => e.name !== mark && e.id !== mark));
}