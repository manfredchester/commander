package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"github.com/urfave/cli/v2"
)

// type Cmd struct {
// 	Path         string　　　//运行命令的路径，绝对路径或者相对路径
// 	Args         []string　　 // 命令参数
// 	Env          []string         //进程环境，如果环境为空，则使用当前进程的环境
// 	Dir          string　　　//指定command的工作目录，如果dir为空，则comman在调用进程所在当前目录中运行
// 	Stdin        io.Reader　　//标准输入，如果stdin是nil的话，进程从null device中读取（os.DevNull），stdin也可以时一个文件，否则的话则在运行过程中再开一个goroutine去
// 　　　　　　　　　　　　　//读取标准输入
// 	Stdout       io.Writer       //标准输出
// 	Stderr       io.Writer　　//错误输出，如果这两个（Stdout和Stderr）为空的话，则command运行时将响应的文件描述符连接到os.DevNull
// 	ExtraFiles   []*os.File
// 	SysProcAttr  *syscall.SysProcAttr
// 	Process      *os.Process    //Process是底层进程，只启动一次
// 	ProcessState *os.ProcessState　　//ProcessState包含一个退出进程的信息，当进程调用Wait或者Run时便会产生该信息．
// }
func main() {
	cmd := exec.Command("ls")
	fmt.Println(cmd)
}

func IfSoftwareExt(name string) (string, error) {
	cmd1 := exec.Command("/usr/bin/rpm", "-qa")
	cmd2 := exec.Command("grep", name)
	var outbuf1 bytes.Buffer
	cmd1.Stdout = &outbuf1
	if err := cmd1.Start(); err != nil {
		return "", err
	}
	if err := cmd1.Wait(); err != nil {
		return "", err
	}
	var outbuf2 bytes.Buffer
	cmd2.Stdin = &outbuf1
	cmd2.Stdout = &outbuf2
	if err := cmd2.Start(); err != nil {
		return "", err
	}
	if err := cmd2.Wait(); err != nil {
		return "", err
	}
	return outbuf2.String(), nil
}

var StopShCmd = &cli.Command{
	Name:  "stopsh",
	Usage: "run sh stop.sh",
	Action: func(context *cli.Context) error {
		logger.Info("command stopsh")

		if !util.HasExist("/data/shepherd/stop.sh") {
			util.CreatFile("/data/shepherd/stop.sh", getStopShConent())
		}
		err2 := os.Remove("/data/shepherd/agent.lock")
		if util.IsNotEmpty(err2) {
			logger.Error(err2)
		}
		// command2 := exec.Command("/usr/bin/rm", "-f", "/data/shepherd/agent.lock")
		// err2 := command2.Start()
		// if util.IsNotEmpty(err2) {
		// 	logger.Error(err2)
		// 	return err2
		// }
		command := exec.Command("/bin/bash", "/data/shepherd/stop.sh")
		//command.Start()
		err := command.Run()
		if util.IsNotEmpty(err) {
			logger.Error(err)
			return err
		}
		return nil
	},
}

func getStopShConent() string {
	buff := bytes.Buffer{}
	buff.WriteString("#!/bin/bash\n")
	buff.WriteString("rm -f /data/shepherd/agent.lock\n")
	buff.WriteString("proc=`ps -ef | grep shepherd | grep -v grep | awk '{print $2}'`\n")
	buff.WriteString("if [ -n \"$proc\" ]; then\n")
	buff.WriteString("  echo \"proc=$proc \"1\n")
	buff.WriteString("  kill -9  $proc\n")
	buff.WriteString("fi\n")
	return buff.String()
}