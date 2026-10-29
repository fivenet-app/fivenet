package properties

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	utilsCache "github.com/fivenet-app/fivenet/v2026/pkg/utils/cache"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/proto"
)

var ErrSystemPropertyNotFound = errors.New("system property not found")

const (
	definitionCacheTTL      = 5 * time.Minute
	definitionCacheCapacity = 256
)

// DefinitionCache caches effective property metadata. It is safe for
// concurrent use and should be provided as a singleton by Fx so all stores
// share the same cache.
type DefinitionCache struct {
	definitions *utilsCache.LRUCache[string, map[string]*propertypb.PropertyDefinition]
	groups      *utilsCache.LRUCache[string, []*propertypb.PropertyGroup]
	enabled     *utilsCache.LRUCache[string, bool]
	loads       singleflight.Group

	mu         sync.RWMutex
	generation map[propertypb.PropertyScope]uint64
}

type ResolvedCatalog struct {
	Definitions map[string]*propertypb.PropertyDefinition
	Groups      []*propertypb.PropertyGroup
}

func NewDefinitionCache() *DefinitionCache {
	return &DefinitionCache{
		definitions: utilsCache.NewLRUCache[string, map[string]*propertypb.PropertyDefinition](
			definitionCacheCapacity,
		),
		groups: utilsCache.NewLRUCache[string, []*propertypb.PropertyGroup](
			definitionCacheCapacity,
		),
		enabled:    utilsCache.NewLRUCache[string, bool](definitionCacheCapacity),
		generation: make(map[propertypb.PropertyScope]uint64),
	}
}

func (c *DefinitionCache) Resolve(
	ctx context.Context,
	db *sql.DB,
	scope propertypb.PropertyScope,
	keys []string,
) (map[string]*propertypb.PropertyDefinition, error) {
	cacheKey := scopeAndKeys(scope, keys)
	if value, ok := c.definitions.Get(cacheKey); ok {
		return cloneDefinitions(value), nil
	}

	loaded, err, _ := c.loads.Do("definitions:"+cacheKey, func() (any, error) {
		if value, ok := c.definitions.Get(cacheKey); ok {
			return value, nil
		}
		generation := c.currentGeneration(scope)
		value, err := ResolveDefinitions(ctx, db, scope, keys)
		if err != nil {
			return nil, err
		}
		if c.currentGeneration(scope) == generation {
			c.definitions.Put(cacheKey, value, definitionCacheTTL)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	return cloneDefinitions(loaded.(map[string]*propertypb.PropertyDefinition)), nil
}

func (c *DefinitionCache) Groups(
	ctx context.Context,
	db *sql.DB,
	scope propertypb.PropertyScope,
	definitions map[string]*propertypb.PropertyDefinition,
) ([]*propertypb.PropertyGroup, error) {
	cacheKey := definitionGroupKey(scope, definitions)
	if value, ok := c.groups.Get(cacheKey); ok {
		return cloneGroups(value), nil
	}

	loaded, err, _ := c.loads.Do("groups:"+cacheKey, func() (any, error) {
		if value, ok := c.groups.Get(cacheKey); ok {
			return value, nil
		}
		generation := c.currentGeneration(scope)
		value, err := ResolvePropertyGroupsFromDB(ctx, db, scope, definitions)
		if err != nil {
			return nil, err
		}
		if c.currentGeneration(scope) == generation {
			c.groups.Put(cacheKey, value, definitionCacheTTL)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	return cloneGroups(loaded.([]*propertypb.PropertyGroup)), nil
}

// ResolveCatalog returns definition and group metadata from one cache service
// so callers do not need to coordinate the two resolution steps themselves.
func (c *DefinitionCache) ResolveCatalog(
	ctx context.Context,
	db *sql.DB,
	scope propertypb.PropertyScope,
	keys []string,
) (*ResolvedCatalog, error) {
	definitions, err := c.Resolve(ctx, db, scope, keys)
	if err != nil {
		return nil, err
	}
	groups, err := c.Groups(ctx, db, scope, definitions)
	if err != nil {
		return nil, err
	}
	return &ResolvedCatalog{Definitions: definitions, Groups: groups}, nil
}

// SystemPropertyEnabled resolves a system property's effective enabled state
// through the same cached override path used for normal definition reads.
func (c *DefinitionCache) SystemPropertyEnabled(
	ctx context.Context,
	db *sql.DB,
	scope propertypb.PropertyScope,
	key string,
) (bool, error) {
	definition, ok := SystemPropertyDefinition(scope, key)
	if !ok ||
		definition.GetSource() != propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_SYSTEM {
		return false, ErrSystemPropertyNotFound
	}
	cacheKey := scopeAndKeys(scope, []string{key})
	if enabled, ok := c.enabled.Get(cacheKey); ok {
		return enabled, nil
	}

	loaded, err, _ := c.loads.Do("enabled:"+cacheKey, func() (any, error) {
		if enabled, ok := c.enabled.Get(cacheKey); ok {
			return enabled, nil
		}
		generation := c.currentGeneration(scope)
		scopeName, ok := ScopeName(scope)
		if !ok {
			return nil, ErrSystemPropertyNotFound
		}
		tOverrides := table.FivenetPropertySystemOverrides
		rows := []model.FivenetPropertySystemOverrides{}
		if err := tOverrides.SELECT(tOverrides.Enabled).
			FROM(tOverrides).
			WHERE(tOverrides.Scope.EQ(mysql.String(scopeName)).AND(tOverrides.PropertyKey.EQ(mysql.String(key)))).
			QueryContext(ctx, db, &rows); err != nil {
			return nil, err
		}
		enabled := systemPropertyEnabledFromOverrides(definition, rows)
		if c.currentGeneration(scope) == generation {
			c.enabled.Put(cacheKey, enabled, definitionCacheTTL)
		}
		return enabled, nil
	})
	if err != nil {
		return false, err
	}
	return loaded.(bool), nil
}

// Invalidate clears all cached metadata for a scope. Definition and group
// changes are invalidated together because group metadata contributes to the
// effective property view.
func (c *DefinitionCache) Invalidate(scope propertypb.PropertyScope) {
	prefix := scopeAndKeys(scope, nil)
	definitionKeys := make([]string, 0)
	c.definitions.Range(func(key string, _ map[string]*propertypb.PropertyDefinition) bool {
		if strings.HasPrefix(key, prefix) {
			definitionKeys = append(definitionKeys, key)
		}
		return true
	})
	for _, key := range definitionKeys {
		c.definitions.Delete(key)
	}
	groupKeys := make([]string, 0)
	c.groups.Range(func(key string, _ []*propertypb.PropertyGroup) bool {
		if strings.HasPrefix(key, prefix) {
			groupKeys = append(groupKeys, key)
		}
		return true
	})
	for _, key := range groupKeys {
		c.groups.Delete(key)
	}
	enabledKeys := make([]string, 0)
	c.enabled.Range(func(key string, _ bool) bool {
		if strings.HasPrefix(key, prefix) {
			enabledKeys = append(enabledKeys, key)
		}
		return true
	})
	for _, key := range enabledKeys {
		c.enabled.Delete(key)
	}
	c.mu.Lock()
	c.generation[scope]++
	c.mu.Unlock()
}

func systemPropertyEnabledFromOverrides(
	definition *propertypb.PropertyDefinition,
	overrides []model.FivenetPropertySystemOverrides,
) bool {
	if len(overrides) == 0 {
		return definition.GetEnabled()
	}
	return overrides[0].Enabled
}

func (c *DefinitionCache) currentGeneration(scope propertypb.PropertyScope) uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.generation[scope]
}

func scopeAndKeys(scope propertypb.PropertyScope, keys []string) string {
	ordered := append([]string(nil), keys...)
	sort.Strings(ordered)
	return string(rune(scope)) + ":" + strings.Join(ordered, "\x00")
}

func definitionGroupKey(
	scope propertypb.PropertyScope,
	definitions map[string]*propertypb.PropertyDefinition,
) string {
	keys := make([]string, 0, len(definitions))
	for key, definition := range definitions {
		if definition == nil {
			keys = append(keys, key+"\x00nil")
			continue
		}
		keys = append(keys, strings.Join([]string{
			key,
			definition.GetGroupKey(),
			strconv.FormatBool(definition.GetEnabled()),
			strconv.FormatInt(int64(definition.GetGroupSortOrder()), 10),
		}, "\x00"))
	}
	sort.Strings(keys)
	return string(rune(scope)) + ":" + strings.Join(keys, "\x01")
}

func cloneDefinitions(
	in map[string]*propertypb.PropertyDefinition,
) map[string]*propertypb.PropertyDefinition {
	out := make(map[string]*propertypb.PropertyDefinition, len(in))
	for key, value := range in {
		out[key] = proto.Clone(value).(*propertypb.PropertyDefinition)
	}
	return out
}

func cloneGroups(in []*propertypb.PropertyGroup) []*propertypb.PropertyGroup {
	out := make([]*propertypb.PropertyGroup, len(in))
	for i, value := range in {
		out[i] = proto.Clone(value).(*propertypb.PropertyGroup)
	}
	return out
}
