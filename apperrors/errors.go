package apperrors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"
)

type (
	// AppError is the error description of application
	AppError struct {
		Err      error   `json:"err"`       // error return by function
		ErrPoint string  `json:"err_point"` // error point stores context of function call, contains file, line and function
		Code     RetCode `json:"-"`         // error code
		Msg      string  `json:"-"`         // error description
	}

	// CallStack stores information about the active subroutines
	CallStack struct {
		ErrRet    *AppError    `json:"-"`         // error return by application
		Stack     []*CallStack `json:"stack"`     // tree structure of CallStack
		ProcTime  int64        `json:"proc_time"` // how long a process takes
		Tag       string       `json:"tag"`       // category of process
		BeginTime int64        `json:"-"`         // when this process began
		parent    *CallStack   // parent node
		desc      string       // description
	}

	// StatInfo stores information of statistic
	StatInfo struct {
		ProcTime int64  `json:"procTime"`
		Ret      string `json:"ret"`
	}

	// CallStackWrapper is the wrapper of CallStack, use to generate error messages
	CallStackWrapper struct {
		Cs             *CallStack          `json:"callstack"`
		Code           RetCode             `json:"code"`
		Msg            string              `json:"msg"`
		Stat           map[string]StatInfo `json:"stat"`
		ErrPath        string              `json:"err_path"`
		ErrPathArr     []string            `json:"err_path_arr"`
		TimeoutPath    string              `json:"timeout_path"`
		TimeoutPathArr []string            `json:"timeout_path_arr"`
	}
)

const (
	SUCCESS     = 0
	ERR_UNKNOWN = -1
)

var (
	tags = []string{"mysql", "redis"}

	errorMap = map[RetCode]string{
		SUCCESS:     "Success",
		ERR_UNKNOWN: "Unknown Error",
	}
)

//应用名,必须初始化
var (
	strAppName string
)

func Init(appName string) {
	strAppName = appName
}

// RegisterTags register tags array, use for prefix matching when generate error message
func RegisterTags(tag []string) {
	tags = append(tags, tag...)
	// sort, for "Maximum matching algorithm"
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))
}

// RegisterError register error code and corresponding error message
func RegisterError(errMap map[RetCode]string) {
	for code, msg := range errMap {
		if _, ok := errorMap[code]; !ok {
			errorMap[code] = msg
		}
	}
}

// BeginCallStack creates a new call stack when a process began
func BeginCallStack(tag string) *CallStack {
	s := new(CallStack)
	s.Tag = tag
	s.BeginTime = GetMSTimeStamp()
	return s
}

// PushBackCallStack push a call stack when a process step into a subroutine
func (s *CallStack) PushBackCallStack(cs *CallStack) {
	if s == cs {
		return
	}

	s.Stack = append(s.Stack, cs)
	if cs.ErrRet != nil {
		s.ErrRet = &AppError{Err: cs.ErrRet.Err, ErrPoint: GetErrPoint(2), Code: cs.ErrRet.Code, Msg: cs.ErrRet.Msg}
	}
	cs.parent = s
}

// PushBackCallStackIgnoreErr push a call stack and ignore recording error
func (s *CallStack) PushBackCallStackIgnoreErr(cs *CallStack) {
	if s == cs {
		return
	}

	s.Stack = append(s.Stack, cs)
	cs.parent = s
}

func shortFileFuncName(fileName, funcName string) (string, string) {
	items := strings.Split(fileName, fmt.Sprintf("%s/", strAppName))
	shortFileName := items[len(items)-1]

	items = strings.Split(funcName, ".")
	shortFuncName := items[len(items)-1]

	return shortFileName, shortFuncName
}

// EndCall ends off a process, calculates time consuming and records description
func (s *CallStack) EndCall(skip int) *CallStack {
	pc, file, _, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()

	shortFileName, shortFuncName := shortFileFuncName(file, funcName)

	s.ProcTime = GetMSTimeStamp() - s.BeginTime
	s.desc = fmt.Sprintf("%s:%s.%s:%d", shortFileName, shortFuncName, s.Tag, s.ProcTime)
	return s
}

// FormatJson encodes call stack to json
func (s *CallStack) FormatJson() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// GetStat samples the time consuming of subroutinecess which has prefix of tags
func (s *CallStack) GetStat(mpStat map[string]StatInfo) {
	if s.Stack != nil {
		for _, st := range s.Stack {
			st.GetStat(mpStat)
		}
	}

	for _, tag := range tags {
		if strings.HasPrefix(s.Tag, tag) {
			st := StatInfo{ProcTime: s.ProcTime}
			if s.ErrRet != nil && s.ErrRet.Err != nil {
				st.Ret = s.ErrRet.Err.Error()
			}
			mpStat[tag] = st
		}
	}
}

// GetErrPath traverses the call stack and returns the error path
func (s *CallStack) GetErrPath(arr []string) string {
	var buffer bytes.Buffer
	if s.ErrRet == nil || s.ErrRet.Err == nil {
		buffer.WriteString("")
	} else {
		if s.parent != nil {
			buffer.WriteString("~")
		}
		path := fmt.Sprintf("%s:%d", s.ErrRet.ErrPoint, s.ProcTime)
		buffer.WriteString(path)
		arr = append(arr, path)
		for _, st := range s.Stack {
			if st.ErrRet != nil {
				buffer.WriteString(st.GetErrPath(arr))
			}
		}
		if s.parent == nil {
			buffer.WriteString("#")
			buffer.WriteString(s.ErrRet.Err.Error())
		}

	}

	return buffer.String()
}

// GetTimeoutPath returns path which is timeout
func (s *CallStack) GetTimeoutPath(arr []string, funcTimeout int64) string {
	if s.ProcTime <= funcTimeout {
		return ""
	}

	var buffer bytes.Buffer
	if s.parent != nil {
		buffer.WriteString("~")
	}
	buffer.WriteString(s.desc)
	arr = append(arr, buffer.String())
	for _, st := range s.Stack {
		buffer.WriteString(st.GetTimeoutPath(arr, funcTimeout))
	}
	return buffer.String()

}

// SetCode sets error code and the corresponding error message
func (err *AppError) SetCode(code RetCode) {
	err.Code = code
	msg, ok := errorMap[code]
	if !ok {
		err.Msg = errorMap[ERR_UNKNOWN]
	} else {
		err.Msg = msg
	}
}

// GetCodeMsg returns the code message
func (err *AppError) GetCodeMsg() (RetCode, string) {
	return err.Code, err.Msg
}

// GetErrPoint returns the error context, file:line:function
func GetErrPoint(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()

	items := strings.Split(file, fmt.Sprintf("%s/", strAppName))
	shortFileName := items[len(items)-1]

	items = strings.Split(funcName, ".")
	shortFuncName := items[len(items)-1]

	s := fmt.Sprintf("%s:%d:%s", shortFileName, line, shortFuncName)
	return s
}

func (err *AppError) Error() string {
	ret := fmt.Sprintf("%s#%s", err.ErrPoint, err.Err.Error())
	return ret
}

// GetProcMsg returns messages of call stack
func (cs *CallStack) GetProcMsg(funcTimeout int64) string {
	if cs == nil {
		return ""
	}

	csw := new(CallStackWrapper)
	csw.ErrPathArr = make([]string, 0)
	csw.TimeoutPathArr = make([]string, 0)
	csw.Stat = make(map[string]StatInfo)
	csw.Cs = cs

	if cs.ErrRet == nil {
		cs.ErrRet = &AppError{Code: SUCCESS, Msg: errorMap[SUCCESS]}
	}
	csw.Code = cs.ErrRet.Code
	csw.Msg = cs.ErrRet.Msg

	cs.GetStat(csw.Stat)
	csw.ErrPath = cs.GetErrPath(csw.ErrPathArr)
	csw.TimeoutPath = cs.GetTimeoutPath(csw.TimeoutPathArr, funcTimeout)

	b, _ := json.Marshal(csw)

	return string(b)
}

func GetMSTimeStamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func CheckError(cs *CallStack) bool {
	if cs.ErrRet != nil && (cs.ErrRet.Err != nil || cs.ErrRet.Code != 0) {
		return false
	}
	return true
}
