package usecase

import (
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

type (
	ConfigData struct {
		tileConfigs        map[models.TileType]map[models.Variant]*TileConfig
		dynamicTileConfigs map[models.TileType]map[models.Variant]*DynamicTileConfig
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Validator       utils.Validator
		Path            string
		InitialMaxDelay int
	}

	// DynamicTileConfig struct is used by GetConfig endpoint to check / hydrate config
	DynamicTileConfig struct {
		Validator utils.Validator
		Builder   builder.DynamicTileBuilder
	}
)

func initConfigData() *ConfigData {
	// TODO

	return &ConfigData{
		tileConfigs:        make(map[models.TileType]map[models.Variant]*TileConfig),
		dynamicTileConfigs: make(map[models.TileType]map[models.Variant]*DynamicTileConfig),
	}
}

func (cu *configUsecase) RegisterTile(tileType models.TileType, variant []models.Variant, version string) {
	// TODO
}

func (cu *configUsecase) EnableTile(
	tileType models.TileType, variant models.Variant, clientConfigValidator utils.Validator, path string, initialMaxDelay int,
) {
	value, exists := cu.configData.tileConfigs[tileType]
	if !exists {
		value = make(map[models.Variant]*TileConfig)
		cu.configData.tileConfigs[tileType] = value
	}

	value[variant] = &TileConfig{
		Path:            path,
		Validator:       clientConfigValidator,
		InitialMaxDelay: initialMaxDelay,
	}
}

func (cu *configUsecase) EnableDynamicTile(
	tileType models.TileType, variant models.Variant, clientConfigValidator utils.Validator, builder builder.DynamicTileBuilder,
) {
	// Used for authorized type
	cu.configData.tileConfigs[tileType] = nil

	value, exists := cu.configData.dynamicTileConfigs[tileType]
	if !exists {
		value = make(map[models.Variant]*DynamicTileConfig)
	}

	value[variant] = &DynamicTileConfig{
		Validator: clientConfigValidator,
		Builder:   builder,
	}
	cu.configData.dynamicTileConfigs[tileType] = value
}
