package rilldata

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/yourorg/datalake-free/internal/domain/dataset"
)

// LoadConfig loads main config and all imported domain files
func LoadConfig(mainPath string) (*Config, error) {
	// 1. Load main config
	mainData, err := os.ReadFile(mainPath)
	if err != nil {
		return nil, fmt.Errorf("read main config: %w", err)
	}

	var mainCfg MainConfig
	if err := yaml.Unmarshal(mainData, &mainCfg); err != nil {
		return nil, fmt.Errorf("parse main config: %w", err)
	}

	// 2. Load all imported datasets
	baseDir := filepath.Dir(mainPath)
	var datasets []DatasetConfig

	for _, importPath := range mainCfg.Imports {
		fullPath := filepath.Join(baseDir, importPath)

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", importPath, err)
		}

		var dsFile DatasetFile
		if err := yaml.Unmarshal(data, &dsFile); err != nil {
			return nil, fmt.Errorf("parse %s: %w", importPath, err)
		}

		datasets = append(datasets, dsFile.Dataset)
	}

	cfg := &Config{
		Version:  mainCfg.Version,
		Datasets: datasets,
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

// ToDomain converts YAML config to domain entities
func (c *Config) ToDomain() []dataset.Dataset {
	datasets := make([]dataset.Dataset, len(c.Datasets))
	for i, ds := range c.Datasets {
		datasets[i] = dataset.Dataset{
			Name:        ds.Name,
			Description: ds.Description,
			Table:       ds.Table,
			Columns:     toColumns(ds.Columns),
			Queries:     toQueries(ds.Queries),
		}
	}
	return datasets
}

func toColumns(cols []ColumnDef) []dataset.Column {
	result := make([]dataset.Column, len(cols))
	for i, c := range cols {
		result[i] = dataset.Column{
			Name:        c.Name,
			Type:        dataset.ParamType(c.Type),
			Description: c.Description,
		}
	}
	return result
}

func toQueries(queries []QueryDef) []dataset.Query {
	result := make([]dataset.Query, len(queries))
	for i, q := range queries {
		result[i] = dataset.Query{
			Name:        q.Name,
			Description: q.Description,
			Method:      dataset.QueryMethod(q.Method),
			SQL:         q.SQL,
			Params:      toParams(q.Params),
		}
	}
	return result
}

func toParams(params []ParamDef) []dataset.Param {
	result := make([]dataset.Param, len(params))
	for i, p := range params {
		result[i] = dataset.Param{
			Name:     p.Name,
			Type:     dataset.ParamType(p.Type),
			Default:  p.Default,
			Required: p.Required,
		}
	}
	return result
}
