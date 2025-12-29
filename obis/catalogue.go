package obis

var (
	catalogue map[string]string
)

func GetFromCatalogue(key string) string {
	if catalogue == nil {
		catalogue = map[string]string{
			// power values
			"1.8.0":  "energy_cons",
			"1.8.1":  "energy_cons_tariff_t1",
			"1.8.2":  "energy_cons_tariff_t2",
			"2.8.0":  "energy_prod",
			"2.8.1":  "energy_prod_tariff_t1",
			"2.8.2":  "energy_prod_tariff_t2",
			"1.7.0":  "power_active_cons",
			"15.7.0": "power_active_abs",
			"16.7.0": "power_active",
			"36.7.0": "power_active_l1",
			"56.7.0": "power_active_l2",
			"76.7.0": "power_active_l3",

			// heat values
			"6.8": "heat_energy_total",

			// device id
			"96.1.0": "device_id",
			"0.0.9":  "device_id",
			"9.20":   "device_id",
			"9.21":   "device_id",
		}
	}
	return catalogue[key]
}
