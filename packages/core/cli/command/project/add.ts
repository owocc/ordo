import { program } from "npm:commander";
import { addProject } from "../../../store/project.store.ts";
import { parsePath } from "../../../utils/path.ts";
import consola from "consola";
import { findAllIDE, findOneIDE } from "../../../store/ide.store.ts";

/**
 * 向项目管理器添加一个新项目
 */
export const addCommand = program.createCommand("add")
    .description(
        "将一个项目添加到管理器，指定项目名称、目录，并可以选择一个可选的 IDE。",
    )
    .argument("[name]", "项目名称")
    .argument("[dir]", "项目目录")
    .argument("[ide]", "IDE 名称或 ID")
    .action(async (name?: string, dir?: string, ideIdOrName?: string) => {
        // 如果未提供名称，提示用户输入
        name ||= await consola.prompt("请输入项目名称：", {
            type: "text",
        });
        if (!name) {
            consola.error("项目名称是必填项！");
            return;
        }

        // 如果未提供目录，提示用户输入
        dir ||= await consola.prompt("请输入项目目录：", {
            type: "text",
        });
        if (!dir) {
            consola.error("项目目录是必填项！");
            return;
        }

        // 如果没有提供 IDE，提示用户选择一个 IDE
        const allIdes = findAllIDE();
        ideIdOrName ||= await consola.prompt(
            "请选择一个 IDE（或者跳过）：",
            {
                type: "select",
                options: [
                    ...allIdes.map((ide) => ({
                        label: ide.name.trim(),
                        value: ide.id.trim(),
                        hint: ide?.desc,
                    })),
                    {
                        label: "无 IDE",
                        value: "",
                        hint: "不设置默认的 IDE",
                    },
                ],
            },
        ) as unknown as string;

        // 解析目录路径并找到 IDE（如果提供）
        const resolvedPath = parsePath(dir);
        const ide = ideIdOrName ? findOneIDE(ideIdOrName) : null;

        // 确保解析的路径和 IDE 是有效的
        if (!resolvedPath) {
            consola.error("无效的项目目录路径！");
            return;
        }

        // 如果提供的 IDE 不存在，提示用户选择的 IDE 无效
        if (ideIdOrName && !ide) {
            consola.error(`没有找到标识符为 "${ideIdOrName}" 的 IDE。`);
            return;
        }

        // 将项目添加到 store 中
        addProject({
            name: name.trim(),
            dir: resolvedPath.trim(),
            relationIdeId: ide?.id.trim(),
        });

        consola.success(`项目 '${name}' 已成功添加！`);
    });