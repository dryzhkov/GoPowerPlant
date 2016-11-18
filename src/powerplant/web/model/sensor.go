package model

type Sensor struct {
	Name         string  `json:"name"`
	SerialNo     string  `json:"serialNo"`
	UnitType     string  `json:"unitType"`
	MinSafeValue float64 `json:"minSafeValue"`
	MaxSafeValue float64 `json:"maxSafeValue"`
}

var sensorDb map[string]Sensor

func init() {
	sensorDb = make(map[string]Sensor)

	sensorDb["sensor1"] = Sensor{
		Name:         "sensor1",
		SerialNo:     "1234",
		UnitType:     "F",
		MinSafeValue: 0,
		MaxSafeValue: 10,
	}

	sensorDb["sensor2"] = Sensor{
		Name:         "sensor2",
		SerialNo:     "2345",
		UnitType:     "F",
		MinSafeValue: 0,
		MaxSafeValue: 10,
	}

	sensorDb["sensor3"] = Sensor{
		Name:         "sensor3",
		SerialNo:     "3456",
		UnitType:     "KPa",
		MinSafeValue: 0,
		MaxSafeValue: 10,
	}

	sensorDb["sensor4"] = Sensor{
		Name:         "sensor4",
		SerialNo:     "4567",
		UnitType:     "KPa",
		MinSafeValue: 0,
		MaxSafeValue: 10,
	}

}

func GetSensorByName(name string) Sensor {
	return sensorDb[name]
}
