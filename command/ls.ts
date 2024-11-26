import { getProjects } from "../lib/config.ts"
import { program } from 'npm:commander'
import consola from 'npm:consola';
import { solidBorder } from "../config/table/border.ts";
import {
    table
} from 'npm:table'

export const lsCommand = program.createCommand('ls')
    .description('Add a project to the manager with specified name and directory.')
    .action(() => {
        const projects = getProjects()
        const tableHeader = ['名称', '目录']
        const data = projects.map(e => [e.name, e.dir])

        consola.log(table([tableHeader, ...data], {
            border: solidBorder
        }))

    })

