package db

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "gorm.io/driver/postgres"
	"net.vikesh/goshop/config"
	"regexp"
	"strings"
)

var db *pgxpool.Pool
var cfg = config.Get()
var txOptions = pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable}
type TxQuery struct {
	Sql string
	Ctx context.Context
}

func Connect() error {
	cfg := config.Get()
	//parseConfig, parseError := pgx.ParseConfig("")
	//if parseError != nil {
	//	return parseError
	//}
	//parseConfig.Config.Host = cfg.GetString(cfg.GetString(config.DBHost))
	//parseConfig.Config.Port = uint16(cfg.GetInt(cfg.GetString(config.DBPort)))
	//parseConfig.Config.User = cfg.GetString(cfg.GetString(config.DBUser))
	//parseConfig.Config.Password = cfg.GetString(cfg.GetString(config.DBPassword))
	//parseConfig.Config.Database = cfg.GetString(cfg.GetString(config.DBName))
	//connectionConfig := &pgxpool.Config{
	//	MaxConns: int32(cfg.GetInt(config.DBMaxConnection)),
	//	MinConns: int32(cfg.GetInt(config.DBMinConnection)),
	//	HealthCheckPeriod: cfg.GetDuration(config.DBHealthCheckPeriod),
	//	LazyConnect: false,
	//	ConnConfig: parseConfig,
	//}
	//pgx.Conn
	conn, connectError := pgxpool.Connect(context.Background(), cfg.GetString(config.ServerDbUrl))
	db = conn
	return connectError
}


func assign(values []pgtype.Value, targets []interface{}) {
	size := len(values)
	for i, v := range values {
		if v != nil && size > i {
			_ = v.AssignTo(targets[i])
		}
	}
}

func present(value pgtype.Value) bool {
	return value != nil && value.Get() != pgtype.Undefined && value.Get() != pgtype.Null
}

func createUrlFromTitle(title string) string {
	compile := regexp.MustCompile("\\s+")
	asciiChars := regexp.MustCompile("[^a-zA-Z0-9\\s+]").ReplaceAllString(title, "")
	return strings.ToLower(compile.ReplaceAllString(asciiChars, "-"))
}

