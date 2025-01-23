import consola from "consola";
import { findAllIDE } from "../store/ide.store.ts";
import type { ProjectAndIDE } from "../types/store.ts";
import { isAbsolute } from "jsr:@std/path";

// 获取前缀命令，针对不同操作系统返回不同的前缀
function getPrefix() {
    if (Deno.build.os === 'darwin' || Deno.build.os === 'linux') {
        return 'open'
    }
    return ''
}

/** 获取所有项目的启动命令 */
export async function getCommands(projects: ProjectAndIDE[]) {
    const ides = findAllIDE();

    // 如果没有找到任何 IDE，返回错误并终止操作
    if (ides.length === 0) {
        consola.error(
            "未找到 IDE。请在继续之前配置至少一个 IDE。",
        );
        return [];
    }

    const commands: string[] = [];

    for (const project of projects) {
        let idePath = project?.ide?.path;

        // 如果项目没有默认 IDE，提示用户选择一个 IDE
        if (!project?.ide) {
            idePath = (await consola.prompt(
                "🤯 该项目没有配置默认 IDE，请选择一个 IDE:",
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

        // 如果 IDE 路径存在并且是绝对路径，构造 IDE 启动命令
        if (idePath && isAbsolute(idePath)) {
            idePath = [getPrefix(), `"${idePath}"`].join(' ')
        }

        // 组合 IDE 路径和项目目录，生成启动命令
        const command = [idePath, project.dir].join(' ')
        commands.push(command);
    }

    return commands;
}