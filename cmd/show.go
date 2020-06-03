package cmd

import (
	"bufio"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	COMMAND_DIR        = ".slc"                                          //文件夹名称
	COMMAND_PATH       = filepath.Join(getHomePath(), COMMAND_DIR)       //文件保存地址
	REMOTE_COMMAND_URL = "https://unpkg.com/linux-command/command/%s.md" //linux命令下载地址
)

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <command>",
		Short: "Show the specified command usage.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("[sorry]: the show command does not accept any arguments")
				return
			}
			showCmd(args[0])
		},
	}
	return cmd
}

func showCmd(cmd string) {
	cmd = strings.ToLower(cmd)
	p := path.Join(getHomePath(), COMMAND_DIR, fmt.Sprintf("%s.md", cmd))
	if !checkFileExist(p) {
		downloadCmdFile(cmd)
	}

	source, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Println("[show command '" + cmd + "' error]")
		return
	}
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	result := markdown.Render(string(source), 80, 6)
	fmt.Println(string(result))
}

func downloadCmdFile(cmd string) (error, int) {
	//检查创建文件夹
	if err := makeCmdDir(); err != nil {
		return err, 0
	}

	//请求远程资源
	resp, err := http.Get(fmt.Sprintf(REMOTE_COMMAND_URL, cmd))
	if err != nil {
		return err, 0
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, resp.StatusCode
	}

	defer resp.Body.Close()

	content := make([]byte, 0)
	reader := bufio.NewReader(resp.Body)
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return err, 0
		}
		if err == io.EOF {
			break
		}
		content = append(content, line...)
		content = append(content, []byte("\n")...)
	}
	p := path.Join(COMMAND_PATH, fmt.Sprintf("%s.md", cmd))
	return ioutil.WriteFile(p, content, 0666), 0
}

//获取用户目录
func getHomePath() string {
	home, _ := homedir.Expand("~")
	return home
}

//判断文件是否存在
func checkFileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func makeCmdDir() error {
	if _, err := os.Lstat(COMMAND_PATH); err != nil && !os.IsExist(err) {
		return os.Mkdir(COMMAND_PATH, 0755)
	}
	return nil
}
