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

package html

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
)

type HTML struct***REMOVED******REMOVED***

func New() *HTML ***REMOVED***
	return &HTML***REMOVED******REMOVED***
***REMOVED***

func (HTML) ParseHTML(ctx context.Context, src string) (Selection, error) ***REMOVED***
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(src))
	if err != nil ***REMOVED***
		return Selection***REMOVED******REMOVED***, err
	***REMOVED***
	return Selection***REMOVED***common.GetRuntime(ctx), doc.Selection***REMOVED***, nil
***REMOVED***

type Selection struct ***REMOVED***
	rt  *goja.Runtime
	sel *goquery.Selection
***REMOVED***

func (s Selection) Add(arg goja.Value) Selection ***REMOVED***
	switch val := arg.Export().(type) ***REMOVED***
	case Selection:
		return Selection***REMOVED***s.rt, s.sel.AddSelection(val.sel)***REMOVED***
	default:
		return Selection***REMOVED***s.rt, s.sel.Add(arg.String())***REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Find(sel string) Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Find(sel)***REMOVED***
***REMOVED***

func (s Selection) Text() string ***REMOVED***
	return s.sel.Text()
***REMOVED***

func (s Selection) Attr(name string, def ...goja.Value) goja.Value ***REMOVED***
	val, exists := s.sel.Attr(name)
	if !exists ***REMOVED***
		if len(def) > 0 ***REMOVED***
			return def[0]
		***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return s.rt.ToValue(val)
***REMOVED***

func (s Selection) Html() goja.Value ***REMOVED***
	val, err := s.sel.Html()
	if err != nil ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return s.rt.ToValue(val)
***REMOVED***

func optionVal(s *goquery.Selection) string ***REMOVED***
	val, exists := s.Attr("value")
	if exists ***REMOVED***
		return val
	***REMOVED***

	val, err := s.Html()
	if err != nil ***REMOVED***
		return ""
	***REMOVED***

	return val
***REMOVED***

func(s Selection) Val() goja.Value ***REMOVED***
	switch goquery.NodeName(s.sel) ***REMOVED***
		case "input":
			return s.Attr("value")

		case "textarea":
			return s.Html()

		case "button":
			return s.Attr("value")

		case "option":
			return s.rt.ToValue(optionVal(s.sel))

		case "select":
			selected := s.sel.First().Find("option[selected]")

			if _, exists := s.sel.Attr("multiple"); exists ***REMOVED***
				return s.rt.ToValue(selected.Map(func(idx int, opt *goquery.Selection) string ***REMOVED*** return optionVal(opt) ***REMOVED***))
			***REMOVED*** else ***REMOVED***
				return s.rt.ToValue(optionVal(selected))
			***REMOVED***

		default:
			return goja.Undefined()
	***REMOVED***
***REMOVED***

func (s Selection) Closest(selector string) Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Closest(selector)***REMOVED***
***REMOVED***

func (s Selection) Children(def ...string) Selection ***REMOVED***
	if len(def) == 0 ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Children()***REMOVED***
	***REMOVED*** else ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.ChildrenFiltered(def[0])***REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Contents() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Contents()***REMOVED***
***REMOVED***

func (s Selection) Each(v goja.Value) Selection ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if isFn ***REMOVED***
		fn := func(idx int, sel *goquery.Selection) ***REMOVED***
			gojaFn(v, s.rt.ToValue(idx), s.rt.ToValue(sel))
		***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Each(fn)***REMOVED***
	***REMOVED*** else ***REMOVED***
		s.rt.Interrupt("Argument to each() must be a function")
		return s
	***REMOVED***
***REMOVED***

func (s Selection) End() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.End()***REMOVED***
***REMOVED***

func (s Selection) buildMatcher(v goja.Value, gojaFn goja.Callable) func (int, *goquery.Selection) bool ***REMOVED***
	return func(idx int, sel *goquery.Selection) bool ***REMOVED***
		fnRes, fnErr := gojaFn(v, s.rt.ToValue(idx), s.rt.ToValue(sel))
		return fnErr == nil && fnRes.ToBoolean()
	***REMOVED***
***REMOVED***

func (s Selection) Filter(v goja.Value) Selection ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if isFn ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.FilterFunction(s.buildMatcher(v, gojaFn))***REMOVED***
	***REMOVED*** else ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Filter(v.String())***REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Is(v goja.Value) bool ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if isFn ***REMOVED***
		return s.sel.IsFunction(s.buildMatcher(v, gojaFn))
	***REMOVED*** else ***REMOVED***
		return s.sel.Is(v.String())
	***REMOVED***
***REMOVED***
