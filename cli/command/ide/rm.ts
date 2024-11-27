import { deleteIDE, findOneIDE } from "../../../store/ide.store.ts";
import { program } from "npm:commander";
import consola from "consola";

/**
 * 删除一个 IDE 通过其标识符（mark）从管理器中。
 */
export const rmIdeCommand = program
    .createCommand("rm")
    .description("通过标识符（mark）从管理器中删除一个 IDE")
    .argument("[mark]", "要删除的 IDE 的标识符（mark）")
    .action(async (mark?: string) => {
        try {
            // 如果没有提供 mark，则提示用户输入
            mark ||= await consola.prompt("请输入要删除的 IDE 的标识符（mark）：", {
                type: "text",
            });

            // 验证 mark 是否有效
            if (!mark) {
                consola.error("IDE 标识符（mark）是必填项！");
                return;
            }

            // 检查指定标识符是否存在
            const ide = findOneIDE(mark);
            if (!ide) {
                consola.error(`没有找到标识符为 "${mark}" 的 IDE。`);
                return;
            }

            // 执行删除操作
            deleteIDE(mark);
            consola.success(`标识符为 "${mark}" 的 IDE 已成功删除。`);
        } catch (error) {
            consola.error(`删除 IDE 失败：${error}`);
        }
    });