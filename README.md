# ORDO

Ordo 是一个简单高效的命令行工具，用于管理开发项目和集成常用 IDE,并且支持跨平台使用!!!😆

## 功能概览

📂 项目管理

    •	添加、打开、列出、移除项目
    •	支持按名称或 ID 操作

🛠️ IDE 管理

    •	添加、列出、移除 IDE
    •	兼容自定义命令或路径

## 使用说明

📂 项目管理

#### 添加项目

将项目添加到 Ordo 的管理列表：

```bash
ordo add [projectName] [projectPath] [defaultIDE]
```

#### 打开项目

通过项目名称或 ID 打开项目：

```bash
ordo open [projectName | projectId]
```

#### 显示所有项目

列出当前已添加的所有项目：

```bash
ordo ls
```

#### 移除项目

通过名称或 ID 从管理列表中删除项目：

```bash
ordo rm [projectName | projectId]
```

🛠️ IDE 管理

#### 添加 IDE

将新的 IDE 添加到 Ordo 的支持列表：

```bash
ordo ide add [IDE Name] [IDE Path | IDE Command]
```

#### 显示所有 IDE

查看已添加的所有 IDE 列表：

```bash
ordo ide ls
```

#### 移除 IDE

通过名称或 ID 删除指定的 IDE：

```bash
ordo ide rm [IDE Name | IDE Id]
```

主要库和工具：

    •	Commander.js：命令行界面构建工具
    •	Consola：增强型日志工具
    •	cli-table3：格式化表格输出

备注:

    •	Ordo 支持任何自定义 IDE，只需通过路径或命令注册即可。
    •	如果遇到问题，Ordo 会以友好的方式提示你修复错误。

Ordo 的作用是为了简化项目管理,让项目列表更加扁平，让开发更高效，工具更便捷！ 😊
