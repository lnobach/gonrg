package cli

import (
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/outjson"
	"github.com/lnobach/gonrg/outtable"
	"github.com/lnobach/gonrg/server"
	"github.com/lnobach/gonrg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	debug      bool
	jsonOut    bool
	device     string
	baudrate   int
	strictMode bool
	rootCmd    = &cobra.Command{
		Use:     "gonrg",
		Short:   "\u26A1\uFE0F gonrg - a simple D0 OBIS energy meter CLI tool or server.",
		Long:    "\u26A1\uFE0F gonrg - a simple D0 OBIS energy meter CLI tool or server.",
		Version: version.GonrgVersion,

		Run: func(cmd *cobra.Command, args []string) {

			if debug {
				log.SetLevel(log.DebugLevel)
			}

			d, err := d0.NewDevice(d0.DeviceConfig{
				Device:   device,
				BaudRate: baudrate,
			})
			if err != nil {
				panic(err)
			}

			mt := time.Now()
			rawdata, err := d.Get()
			if err != nil {
				panic(err)
			}

			log.WithField("raw", rawdata).Debugf("Raw data from device %s", device)

			p, err := d0.NewParser(d0.ParseConfig{
				StrictMode: strictMode,
			})
			if err != nil {
				panic(err)
			}

			result, err := p.GetOBISList(rawdata, mt)
			if err != nil {
				panic(err)
			}

			if jsonOut {
				outjson.PrintJSON(result)
			} else {
				outtable.PrintTable(result)
			}

		},
	}
	config    string
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "run in server mode given a config",
		Long:  "run in server mode given a config",

		Run: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}

			conf, err := server.ConfigFromFile(config)
			if err != nil {
				panic(err)
			}

			s, err := server.NewServer(conf, debug)
			if err != nil {
				panic(err)
			}
			err = s.ListenAndServe()
			if err != nil {
				panic(err)
			}
		},
	}
)

func Start() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "D", false, "set debug log level")
	rootCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output json instead of pretty table")
	rootCmd.Flags().StringVarP(&device, "device", "d", "/dev/ttyUSB0", "device to read from")
	rootCmd.Flags().IntVarP(&baudrate, "baudrate", "b", 0, "baud rate, 0 means choose best option")
	rootCmd.Flags().BoolVarP(&strictMode, "strict", "S", false, "strict mode for parsing - fail fast")
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&config, "config", "C", "/etc/gonrg.yaml", "Config to use for server")
	serverCmd.Flags().BoolVarP(&debug, "debug", "D", false, "set debug log level")

}
