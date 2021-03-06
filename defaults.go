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
	"net"
	"time"
)

var defaultBlocklists = []string{
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
	"https://mirror1.malwaredomains.com/files/justdomains",
	"http://sysctl.org/cameleon/hosts",
	"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist",
	"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt",
	"https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt",
	"https://hosts-file.net/ad_servers.txt",
}

const defaultIPv4ResolutionIP = "127.0.0.1"
const defaultIPv6ResolutionIP = "::1"

var defaultConfigWithoutRules = adsPluginConfig{
	BlocklistURLs:       []string{},
	BlacklistRules:      []string{},
	WhitelistRules:      []string{},
	RegexBlacklistRules: []string{},
	RegexWhitelistRules: []string{},

	TargetIP:   net.ParseIP(defaultIPv4ResolutionIP),
	TargetIPv6: net.ParseIP(defaultIPv6ResolutionIP),

	EnableAutoUpdate:              true,
	BlocklistRenewalInterval:      time.Hour * 24,
	BlocklistRenewalRetryInterval: time.Minute,
	BlocklistRenewalRetryCount:    5,

	EnableBlocklistPersistence: false,
	BlocklistPersistencePath:   "",
}