// Code generated by "stringer -type ThermostatMode"; DO NOT EDIT

package thermostat

import "fmt"

const _ThermostatMode_name = "THERMOSTAT_OFFTHERMOSTAT_HEATTHERMOSTAT_COOLTHERMOSTAT_AUTO"

var _ThermostatMode_index = [...]uint8{0, 14, 29, 44, 59}

func (i ThermostatMode) String() string {
	if i < 0 || i >= ThermostatMode(len(_ThermostatMode_index)-1) {
		return fmt.Sprintf("ThermostatMode(%d)", i)
	}
	return _ThermostatMode_name[_ThermostatMode_index[i]:_ThermostatMode_index[i+1]]
}
