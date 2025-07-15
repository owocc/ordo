import { findManyProjectAndIDE } from "../../../store/project.store.ts";
import { program } from "npm:commander";
import consola from "npm:consola";
import Table from "npm:cli-table3";
import { blankBorder } from "../../../lib/config/tableBorder.ts";

/**
 * 列出所有已管理的项目及其详细信息。
 */
export const lsCommand = program.createCommand("ls")
    .description("列出所有已管理的项目及其详细信息。")
    .action(() => {
        // 获取所有项目和其对应的 IDE 信息
        const projects = findManyProjectAndIDE();

        // 如果没有找到任何项目，提示用户
        if (projects.length === 0) {
            consola.info(
                "没有找到任何项目。请使用 'add' 命令添加一个项目。",
            );
            return;
        }

        // 将项目的数据进行格式化，包含项目 ID、名称、目录和默认 IDE
        const data = projects.map((project) => [
            project.id,
            project.name,
            project.dir,
            project?.ide?.name || "None", // 如果没有关联的 IDE，则显示 "None"
        ]);

        // 使用 Table 库创建一个表格来展示项目列表
        const table = new Table({
            head: ["ID", "项目名称", "项目目录", "默认 IDE"], // 表头内容
            chars: blankBorder, // 自定义表格边框样式
        });

        // 将项目数据填充到表格中
        data.forEach((row) => {
            table.push(row);
        });

        // 输出找到的项目数量，并显示表格
        consola.success(`找到 ${projects.length} 个项目：`);
        console.log(table.toString());
    });