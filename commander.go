package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"test/zhlog"
)

// LookPath在环境变量中查找科执行二进制文件，如果file中包含一个斜杠，则直接根据绝对路径或者相对本目录的相对路径去查找
func LookPath() {
	f, err := exec.LookPath("/root/zhxu")
	zhlog.Assert(err)
	fmt.Println("file:", f)

}

//command返回cmd结构来执行带有相关参数的命令,它仅仅设定cmd结构中的Path和Args参数
// 如果name参数中不包含路径分隔符，command使用LookPath来解决路径问题，否则的话就直接使用name
// Args直接跟在command命令之后，所以在Args中不许要添加命令．
func Cmmd() {
	cmd := exec.Command("tr", "a-z", "A-Z")
	cmd.Stdin = strings.NewReader("some inpuit")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	zhlog.Assert(err)
	fmt.Println("out string{}:", out.String())
}

//运行命令，并返回标准输出和标准错误
// 注意：Output()和CombinedOutput()不能够同时使用，因为command的标准输出只能有一个，同时使用的话便会定义了两个，便会报错
func CombinedOutput() {
	cmd := exec.Command("ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

// StdinPipe返回一个连接到command标准输入的管道pipe
func StdinPipe() {
	cmd := exec.Command("ls")
	stdin, err := cmd.StdinPipe()
	zhlog.Assert(err)
	_, err = stdin.Write([]byte("tep.txt"))
	zhlog.Assert(err)
	stdin.Close()
	cmd.Stdout = os.Stdout
	cmd.Start()
}

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

func cmder() {
	cmd := NewZhCommand()
	a, b, e := cmd.Exec(os.Getenv("SHELL"), "ls /data/mysql")
	fmt.Println("===:", a, b, e)
	// LookPath()
	Cmmd()
	CombinedOutput()
	StdinPipe()
}

type Zhcommander interface {
	// 执行命令并返回进程的pid，命令行结构，错误消息
	Exec(sysType string, args ...string) (int, string, error)

	// 异步执行命名并通过chan返回结果,error 用panic
	ExecAsync(sysType string, stdout chan string, args ...string) int

	// ExecIgnore(args ...string) error
}

func NewZhCommand() Zhcommander {
	var cmd Zhcommander

	switch runtime.GOOS {
	case "linux":
		cmd = Zhcmd()
	case "windows":
		// cmd = ZhWindows()
	case "mac":

	}
	return cmd
}

type LinuxCommander struct {
}

func Zhcmd() *LinuxCommander {
	return &LinuxCommander{}
}

func (e *LinuxCommander) Exec(sysType string, args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.Command(sysType, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	outpip, err := cmd.StdoutPipe()
	defer outpip.Close()
	if err != nil {
		return 0, "", err
	}

	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}

	out, err := ioutil.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}

	return cmd.Process.Pid, string(out), nil
}

func (e *LinuxCommander) ExecAsync(sysType string, stdout chan string, args ...string) int {
	var pidChan = make(chan int, 1)
	go func() {
		args = append([]string{"-c"}, args...)
		cmd := exec.Command(sysType, args...)

		outpip, err := cmd.StdoutPipe()
		zhlog.Assert(err)
		err = cmd.Start()
		zhlog.Assert(err)

		pidChan <- cmd.Process.Pid
		out, err := ioutil.ReadAll(outpip)
		zhlog.Assert(err)

		stdout <- string(out)
	}()
	return <-pidChan
}
