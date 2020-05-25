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
		Use:                   "protog <name> [-ofop] [-m Message[field:type,field:type,...]]",
		Short:                 "protog is a protobuf file generator for the command line",
		Example:               "protog Greet.v1 -m HelloRequest[data:string]",
		Version:               "1.0.0",
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

			messages, err := flags.GetStringArray("message")
			fatalErr(err)

			msgs := map[string]interface{}{}
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

			in := map[string]interface{}{
				"syntax":  "proto3",
				"package": name,
			}

			if len(msgs) > 0 {
				in["message"] = msgs
			}

			data, err := protog.Encode(in)
			fatalErr(err)

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
	rootCmd.Flags().StringArrayP("message", "m", nil, "message related arguments")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
