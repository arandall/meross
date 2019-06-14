package mdp

import (
	"encoding/json"
	"sort"
)

const Ability_DeviceAbility = "Appliance.System.Ability"

type deviceAbility struct {
	PayloadVersion int                    `json:"payloadVersion"`
	Abilities      map[string]interface{} `json:"ability"`
}

// DeviceAbility lists the name spaces available on the device in alphabetical order.
type DeviceAbility struct {
	PayloadVersion int
	Abilities      []string
}

func (a *DeviceAbility) UnmarshalJSON(data []byte) error {
	var da deviceAbility
	err := json.Unmarshal(data, &da)
	if err != nil {
		return err
	}
	abilities := make([]string, len(da.Abilities))
	index := 0
	for ability, _ := range da.Abilities {
		abilities[index] = ability
		index++
	}

	sort.Strings(abilities)
	*a = DeviceAbility{
		Abilities:      abilities,
		PayloadVersion: da.PayloadVersion,
	}
	return nil
}

func (a *DeviceAbility) MarshalJSON() ([]byte, error) {
	abilities := make(map[string]interface{})
	for _, ability := range a.Abilities {
		abilities[ability] = struct{}{}
	}
	return json.Marshal(&deviceAbility{
		Abilities:      abilities,
		PayloadVersion: a.PayloadVersion,
	})
}

func GetAbilities() *Packet {
	return NewPacket(Ability_DeviceAbility, Method_GET, []byte(`{}`))
}
