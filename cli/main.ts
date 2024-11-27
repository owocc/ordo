import { program } from 'npm:commander'
import { addCommand } from './command/add.ts'
import { openCommand } from './command/open.ts'
import { lsCommand } from "./command/ls.ts"
import {rmCommand} from "./command/rm.ts";

program
    .name('Ordo')
    .description('Ordo is a project management tool that helps you efficiently organize and manage your projects.')
    .addCommand(addCommand)
    .addCommand(openCommand)
    .addCommand(lsCommand)
    .addCommand(rmCommand)
program.parse()