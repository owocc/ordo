import { findManyProjectAndIDE } from "../../store/project.store.ts";
import { program } from "npm:commander";
import consola from "npm:consola";
import Table from "npm:cli-table3";

export const lsCommand = program.createCommand("ls")
    .description("List all managed projects with their details.")
    .action(() => {
        const projects = findManyProjectAndIDE();

        if (projects.length === 0) {
            consola.info(
                "No projects found. Use the 'add' command to add a project.",
            );
            return;
        }

        const data = projects.map((project) => [
            project.id,
            project.name,
            project.dir,
            project?.ide?.name || "None",
        ]);

        const table = new Table({
            head: ["ID", "Name", "Project Directory", "Default Editor"],
        });

        data.forEach((row) => {
            table.push(row);
        });

        consola.success(`Found ${projects.length} project(s):`);
        console.log(table.toString());
    });
