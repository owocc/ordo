import { program } from 'npm:commander'
import { addCommand } from './command/project/add.ts'
import { openCommand } from './command/project/open.ts'
import { lsCommand } from "./command/project/ls.ts"
import { rmCommand } from "./command/project/rm.ts";
import { ideCommander } from './command/ide/index.ts'

program
    .name('Ordo')
    .description('Ordo is a project management tool that helps you efficiently organize and manage your projects.')
    .addCommand(addCommand)
    .addCommand(openCommand)
    .addCommand(lsCommand)
    .addCommand(rmCommand)
    .addCommand(ideCommander)
program.parse()

