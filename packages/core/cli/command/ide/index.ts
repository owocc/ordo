import { program } from 'npm:commander'
import { addIdeCommand } from './add.ts'
import { lsIDECommand } from "./ls.ts"
import { rmIdeCommand } from "./rm.ts"
export const ideCommander = program.createCommand('ide')
    .addCommand(lsIDECommand)
    .addCommand(addIdeCommand)
    .addCommand(rmIdeCommand)