# 使用golang开发selpg命令行程序


----------


根据[selpg](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)设计编写

参考博客[服务计算——selpg命令行程序][1]


----------


引用包

    import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
	flag "github.com/spf13/pflag"
)

结构体

    type selpg_args struct {
	startPage  int
	endPage    int
	inFile string
	pageLen    int
	pageType   bool
	printDest string
}

函数

    func main()//主程序
    func process_args()//检查参数
    func process_input()//处理参数并输出

  [1]: https://blog.csdn.net/qq_26003929/article/details/78262002