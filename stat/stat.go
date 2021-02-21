package stat

// Inter 统计接口
type Inter interface {
	AddReportHeadItem(itemName string)
	AddReportBodyRowItem(itemName string)
	AddReportBodyColItem(itemName string)
	AddReportTailItem(itemName string)
	AddReportErrorItem(itemName string)
	AddReportIPError()

	IncStat(itemName string, val uint)
	IncKey(itemName string)
	IncStatByTab(rowName string, colName string, val uint)
	IncErrnoStat(errno int, val uint)
	IncErrnoStatByItem(itemName string, errno int, val uint)
	IncErrnoIP(ip uint, errno int, val uint)

	SetStat(itemName string, val uint)
	GetStat(itemName string) int
	GetStatValueByTab(itemName string, colName string)
	TimeStatGet(rowName string) (count uint, avgDelay float32, maxDelay float32, upDelay uint, upDelay2 uint, upDelay3 uint)

	NoCheckAndPrint()
	Print()
	PrintHeader()
	PrintBody()
	PrintTail()
	PrintRowError()
	PrintIPError()

	ClearAll()

	Reset()
}
