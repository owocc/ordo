import Conf from 'npm:conf'
import type { Project } from '../model/project.ts'
import type { Config } from "../model/config.ts"
export const config = new Conf<Config>({
    projectName: 'ordo',
})

export const getProjects = () => {
    return config.get("projects") || []
}

export const saveProjects = (projects: Project[]) => {
    config.set('projects', projects)
}

export const addProject = (project: Project) => {
    const projects = getProjects()
    projects.push(
        project
    )
    saveProjects(projects)
}


export const getProjectByName = (name: string) => {
    const projects = getProjects()
    return projects.find(e => e.name == name)
}