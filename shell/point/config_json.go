// +build json

package point

import (
	"encoding/json"
	"errors"
	"io"

	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
)

const (
	DefaultRefreshMinute = int(9999)
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Port            v2net.Port                `json:"port"` // Port of this Point server.
		LogConfig       *log.Config               `json:"log"`
		RouterConfig    *router.Config            `json:"routing"`
		DNSConfig       *dns.Config               `json:"dns"`
		InboundConfig   *InboundConnectionConfig  `json:"inbound"`
		OutboundConfig  *OutboundConnectionConfig `json:"outbound"`
		InboundDetours  []*InboundDetourConfig    `json:"inboundDetour"`
		OutboundDetours []*OutboundDetourConfig   `json:"outboundDetour"`
		Transport       *transport.Config         `json:"transport"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Point: Failed to parse config: " + err.Error())
	}
	this.Port = jsonConfig.Port
	this.LogConfig = jsonConfig.LogConfig
	this.RouterConfig = jsonConfig.RouterConfig

	if jsonConfig.InboundConfig == nil {
		return errors.New("Point: Inbound config is not specified.")
	}
	this.InboundConfig = jsonConfig.InboundConfig

	if jsonConfig.OutboundConfig == nil {
		return errors.New("Point: Outbound config is not specified.")
	}
	this.OutboundConfig = jsonConfig.OutboundConfig
	this.InboundDetours = jsonConfig.InboundDetours
	this.OutboundDetours = jsonConfig.OutboundDetours
	if jsonConfig.DNSConfig == nil {
		jsonConfig.DNSConfig = &dns.Config{
			NameServers: []*v2net.DestinationPB{{
				Network: v2net.Network_UDP,
				Address: &v2net.AddressPB{
					Address: &v2net.AddressPB_Domain{
						Domain: "localhost",
					},
				},
				Port: 53,
			}},
		}
	}
	this.DNSConfig = jsonConfig.DNSConfig
	this.TransportConfig = jsonConfig.Transport
	return nil
}

func (this *InboundConnectionConfig) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Port          uint16                 `json:"port"`
		Listen        *v2net.AddressPB       `json:"listen"`
		Protocol      string                 `json:"protocol"`
		StreamSetting *internet.StreamConfig `json:"streamSettings"`
		Settings      json.RawMessage        `json:"settings"`
		AllowPassive  bool                   `json:"allowPassive"`
	}

	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Point: Failed to parse inbound config: " + err.Error())
	}
	this.Port = v2net.Port(jsonConfig.Port)
	this.ListenOn = v2net.AnyIP
	if jsonConfig.Listen != nil {
		if jsonConfig.Listen.AsAddress().Family().IsDomain() {
			return errors.New("Point: Unable to listen on domain address: " + jsonConfig.Listen.AsAddress().Domain())
		}
		this.ListenOn = jsonConfig.Listen.AsAddress()
	}
	if jsonConfig.StreamSetting != nil {
		this.StreamSettings = jsonConfig.StreamSetting
	}

	this.Protocol = jsonConfig.Protocol
	this.Settings = jsonConfig.Settings
	this.AllowPassiveConnection = jsonConfig.AllowPassive
	return nil
}

func (this *OutboundConnectionConfig) UnmarshalJSON(data []byte) error {
	type JsonConnectionConfig struct {
		Protocol      string                 `json:"protocol"`
		SendThrough   *v2net.AddressPB       `json:"sendThrough"`
		StreamSetting *internet.StreamConfig `json:"streamSettings"`
		Settings      json.RawMessage        `json:"settings"`
	}
	jsonConfig := new(JsonConnectionConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Point: Failed to parse outbound config: " + err.Error())
	}
	this.Protocol = jsonConfig.Protocol
	this.Settings = jsonConfig.Settings

	if jsonConfig.SendThrough != nil {
		address := jsonConfig.SendThrough.AsAddress()
		if address.Family().IsDomain() {
			return errors.New("Point: Unable to send through: " + address.String())
		}
		this.SendThrough = address
	}
	if jsonConfig.StreamSetting != nil {
		this.StreamSettings = jsonConfig.StreamSetting
	}
	return nil
}

func (this *InboundDetourAllocationConfig) UnmarshalJSON(data []byte) error {
	type JsonInboundDetourAllocationConfig struct {
		Strategy    string `json:"strategy"`
		Concurrency int    `json:"concurrency"`
		RefreshMin  int    `json:"refresh"`
	}
	jsonConfig := new(JsonInboundDetourAllocationConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Point: Failed to parse inbound detour allocation config: " + err.Error())
	}
	this.Strategy = jsonConfig.Strategy
	this.Concurrency = jsonConfig.Concurrency
	this.Refresh = jsonConfig.RefreshMin
	if this.Strategy == AllocationStrategyRandom {
		if this.Refresh == 0 {
			this.Refresh = 5
		}
		if this.Concurrency == 0 {
			this.Concurrency = 3
		}
	}
	if this.Refresh == 0 {
		this.Refresh = DefaultRefreshMinute
	}
	return nil
}

func (this *InboundDetourConfig) UnmarshalJSON(data []byte) error {
	type JsonInboundDetourConfig struct {
		Protocol      string                         `json:"protocol"`
		PortRange     *v2net.PortRange               `json:"port"`
		ListenOn      *v2net.AddressPB               `json:"listen"`
		Settings      json.RawMessage                `json:"settings"`
		Tag           string                         `json:"tag"`
		Allocation    *InboundDetourAllocationConfig `json:"allocate"`
		StreamSetting *internet.StreamConfig         `json:"streamSettings"`
		AllowPassive  bool                           `json:"allowPassive"`
	}
	jsonConfig := new(JsonInboundDetourConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Point: Failed to parse inbound detour config: " + err.Error())
	}
	if jsonConfig.PortRange == nil {
		log.Error("Point: Port range not specified in InboundDetour.")
		return common.ErrBadConfiguration
	}
	this.ListenOn = v2net.AnyIP
	if jsonConfig.ListenOn != nil {
		if jsonConfig.ListenOn.AsAddress().Family().IsDomain() {
			return errors.New("Point: Unable to listen on domain address: " + jsonConfig.ListenOn.AsAddress().Domain())
		}
		this.ListenOn = jsonConfig.ListenOn.AsAddress()
	}
	this.Protocol = jsonConfig.Protocol
	this.PortRange = *jsonConfig.PortRange
	this.Settings = jsonConfig.Settings
	this.Tag = jsonConfig.Tag
	this.Allocation = jsonConfig.Allocation
	if this.Allocation == nil {
		this.Allocation = &InboundDetourAllocationConfig{
			Strategy: AllocationStrategyAlways,
			Refresh:  DefaultRefreshMinute,
		}
	}
	if jsonConfig.StreamSetting != nil {
		this.StreamSettings = jsonConfig.StreamSetting
	}
	this.AllowPassiveConnection = jsonConfig.AllowPassive
	return nil
}

func (this *OutboundDetourConfig) UnmarshalJSON(data []byte) error {
	type JsonOutboundDetourConfig struct {
		Protocol      string                 `json:"protocol"`
		SendThrough   *v2net.AddressPB       `json:"sendThrough"`
		Tag           string                 `json:"tag"`
		Settings      json.RawMessage        `json:"settings"`
		StreamSetting *internet.StreamConfig `json:"streamSettings"`
	}
	jsonConfig := new(JsonOutboundDetourConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Point: Failed to parse outbound detour config: " + err.Error())
	}
	this.Protocol = jsonConfig.Protocol
	this.Tag = jsonConfig.Tag
	this.Settings = jsonConfig.Settings

	if jsonConfig.SendThrough != nil {
		address := jsonConfig.SendThrough.AsAddress()
		if address.Family().IsDomain() {
			return errors.New("Point: Unable to send through: " + address.String())
		}
		this.SendThrough = address
	}

	if jsonConfig.StreamSetting != nil {
		this.StreamSettings = jsonConfig.StreamSetting
	}
	return nil
}

func JsonLoadConfig(input io.Reader) (*Config, error) {
	jsonConfig := &Config{}
	decoder := json.NewDecoder(input)
	err := decoder.Decode(jsonConfig)
	if err != nil {
		log.Error("Point: Failed to load server config: ", err)
		return nil, err
	}

	return jsonConfig, err
}

func init() {
	configLoader = JsonLoadConfig
}
