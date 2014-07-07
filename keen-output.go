package heka_clever_plugins

import (
	"encoding/json"
	"fmt"
	"github.com/inconshreveable/go-keen"
	"github.com/mozilla-services/heka/pipeline"
)

type KeenOutput struct {
	client *keen.Client
}

type KeenOutputConfig struct {
	ApiKey       string `toml:"api_key"`
	ProjectToken string `toml:"project_token"`
}

func (ko *KeenOutput) ConfigStruct() interface{} {
	return &KeenOutputConfig{}
}

func (ko *KeenOutput) Init(rawConf interface{}) error {
	config := rawConf.(*KeenOutputConfig)
	ko.client = &keen.Client{ApiKey: config.ApiKey, ProjectToken: config.ProjectToken}
	return nil
}

func (ko *KeenOutput) Run(or pipeline.OutputRunner, h pipeline.PluginHelper) error {
	var (
		err  error
		pack *pipeline.PipelinePack
	)

	for pack = range or.InChan() {
		data, ok := pack.Message.GetFieldValue("Data")
		if !ok {
			or.LogError(fmt.Errorf("Could not get field value 'Data' from message %#v", data))
			continue
		}

		event := make(map[string]interface{})
		err = json.Unmarshal([]byte(data.(string)), &event)
		if err != nil {
			or.LogError(err)
			continue
		}
		err = ko.client.AddEvent("job-finished", event)
		if err != nil {
			or.LogError(err)
			continue
		}
		pack.Recycle()
	}
	return nil
}

func init() {
	pipeline.RegisterPlugin("KeenOutput", func() interface{} {
		return new(KeenOutput)
	})
}
