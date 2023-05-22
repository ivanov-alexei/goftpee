package server

const (
	ResponseTemplate = "%d %s"
)

const (
	// 2xx
	StatusCommandOk    		= 200
	StatusServiceReady 		= 220
	StatusCloseControlConn 	= 221
	StatusPassiveMode  		= 227
	StatusUserLoggedIn 		= 230
	StatusFileActionOk 		= 250

	// 3xx
	StatusUserOk = 331

	// 4xx

	// 5xx
	StatusSyntaxError           = 500
	StatusCommandNotImplemented = 502
	StatusNotLoggedIn           = 530
)
