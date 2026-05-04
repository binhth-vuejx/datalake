package rilldata

// MainConfig is the root config file that imports domain files
type MainConfig struct {
	Version string   `yaml:"version"`
	Imports []string `yaml:"imports"`
}

// DatasetFile represents a single domain YAML file
type DatasetFile struct {
	Dataset DatasetConfig `yaml:"dataset"`
}

// Config is the merged configuration from all imported files
type Config struct {
	Version  string
	Datasets []DatasetConfig
}

// YAML schema types (internal to infrastructure)
type DatasetConfig struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Table       string        `yaml:"table"`
	Columns     []ColumnDef   `yaml:"columns"`
	Queries     []QueryDef    `yaml:"queries"`
}

type ColumnDef struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

type QueryDef struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Method      string     `yaml:"method"`
	SQL         string     `yaml:"sql"`
	Params      []ParamDef `yaml:"params"`
}

type ParamDef struct {
	Name     string      `yaml:"name"`
	Type     string      `yaml:"type"`
	Default  interface{} `yaml:"default"`
	Required bool        `yaml:"required"`
}
