package conf

// Config 配置文件不能有下划线或横线，否则不能解析
type Config struct {
	Http *Http `yaml:"http"`

	Debug *Debug `yaml:"debug"`

	Database *Database `yaml:"database"`

	Redis *Redis `yaml:"redis"`

	JWKS *JWKS `yaml:"jwks"`

	Metrics *Metrics `yaml:"metrics"`

	OpenAI *OpenAI `yaml:"openai"`

	LLM *LLM `yaml:"llm"`

	S3 *S3 `yaml:"s3"`

	Milvus *Milvus `yaml:"milvus"`

	Kafka *Kafka `yaml:"kafka"`

	Account *Account `yaml:"account"`

	ThirdParty *ThirdParty `yaml:"third_party" mapstructure:"third_party"`
}

type ThirdParty struct {
	JinaAIKey  string `yaml:"jina_ai_key" mapstructure:"jina_ai_key"`
	BingAPIKey string `yaml:"bing_api_key" mapstructure:"bing_api_key"`
}

type Http struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Url  string `yaml:"url"`
}

type Debug struct {
	Enabled bool `yaml:"enabled"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWKS struct {
	Url string `yaml:"url"`
}

type Metrics struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
}

type OpenAI struct {
	ApiKey            string   `yaml:"api_key" mapstructure:"api_key"`
	BaseUrl           string   `yaml:"base_url" mapstructure:"base_url"`
	InternalBaseUrl   string   `yaml:"internal_base_url" mapstructure:"internal_base_url"`
	Model             string   `yaml:"model" mapstructure:"model"`
	VisionModel       string   `yaml:"vision_model" mapstructure:"vision_model"`
	LongContextModel  string   `yaml:"long_context_model" mapstructure:"long_context_model"`
	MemoryModel       string   `yaml:"memory_model" mapstructure:"memory_model"`
	EmbeddingModel    string   `yaml:"embedding_model" mapstructure:"embedding_model"`
	EmbeddingDim      int      `yaml:"embedding_dim" mapstructure:"embedding_dim"`
	EmbeddingMaxToken int      `yaml:"embedding_max_token"  mapstructure:"embedding_max_token"`
	DallEModel        string   `yaml:"dall_e_model" mapstructure:"dall_e_model"`
	AllowedModels     []string `yaml:"allowed_models" mapstructure:"allowed_models"`
}

type S3 struct {
	Endpoint         string `yaml:"endpoint" mapstructure:"endpoint"`
	ExternalEndpoint string `yaml:"external_endpoint" mapstructure:"external_endpoint"`
	AccessKey        string `yaml:"access_key" mapstructure:"access_key"`
	SecretKey        string `yaml:"secret_key" mapstructure:"secret_key"`
	Bucket           string `yaml:"bucket" mapstructure:"bucket"`
	UseSSL           bool   `yaml:"use_ssl" mapstructure:"use_ssl"`
	Region           string `yaml:"region" mapstructure:"region"`
}

type LLM struct {
	MaxTokens                  int    `yaml:"max_tokens" mapstructure:"max_tokens"`
	ContextOptimizeActiveCount int    `yaml:"context_optimize_active_count" mapstructure:"context_optimize_active_count"`
	PrimarySystemPrompt        string `yaml:"primary_system_prompt" mapstructure:"primary_system_prompt"`
	DefaultSystemPrompt        string `yaml:"default_system_prompt" mapstructure:"default_system_prompt"`
	//Temperature float64 `yaml:"temperature"`
	//TopP float64 `yaml:"top_p"`
	//N int `yaml:"n"`
}

type Milvus struct {
	Host                   string `yaml:"host" mapstructure:"host"`
	Port                   int    `yaml:"port" mapstructure:"port"`
	DBName                 string `yaml:"db_name" mapstructure:"db_name"`
	MemoryCollection       string `yaml:"memory_collection" mapstructure:"memory_collection"`
	DocumentCollection     string `yaml:"document_collection" mapstructure:"document_collection"`
	MessageBlockCollection string `yaml:"message_block_collection" mapstructure:"message_block_collection"`
	User                   string `yaml:"user" mapstructure:"user"`
	Password               string `yaml:"password" mapstructure:"password"`
}

type Kafka struct {
	BootstrapServers KafkaBootstrapServers `yaml:"bootstrap_servers" mapstructure:"bootstrap_servers"`
	Topic            string                `yaml:"topic" mapstructure:"topic"`
	GroupId          string                `yaml:"group_id" mapstructure:"group_id"`
	Username         string                `yaml:"username" mapstructure:"username"`
	Password         string                `yaml:"password" mapstructure:"password"`
}

type KafkaBootstrapServers []string

type Account struct {
	Host           string `yaml:"host" mapstructure:"host"`
	ApplicationKey string `yaml:"application_key" mapstructure:"application_key"`
	Unit           string `yaml:"unit" mapstructure:"unit"`
	UnitStart      int64  `yaml:"unit_start" mapstructure:"unit_start"`
}
