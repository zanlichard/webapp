// common key
const (
       STAT_IN			        = "MsgIn"
       STAT_OUT				= "MsgOut"
       INDEGREE_Recive 		        = "InDegree_Recive(MB)"
       INDEGREE_Send			= "InDegree_Send(MB)"
       OUTDEGREE_Recive 		= "OutDegree_Recive(MB)"
       OUTDEGREE_Send			= "OutDegree_Send(MB)"
	   
)

// application key
const (
	StatRedisGet				= "StatRedisGet"
)
	
// base initialize
	logconfig := make(stat.LoggerParm)
	logconfig.level = "info"
	logconfig.path = "./stat"
	logconfig.namePrefix = "test"
	logconfig.filename = "stat.log"
	logconfig.maxfilesize = 10000
	logconfig.maxdays = 7
	logconfig.maxlines = 10000
	logconfig.chanlen = 10000
	stat.Init(logconfig, 60)
	stat.StatProc()
	
// application initialize 
	stat.SetDelayUp(20,50,100)

	
// using method add stat data
    stat.Push(elem)

// exit must call
    stat.Exit()







