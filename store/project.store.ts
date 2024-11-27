import { config } from "../lib/config.ts";
import { Project, ProjectAndIDE } from "../types/store.ts";
import { findAllIDE } from "./ide.store.ts";
import { v4 as uuid } from "npm:uuid";

// + [ Base Query Functions ] +

/** Find All Projects*/
export function findAllProjects() {
    return config.get("projects") || [];
}

/** Find Many Project */
export function findManyProject(condition?: (item: Project) => boolean) {
    const allProjects = findAllProjects();
    if (condition && typeof condition === "function") {
        return allProjects.filter(condition);
    } else {
        return allProjects;
    }
}

export function findManyProjectAndIDE(
    condition?: (item: Project) => boolean,
): ProjectAndIDE[] {
    const projects = findManyProject(condition);
    const ides = findAllIDE();
    // create map
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

/** Save a Projects */
export function saveProjects(projects: Project[]) {
    config.set("projects", projects);
}

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

export function findOneProject(idOrName: string) {
    const projects = findManyProject((item) => {
        return item.id === idOrName || item.name === idOrName;
    });
    return projects.length >= 1 ? projects[0] : undefined;
}

export function getProjectAndIDE(idOrNames: string[]) {
    return findManyProjectAndIDE((project) =>
        idOrNames.some((idOrName) =>
            idOrName === project.name || idOrName === project.id
        )
    );
}

export function deleteProject(mark: string) {
    const projects = findManyProject();
    saveProjects(projects.filter((e) => e.name !== mark && e.id !== mark));
}
