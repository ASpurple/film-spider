package app

// Addr 播放地址
type Addr struct {
	ID     int
	URL    string
	Remark string
	FID    int
}

// Film 电影信息
type Film struct {
	ID       int
	Name     string
	Director string
	Actor    string
	Type     string
	Area     string
	Lang     string
	Year     string
	Summary  string
	Img      string
	Site     string
	Amd      string
	Stat     int
	Addrs    []Addr
}

// SiteConf 站点配置
type SiteConf struct {
	Host         string
	Name         string
	ListURL      string
	DetailURL    string
	HrefSep      string
	NameSep1     string
	NameSep2     string
	DirectorSep1 string
	DirectorSep2 string
	ActorSep1    string
	ActorSep2    string
	TypeSep1     string
	TypeSep2     string
	AreaSep1     string
	AreaSep2     string
	LangSep1     string
	LangSep2     string
	YearSep1     string
	YearSep2     string
	SummarySep1  string
	SummarySep2  string
	ImgSep1      string
	ImgSep2      string
	AddrSep1     string
	AddrSep2     string
}
