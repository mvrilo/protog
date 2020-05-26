package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mvrilo/protog"
	"github.com/spf13/cobra"
)

func fatal(err string) {
	fmt.Fprintf(os.Stderr, err)
	os.Exit(1)
}

func fatalErr(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:                   "protog <name> [-fhmops] [-m Message[field:type,field:type,...]] [-s ServiceName[MethodName:In:Out]]",
		Short:                 "protog is a protobuf file generator for the command line",
		Example:               "protog Greet.v1 -m HelloRequest[data:string]",
		Version:               "0.0.1",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Usage()
				return
			}

			name := args[0]
			if name == "" {
				fatal("proto name missing")
			}

			args = args[1:]
			flags := cmd.Flags()
			output, err := flags.GetString("output")
			fatalErr(err)

			force, err := flags.GetBool("force")
			fatalErr(err)

			messages, err := flags.GetStringArray("message")
			fatalErr(err)

			msgs := make(map[string]interface{})
			for _, message := range messages {
				msgSplit := strings.Split(message, "[")
				if len(msgSplit) != 2 {
					fatal("error parsing message")
				}

				msgName := msgSplit[0]
				msgFields := strings.Replace(msgSplit[1], "]", "", -1)
				fields := map[string]string{}
				for _, field := range strings.Split(msgFields, ",") {
					parsed := strings.Split(field, ":")
					if len(parsed) != 2 {
						fatal("error parsing message fields")
					}
					fName := parsed[0]
					fType := parsed[1]
					fields[fName] = fType
				}
				msgs[msgName] = fields
			}

			services, err := flags.GetStringArray("service")
			fatalErr(err)

			svcs := make(map[string]interface{})
			for _, service := range services {
				svcSplit := strings.Split(service, "[")
				if len(svcSplit) != 2 {
					fatal("error parsing service")
				}

				svcName := svcSplit[0]
				svcMethods := strings.Replace(svcSplit[1], "]", "", -1)
				methods := make(map[string]interface{})

				for _, method := range strings.Split(svcMethods, ",") {
					parsed := strings.Split(method, ":")
					if len(parsed) < 1 {
						fatal("error parsing service methods")
					}

					methodName := parsed[0]
					var in, out string
					if len(parsed) > 1 {
						in = parsed[1]
					}
					if len(parsed) > 2 {
						out = parsed[2]
					}

					methods[methodName] = map[string]string{
						"in":  in,
						"out": out,
					}
				}

				svcs[svcName] = methods
			}

			in := map[string]interface{}{
				"syntax":  "proto3",
				"package": name,
			}

			if len(msgs) > 0 {
				in["message"] = msgs
			}

			if len(svcs) > 0 {
				in["service"] = svcs
			}

			data, err := protog.Encode(in)
			fatalErr(err)

			dryrun, err := flags.GetBool("dryrun")
			fatalErr(err)

			if dryrun {
				println(string(data))
				return
			}

			var filename string
			if !strings.Contains(name, ".proto") {
				filename = strings.ToLower(name) + ".proto"
			}

			path := output + "/" + filename
			if !force {
				_, err := os.Stat(path)
				if err == nil || !os.IsNotExist(err) {
					fatal("file already exists, pass -f if you want to overwrite it")
				}
			}

			filepath := output + "/" + filename
			f, err := os.Create(filepath)
			fatalErr(err)

			f.Write(data)
			fatalErr(f.Close())
		},
	}

	rootCmd.Flags().BoolP("force", "f", false, "overwrite the file if it already exists")
	rootCmd.Flags().BoolP("dryrun", "d", false, "prints the generated proto in the stdout")
	rootCmd.Flags().StringP("output", "o", ".", "output dir for the generated proto")
	rootCmd.Flags().StringP("package", "p", "", "package name")
	rootCmd.Flags().StringArrayP("message", "m", nil, "message data")
	rootCmd.Flags().StringArrayP("service", "s", nil, "service data")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
