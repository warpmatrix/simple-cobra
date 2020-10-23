<!-- omit in toc -->
# Specification

<!-- omit in toc -->
## Table of Contents

- [`cobra` 的简单介绍](#cobra-的简单介绍)
- [`simple cobra` 的设计说明](#simple-cobra-的设计说明)
- [`simple cobra` 的单元测试](#simple-cobra-的单元测试)
- [`simple cobra` 的功能测试](#simple-cobra-的功能测试)

## `cobra` 的简单介绍

`cobra` 可以生成一个简单的带子命令的命令行程序，使用 `go get -u github.com/spf13/cobra` 指令可以安装 `cobra` 库，安装过程可能需要手动通过 github 镜像安装对应的依赖库。使用 `go get github.com/spf13/cobra/cobra` 可以安装 `cobra` 应用程序。

完成上述两个安装步骤后，可以使用 `cobra` 指令，快速创建新的应用程序：

```bash
mkdir -p newApp && cd newApp
cobra init --pkg-name github.com/userid/newApp
```

执行初始化指令后，新的应用程序工作目录结构如下：

```plaintext
.
|---cmd
|   `root.go
`---main.go
```

在 `cmd` 文件夹中，存放了应用程序的各条子命令。可以使用 `cobra add` 指令给应用程序增加新的子命令，如执行 `cobra add config`。此时，文件结构变为：

```plaintext
.
|---cmd
|---|---config.go
|   `---root.go
`---main.go
```

此时，可以执行 `go run main.go config` 指令，得到输出的结果为：`config called`。`config` 子命令成功被加到了新的应用程序中。

## `simple cobra` 的设计说明

`simple cobra` 在 `cobra` 的基础上进行了简化，去除了一些用处不大的功能，实现了最基础的用例需求——支持简单带子命令的命令行程序开发。

<!-- TODO: 设计说明 -->

## `simple cobra` 的单元测试

<!-- TODO: 单元测试 -->

## `simple cobra` 的功能测试

<!-- TODO: 功能测试 -->
