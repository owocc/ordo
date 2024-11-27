import consola from "consola";
import { findAllIDE } from "../store/ide.store.ts";
import type { ProjectAndIDE } from "../types/store.ts";
import { isAbsolute } from "jsr:@std/path";

/** Get all project start command */
export async function getCommands(projects: ProjectAndIDE[]) {
    const ides = findAllIDE();
    if (ides.length === 0) {
        consola.error(
            "No IDEs found. Please configure at least one IDE before proceeding.",
        );
        return [];
    }
    const commands: string[] = [];

    for (const project of projects) {
        let idePath = project?.ide?.path;
        if (!project?.ide) {
            idePath = (await consola.prompt(
                "🤯 This project does not have a default IDE configured. Please select an IDE:",
                {
                    type: "select",
                    options: ides.map((e) => ({
                        label: e.name,
                        value: e.path,
                        hint: e?.desc,
                    })),
                },
            )) as unknown as string;
        }
        if (idePath && isAbsolute(idePath)) {
            idePath = `"${idePath}"`;
        }
        const command = `${idePath} ${project.dir}`;
        commands.push(command);
    }

    return commands;
}
