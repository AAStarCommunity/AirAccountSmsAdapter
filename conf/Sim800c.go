package conf

var port string
var baud int

// GetSim800c 获取Sim800c通信方式
func GetSim800c() (string, int) {
	if len(port) == 0 || baud == 0 {
		port = Get().Sim800C.Port
		baud = Get().Sim800C.Baud
	}

	return port, baud
}
