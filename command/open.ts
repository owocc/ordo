import { program } from 'npm:commander'
import { getProjectByName } from '../lib/config.ts'
import { exec } from "node:child_process";

interface OpenOptions {
    name: string
}
export const openCommand = program.createCommand('open')
    .option('-n, --name [string]', 'open you project')
    .action((opts: OpenOptions) => {
        const {
            name
        } = opts
        const project = getProjectByName(name)
        if (project) {
            exec(`code ${project.dir}`)
        }

    })
