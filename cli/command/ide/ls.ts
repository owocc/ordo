import { program } from "npm:commander";
import { findAllIDE } from "../../../store/ide.store.ts";
import Table from "npm:cli-table3";
import { blankBorder } from "../../../lib/config/tableBorder.ts";
import consola from "npm:consola";

/**
 * 列出所有已注册的 IDE 及其安装路径。
 */
export const lsIDECommand = program.createCommand('ls')
    .description('列出所有已注册的 IDE 及其对应的安装路径。')
    .action(() => {
        const allIDEs = findAllIDE();

        // 如果没有找到任何 IDE，则提示用户
        if (allIDEs.length === 0) {
            consola.info(
                "没有找到任何 IDE。请使用 'add' 命令注册一个 IDE。",
            );
            return;
        }

        // 将所有 IDE 的信息格式化为表格数据
        const data = allIDEs.map((ide) => [
            ide.id,
            ide.name,
            ide.path,
        ]);

        // 使用 Table 库创建并显示 IDE 列表
        const table = new Table({
            head: ["ID", "IDE 名称", "安装路径"],
            chars: blankBorder, // 自定义表格边框样式
        });

        data.forEach((row) => {
            table.push(row);
        });

        // 输出找到的 IDE 数量，并显示表格
        consola.success(`找到 ${allIDEs.length} 个 IDE：`);
        console.log(table.toString());
    });
