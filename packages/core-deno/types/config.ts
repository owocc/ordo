import type { IDE, Project } from "./store.ts";
export interface Config {
    projects: Project[];
    ides: IDE[];
    defaultIde?: IDE;
}
