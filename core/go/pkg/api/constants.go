package api

const (
	ProtocolVersion     = "1.0.0"
	MethodHello         = "device.hello"
	MethodGetInfo       = "device.get_info"
	MethodSetApConfig   = "wifi.set_ap_config"
	MethodGetClients    = "wifi.get_clients"
	MethodForwardStart  = "forward.start"
	MethodForwardStop   = "forward.stop"
	MethodGetStats      = "forward.get_stats"
	MethodSubscribeLog  = "log.subscribe"
	MethodPing          = "diag.ping"
	MethodOtaBegin      = "ota.begin"
	MethodOtaChunk      = "ota.chunk"
	MethodOtaCommit     = "ota.commit"
)

const (
	ChannelControl = 1
	ChannelData    = 2
	ChannelOta     = 3
)

const (
	ErrorOK                  = 0
	ErrorCodeProtocolBase    = 1000
	ErrorCodeDeviceBase      = 2000
	ErrorCodeRuntimeBase     = 3000
	ErrorCodeOtaBase         = 4000
	ErrorCodeInternalBase    = 5000
)
