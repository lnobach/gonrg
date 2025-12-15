package obis

var (
	catalogue map[string]string
)

func GetFromCatalogue(key string) string {
	if catalogue == nil {
		catalogue = map[string]string{
			// power values
			"1.8.0": "energy_cons",
			"2.8.0": "energy_prod",

			// heat values
			"6.8": "heat_energy_total",
		}
	}
	return catalogue[key]
}
