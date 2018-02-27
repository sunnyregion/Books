package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type books struct {
	name  string
	size  string
	path  string
	bytes int64
}

const (
	_          = iota
	KB float64 = 1 << (10 * iota)
	MB
	GB
)

func fileSize(i int64) string {
	b := float64(i)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fBytes", b)
}

func walkDir(dirPth, suffix string) (files []books, content string, err error) {
	files = make([]books, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		// 忽略目录
		if fi.IsDir() {
			switch len(strings.Split(filename, "\\")) {
			case 1:
				content += "## "
			case 2:
				content += "* "
			case 3:
				content += "  * "
			}
			content += fi.Name() + "\r\n"
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			input := books{
				name:  fi.Name(),
				size:  fileSize(fi.Size()),
				path:  strings.Replace(filename, "\\", "/", -1),
				bytes: fi.Size(),
			}
			switch len(strings.Split(filename, "\\")) {
			case 3:
				content += "  * "
			case 4:
				content += "    * "
			}
			//[静态资源](/静态资源/README.md)
			content += fmt.Sprintf("[%s [%s]](./%s)", strings.Replace(strings.Split(input.name, ".")[0], "_", " ", -1), input.size, input.path) + "\r\n"
			files = append(files, input)
		}
		return nil
	})
	return files, content, err
}

func writeMarkDown(fileName, content string) {
	// open output file
	fo, err := os.Create(fileName + ".md")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// make a write buffer
	w := bufio.NewWriter(fo)
	w.WriteString(content)
	w.Flush()
}

func main() {
	var count int
	var sum int64
	var content string
	content = `# Books

把读过的一些书分享出来给大家，涉及新思想、新科技、算法、人工智能、语言编程类等等，每本书都是严格挑选，这个库也会持续不断的更新。

如果您有好书也不妨提交一个issue，共享出来。这些资源来自于互联网共享于互联网（from the Internet, for the Internet）；如果涉嫌侵权，您也可以提交一个issue告知，我会及时删除。在此，对这些书的作者（或译者）们表示感谢！
	
如果您觉得这个资源还不错，不妨给我来杯咖啡，我会把赞助者的名单列在下面以示感谢！

![微信支付](./weixin.jpg)

`
	if f, c, e := walkDir("./图书目录", ""); e == nil {
		count = len(f)
		for _, v := range f {
			sum += v.bytes
		}
		content += c
	}
	content += "\r\n" + fmt.Sprintln("共", count, "本书，计", fileSize(sum))
	content += "\r\n最后更新时间：" + time.Now().Format("2006年1月2日 15:04:05") + "\r\n"
	content += `
## 鸣谢

项目赞助者：
`
	writeMarkDown("README", content)
}
