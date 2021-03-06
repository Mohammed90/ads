// Copyright 2018 - 2019 Christian Müller <dev@c-mueller.xyz>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ads

import (
	"fmt"
	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

const default_Corefile = `ads`

const valid_Flag_Corefile = `ads {
  log
  default-lists
  disable-auto-update
}`

const valid_Target_Corefile = `ads {
  target 10.10.10.10
}`
const invalid_Target_MissingIP_Corefile = `ads {
  target
}`
const invalid_Target_InvalidIP_Corefile = `ads {
  target not-a-valid-ip
}`

const valid_List_Corefile = `ads {
  list %s/list.txt
}`
const invalid_List_MissingURL_Corefile = `ads {
  list
}`
const invalid_List_InvalidURL_Corefile = `ads {
  list not-a-valid-url
}`

const valid_Interval_Corefile = `ads {
  auto-update-interval 24h
}`
const invalid_Interval_MissingDuration_Corefile = `ads {
  auto-update-interval
}`
const invalid_Interval_InvalidDuration_Corefile = `ads {
  auto-update-interval not-a-parsable-duration-string
}`

const valid_Whitelist_Single = `ads {
  whitelist test.com
}`
const valid_Whitelist_Multi = `ads {
  whitelist test.com
  whitelist test.xyz
}`
const valid_Blacklist_Single = `ads {
  blacklist test.com
}`
const valid_Blacklist_Multi = `ads {
  blacklist test.com
  blacklist test.xyz
}`

const invalid_Whitelist_Single = `ads {
  whitelist
}`
const invalid_Whitelist_Multi = `ads {
  whitelist test.com
  whitelist
}`
const invalid_Blacklist_Single = `ads {
  blacklist
}`
const invalid_Blacklist_Multi = `ads {
  blacklist test.com
  blacklist
}`

const valid_Regex_Whitelist_Single = `ads {
  whitelist-regex (^|\.)local\.c-mueller\.de$
}`
const valid_Regex_Whitelist_Multi = `ads {
  whitelist-regex (^|\.)local\.c-mueller\.de$
  whitelist-regex (^|\.)local\.c-mueller\.xyz$
}`
const valid_Regex_Blacklist_Single = `ads {
  blacklist-regex (^|\.)local\.c-mueller\.de$
}`
const valid_Regex_Blacklist_Multi = `ads {
  blacklist-regex (^|\.)local\.c-mueller\.de$
  blacklist-regex (^|\.)local\.c-mueller\.xyz$
}`

const invalid_Regex_Whitelist_Single = `ads {
  whitelist-regex
}`
const invalid_Regex_Whitelist_Multi = `ads {
  whitelist-regex (^|\.)local\.c-mueller\.de$
  whitelist-regex
}`
const invalid_Regex_Blacklist_Single = `ads {
  blacklist-regex
}`
const invalid_Regex_Blacklist_Multi = `ads {
  blacklist-regex (^|\.)local\.c-mueller\.de$
  blacklist-regex
}`

func TestSetup_Initialisation(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", default_Corefile)
	assert.NoError(t, setup(c))
}

func TestSetup_Defaults(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", default_Corefile)
	assert.NoError(t, setup(c))
}

func TestValidFlags(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", valid_Flag_Corefile)
	assert.NoError(t, setup(c))
}

func TestSetup_ValidInterval(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", valid_Interval_Corefile)
	assert.NoError(t, setup(c))
}

func TestSetup_InvalidInterval(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", invalid_Interval_InvalidDuration_Corefile)
	assert.Error(t, setup(c))

	c = caddy.NewTestController("dns", invalid_Interval_MissingDuration_Corefile)
	assert.Error(t, setup(c))
}

func TestSetup_ValidTarget(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", valid_Target_Corefile)
	assert.NoError(t, setup(c))
}

func TestSetup_InvalidTarget(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", invalid_Target_InvalidIP_Corefile)
	assert.Error(t, setup(c))

	c = caddy.NewTestController("dns", invalid_Target_MissingIP_Corefile)
	assert.Error(t, setup(c))
}

func TestSetup_ValidList(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	srv := initTestServer(t)
	defer srv.Close()

	c := caddy.NewTestController("dns", fmt.Sprintf(valid_List_Corefile, srv.URL))
	assert.NoError(t, setup(c))
}

func TestSetup_InvalidList(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	c := caddy.NewTestController("dns", invalid_List_MissingURL_Corefile)
	assert.Error(t, setup(c))

	c = caddy.NewTestController("dns", invalid_List_InvalidURL_Corefile)
	assert.Error(t, setup(c))
}

func TestSetup_ValidWhiteAndBlacklist(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	cfs := []string{
		valid_Blacklist_Multi,
		valid_Blacklist_Single,
		valid_Whitelist_Multi,
		valid_Whitelist_Single,
	}

	for _, v := range cfs {
		c := caddy.NewTestController("dns", v)
		assert.NoError(t, setup(c))
	}
}

func TestSetup_InvalidWhiteAndBlacklist(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	cfs := []string{
		invalid_Blacklist_Multi,
		invalid_Blacklist_Single,
		invalid_Whitelist_Multi,
		invalid_Whitelist_Single,
	}

	for _, v := range cfs {
		c := caddy.NewTestController("dns", v)
		assert.Error(t, setup(c))
	}
}

func TestSetup_ValidRegexpWhiteAndBlacklist(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	cfs := []string{
		valid_Regex_Blacklist_Multi,
		valid_Regex_Blacklist_Single,
		valid_Regex_Whitelist_Multi,
		valid_Regex_Whitelist_Single,
	}

	for _, v := range cfs {
		c := caddy.NewTestController("dns", v)
		assert.NoError(t, setup(c))
	}
}

func TestSetup_InvalidRegexpWhiteAndBlacklist(t *testing.T) {
	s := updateDefaultBlocklists(t)
	defer s.Close()

	cfs := []string{
		invalid_Regex_Blacklist_Multi,
		invalid_Regex_Blacklist_Single,
		invalid_Regex_Whitelist_Multi,
		invalid_Regex_Whitelist_Single,
	}

	for _, v := range cfs {
		c := caddy.NewTestController("dns", v)
		assert.Error(t, setup(c))
	}
}

func updateDefaultBlocklists(t *testing.T) *httptest.Server {
	srv := initTestServer(t)

	defaultBlocklists = []string{fmt.Sprintf("%s/my-test-list.txt", srv.URL)}

	return srv
}
