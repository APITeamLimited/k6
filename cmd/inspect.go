/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
)

func getInspectCmd(gs *globalState) *cobra.Command ***REMOVED***
	var addExecReqs bool

	// inspectCmd represents the inspect command
	inspectCmd := &cobra.Command***REMOVED***
		Use:   "inspect [file]",
		Short: "Inspect a script or archive",
		Long:  `Inspect a script or archive.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			test, err := loadTest(gs, cmd, args, nil)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			// At the moment, `k6 inspect` output can take 2 forms: standard
			// (equal to the lib.Options struct) and extended, with additional
			// fields with execution requirements.
			inspectOutput := interface***REMOVED******REMOVED***(test.initRunner.GetOptions())

			if addExecReqs ***REMOVED***
				inspectOutput, err = addExecRequirements(gs, cmd, test)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***

			data, err := json.MarshalIndent(inspectOutput, "", "  ")
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			printToStdout(gs, string(data))

			return nil
		***REMOVED***,
	***REMOVED***

	inspectCmd.Flags().SortFlags = false
	inspectCmd.Flags().AddFlagSet(runtimeOptionFlagSet(false))
	inspectCmd.Flags().StringVarP(&gs.flags.testType, "type", "t",
		gs.flags.testType, "override file `type`, \"js\" or \"archive\"")
	inspectCmd.Flags().BoolVar(&addExecReqs,
		"execution-requirements",
		false,
		"include calculations of execution requirements for the test")

	return inspectCmd
***REMOVED***

func addExecRequirements(gs *globalState, cmd *cobra.Command, test *loadedTest) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// we don't actually support CLI flags here, so we pass nil as the getter
	if err := test.consolidateDeriveAndValidateConfig(gs, cmd, nil); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	et, err := lib.NewExecutionTuple(test.derivedConfig.ExecutionSegment, test.derivedConfig.ExecutionSegmentSequence)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	executionPlan := test.derivedConfig.Scenarios.GetFullExecutionRequirements(et)
	duration, _ := lib.GetEndOffset(executionPlan)

	return struct ***REMOVED***
		lib.Options
		TotalDuration types.NullDuration `json:"totalDuration"`
		MaxVUs        uint64             `json:"maxVUs"`
	***REMOVED******REMOVED***
		test.derivedConfig.Options,
		types.NewNullDuration(duration, true),
		lib.GetMaxPossibleVUs(executionPlan),
	***REMOVED***, nil
***REMOVED***
