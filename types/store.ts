export interface Project {
    name: string;
    dir: string;
    tag?: string[];
    id: string;
    relationIdeId?: string;
}

export interface IDE {
    name: string;
    path: string;
    args?: string;
    id: string;
    desc?:string
}

export interface ProjectAndIDE extends Project {
    ide?: IDE;
}
