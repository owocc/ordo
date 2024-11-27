import { program } from "npm:commander";
import { getProjectAndIDE } from "../../store/project.store.ts";
import { exec } from "node:child_process";
import consola from "consola";
import { getCommands } from "../../utils/command.ts";

export const openCommand = program.createCommand("open")
    .argument(
        "[idOrName...]",
        "Specify one or more project names or IDs to open",
    )
    .description("Open specified projects in their associated IDEs")
    .action(async (idOrName: string[]) => {
        if (idOrName.length === 0) {
            consola.warn(
                "No project names or IDs provided. Please specify at least one project.",
            );
            return;
        }

        const projects = getProjectAndIDE(idOrName);

        if (projects.length === 0) {
            consola.warn(
                `No matching projects found for: ${
                    idOrName.join(", ")
                }. Please check your input and try again.`,
            );
            return;
        }

        const commands = await getCommands(projects);
        if (commands.length === 0) {
            consola.warn(
                "No commands generated. Please ensure your IDE configuration is correct.",
            );
            return;
        }

        commands.forEach((command) => {
            consola.info(`Executing command: ${command}`);
            exec(command, (error) => {
                if (error) {
                    consola.error(`Failed to execute: ${command}`);
                } else {
                    consola.success(
                        `Successfully opened project with command: ${command}`,
                    );
                }
            });
        });
    });
