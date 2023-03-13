package mysql

type Config struct {
	Url        string `json:"Url"`
	User       string `json:"User"`
	Password   string `json:"Password"`
	Dbname     string `json:"Dbname"`
	Charset    string `json:"Charset"`
	MaxIdle    int    `json:"MaxIdle"`
	MaxConn    int    `json:"MaxConn"`
	ShowSql    bool   `json:"ShowSql"`
	SqlLogFile string `json:"SqlLogFile"`
	LogLv      int    `json:"LogLv"`
}
