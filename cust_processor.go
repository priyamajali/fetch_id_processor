// custom_processor.go
package resourceProcessor

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.opentelemetry.io/collector/processor"
	"net/http"
	"io"
	// "go.opentelemetry.io/collector/config/configmodels"
	

	

	"go.uber.org/zap"
	
)
var processorCapabilities = consumer.Capabilities{MutatesData: true}

// type Config struct {

// 	// AttributesActions specifies the list of actions to be applied on resource attributes.
// 	// The set of actions are {INSERT, UPDATE, UPSERT, DELETE, HASH, EXTRACT}.
// 	// AttributesActions = "trial"
	
// }

type Config struct {
	LabelName  string `mapstructure:"labelname"`
	LabelValue string `mapstructure:"labelvalue"`
}

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	
	return nil
}
type resourceProcessor struct {
	logger   *zap.Logger
	config   *Config
	
}

func (rp *resourceProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	rms := md.ResourceMetrics()
	// fmt.Printf("Name: %s, Value: %s\n", cfg.Name, cfg.Value)
	fmt.Printf("Name: %s, Value: %s\n", rp.config.LabelName, rp.config.LabelValue)
	
	for i := 0; i < rms.Len(); i++ {
		// rp.logger.Info("Hello, World!")
		fmt.Printf("Name: %s, Value: %s\n", rp.config.LabelName, rp.config.LabelValue)

	}
	apiValue, err := rp.fetchAPIValue(ctx, "http://172.233.224.248:5000/")
	if err != nil {
		rp.logger.Error("Failed to fetch API value", zap.Error(err))
	} else {
		rp.logger.Info("API Value:", zap.String("value", apiValue))
		fmt.Println("API Value: %s\n",  apiValue)
	}
	return md, nil
}

func (rp *resourceProcessor) fetchAPIValue(ctx context.Context, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	// Note: This is a simplified example; you might need to adjust based on your API response format
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}


// func NewFactory() component.ProcessorFactory {
// 	return processor.NewFactory(
// 		"resource",
// 		// createDefaultConfig,
// 		processorhelper.WithMetrics(createMetricsProcessor))
// }

func NewFactory() processor.Factory {
	return processor.NewFactory(
		"custom_processor",
		createDefaultConfig,
		processor.WithMetrics(createMetricsProcessor, component.StabilityLevelStable),
	)
}

// func createDefaultConfig()  {
// 	return 
// }

func createDefaultConfig() component.Config {
	return &Config{
		LabelName:  "defaultName",
		LabelValue: "defaultValue",
	}
}


func createMetricsProcessor(
	ctx context.Context,
	set processor.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Metrics) (processor.Metrics, error) {
	
	myConfig := cfg.(*Config)
	proc := &resourceProcessor{logger: set.Logger,config: myConfig}
	return processorhelper.NewMetricsProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		proc.processMetrics,
		processorhelper.WithCapabilities(processorCapabilities))
}

