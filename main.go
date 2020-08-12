package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/schollz/cowyo/server"

	cli "github.com/urfave/cli/v2"
)

var version string
var pathToData string

func main() {
	app := cli.NewApp()
	app.Name = "cowyo"
	app.Usage = "a simple wiki"
	app.Version = version
	app.Compiled = time.Now()
	app.Action = func(c *cli.Context) error {
		pathToData = c.String("data")
		os.MkdirAll(pathToData, 0755)
		host := c.String("host")
		crtFlag := c.String("cert") // crt flag
		keyFlag := c.String("key")  // key flag
		if host == "" {
			host = GetLocalIP()
		}
		TLS := false
		if crtFlag != "" && keyFlag != "" {
			TLS = true
		}
		if TLS {
			fmt.Printf("\nRunning cowyo server (version %s) at https://%s:%s\n\n", version, host, c.String("port"))
		} else {
			fmt.Printf("\nRunning cowyo server (version %s) at http://%s:%s\n\n", version, host, c.String("port"))
		}

		server.Serve(
			pathToData,
			c.String("host"),
			c.String("port"),
			c.String("cert"),
			c.String("key"),
			TLS,
			c.String("css"),
			c.String("default-page"),
			c.String("lock"),
			c.Int("debounce"),
			c.Bool("diary"),
			c.String("cookie-secret"),
			c.String("access-code"),
			c.Bool("allow-insecure-markup"),
			c.Bool("allow-file-uploads"),
			c.Uint("max-upload-mb"),
			c.Uint("max-document-length"),
			logger(c.Bool("debug")),
		)
		return nil
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "data",
			Value: "data",
			Usage: "data folder to use",
		},
		&cli.StringFlag{
			Name:  "olddata",
			Value: "",
			Usage: "data folder for migrating",
		},
		&cli.StringFlag{
			Name:  "host",
			Value: "",
			Usage: "host to use",
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "8050",
			Usage:   "port to use",
		},
		&cli.StringFlag{
			Name:  "cert",
			Value: "",
			Usage: "absolute path to SSL public sertificate",
		},
		&cli.StringFlag{
			Name:  "key",
			Value: "",
			Usage: "absolute path to SSL private key",
		},
		&cli.StringFlag{
			Name:  "css",
			Value: "",
			Usage: "use a custom CSS file",
		},
		&cli.StringFlag{
			Name:  "default-page",
			Value: "",
			Usage: "show default-page/read instead of editing (default: show random editing)",
		},
		&cli.BoolFlag{
			Name:  "allow-insecure-markup",
			Usage: "Skip HTML sanitization",
		},
		&cli.StringFlag{
			Name:  "lock",
			Value: "",
			Usage: "password to lock editing all files (default: all pages unlocked)",
		},
		&cli.IntFlag{
			Name:  "debounce",
			Value: 500,
			Usage: "debounce time for saving data, in milliseconds",
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "turn on debugging",
		},
		&cli.BoolFlag{
			Name:  "diary",
			Usage: "turn diary mode (doing New will give a timestamped page)",
		},
		&cli.StringFlag{
			Name:  "access-code",
			Value: "",
			Usage: "Secret code to login with before accessing any wiki stuff",
		},
		&cli.StringFlag{
			Name:  "cookie-secret",
			Value: "secret",
			Usage: "random data to use for cookies; changing it will invalidate all sessions",
		},
		&cli.BoolFlag{
			Name:  "allow-file-uploads",
			Usage: "Enable file uploads",
		},
		&cli.UintFlag{
			Name:  "max-upload-mb",
			Value: 2,
			Usage: "Largest file upload (in mb) allowed",
		},
		&cli.UintFlag{
			Name:  "max-document-length",
			Value: 100000000,
			Usage: "Largest wiki page (in characters) allowed",
		},
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "migrate from the old cowyo",
			Action: func(c *cli.Context) error {
				pathToData = c.String("data")
				pathToOldData := c.String("olddata")
				if len(pathToOldData) == 0 {
					fmt.Printf("You need to specify folder with -olddata")
					return nil
				}
				os.MkdirAll(pathToData, 0755)
				if !exists(pathToOldData) {
					fmt.Printf("Can not find '%s', does it exist?", pathToOldData)
					return nil
				}
				server.Migrate(pathToOldData, pathToData, logger(c.Bool("debug")))
				return nil
			},
		},
	}

	app.Run(os.Args)
}

// GetLocalIP returns the local ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	bestIP := ""
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return bestIP
}

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func logger(debug bool) *lumber.ConsoleLogger {
	if !debug {
		return lumber.NewConsoleLogger(lumber.WARN)
	}
	return lumber.NewConsoleLogger(lumber.TRACE)

}
