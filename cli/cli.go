package cli

import (
	"fmt"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/outjson"
	"github.com/lnobach/gonrg/outtable"
	"github.com/lnobach/gonrg/server"
	"github.com/lnobach/gonrg/sml"
	"github.com/lnobach/gonrg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	debug         bool
	smlmode       bool
	jsonOut       bool
	device        string
	baudrate      int
	baudrate_read int
	responsedelay time.Duration
	d0timeout     time.Duration
	strictMode    bool
	dev_opts      *[]string
	rootCmd       = &cobra.Command{
		Use:     "gonrg",
		Short:   "\u26A1\uFE0F gonrg - a simple D0 OBIS/SML energy meter CLI tool or server.",
		Long:    "\u26A1\uFE0F gonrg - a simple D0 OBIS/SML energy meter CLI tool or server.",
		Version: version.GonrgVersion,

		RunE: func(cmd *cobra.Command, args []string) error {

			setLogLevel()

			cfg := d0.DeviceConfig{
				Device:        device,
				BaudRate:      baudrate,
				BaudRateRead:  baudrate_read,
				ResponseDelay: responsedelay,
				D0Timeout:     d0timeout,
				DeviceOptions: *dev_opts,
			}

			var err error
			var d d0.Device

			if smlmode {
				d, err = sml.NewDevice(cfg)
			} else {
				d, err = d0.NewDevice(cfg)
			}
			if err != nil {
				return fmt.Errorf("error creating device: %w", err)
			}

			mt := time.Now()
			rawdata, err := d.Get()
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("error retrieving data from device: %w", err)
			}

			pcfg := d0.ParseConfig{
				StrictMode: strictMode,
			}
			if err != nil {
				return fmt.Errorf("error creating parser: %w", err)
			}

			result, err := d0.ParseOBISList(&pcfg, rawdata, mt)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("error parsing data: %w", err)
			}

			if jsonOut {
				outjson.PrintJSON(result)
			} else {
				outtable.PrintTable(result)
			}
			return nil
		},
	}
	config    string
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "run in server mode given a config",
		Long:  "run in server mode given a config",

		RunE: func(cmd *cobra.Command, args []string) error {

			setLogLevel()

			conf, err := server.ConfigFromFile(config)
			if err != nil {
				return fmt.Errorf("error reading configuration: %w", err)
			}

			s, err := server.NewServer(conf, debug)
			if err != nil {
				return fmt.Errorf("error creating server: %w", err)
			}
			err = s.ListenAndServe()
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("error while serving: %w", err)
			}
			return nil
		},
	}
)

func setLogLevel() {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func Start() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "D", false, "set debug log level")
	rootCmd.Flags().BoolVarP(&smlmode, "sml", "s", false, "connect to an SML device rather than a plain OBIS device")
	rootCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output json instead of pretty table")
	rootCmd.Flags().StringVarP(&device, "device", "d", "/dev/ttyUSB0", "device to read from")
	rootCmd.Flags().IntVarP(&baudrate, "baudrate", "b", 0, "baud rate, 0 means choose best option")
	rootCmd.Flags().IntVarP(&baudrate_read, "baudrate-read", "r", 0, "(non-SML) baud rate for reading, 0 means same like baudrate")
	rootCmd.Flags().DurationVarP(&responsedelay, "response-delay", "l", 0, "(non-SML) wait before expecting response")
	rootCmd.Flags().DurationVarP(&d0timeout, "d0-timeout", "t", 0, "read timeout of the d0 serial connection")
	dev_opts = rootCmd.Flags().StringSliceP("device-option", "o", []string{}, "device option")
	rootCmd.Flags().BoolVarP(&strictMode, "strict", "S", false, "strict mode for parsing - fail fast")
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&config, "config", "C", "/etc/gonrg.yaml", "config to use for server")
	serverCmd.Flags().BoolVarP(&debug, "debug", "D", false, "set debug log level")

}
