import { deleteProject } from "../../store/project.store.ts";
import { program } from "npm:commander";

const rmIdeCommand = program.createCommand("ide")
    .argument("[mark]", "rm mark")
    .action(() => {
    });

export const rmCommand = program.createCommand("rm")
    .description(
        "Remove project ",
    )
    .argument("[mark]", "rm mark")
    .action((mark: string) => {
        deleteProject(mark);
    })
    .addCommand(rmIdeCommand);
