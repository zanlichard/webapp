package stat

// statex
// 内部提供最终的输出控制,输出的控制不交给外部完成
// 为了避免使用double buffer机制,数据接收和重置放到协程中统一事件处理
// 定义外部与内部协程之间的通信数据的格式
// 通信协议包括 错误码、延时、成功与失败
import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/zanlichard/beegoe/logs"
)

const (
	errNum       = 5
	ipNum        = 3
	strFmt       = "%-13s"
	ipFmt        = "%-17s"
	perFmt       = "%-10f%%"
	miniFloatFmt = "%-5f"
	miniIntFmt   = "%-4d"
	miniStrFmt   = "%-4s"
	spt          = " | "

	statDelayEndCount  = "tcount"
	statDelayTotalTime = "de_total_ms"
	statDelayMax       = "de_max_s"
	statMaxIP          = "max_ip"
	statDelayUp        = "de_up"
	statDelayUp2       = "de_up_2"
	statDelayUp3       = "de_up_3"
	statDelayUp4       = "de_up_4" // 超时前一个分段

	statIn          = "MsgIn"
	statOut         = "MsgOut"
	indegreeRecive  = "INDEGREE_Recive(B)"
	indegreeSend    = "INDEGREE_Send(B)"
	outdegreeRecive = "OUTDEGREE_Recive(B)"
	outdegreeSend   = "OUTDEGREE_Send(B)"
)

// LoggerParam 日志参数
type LoggerParam struct {
	Level       string
	Path        string
	NamePrefix  string
	Filename    string
	Maxfilesize int
	Maxdays     int
	Maxlines    int
	Chanlen     int64
}

// Item 统一单元元素
type Item struct {
	Name      string // 统计的接口名
	Delay     uint   // 接口执行的延时,单位ms
	Errcode   int    // 当次接口请求的错误码,0--成功`
	Ipsrc     net.IP // 请求的来源ip
	Payload   uint   // 请求的载荷
	Direction int    // 上行or下行    1 ---- 上行   0 ----下行
	InOrOut   int    // 入度还是出度   1 ---- in  0 ----out
}

// Mystat 监控统计对象
type Mystat struct {
	sampleChan    chan Item // 监控统计channel
	ctrlChan      chan int  // 推出控制
	vHeadItems    *list.List
	vBodyRowItems *list.List
	vBodyColItems *list.List
	vTailItems    *list.List
	statlog       *logs.BeeLogger
	timeout       uint
	statGap       time.Duration
	delayUp       uint
	delayUp2      uint
	delayUp3      uint
	mapErrNum     map[string]map[int]uint //request,errnode,count
	mapErrIP      map[int]map[int]uint    //retcode,ip,count
	isHadIPErr    bool
	countMap      map[string]uint // 请求每段的几数 requests,count
	errnoCountMap map[int]uint    // errno,count
	timeoutMap    map[string]uint

	IsClearFlag bool // 清除标记
}

// GStat 全局唯一统计对象
var GStat *Mystat

// Init 初始化调用,线程不安全
func Init(logconfig LoggerParam, statgap time.Duration) {
	Logger := logs.NewLogger(logconfig.Chanlen)
	logConfig := fmt.Sprintf(`{
				"filename":"%s/%s_%s",
				"maxlines":%d,
				"maxsize":%d,
				"maxDays":%d,
				"blankprefix":true,
				"perm": "664",
				"rotateperm":"444"
		  }`,
		logconfig.Path,
		logconfig.NamePrefix,
		logconfig.Filename,
		logconfig.Maxlines,
		logconfig.Maxfilesize,
		logconfig.Maxdays)

	var level int
	switch logconfig.Level {
	case "debug":
		level = logs.LevelDebug
	case "info":
		level = logs.LevelInformational
	case "notice":
		level = logs.LevelNotice
	case "warn":
		level = logs.LevelWarning
	case "error":
		level = logs.LevelError
	case "critical":
		level = logs.LevelCritical
	case "alert":
		level = logs.LevelAlert
	case "emergency":
		level = logs.LevelEmergency
	}

	// 初始化输出日志
	Logger.SetLogger("file", logConfig)
	Logger.SetLevel(level)
	Logger.BlankPrefix()
	Logger.Async()

	GStat = new(Mystat)
	GStat.statGap = statgap
	GStat.statlog = Logger
	GStat.vHeadItems = list.New()
	GStat.vBodyRowItems = list.New()
	GStat.vBodyColItems = list.New()
	GStat.vTailItems = list.New()
	GStat.countMap = make(map[string]uint)
	GStat.errnoCountMap = make(map[int]uint)
	GStat.mapErrIP = make(map[int]map[int]uint)
	GStat.mapErrNum = make(map[string]map[int]uint)
	GStat.timeoutMap = make(map[string]uint)

	GStat.sampleChan = make(chan Item, 4096) // 1024 有点小, 高并发有可能阻塞
	GStat.ctrlChan = make(chan int)
	GStat.AddReportHeadItem(statIn)
	GStat.AddReportHeadItem(statOut)
	GStat.AddReportTailItem(indegreeRecive)
	GStat.AddReportTailItem(indegreeSend)
	GStat.IsClearFlag = false
}

// PushStat 把函数统计输入统计
func PushStat(itemName string, procTime int, requestIP net.IP, payload int, retcode int) {
	statItem := new(Item)
	statItem.Name = itemName
	statItem.Delay = uint(procTime)
	statItem.Errcode = retcode
	statItem.Ipsrc = requestIP
	statItem.Payload = uint(payload)
	statItem.Direction = 1
	statItem.InOrOut = 1
	for i := 0; i < 5; i++ {
		if GStat.IsClearFlag {
			<-time.After(200 * time.Millisecond)
		} else {
			break
		}
	}
	GStat.sampleChan <- *statItem
}

// SetDelayUp 设置分段统计延时
// 初始化调用,线程不安全
func SetDelayUp(delayUp uint, delayUp2 uint, delayUp3 uint) {
	if delayUp >= delayUp2 || delayUp2 >= delayUp3 {
		panic(errors.New("delay up error"))
	}
	GStat.delayUp = delayUp
	GStat.delayUp2 = delayUp2
	GStat.delayUp3 = delayUp3
}

// SetTimeOut 设置客户端请求超时时间
// 初始化调用,线程不安全
// 和delayup参数意义相同, 起一个分段作用,只不过是最大分段
func (stat *Mystat) SetTimeOut(timeout uint) {
	if timeout <= stat.delayUp3 {
		panic(errors.New("timeout error"))
	}
	stat.timeout = timeout
}

// Proc 统计统一处理函数
func Proc() {
	go func() {
		t1 := time.NewTimer(time.Second * GStat.statGap)
		for {
			select {
			case <-t1.C:
				GStat.NoCheckAndPrint()
				t1.Reset(time.Second * GStat.statGap)

			case <-GStat.ctrlChan:
				//退出
				return

			case elem := <-GStat.sampleChan:
				GStat.addSampleStat(&elem)
			}
		}
	}()

}

// Exit 推出统计
func Exit() {
	GStat.ctrlChan <- 0
	close(GStat.sampleChan)
	close(GStat.ctrlChan)
}

func (stat *Mystat) addSampleStat(elem *Item) {
	stat.IncKey(elem.Name + statDelayEndCount)             // 延时的总请求数
	stat.IncStat(elem.Name+statDelayTotalTime, elem.Delay) // 延时总和
	if elem.Direction == 1 {                               // 上行 (现在注意统计上行)
		stat.IncStat(statIn, 1)
	} else {
		stat.IncStat(statOut, 1)
	}

	if elem.InOrOut == 1 {
		if elem.Direction == 1 {
			stat.IncStat(indegreeRecive, elem.Payload)
		} else {
			stat.IncStat(indegreeSend, elem.Payload)
		}
	} else {
		if elem.Direction == 1 {
			stat.IncStat(outdegreeRecive, elem.Payload)
		} else {
			stat.IncStat(outdegreeSend, elem.Payload)
		}
	}

	// 修改最大延时
	max := stat.GetStat(elem.Name + statDelayMax)
	if elem.Delay > max {
		stat.SetStat(elem.Name+statDelayMax, elem.Delay)               //最大延时
		stat.SetStat(elem.Name+statMaxIP, uint(inet_aton(elem.Ipsrc))) //最大延时的ip
	}

	//分段延时统计
	if elem.Delay < stat.delayUp {
		stat.IncKey(elem.Name + statDelayUp)
	} else if elem.Delay >= stat.delayUp && elem.Delay < stat.delayUp2 {
		stat.IncKey(elem.Name + statDelayUp2)
	} else if elem.Delay >= stat.delayUp2 && elem.Delay < stat.delayUp3 {
		stat.IncKey(elem.Name + statDelayUp3)
	} else if elem.Delay >= stat.delayUp3 && elem.Delay < stat.timeout {
		stat.IncKey(elem.Name + statDelayUp4)
	} else {
		stat.IncTimeout(elem.Name)
	}

	stat.IncStat(elem.Name, 1) //请求数累计
	if elem.Errcode != 0 {
		stat.IncErrnoIP(elem.Ipsrc, elem.Errcode, 1)        // 错误码和ip统计
		stat.IncErrnoStatByItem(elem.Name, elem.Errcode, 1) // 请求相关的错误码统计
	}

}

// AddReportHeadItem 添加监控报告头部元素
func (stat *Mystat) AddReportHeadItem(itemName string) {
	stat.vHeadItems.PushBack(itemName)
}

// AddReportBodyRowItem 添加监控报告主体行元素
func (stat *Mystat) AddReportBodyRowItem(itemName string) {
	stat.vBodyRowItems.PushBack(itemName)
}

// AddReportBodyColItem 添加监控报告主体列元素
func (stat *Mystat) AddReportBodyColItem(itemName string) {
	stat.vBodyColItems.PushBack(itemName)
}

// AddReportTailItem 添加监控报告尾部元素
func (stat *Mystat) AddReportTailItem(itemName string) {
	stat.vTailItems.PushBack(itemName)
}

// AddReportErrorItem 添加监控报告错误元素
func (stat *Mystat) AddReportErrorItem(itemName string) {
	delete(stat.mapErrNum, itemName)
	stat.mapErrNum[itemName] = make(map[int]uint)
}

// AddReportIPError 添加监控报告IP地址错误
func (stat *Mystat) AddReportIPError() {
	stat.isHadIPErr = true
}

// IncTimeout 增加单位元素超时计数
func (stat *Mystat) IncTimeout(itemName string) {
	count, ok := stat.timeoutMap[itemName]
	if !ok {
		stat.timeoutMap[itemName] = 1
	} else {
		stat.timeoutMap[itemName] = count + 1
	}

}

// GetTimeout 获取单位元素超时计数
func (stat *Mystat) GetTimeout(itemName string) uint {
	val, ok := stat.timeoutMap[itemName]
	if !ok {
		return 0
	}
	return val
}

// SetStat 设置主体数据单位元素列统计值
func (stat *Mystat) SetStat(itemName string, val uint) {
	stat.countMap[itemName] = val
}

// IncKey 增加主体数据单位元素统计
func (stat *Mystat) IncKey(itemName string) {
	stat.IncStat(itemName, 1)
}

// IncStatByTab 增加主体数据单位元素列统计
func (stat *Mystat) IncStatByTab(rowName string, colName string, val uint) {
	stat.IncStat(rowName+colName, val)
}

// IncStat 增加主体数据单位元素列统计
func (stat *Mystat) IncStat(itemName string, val uint) {
	count, ok := stat.countMap[itemName]
	if !ok {
		stat.countMap[itemName] = val
	} else {
		stat.countMap[itemName] = count + val
	}
}

// IncErrnoStat 增加错误计数
func (stat *Mystat) IncErrnoStat(errno int, val uint) {
	count, ok := stat.errnoCountMap[errno]
	if !ok {
		stat.errnoCountMap[errno] = val
	} else {
		stat.errnoCountMap[errno] = count + val
	}
}

// IncErrnoStatByItem 增加单位元素错误计数
func (stat *Mystat) IncErrnoStatByItem(itemName string, errno int, val uint) {
	stat.IncErrnoStat(errno, val) //错误码统计
	errcodeMap, ok := stat.mapErrNum[itemName]
	if !ok {
		return
	}
	//到底是引用还是值
	count, ok2 := errcodeMap[errno]
	if !ok2 {
		errcodeMap[errno] = val
	} else {
		errcodeMap[errno] = val + count
	}
}

// IncErrnoIP 增加错误IP计数
func (stat *Mystat) IncErrnoIP(ip net.IP, errno int, val uint) {
	ipMap, ok := stat.mapErrIP[errno]
	ipint := inet_aton(ip)
	if !ok {
		ipMap := make(map[int]uint)
		ipMap[ipint] = val
		stat.mapErrIP[errno] = ipMap
	} else {
		count, ok := ipMap[ipint]
		if !ok {
			ipMap[ipint] = val
		} else {
			ipMap[ipint] = val + count
		}
	}
}

// GetStat 获取单位元素计数
func (stat *Mystat) GetStat(itemName string) uint {
	val, ok := stat.countMap[itemName]
	if !ok {
		return 0
	}
	return val
}

// TimeStatGet 单位时间内的计数获取
func (stat *Mystat) TimeStatGet(rowName string) (count uint, avgDelay, maxDelay float32, upDelay, upDelay2, upDelay3, upDelay4 uint) {
	ok := false
	iMaxDelay := uint(0)
	count, ok = stat.countMap[rowName+statDelayEndCount]
	if !ok {
		// fmt.Printf("no key:%s\n", rowName+statDelayEndCount)
	}
	delayTotalTime := uint(0)
	delayTotalTime, ok = stat.countMap[rowName+statDelayTotalTime]
	if ok {
		avgDelay = float32(delayTotalTime) / float32(count)
		// fmt.Printf("no key:%s\n", rowName+statDelayTotalTime)
	}
	iMaxDelay, ok = stat.countMap[rowName+statDelayMax]
	maxDelay = float32(iMaxDelay)
	if !ok {
		// fmt.Printf("no key:%s\n", rowName+statDelayMax)
	}
	upDelay, ok = stat.countMap[rowName+statDelayUp]
	if !ok {
		// fmt.Printf("no key:%s\n", rowName+statDelayUp2)
	}
	upDelay2, ok = stat.countMap[rowName+statDelayUp2]
	if !ok {
		// fmt.Printf("no key:%s\n", rowName+statDelayUp2)
	}
	upDelay3, ok = stat.countMap[rowName+statDelayUp3]
	if !ok {
		// fmt.Printf("no key:%s\n", rowName+statDelayUp3)
	}
	upDelay4, ok = stat.countMap[rowName+statDelayUp4]
	if !ok {
		// fmt.Printf("no key:%s\n", rowName+statDelayUp4)
	}
	return
}

// GetStatValueByTab 获取列元素的值
func (stat *Mystat) GetStatValueByTab(itemName string, colName string) uint {
	count, ok := stat.countMap[itemName+colName]
	if !ok {
		return 0
	}
	return count
}

// Print 打印
func (stat *Mystat) Print() {
	stat.PrintHeader()
	stat.PrintBody()
	stat.PrintRowError()
	stat.PrintIPError()
	stat.PrintTail()
	stat.statlog.Info("")
	return

}

// NoCheckAndPrint 清空缓存并打印
func (stat *Mystat) NoCheckAndPrint() {
	stat.Print()
	stat.Reset()
}

// PrintHeader 打印头部
func (stat *Mystat) PrintHeader() {
	t2 := time.Now()
	stat.statlog.Info("Statistic in %ds,  CTime: %s", stat.statGap, t2.Format("2006-01-02 15:04:05"))
	stat.statlog.Info("\n---------------------\nHead Information\n---------------------")

	line := fmt.Sprintf("%18s", "")
	line1 := fmt.Sprintf("%-18s", "total:")
	line2 := fmt.Sprintf("%-18s", "count /1s:")
	name := ""
	value := uint(0)
	for e := stat.vHeadItems.Front(); e != nil; e = e.Next() {
		name = e.Value.(string)
		value = stat.GetStat(name)
		line = fmt.Sprintf("%s|%9s", line, name)
		line1 = fmt.Sprintf("%s|%9d", line1, value)
		line2 = fmt.Sprintf("%s|%9d", line2, value/uint(stat.statGap))
	}
	stat.statlog.Info("%s", line)
	stat.statlog.Info("%s", line1)
	stat.statlog.Info("%s", line2)

}

// PrintBody 打印主体数据
func (stat *Mystat) PrintBody() {
	stat.statlog.Info("\n---------------------\nOperation Information\n---------------------")

	line := fmt.Sprintf("%-18s", "Op")

	for e := stat.vBodyColItems.Front(); e != nil; e = e.Next() {
		line = fmt.Sprintf("%s|%8s ", line, e.Value)
	}

	line = fmt.Sprintf("%s|%9s|%9s|%9s|%15s|%3d(ms)<|%3d(ms)<|%3d(ms)<|%3d(ms)<|%3d(ms)>|", line, "tcount", "avg_de_ms", "de_max_ms", "max_ip",
		stat.delayUp,
		stat.delayUp2,
		stat.delayUp3,
		stat.timeout, stat.timeout)

	stat.statlog.Info("%s", line)
	for eRow := stat.vBodyRowItems.Front(); eRow != nil; eRow = eRow.Next() {
		rowTotalCount := uint(0)
		name := eRow.Value.(string)
		line = fmt.Sprintf("%-18s", name+":")
		colname := ""
		value := uint(0)
		for eCol := stat.vBodyColItems.Front(); eCol != nil; eCol = eCol.Next() {
			colname = eCol.Value.(string)
			value = stat.GetStatValueByTab(name, colname)
			rowTotalCount += value
			line = fmt.Sprintf("%s|%9d ", line, value)
		}

		tcount, avg, max, up, up2, up3, up4 := stat.TimeStatGet(name)
		maxIP := stat.GetStat(name + statMaxIP)
		if rowTotalCount == 0 && tcount == 0 {
			continue
		}

		line = fmt.Sprintf("%s|%9d|%9.3f|%9.3f|%15s|%8d|%8d|%8d|%8d",
			line,
			tcount,
			avg,
			max,
			inet_ntoa(int(maxIP)).String(),
			up,
			up2,
			up3, up4)

		line = fmt.Sprintf("%s|%8d|", line, stat.GetTimeout(name))

		stat.statlog.Info("%s", line)
	}
}

// PrintTail 打印尾部
func (stat *Mystat) PrintTail() {
	if stat.vTailItems.Len() > 0 {
		stat.statlog.Info("\n---------------------\nTail Information\n---------------------")
		name := ""
		for e := stat.vTailItems.Front(); e != nil; e = e.Next() {
			name = e.Value.(string)
			stat.statlog.Info("%-17s | %8d", name+"#", stat.GetStat(name))
		}
	}
}

// PrintRowError 打印行错误
func (stat *Mystat) PrintRowError() {
	stat.statlog.Info("\n---------------------\nError Information\n---------------------")
	str := fmt.Sprintf("%-17s", "Op")
	format1 := ""
	format2 := ""
	for i := 0; i < errNum; i++ {
		format1 = fmt.Sprintf("%s%d", "Err", i+1)
		format2 = fmt.Sprintf(strFmt, format1)
		str += spt
		str += format2
	}

	format1 = fmt.Sprintf("%s", "total count")
	format2 = fmt.Sprintf(strFmt, format1)
	str += spt
	str += format2
	stat.statlog.Info("%s", str)

	allcount := 0
	count := 0
	for k, v := range stat.mapErrNum {
		count = 0
		if len(v) == 0 {
			continue
		}
		str = ""
		format1 = fmt.Sprintf("%-17s", k+"_E")
		str += format1
		topkArray := GetTopn(v, errNum)
		for _, item := range topkArray {
			format1 = fmt.Sprintf("%d/%d", item, v[item])
			format2 = fmt.Sprintf(strFmt, format1)
			count += int(v[item])
			str += spt
			str += format2
		}
		for i := len(topkArray); i < errNum; i++ {
			format1 = fmt.Sprintf("%d/%d", 0, 0)
			format2 = fmt.Sprintf(strFmt, format1)
			str += spt
			str += format2
		}

		allcount += count
		format1 = fmt.Sprintf("%d", count)
		format2 = fmt.Sprintf(strFmt, format1)
		str += spt
		str += format2
		stat.statlog.Info("%s", str)

	}
	stat.statlog.Info("---------------------")
	str = fmt.Sprintf("%-17s", "TOTAL")
	topnErrnoArray := GetTopn(stat.errnoCountMap, errNum)
	for _, val := range topnErrnoArray {
		format1 = fmt.Sprintf("%d/%d", val, stat.errnoCountMap[val])
		format2 = fmt.Sprintf(strFmt, format1)
		str += spt
		str += format2
	}

	for i := len(topnErrnoArray); i < errNum; i++ {
		format1 = fmt.Sprintf("%d/%d", 0, 0)
		format2 = fmt.Sprintf(strFmt, format1)
		str += spt
		str += format2
	}
	format1 = fmt.Sprintf("%d", allcount)
	format2 = fmt.Sprintf(strFmt, format1)
	str += spt
	str += format2
	stat.statlog.Info("%s", str)
}

// PrintIPError 打印IP错误
func (stat *Mystat) PrintIPError() {
	if !stat.isHadIPErr {
		return
	}
	stat.statlog.Info("\n---------------------\nIP Information\n---------------------")
	str := ""
	format1 := fmt.Sprintf("%-17s", "retcode")
	format2 := ""
	str += format1

	for i := 0; i < ipNum; i++ {
		format1 = fmt.Sprintf("%s%d", "ip", i+1)
		format2 = fmt.Sprintf(ipFmt, format1)
		str += spt
		str += format2
	}
	stat.statlog.Info("%s", str)
	for retcode, v := range stat.mapErrIP {
		count := 0
		if len(v) == 0 {
			continue
		}
		format1 = fmt.Sprintf("%-17d", retcode)
		str = ""
		str += format1
		topkArrayIP := GetTopn(v, ipNum)
		for _, ip := range topkArrayIP {
			format1 = fmt.Sprintf("%s/%d", inet_ntoa(ip).String(), v[ip])
			format2 = fmt.Sprintf(ipFmt, format1)
			count += int(v[ip])
			str += spt
			str += format2

		}
		for i := len(topkArrayIP); i < ipNum; i++ {
			format1 = fmt.Sprintf("%d/%d", 0, 0)
			format2 = fmt.Sprintf(ipFmt, format1)
			str += spt
			str += format2
		}
		stat.statlog.Info("%s", str)

	}

}

// ClearAll 清空计数
func (stat *Mystat) ClearAll() {
	stat.Reset()
}

// Reset 重置计数
func (stat *Mystat) Reset() {
	stat.IsClearFlag = true
	wg := &sync.WaitGroup{}
	wg.Add(5)
	go func(iwg *sync.WaitGroup) {
		for k := range stat.timeoutMap {
			delete(stat.timeoutMap, k)
		}
		iwg.Done()
	}(wg)

	go func(iwg *sync.WaitGroup) {
		for k := range stat.countMap {
			delete(stat.countMap, k)
		}
		iwg.Done()
	}(wg)

	go func(iwg *sync.WaitGroup) {
		for k := range stat.errnoCountMap {
			delete(stat.errnoCountMap, k)
		}
		iwg.Done()
	}(wg)

	go func(iwg *sync.WaitGroup) {
		for _, v := range stat.mapErrNum {
			for k1 := range v {
				delete(v, k1)
			}
			//			delete(stat.mapErrNum, k)
		}
		iwg.Done()
	}(wg)

	go func(iwg *sync.WaitGroup) {
		for k, v := range stat.mapErrIP {
			for k1 := range v {
				delete(v, k1)
			}
			delete(stat.mapErrIP, k)
		}
		iwg.Done()
	}(wg)

	wg.Wait()
	stat.IsClearFlag = false
}
