# go-version

#### 介绍

设置及显示版本信息的库。

#### 使用说明

示例：

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chen-zyc/go-version"
)

func main() {
	version.Version = "v1.0.0"
	version.Branch = "main"
	version.Hash = "abc123"
	version.BuildTime = time.Now().Format(time.RFC3339)
	version.Comment = "example"

	v, err := version.New("demo", "v0.99.99")
	checkErr(err)

	// 或者自定义模板。
	v.Template = version.DefaultTemplate2
	// 默认是 os.Getwd() 的值。
	v.WorkDir = "/tmp"

	s, err := v.String()
	checkErr(err)

	fmt.Println(s)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
```

输出示例：

```
         PROGRAM : demo
         VERSION : v1.0.0
INTERNAL_VERSION : v0.99.99
          BRANCH : main
            HASH : abc123
      BUILD_TIME : 2021-08-19T15:58:46+08:00
      GO_VERSION : go1.16.7
        WORK_DIR : /tmp
        RUN_TIME : 2021-08-19T15:58:46+08:00
         COMMENT : example
```

或者配合 Makefile 使用：

```makefile
VERPKG=github.com/chen-zyc/go-version
VERSION=v0.1.0
BRANCH=`git rev-parse --abbrev-ref HEAD`
HASH=`git rev-list -1 HEAD`
BUILDTIME=`date`
COMMENT=build from Makefile
LDFLAGS="-X '${VERPKG}.Version=${VERSION}' -X '${VERPKG}.Branch=${BRANCH}' -X '${VERPKG}.Hash=${HASH}' -X '${VERPKG}.BuildTime=${BUILDTIME}' -X '${VERPKG}.Comment=${COMMENT}'"

build:
	go build -v -ldflags ${LDFLAGS} example.com
```

```go
func printVersion() {
	ver, err := version.New("demo", "v0.99.99")
	checkError(err)

	s, err := ver.String()
	checkError(err)

	fmt.Println(s)
}
```