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

func parseOptions(options []string) (opts [][]string) {
	opts = make([][]string, 0)

	for _, opt := range options {
		parsed := strings.Split(opt, ":")
		if len(parsed) != 2 {
			fatal("error parsing options")
		}
		opts = append(opts, parsed)
	}

	return
}

func parseMessages(messages []string) (msgs map[string]interface{}) {
	msgs = make(map[string]interface{})

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

	return
}
func parseServices(services []string) (svcs map[string]interface{}) {
	svcs = make(map[string]interface{})

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
				if parsed[1][0] == '+' {
					in = "stream "
					in += parsed[1][1:]
				} else {
					in = parsed[1]
				}
			}
			if len(parsed) > 2 {
				if parsed[2][0] == '+' {
					out = "stream "
					out += parsed[2][1:]
				} else {
					out = parsed[2]
				}
			}

			methods[methodName] = map[string]string{
				"in":  in,
				"out": out,
			}
		}

		svcs[svcName] = methods
	}

	return
}

func main() {
	rootCmd := &cobra.Command{
		Use:                   "protog <name> [-dhfomnsp] [-n option_name:option_value] [-m MessageName[field:type,field:type,...]] [-s ServiceName[MethodName:In:Out]]",
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

			in := map[string]interface{}{
				"syntax":  "proto3",
				"package": name,
			}

			options, err := flags.GetStringSlice("option")
			fatalErr(err)

			opts := parseOptions(options)
			if len(opts) > 0 {
				in["option"] = opts
			}

			messages, err := flags.GetStringArray("message")
			fatalErr(err)

			msgs := parseMessages(messages)
			if len(msgs) > 0 {
				in["message"] = msgs
			}

			services, err := flags.GetStringArray("service")
			fatalErr(err)
			svcs := parseServices(services)
			if len(svcs) > 0 {
				in["service"] = svcs
			}

			data, err := protog.Encode(in)
			fatalErr(err)

			dryrun, err := flags.GetBool("dryrun")
			fatalErr(err)

			if dryrun {
				fmt.Println(string(data))
				return
			}

			var filename string
			if !strings.Contains(name, ".proto") {
				filename = strings.ToLower(name) + ".proto"
			}

			filepath := output + "/" + filename
			force, err := flags.GetBool("force")
			fatalErr(err)

			if !force {
				_, err := os.Stat(filepath)
				if err == nil || !os.IsNotExist(err) {
					fatal("file already exists, pass -f if you want to overwrite it")
				}
			}

			f, err := os.Create(filepath)
			fatalErr(err)

			f.Write(data)
			fatalErr(f.Close())
		},
	}

	rootCmd.Flags().BoolP("force", "f", false, "overwrite the file if it already exists")
	rootCmd.Flags().BoolP("dryrun", "d", false, "prints the generated proto to stdout")
	rootCmd.Flags().StringP("output", "o", ".", "output dir for the generated proto")
	rootCmd.Flags().StringP("package", "p", "", "package name")
	rootCmd.Flags().StringSliceP("option", "n", nil, "add an option")
	rootCmd.Flags().StringArrayP("message", "m", nil, "add a message and its fields")
	rootCmd.Flags().StringArrayP("service", "s", nil, "add a service and its methods")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
