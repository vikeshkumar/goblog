package config

const (
	//ui configurations
	TemplateRoot       string = "template.root"
	Favicon            string = "site.favicon"
	Robots             string = "site.robots"
	TemplateSuffix     string = "template.suffix"
	StaticResourcePath string = "static.path"
	StaticPathPrefix   string = "static.prefix"
	AdminResourcePath  string = "admin.path"
	AdminPathPrefix    string = "admin.prefix"
	UploadResourcePath string = "uploads.path"
	UploadPathPrefix   string = "uploads.prefix"

	//server configurations
	ServerListen        string = "server.listenOn"
	ServerWriteTimeout  string = "server.writetimeout"
	ServerReadTimeout   string = "server.readtimeout"
	ServerAppUrl        string = "server.appUrl"
	ServerUploadPath    string = "server.uploadPath"
	ServerExecDir       string = "server.execPathDir"
	ServerDbUrl         string = "db.url"
	DBHost              string = "db.host"
	DBPort              string = "db.port"
	DBUser              string = "db.user"
	DBPassword          string = "db.password"
	DBName              string = "db.name"
	DBMaxConnection     string = "db.maxconnection"
	DBMinConnection     string = "db.minconnection"
	DBHealthCheckPeriod string = "db.healthcheckperiod"
	FileUploadDirectory string = "server.fileUploadDirectory"

	BcryptCost = "bcrypt.cost"

	ApiCookieName     = "api.cookie.name"
	ApiCookieDomain   = "api.cookie.domain"
	ApiCookieValidity = "api.cookie.validity"
	ApiCookieSecure   = "api.cookie.secure"
)
