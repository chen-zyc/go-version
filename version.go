package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"text/template"
	"time"
)

// 编译时传入具体值
var (
	Version   string
	Branch    string
	Hash      string
	BuildTime string
	Comment   string
)

const DefaultTemplate = `
{{.ProgramName}} ({{.Comment}})
	VERSION          : {{.Version}}
	INTERNAL_VERSION : {{.InternalVersion}}
	BRANCH           : {{.Branch}}
	HASH             : {{.Hash}}
	BUILD_TIME       : {{.BuildTime}}
	GO_VERSION       : {{.GoVersion}}
	WORK_DIR         : {{.WorkDir}}
	RUN_TIME         : {{.RunTime}}
`

const DefaultTemplate2 = `
{{.ProgramName}} ({{.Comment}})
	         VERSION : {{.Version}}
	INTERNAL_VERSION : {{.InternalVersion}}
	          BRANCH : {{.Branch}}
	            HASH : {{.Hash}}
	      BUILD_TIME : {{.BuildTime}}
	      GO_VERSION : {{.GoVersion}}
	        WORK_DIR : {{.WorkDir}}
	        RUN_TIME : {{.RunTime}}
`

// V 保存程序版本及相关的值。
type V struct {
	// 编译时传递的值
	Version   string `json:"version"`    // 程序版本
	Branch    string `json:"branch"`     // git 分支
	Hash      string `json:"hash"`       // git 哈希值
	BuildTime string `json:"build_time"` // 构建时间
	Comment   string `json:"comment"`    // 注释

	// 运行时的值
	ProgramName     string `json:"program_name"`     // 程序名称
	InternalVersion string `json:"internal_version"` // 程序内部的版本号
	GoVersion       string `json:"go_version"`       // go 版本
	WorkDir         string `json:"work_dir"`         // 工作目录
	Template        string `json:"-"`                // 渲染模板，如果为空，则使用默认模板。
	RunTime         string `json:"run_time"`         // 运行时间
}

func New(programName, internalVersion string) (v *V, err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return
	}
	v = &V{
		Version:         Version,
		Branch:          Branch,
		Hash:            Hash,
		BuildTime:       BuildTime,
		Comment:         Comment,
		ProgramName:     programName,
		InternalVersion: internalVersion,
		GoVersion:       runtime.Version(),
		WorkDir:         workDir,
		RunTime:         time.Now().Format(time.RFC3339),
	}
	return
}

func (v V) String() (string, error) {
	tpl := v.Template
	if tpl == "" {
		tpl = DefaultTemplate
	}
	t, err := template.New(v.ProgramName + "-version").Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := strings.Builder{}
	if err = t.Execute(&buf, v); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (v V) JSON() ([]byte, error) {
	return json.Marshal(v)
}

func (v V) IndentJSON() ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}

// MustPrintVersionAndExit 向标准输出中输出版本信息并退出。
func MustPrintVersionAndExit(name, internalVersion string) {
	v, err := New(name, internalVersion)
	if err != nil {
		panic(err)
	}
	s, err := v.String()
	if err != nil {
		panic(err)
	}
	fmt.Print(s)
	os.Exit(0)
}

// Handler 以 json 的形式输出版本信息.
func Handler(name, internalVersion string) http.HandlerFunc {
	v, err := New(name, internalVersion)
	if err != nil {
		panic(err)
	}
	jsonData, err := v.JSON()
	if err != nil {
		panic(err)
	}
	indentJSON, err := v.IndentJSON()
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		var err error
		if req.URL.Query().Get("pretty") != "" {
			_, err = w.Write(indentJSON)
		} else {
			_, err = w.Write(jsonData)
		}
		if err != nil {
			fmt.Println("Failed to write:", err)
		}
	}
}
