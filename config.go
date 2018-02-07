package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type brconfig struct {
	NetInterface string          `toml:"net_interface"`
	Devices      []bonjourDevice `toml:"devices"`
}

type bonjourDevice struct {
	MacAddress  string `toml:"mac_address"`
	OriginPool  int    `toml:"origin_pool"`
	SharedPools []int  `toml:"shared_pools"`
}

func readConfig(path string) (cfg brconfig, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return brconfig{}, err
	}
	_, err = toml.Decode(string(content), &cfg)
	return cfg, err
}

func mapByPool(devices []bonjourDevice) map[int]([]int) {
	seen := make(map[int]map[int]bool)
	poolsMap := make(map[int]([]int))
	for _, device := range devices {
		for _, pool := range device.SharedPools {
			if _, ok := seen[pool]; !ok {
				seen[pool] = make(map[int]bool)
			}
			if _, ok := seen[pool][device.OriginPool]; !ok {
				seen[pool][device.OriginPool] = true
				poolsMap[pool] = append(poolsMap[pool], device.OriginPool)
			}
		}
	}
	return poolsMap
}

func mapByAddress(devices []bonjourDevice) map[string]([]int) {
	addressMap := make(map[string]([]int))
	for _, device := range devices {
		addressMap[device.MacAddress] = device.SharedPools
	}
	return addressMap
}
