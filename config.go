package railgun

var Config = conf{}

type conf struct {
	Parser struct {
		Proxy struct {
			Enabled bool
			CN      string
			HK      string
			TW      string
			TH      string
		}

		DefaultRegion Region
	}
}

func (c *conf) setProxy(t Region, addr string) {
	switch t {
	case REGION_CN:
		c.Parser.Proxy.CN = addr
	case REGION_HK:
		c.Parser.Proxy.HK = addr
	case REGION_TW:
		c.Parser.Proxy.TW = addr
	case REGION_TH:
		c.Parser.Proxy.TH = addr
	}
}

type ProxyHostname string

func (c *conf) getProxy(t Region) ProxyHostname {
	switch t {
	case REGION_CN:
		return ProxyHostname(c.Parser.Proxy.CN)
	case REGION_HK:
		return ProxyHostname(c.Parser.Proxy.HK)
	case REGION_TW:
		return ProxyHostname(c.Parser.Proxy.TW)
	case REGION_TH:
		return ProxyHostname(c.Parser.Proxy.TH)
	default:
		return ""
	}
}

func (c *conf) getDefaultRegion() Region {
	return c.Parser.DefaultRegion
}

func (p ProxyHostname) GetProxyURL() string {
	if p == `` {
		return ``
	}
	return stringBuilder("https://", string(p))
}
