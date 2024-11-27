import { config } from "../lib/config.ts";
import { Project, ProjectAndIDE } from "../types/store.ts";
import { findAllIDE } from "./ide.store.ts";
import { v4 as uuid } from "npm:uuid";

// + [ 基础查询函数 ] +

/** 查找所有项目 */
export function findAllProjects() {
    return config.get("projects") || [];
}

/** 查找多个项目 */
export function findManyProject(condition?: (item: Project) => boolean) {
    const allProjects = findAllProjects();
    if (condition && typeof condition === "function") {
        return allProjects.filter(condition);
    } else {
        return allProjects;
    }
}

/** 查找多个项目及其对应的 IDE */
export function findManyProjectAndIDE(
    condition?: (item: Project) => boolean,
): ProjectAndIDE[] {
    const projects = findManyProject(condition);
    const ides = findAllIDE();
    // 创建 IDE 映射
    const ideMap = new Map(ides.map((ide) => [ide.id, ide]));
    return projects.map((project) => {
        return {
            ...project,
            ide: project?.relationIdeId
                ? ideMap.get(project.relationIdeId)
                : undefined,
        };
    });
}

/** 保存项目 */
export function saveProjects(projects: Project[]) {
    config.set("projects", projects);
}

/** 添加新的项目 */
export function addProject(project: Omit<Project, "id">) {
    const projects = findAllProjects();
    projects.push(
        {
            ...project,
            id: uuid(),
        },
    );
    saveProjects(projects);
}

/** 查找单个项目 */
export function findOneProject(idOrName: string) {
    const projects = findManyProject((item) => {
        return item.id === idOrName || item.name === idOrName;
    });
    return projects.length >= 1 ? projects[0] : undefined;
}

/** 根据项目名称或 ID 获取项目及其 IDE */
export function getProjectAndIDE(idOrNames: string[]) {
    return findManyProjectAndIDE((project) =>
        idOrNames.some((idOrName) =>
            idOrName === project.name || idOrName === project.id
        )
    );
}

/** 删除项目 */
export function deleteProject(mark: string) {
    const projects = findAllProjects();
    saveProjects(projects.filter((e) => e.name !== mark && e.id !== mark));
}