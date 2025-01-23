import { deleteProject } from "../../../store/project.store.ts";
import { program } from "npm:commander";
import consola from "npm:consola";

export const rmCommand = program.createCommand("rm")
    .description("从管理器中删除一个项目")
    .argument("[mark]", "要删除的项目标识符（标记）")
    .action(async (mark?: string) => {
        // 如果没有提供标记，提示用户输入
        mark ||= await consola.prompt("请输入要删除的项目标记:", {
            type: "text",
        });

        // 验证标记
        if (!mark) {
            consola.error("项目标记是必需的！");
            return;
        }

        try {
            deleteProject(mark);
            consola.success(`项目标记为 "${mark}" 的项目已成功删除。`);
        } catch (error) {
            consola.error(`删除项目失败: ${error}`);
        }
    });