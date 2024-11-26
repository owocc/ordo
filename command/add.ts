import { program } from 'npm:commander'
import { resolve, isAbsolute } from "jsr:@std/path";
import { cwd } from "node:process";
import { addProject } from '../lib/config.ts'

interface AddOptions {
    dir: string
    name: string
}
// opts.dir
const parsePath = (path: string) => {
    if (isAbsolute(path)) {
        return path
    }
    return resolve(cwd(), path)
}

export const addCommand = program.createCommand('add')
    .description('Add a project to the manager with specified name and directory.')
    .option('-n, --name [string]', 'Specify the name of the project')
    .option('-d, --dir [dirs]', 'Specify the directory of the project to add to the manager')
    .action((opts: AddOptions) => {
        const {
            dir,
            name
        } = opts
        const path = parsePath(dir)
        addProject({
            dir: path,
            name,
        })
    })

