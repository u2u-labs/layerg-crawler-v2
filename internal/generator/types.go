package generator

// CrawlerConfig represents the subgraph.yaml configuration
type CrawlerConfig struct {
	Network struct {
		ChainId  string   `yaml:"chainId"`
		Endpoint []string `yaml:"endpoint"`
	} `yaml:"network"`
	DataSources []struct {
		Kind       string         `yaml:"kind"`
		Options    ContractConfig `yaml:"options"`
		StartBlock int64          `yaml:"startBlock"`
		Mapping    struct {
			Handlers []struct {
				Kind    string `yaml:"kind"`
				Handler string `yaml:"handler"`
				Filter  struct {
					Function string   `yaml:"function,omitempty"`
					Topics   []string `yaml:"topics,omitempty"`
				} `yaml:"filter"`
			} `yaml:"handlers"`
		} `yaml:"mapping"`
	} `yaml:"dataSources"`
}

type ContractConfig struct {
	Address string `yaml:"address"`
	Abi     string `yaml:"abi"`
}
