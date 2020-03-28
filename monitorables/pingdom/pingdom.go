//+build !faker

package pingdom

import (
	"fmt"
	"net/url"

	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"

	coreModels "github.com/monitoror/monitoror/models"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	pingdomDelivery "github.com/monitoror/monitoror/monitorables/pingdom/api/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	pingdomRepository "github.com/monitoror/monitoror/monitorables/pingdom/api/repository"
	pingdomUsecase "github.com/monitoror/monitoror/monitorables/pingdom/api/usecase"
	pingdomConfig "github.com/monitoror/monitoror/monitorables/pingdom/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*pingdomConfig.Pingdom
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.VariantName]*pingdomConfig.Pingdom)

	// Load core config from env
	pkgMonitorable.LoadConfig(&monitorable.config, pingdomConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.PingdomCheckTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.PingdomChecksTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "Pingdom"
}

func (m *Monitorable) GetVariants() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(variant coreModels.VariantName) (bool, error) {
	conf := m.config[variant]

	// No configuration set
	if conf.URL == pingdomConfig.Default.URL && conf.Token == "" {
		return false, nil
	}

	// Error in URL
	if _, err := url.Parse(conf.URL); err != nil {
		return false, fmt.Errorf(`%s contains invalid URL: "%s"`, pkgMonitorable.BuildMonitorableEnvKey(conf, variant, "URL"), conf.URL)
	}

	// Error in Token
	if conf.Token == "" {
		return false, fmt.Errorf(`%s is required, no value found`, pkgMonitorable.BuildMonitorableEnvKey(conf, variant, "TOKEN"))
	}

	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	conf := m.config[variant]

	repository := pingdomRepository.NewPingdomRepository(conf)
	usecase := pingdomUsecase.NewPingdomUsecase(repository, m.store.CacheStore, conf.CacheExpiration)
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/pingdom", variant)
	route := routeGroup.GET("/pingdom", delivery.GetCheck)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.PingdomCheckTileType, variant,
		&pingdomModels.CheckParams{}, route.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableDynamicTile(api.PingdomChecksTileType, variant,
		&pingdomModels.ChecksParams{}, usecase.Checks)
}
