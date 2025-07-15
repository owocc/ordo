import { addIDE } from "../../../store/ide.store.ts";
import { program } from "npm:commander";
import consola from "consola";

/**
 * 添加一个新的 IDE 到管理器。
 */
export const addIdeCommand = program.createCommand("add")
    .argument("[name]", "IDE 名称") // IDE 的名称
    .argument("[pathOrCommand]", "IDE 路径或命令") // IDE 的路径或命令
    .description("将一个新的 IDE 添加到管理器中。")
    .action(async (name?: string, pathOrCommand?: string) => {
        // 如果未提供 IDE 名称，提示用户输入
        name ||= await consola.prompt("请输入 IDE 名称:", {
            type: "text",
            placeholder: "例如: VSCode",
        });
        if (!name) {
            consola.error("IDE 名称是必填的！");
            return;
        }

        // 如果未提供 IDE 路径或命令，提示用户输入
        pathOrCommand ||= await consola.prompt(
            "请输入 IDE 路径或命令:",
            {
                type: "text",
                placeholder: "例如: /usr/bin/code 或 idea",
            },
        );
        if (!pathOrCommand) {
            consola.error("IDE 路径或命令是必填的！");
            return;
        }

        // 调用添加 IDE 的函数
        addIDE({
            name,
            path: pathOrCommand,
        });

        // 显示成功消息
        consola.success(`IDE '${name}' 已成功添加！`);
    });