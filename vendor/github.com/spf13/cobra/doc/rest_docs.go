//Copyright 2015 Red Hat Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package doc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func printOptionsReST(buf *bytes.Buffer, cmd *cobra.Command, name string) error ***REMOVED***
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() ***REMOVED***
		buf.WriteString("Options\n")
		buf.WriteString("~~~~~~~\n\n::\n\n")
		flags.PrintDefaults()
		buf.WriteString("\n")
	***REMOVED***

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() ***REMOVED***
		buf.WriteString("Options inherited from parent commands\n")
		buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n\n::\n\n")
		parentFlags.PrintDefaults()
		buf.WriteString("\n")
	***REMOVED***
	return nil
***REMOVED***

// linkHandler for default ReST hyperlink markup
func defaultLinkHandler(name, ref string) string ***REMOVED***
	return fmt.Sprintf("`%s <%s.rst>`_", name, ref)
***REMOVED***

// GenReST creates reStructured Text output.
func GenReST(cmd *cobra.Command, w io.Writer) error ***REMOVED***
	return GenReSTCustom(cmd, w, defaultLinkHandler)
***REMOVED***

// GenReSTCustom creates custom reStructured Text output.
func GenReSTCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string, string) string) error ***REMOVED***
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	short := cmd.Short
	long := cmd.Long
	if len(long) == 0 ***REMOVED***
		long = short
	***REMOVED***
	ref := strings.Replace(name, " ", "_", -1)

	buf.WriteString(".. _" + ref + ":\n\n")
	buf.WriteString(name + "\n")
	buf.WriteString(strings.Repeat("-", len(name)) + "\n\n")
	buf.WriteString(short + "\n\n")
	buf.WriteString("Synopsis\n")
	buf.WriteString("~~~~~~~~\n\n")
	buf.WriteString("\n" + long + "\n\n")

	if cmd.Runnable() ***REMOVED***
		buf.WriteString(fmt.Sprintf("::\n\n  %s\n\n", cmd.UseLine()))
	***REMOVED***

	if len(cmd.Example) > 0 ***REMOVED***
		buf.WriteString("Examples\n")
		buf.WriteString("~~~~~~~~\n\n")
		buf.WriteString(fmt.Sprintf("::\n\n%s\n\n", indentString(cmd.Example, "  ")))
	***REMOVED***

	if err := printOptionsReST(buf, cmd, name); err != nil ***REMOVED***
		return err
	***REMOVED***
	if hasSeeAlso(cmd) ***REMOVED***
		buf.WriteString("SEE ALSO\n")
		buf.WriteString("~~~~~~~~\n\n")
		if cmd.HasParent() ***REMOVED***
			parent := cmd.Parent()
			pname := parent.CommandPath()
			ref = strings.Replace(pname, " ", "_", -1)
			buf.WriteString(fmt.Sprintf("* %s \t - %s\n", linkHandler(pname, ref), parent.Short))
			cmd.VisitParents(func(c *cobra.Command) ***REMOVED***
				if c.DisableAutoGenTag ***REMOVED***
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				***REMOVED***
			***REMOVED***)
		***REMOVED***

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children ***REMOVED***
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() ***REMOVED***
				continue
			***REMOVED***
			cname := name + " " + child.Name()
			ref = strings.Replace(cname, " ", "_", -1)
			buf.WriteString(fmt.Sprintf("* %s \t - %s\n", linkHandler(cname, ref), child.Short))
		***REMOVED***
		buf.WriteString("\n")
	***REMOVED***
	if !cmd.DisableAutoGenTag ***REMOVED***
		buf.WriteString("*Auto generated by spf13/cobra on " + time.Now().Format("2-Jan-2006") + "*\n")
	***REMOVED***
	_, err := buf.WriteTo(w)
	return err
***REMOVED***

// GenReSTTree will generate a ReST page for this command and all
// descendants in the directory given.
// This function may not work correctly if your command names have `-` in them.
// If you have `cmd` with two subcmds, `sub` and `sub-third`,
// and `sub` has a subcommand called `third`, it is undefined which
// help output will be in the file `cmd-sub-third.1`.
func GenReSTTree(cmd *cobra.Command, dir string) error ***REMOVED***
	emptyStr := func(s string) string ***REMOVED*** return "" ***REMOVED***
	return GenReSTTreeCustom(cmd, dir, emptyStr, defaultLinkHandler)
***REMOVED***

// GenReSTTreeCustom is the the same as GenReSTTree, but
// with custom filePrepender and linkHandler.
func GenReSTTreeCustom(cmd *cobra.Command, dir string, filePrepender func(string) string, linkHandler func(string, string) string) error ***REMOVED***
	for _, c := range cmd.Commands() ***REMOVED***
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() ***REMOVED***
			continue
		***REMOVED***
		if err := GenReSTTreeCustom(c, dir, filePrepender, linkHandler); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".rst"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := GenReSTCustom(cmd, f, linkHandler); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// adapted from: https://github.com/kr/text/blob/main/indent.go
func indentString(s, p string) string ***REMOVED***
	var res []byte
	b := []byte(s)
	prefix := []byte(p)
	bol := true
	for _, c := range b ***REMOVED***
		if bol && c != '\n' ***REMOVED***
			res = append(res, prefix...)
		***REMOVED***
		res = append(res, c)
		bol = c == '\n'
	***REMOVED***
	return string(res)
***REMOVED***
