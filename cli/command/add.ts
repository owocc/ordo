import { program } from "npm:commander";
import { addProject } from "../../store/project.store.ts";
import { parsePath } from "../../utils/path.ts";
import consola from "consola";
import { addIde, findAllIDE, findOneIDE } from "../../store/ide.store.ts";

/**
 * Adds a new IDE to the manager.
 */
const addIdeCommand = program.createCommand("ide")
    .argument("[name]", "IDE Name")
    .argument("[pathOrCommand]", "IDE Path or Command")
    .description("Add a new IDE to the manager.")
    .action(async (name?: string, pathOrCommand?: string) => {
        name ||= await consola.prompt("Please input the IDE name:", {
            type: "text",
            placeholder: "e.g., VSCode",
        });
        if (!name) {
            consola.error("IDE name is required!");
            return;
        }

        pathOrCommand ||= await consola.prompt(
            "Please input the IDE path or command:",
            {
                type: "text",
                placeholder: "e.g., /usr/bin/code or idea",
            },
        );
        if (!pathOrCommand) {
            consola.error("IDE path or command is required!");
            return;
        }

        addIde({
            name,
            path: pathOrCommand,
        });
        consola.success(`IDE '${name}' added successfully!`);
    });

/**
 * Adds a new project to the manager.
 */
export const addCommand = program.createCommand("add")
    .description(
        "Add a project to the manager with a specified name, directory, and optional IDE.",
    )
    .argument("[name]", "Project Name")
    .argument("[dir]", "Project Directory")
    .argument("[ide]", "IDE Name or ID")
    .action(async (name?: string, dir?: string, ideIdOrName?: string) => {
        // Prompt for missing name
        name ||= await consola.prompt("Please input your project name:", {
            type: "text",
        });
        if (!name) {
            consola.error("Project name is required!");
            return;
        }

        // Prompt for missing directory
        dir ||= await consola.prompt("Please input your project directory:", {
            type: "text",
        });
        if (!dir) {
            consola.error("Project directory is required!");
            return;
        }

        // Resolve IDE if not provided
        const allIdes = findAllIDE();
        ideIdOrName ||= await consola.prompt(
            "Please select an IDE (or skip):",
            {
                type: "select",
                options: [
                    ...allIdes.map((ide) => ({
                        label: ide.name.trim(),
                        value: ide.id.trim(),
                        hint: ide?.desc,
                    })),
                    {
                        label: "No IDE",
                        value: "",
                        hint: "Do not set a default IDE",
                    },
                ],
            },
        ) as unknown as string;

        // Resolve directory path and IDE
        const resolvedPath = parsePath(dir);
        const ide = ideIdOrName ? findOneIDE(ideIdOrName) : null;

        // Add the project to the store
        addProject({
            name: name.trim(),
            dir: resolvedPath.trim(),
            relationIdeId: ide?.id.trim(),
        });

        consola.success(`Project '${name}' added successfully!`);
    })
    .addCommand(addIdeCommand);
