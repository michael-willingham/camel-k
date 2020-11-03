/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newCmdInspect(rootCmdOptions *RootCmdOptions) (*cobra.Command, *inspectCmdOptions) {
	options := inspectCmdOptions{
		RootCmdOptions: rootCmdOptions,
	}

	cmd := cobra.Command{
		Use:   "inspect [files to inspect]",
		Short: "Generate dependencies list given integration files.",
		Long: `Output dependencies for a list of integration files. By default this command returns the
top level dependencies only. When --all-dependencies is enabled, the transitive dependencies
will be generated by calling Maven and then printed in the selected output format.`,
		PreRunE: decode(&options),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := options.validate(args); err != nil {
				return err
			}
			if err := options.run(args); err != nil {
				fmt.Println(err.Error())
			}

			return nil
		},
		Annotations: map[string]string{
			offlineCommandLabel: "true",
		},
	}

	cmd.Flags().Bool("all-dependencies", false, "Compute transitive dependencies and move them to directory pointed to by the --dependencies-directory flag.")
	cmd.Flags().StringArrayP("dependency", "d", nil, `Additional top-level dependency with the format:
<type>:<dependency-name>
where <type> is one of {`+strings.Join(acceptedDependencyTypes, "|")+`}.`)
	cmd.Flags().StringP("output", "o", "", "Output format. One of: json|yaml")

	return &cmd, &options
}

type inspectCmdOptions struct {
	*RootCmdOptions
	AllDependencies        bool     `mapstructure:"all-dependencies"`
	OutputFormat           string   `mapstructure:"output"`
	AdditionalDependencies []string `mapstructure:"dependencies"`
}

func (command *inspectCmdOptions) validate(args []string) error {
	return validateIntegrationForDependencies(args, command.AdditionalDependencies)
}

func (command *inspectCmdOptions) run(args []string) error {
	// Fetch dependencies.
	dependencies, err := getDependencies(args, command.AdditionalDependencies, command.AllDependencies)
	if err != nil {
		return err
	}

	// Print dependencies.
	err = outputDependencies(dependencies, command.OutputFormat)
	if err != nil {
		return err
	}

	return nil
}
