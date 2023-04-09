package mongodb

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type (
	Raw    bson.Raw                 // Raw is a raw sql that will not be treated as argument but as a direct sql part.
	Value  = *gvar.Var              // Value is the field value type.
	Record map[string]Value         // Record is the row record of the table.
	Result []Record                 // Result is the row record array.
	Map    = map[string]interface{} // Map is alias of map[string]interface{}, which is the most common usage map type.
	List   = []Map                  // List is type of map array.
)

type queryType int

const (
	defaultModelSafe                      = false
	defaultCharset                        = `utf8`
	defaultProtocol                       = `tcp`
	queryTypeNormal           queryType   = 0
	queryTypeCount            queryType   = 1
	queryTypeValue            queryType   = 2
	unionTypeNormal                       = 0
	unionTypeAll                          = 1
	defaultMaxIdleConnCount               = 10               // Max idle connection count in pool.
	defaultMaxOpenConnCount               = 0                // Max open connection count in pool. Default is no limit.
	defaultMaxConnLifeTime                = 30 * time.Second // Max lifetime for per connection in pool in seconds.
	ctxTimeoutTypeExec                    = 0
	ctxTimeoutTypeQuery                   = 1
	ctxTimeoutTypePrepare                 = 2
	cachePrefixTableFields                = `TableFields:`
	cachePrefixSelectCache                = `SelectCache:`
	commandEnvKeyForDryRun                = "gf.gdb.dryrun"
	modelForDaoSuffix                     = `ForDao`
	dbRoleSlave                           = `slave`
	ctxKeyForDB               gctx.StrKey = `CtxKeyForDB`
	ctxKeyCatchSQL            gctx.StrKey = `CtxKeyCatchSQL`
	ctxKeyInternalProducedSQL gctx.StrKey = `CtxKeyInternalProducedSQL`

	// type:[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	linkPattern = `(\w+):([\w\-]*):(.*?)@(\w+?)\((.+?)\)/{0,1}([^\?]*)\?{0,1}(.*)`
)

var (
	// instances is the management map for instances.
	instances = gmap.NewStrAnyMap(true)

	// lastOperatorRegPattern is the regular expression pattern for a string
	// which has operator at its tail.
	lastOperatorRegPattern = `[<>=]+\s*$`

	// regularFieldNameRegPattern is the regular expression pattern for a string
	// which is a regular field name of table.
	regularFieldNameRegPattern = `^[\w\.\-]+$`

	// regularFieldNameWithoutDotRegPattern is similar to regularFieldNameRegPattern but not allows '.'.
	// Note that, although some databases allow char '.' in the field name, but it here does not allow '.'
	// in the field name as it conflicts with "db.table.field" pattern in SOME situations.
	regularFieldNameWithoutDotRegPattern = `^[\w\-]+$`

	// allDryRun sets dry-run feature for all database connections.
	// It is commonly used for command options for convenience.
	allDryRun = false

	// tableFieldsMap caches the table information retrieved from database.
	tableFieldsMap = gmap.NewStrAnyMap(true)
)

// Core is the base struct for database management.
type Core struct {
	//db            DB              // DB interface object.
	ctx           context.Context // Context for chaining operation only. Do not set a default value in Core initialization.
	group         string          // Configuration group name.
	schema        string          // Custom schema for this object.
	debug         *gtype.Bool     // Enable debug mode for the database, which can be changed in runtime.
	cache         *gcache.Cache   // Cache manager, SQL result cache only.
	links         *gmap.StrAnyMap // links caches all created links by node.
	logger        glog.ILogger    // Logger for logging functionality.
	config        *ConfigNode     // Current config node.
	dynamicConfig dynamicConfig   // Dynamic configurations, which can be changed in runtime.
}

type dynamicConfig struct {
	MaxIdleConnCount int
	MaxOpenConnCount int
	MaxConnLifeTime  time.Duration
}
