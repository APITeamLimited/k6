/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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

package har

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/pretty"
	"net/url"
	"sort"
	"strings"
)

func Convert(h HAR, enableChecks bool, returnOnFailedCheck bool, batchTime uint, nobatch bool, correlate bool, only, skip []string) (string, error) ***REMOVED***
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	if returnOnFailedCheck && !enableChecks ***REMOVED***
		return "", errors.Errorf("return on failed check requires --enable-status-code-checks")
	***REMOVED***

	if correlate && !nobatch ***REMOVED***
		return "", errors.Errorf("correlation requires --no-batch")
	***REMOVED***

	if enableChecks ***REMOVED***
		fmt.Fprint(w, "import ***REMOVED*** group, check, sleep ***REMOVED*** from 'k6';\n")
	***REMOVED*** else ***REMOVED***
		fmt.Fprint(w, "import ***REMOVED*** group, sleep ***REMOVED*** from 'k6';\n")
	***REMOVED***
	fmt.Fprint(w, "import http from 'k6/http';\n\n")

	fmt.Fprintf(w, "// Version: %v\n", h.Log.Version)
	fmt.Fprintf(w, "// Creator: %v\n", h.Log.Creator.Name)
	if h.Log.Browser != nil ***REMOVED***
		fmt.Fprintf(w, "// Browser: %v\n", h.Log.Browser.Name)
	***REMOVED***
	if h.Log.Comment != "" ***REMOVED***
		fmt.Fprintf(w, "// %v\n", h.Log.Comment)
	***REMOVED***

	// recordings include redirections as separate requests, and we dont want to trigger them twice
	fmt.Fprint(w, "\nexport let options = ***REMOVED*** maxRedirects: 0 ***REMOVED***;\n\n")

	fmt.Fprint(w, "export default function() ***REMOVED***\n\n")

	pages := h.Log.Pages
	sort.Sort(PageByStarted(pages))

	// Grouping by page and URL filtering
	pageEntries := make(map[string][]*Entry)
	for _, e := range h.Log.Entries ***REMOVED***

		// URL filtering
		u, err := url.Parse(e.Request.URL)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if !IsAllowedURL(u.Host, only, skip) ***REMOVED***
			continue
		***REMOVED***

		// Avoid multipart/form-data requests until k6 scripts can support binary data
		if e.Request.PostData != nil && strings.HasPrefix(e.Request.PostData.MimeType, "multipart/form-data") ***REMOVED***
			continue
		***REMOVED***

		// Create new group o adding page to a existing one
		if _, ok := pageEntries[e.Pageref]; !ok ***REMOVED***
			pageEntries[e.Pageref] = append([]*Entry***REMOVED******REMOVED***, e)
		***REMOVED*** else ***REMOVED***
			pageEntries[e.Pageref] = append(pageEntries[e.Pageref], e)
		***REMOVED***
	***REMOVED***

	for i, page := range pages ***REMOVED***

		entries := pageEntries[page.ID]
		fmt.Fprintf(w, "\tgroup(\"%s - %s\", function() ***REMOVED***\n", page.ID, page.Title)

		sort.Sort(EntryByStarted(entries))

		if nobatch ***REMOVED***
			var recordedRedirectURL string
			var recordedRestID string
			previousResponse := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***

			fmt.Fprint(w, "\t\tlet res, redirectUrl, json;\n")

			for entryIndex, e := range entries ***REMOVED***

				var params []string
				var cookies []string
				var body string

				fmt.Fprintf(w, "\t\t// Request #%d\n", entryIndex)

				if e.Request.PostData != nil ***REMOVED***
					body = e.Request.PostData.Text
				***REMOVED***

				for _, c := range e.Request.Cookies ***REMOVED***
					cookies = append(cookies, fmt.Sprintf(`%q: %q`, c.Name, c.Value))
				***REMOVED***
				if len(cookies) > 0 ***REMOVED***
					params = append(params, fmt.Sprintf("\"cookies\": ***REMOVED***\n\t\t\t\t%s\n\t\t\t***REMOVED***", strings.Join(cookies, ",\n\t\t\t\t\t")))
				***REMOVED***

				if headers := buildK6Headers(e.Request.Headers); len(headers) > 0 ***REMOVED***
					params = append(params, fmt.Sprintf("\"headers\": ***REMOVED***\n\t\t\t\t\t%s\n\t\t\t\t***REMOVED***", strings.Join(headers, ",\n\t\t\t\t\t")))
				***REMOVED***

				fmt.Fprintf(w, "\t\tres = http.%s(", strings.ToLower(e.Request.Method))

				if correlate && recordedRedirectURL != "" ***REMOVED***
					if recordedRedirectURL != e.Request.URL ***REMOVED***
						return "", errors.Errorf("The har file contained a redirect but the next request did not match that redirect. Possibly a misbehaving client or concurrent requests?")
					***REMOVED***
					fmt.Fprintf(w, "redirectUrl")
					recordedRedirectURL = ""
				***REMOVED*** else ***REMOVED***
					if recordedRestID != "" && strings.Contains(e.Request.URL, recordedRestID) ***REMOVED***
						fmt.Fprintf(w, "`%s`", strings.Replace(e.Request.URL, recordedRestID, "$***REMOVED***restID***REMOVED***", -1))
					***REMOVED*** else ***REMOVED***
						fmt.Fprintf(w, "%q", e.Request.URL)
					***REMOVED***
				***REMOVED***

				if e.Request.Method != "GET" ***REMOVED***
					if correlate && e.Request.PostData != nil && strings.Contains(e.Request.PostData.MimeType, "json") ***REMOVED***
						requestMap := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***

						escapedPostdata := strings.Replace(e.Request.PostData.Text, "$", "\\$", -1)

						if err := json.Unmarshal([]byte(escapedPostdata), &requestMap); err != nil ***REMOVED***
							return "", err
						***REMOVED***

						if len(previousResponse) != 0 ***REMOVED***
							traverseMaps(requestMap, previousResponse, nil)
						***REMOVED***
						requestText, err := json.Marshal(requestMap)
						if err == nil ***REMOVED***
							prettyJSONString := string(pretty.PrettyOptions(requestText, &pretty.Options***REMOVED***Width: 999999, Prefix: "\t\t\t", Indent: "\t", SortKeys: true***REMOVED***)[:])
							fmt.Fprintf(w, ",\n\t\t\t`%s`", strings.TrimSpace(prettyJSONString))
						***REMOVED*** else ***REMOVED***
							return "", err
						***REMOVED***

					***REMOVED*** else ***REMOVED***
						fmt.Fprintf(w, ",\n\t\t%q", body)
					***REMOVED***
				***REMOVED***

				if len(params) > 0 ***REMOVED***
					fmt.Fprintf(w, ",\n\t\t\t***REMOVED***\n\t\t\t\t%s\n\t\t\t***REMOVED***", strings.Join(params, ",\n\t\t\t"))
				***REMOVED***

				fmt.Fprintf(w, "\n\t\t)\n")

				if e.Response != nil ***REMOVED***
					// the response is nil if there is a failed request in the recording, or if responses were not recorded
					if enableChecks ***REMOVED***
						if e.Response.Status > 0 ***REMOVED***
							if returnOnFailedCheck ***REMOVED***
								fmt.Fprintf(w, "\t\tif (!check(res, ***REMOVED***\"status is %v\": (r) => r.status === %v ***REMOVED***)) ***REMOVED*** return ***REMOVED***;\n", e.Response.Status, e.Response.Status)
							***REMOVED*** else ***REMOVED***
								fmt.Fprintf(w, "\t\tcheck(res, ***REMOVED***\"status is %v\": (r) => r.status === %v ***REMOVED***);\n", e.Response.Status, e.Response.Status)
							***REMOVED***
						***REMOVED***
					***REMOVED***

					responseMimeType := e.Response.Content.MimeType
					if correlate &&
						strings.Index(responseMimeType, "application/") == 0 &&
						strings.Index(responseMimeType, "json") == len(responseMimeType)-4 ***REMOVED***
						if err := json.Unmarshal([]byte(e.Response.Content.Text), &previousResponse); err != nil ***REMOVED***
							return "", err
						***REMOVED***
						fmt.Fprint(w, "\t\tjson = JSON.parse(res.body);\n")
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			batches := SplitEntriesInBatches(entries, batchTime)

			fmt.Fprint(w, "\t\tlet req, res;\n")

			for j, batchEntries := range batches ***REMOVED***

				fmt.Fprint(w, "\t\treq = [")
				for k, e := range batchEntries ***REMOVED***
					r, err := buildK6RequestObject(e.Request)
					if err != nil ***REMOVED***
						return "", err
					***REMOVED***
					fmt.Fprintf(w, "%v", r)
					if k != len(batchEntries)-1 ***REMOVED***
						fmt.Fprint(w, ",")
					***REMOVED***
				***REMOVED***
				fmt.Fprint(w, "];\n")
				fmt.Fprint(w, "\t\tres = http.batch(req);\n")

				if enableChecks ***REMOVED***
					for k, e := range batchEntries ***REMOVED***
						if e.Response.Status > 0 ***REMOVED***
							if returnOnFailedCheck ***REMOVED***
								fmt.Fprintf(w, "\t\tif (!check(res, ***REMOVED***\"status is %v\": (r) => r.status === %v ***REMOVED***)) ***REMOVED*** return ***REMOVED***;\n", e.Response.Status, e.Response.Status)
							***REMOVED*** else ***REMOVED***
								fmt.Fprintf(w, "\t\tcheck(res[%v], ***REMOVED***\"status is %v\": (r) => r.status === %v ***REMOVED***);\n", k, e.Response.Status, e.Response.Status)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***

				if j != len(batches)-1 ***REMOVED***
					lastBatchEntry := batchEntries[len(batchEntries)-1]
					firstBatchEntry := batches[j+1][0]
					t := firstBatchEntry.StartedDateTime.Sub(lastBatchEntry.StartedDateTime).Seconds()
					fmt.Fprintf(w, "\t\tsleep(%.2f);\n", t)
				***REMOVED***
			***REMOVED***

			if i == len(pages)-1 ***REMOVED***
				// Last page; add random sleep time at the group completion
				fmt.Fprint(w, "\t\t// Random sleep between 2s and 4s\n")
				fmt.Fprint(w, "\t\tsleep(Math.floor(Math.random()*3+2));\n")
			***REMOVED*** else ***REMOVED***
				// Add sleep time at the end of the group
				nextPage := pages[i+1]
				lastEntry := entries[len(entries)-1]
				t := nextPage.StartedDateTime.Sub(lastEntry.StartedDateTime).Seconds()
				if t < 0.01 ***REMOVED***
					t = 0.5
				***REMOVED***
				fmt.Fprintf(w, "\t\tsleep(%.2f);\n", t)
			***REMOVED***
		***REMOVED***

		fmt.Fprint(w, "\t***REMOVED***);\n")
	***REMOVED***

	fmt.Fprint(w, "\n***REMOVED***\n")
	if err := w.Flush(); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return b.String(), nil
***REMOVED***

func buildK6RequestObject(req *Request) (string, error) ***REMOVED***
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	fmt.Fprint(w, "***REMOVED***\n")

	method := strings.ToLower(req.Method)
	if method == "delete" ***REMOVED***
		method = "del"
	***REMOVED***
	fmt.Fprintf(w, `"method": %q, "url": %q`, method, req.URL)

	if req.PostData != nil && method != "get" ***REMOVED***
		postParams, plainText, err := buildK6Body(req)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED*** else if len(postParams) > 0 ***REMOVED***
			fmt.Fprintf(w, `, "body": ***REMOVED*** %s ***REMOVED***`, strings.Join(postParams, ", "))
		***REMOVED*** else if plainText != "" ***REMOVED***
			fmt.Fprintf(w, `, "body": %q`, plainText)
		***REMOVED***
	***REMOVED***

	var params []string
	var cookies []string
	for _, c := range req.Cookies ***REMOVED***
		cookies = append(cookies, fmt.Sprintf(`%q: %q`, c.Name, c.Value))
	***REMOVED***
	if len(cookies) > 0 ***REMOVED***
		params = append(params, fmt.Sprintf(`"cookies": ***REMOVED*** %s ***REMOVED***`, strings.Join(cookies, ", ")))
	***REMOVED***

	if headers := buildK6Headers(req.Headers); len(headers) > 0 ***REMOVED***
		params = append(params, fmt.Sprintf(`"headers": ***REMOVED*** %s ***REMOVED***`, strings.Join(headers, ", ")))
	***REMOVED***

	if len(params) > 0 ***REMOVED***
		fmt.Fprintf(w, `, "params": ***REMOVED*** %s ***REMOVED***`, strings.Join(params, ", "))
	***REMOVED***

	fmt.Fprint(w, "***REMOVED***")
	if err := w.Flush(); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var buffer bytes.Buffer
	err := json.Indent(&buffer, b.Bytes(), "\t\t", "\t")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return buffer.String(), nil
***REMOVED***

func buildK6Headers(headers []Header) []string ***REMOVED***
	var h []string
	if len(headers) > 0 ***REMOVED***
		m := make(map[string]Header)
		for _, header := range headers ***REMOVED***
			name := strings.ToLower(header.Name)
			_, exists := m[name]
			// Avoid SPDY's, duplicated or cookie headers
			if !exists && name[0] != ':' && name != "cookie" ***REMOVED***
				m[strings.ToLower(header.Name)] = header
				h = append(h, fmt.Sprintf("%q: %q", header.Name, header.Value))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return h
***REMOVED***

func buildK6Body(req *Request) ([]string, string, error) ***REMOVED***
	var postParams []string
	if req.PostData.MimeType == "application/x-www-form-urlencoded" && len(req.PostData.Params) > 0 ***REMOVED***
		for _, p := range req.PostData.Params ***REMOVED***
			n, err := url.QueryUnescape(p.Name)
			if err != nil ***REMOVED***
				return postParams, "", err
			***REMOVED***
			v, err := url.QueryUnescape(p.Value)
			if err != nil ***REMOVED***
				return postParams, "", err
			***REMOVED***
			postParams = append(postParams, fmt.Sprintf(`%q: %q`, n, v))
		***REMOVED***
		return postParams, "", nil
	***REMOVED***
	return postParams, req.PostData.Text, nil
***REMOVED***

func traverseMaps(request map[string]interface***REMOVED******REMOVED***, response map[string]interface***REMOVED******REMOVED***, path []interface***REMOVED******REMOVED***) ***REMOVED***
	if response == nil ***REMOVED***
		// previous call reached a leaf in the response map so there's no point continuing
		return
	***REMOVED***
	for key, val := range request ***REMOVED***
		responseVal := response[key]
		if responseVal == nil ***REMOVED***
			// no corresponding value in response map (and the type conversion below would fail so we need an early exit)
			continue
		***REMOVED***
		newPath := append(path, key)
		switch concreteVal := val.(type) ***REMOVED***
		case map[string]interface***REMOVED******REMOVED***:
			traverseMaps(concreteVal, responseVal.(map[string]interface***REMOVED******REMOVED***), newPath)
		case []interface***REMOVED******REMOVED***:
			traverseArrays(concreteVal, responseVal.([]interface***REMOVED******REMOVED***), newPath)
		default:
			if responseVal == val ***REMOVED***
				request[key] = jsObjectPath(newPath)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func traverseArrays(requestArray []interface***REMOVED******REMOVED***, responseArray []interface***REMOVED******REMOVED***, path []interface***REMOVED******REMOVED***) ***REMOVED***
	for i, val := range requestArray ***REMOVED***
		newPath := append(path, i)
		if len(responseArray) <= i ***REMOVED***
			// requestArray had more entries than responseArray
			break
		***REMOVED***
		responseVal := responseArray[i]
		switch concreteVal := val.(type) ***REMOVED***
		case map[string]interface***REMOVED******REMOVED***:
			traverseMaps(concreteVal, responseVal.(map[string]interface***REMOVED******REMOVED***), newPath)
		case []interface***REMOVED******REMOVED***:
			traverseArrays(concreteVal, responseVal.([]interface***REMOVED******REMOVED***), newPath)
		case string:
			if responseVal == val ***REMOVED***
				requestArray[i] = jsObjectPath(newPath)
			***REMOVED***
		default:
			panic(jsObjectPath(newPath))
		***REMOVED***
	***REMOVED***
***REMOVED***

func jsObjectPath(path []interface***REMOVED******REMOVED***) string ***REMOVED***
	s := "$***REMOVED***json"
	for _, val := range path ***REMOVED***
		// this may cause issues with non-array keys with numeric values. test this later.
		switch concreteVal := val.(type) ***REMOVED***
		case int:
			s = s + "[" + fmt.Sprint(concreteVal) + "]"
		case string:
			s = s + "." + concreteVal
		***REMOVED***
	***REMOVED***
	s = s + "***REMOVED***"
	return s
***REMOVED***
