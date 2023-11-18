package conf

var host string

// GetAirCenterHost 获取AirCenter服务地址
func GetAirCenterHost() string {
	if len(host) == 0 {
		host = Get().AirCenter.Host
	}

	return host
}
