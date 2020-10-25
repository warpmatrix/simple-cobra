<!-- omit in toc -->
# Simple Cobra

simple cobra 实现了支持简单带子命令的命令行程序开发的功能。支持导入子命令目录，为开发的应用程序增加需要的自定义子命令。所有的子命令文件夹中，一个 go 文件对应应用程序中的一个子命令。可以通过增加或删除 go 文件，对直接对应用程序中的子命令进行改动。

<!-- omit in toc -->
## Table of Contents

- [1. Quick Start](#1-quick-start)
- [2. 使用示例](#2-使用示例)
  - [2.1. 环境创建](#21-环境创建)
  - [2.2. 根指令文件](#22-根指令文件)
  - [2.3. 子指令文件](#23-子指令文件)
  - [2.4. 执行指令](#24-执行指令)

## 1. Quick Start

使用下述的指令可以快速安装 simple cobra 包：

```bash
go get github.com/warpmatrix/simple-cobra
```

可以在对应的 `command.go` 查看 `Command` 结构的相应字段接口：

- `Use`、`Short`、`Long` 用于 help 中描述该指令的用法
- `Run`、`RunE` 执行指令时所执行的函数，可以根据字段选择有无返回值。两个字段填写其一即可，若两个字段都填写，则 `RunE` 有更高的优先级
- `helpFunc`、`helpCommand` 为应用程序默认添加的 `help` 指令和执行 `help` 指令默认调用的函数
- `parent`、`commands` 记录指令的父指令和子指令信息

使用该包利用上述接口，完成对应指令所需接口的实现，即可支持带子命令的命令行程序开发。通常通过以下的两种方法使用对应接口：

- 在子命令的 `init` 函数调用 `AddCommand` 方法可以为指令添加子命令
- 通常通过根指令调用 `Execute` 方法，可以执行的输入指令

## 2. 使用示例

下面提供一个简单的使用示例，创建一个名为 newApp 的项目。这个项目包含的命令如下：

- 根命令包含子命令 server
- server 子命令中包含 config 子命令

### 2.1. 环境创建

首先使用如下指令，创建项目环境：

```bash
mkdir newApp
cd newApp
mkdir cmd
```

接着创建对应的指令文件，并填写相应的代码，配置指令行为：

```bash
cd cmd
touch root.go server.go config.go
cd ..
```

### 2.2. 根指令文件

`root.go` 文件中记录的是应用程序根指令的信息，调用根命令的 `Execute` 方法即可执行用户输入的指令。具体代码填写的示例如下：指令填写 `Use`、`Short`、`Long` 字段用于 help 输出中。`Run` 和 `RunE` 字段表示执行指令所调用的函数，若都不填写则输出 help 的信息。

```go
package cmd

import (
    "fmt"
    "os"

    cobra "github.com/warpmatrix/simple-cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "test [server [config]]",
    Short: "A brief description of your application",
    Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    // Run: func(cmd *cobra.Command, args []string) { fmt.Println("main run") },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
//
// Execute 函数执行涉及所有子命令的指令，该函数通过 main 函数只被调用一次
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

### 2.3. 子指令文件

为指令添加子指令，只需要增加对应的子指令文件，在 `init` 函数中调用父指令的 `AddCommand` 方法即可实现。

下面代码块中的内容为 `server.go` 文件的代码：

```go
package cmd

import (
    "fmt"

    cobra "github.com/warpmatrix/simple-cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("server called, args:", args)
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)

    // Here you will define your flags and configuration settings.
}
```

下面的代码块中的内容为 `config.go` 文件的代码：

```go
package cmd

import (
    "fmt"

    cobra "github.com/warpmatrix/simple-cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("config called, args:", args)
    },
}

func init() {
    serverCmd.AddCommand(configCmd)

    // Here you will define your flags and configuration settings.
}
```

### 2.4. 执行指令

完成上述的所有步骤，即实现了项目的命令需求：

- 根命令包含子命令 server
- server 子命令中包含 config 子命令

我们可以编写 `main.go` 执行用户输入的指令：

```bash
touch main.go
```

给 `main.go` 文件填充以下内容：

```go
package main

// import "domain-name/your-id/your-repo/cmd"
import "github.com/warpmatrix/newApp/cmd"

func main() {
    cmd.Execute()
}
```

编译、执行的结果如下，可以看到在没有定义执行函数时，执行指令输出对应的 help 信息：

```bash
go build
./newApp
```

![run-newApp](images/run-newApp.png)

并且可以实现 help 输出子命令信息以及子命令的执行：

```bash
./newApp help server
./newApp server config arg1 arg2
```

![run-sub-command](images/run-subcmd.png)
