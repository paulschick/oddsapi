// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

const DefaultSports = "upcoming"

type ParameterValue interface {
	Valid() bool
	String() string
}

func ParamPointer(p ParameterValue) *ParameterValue {
	return &p
}

type OddsFormat string

const (
	AmericanOddsFormat OddsFormat = "american"
	DecimalOddsFormat  OddsFormat = "decimal"
	DefaultOddsFormat             = DecimalOddsFormat
)

func (o OddsFormat) Valid() bool {
	switch o {
	case AmericanOddsFormat, DecimalOddsFormat:
		return true
	}
	return false
}

func (o OddsFormat) String() string {
	return string(o)
}

type DateFormat string

const (
	DateFormatIso     DateFormat = "iso"
	DefaultDateFormat            = DateFormatIso
	DateFormatUnix    DateFormat = "unix"
)

func (d DateFormat) Valid() bool {
	switch d {
	case DateFormatIso, DateFormatUnix:
		return true
	}
	return false
}

func (d DateFormat) String() string {
	return string(d)
}

type Region string

const (
	RegionUs        Region = "us"
	DefaultRegion          = RegionUs
	RegionUs2       Region = "us2"
	RegionUk        Region = "uk"
	RegionAustralia Region = "au"
	RegionEurope    Region = "eu"
)

func (r Region) Valid() bool {
	switch r {
	case RegionUs, RegionUs2, RegionUk, RegionAustralia, RegionEurope:
		return true
	}
	return false
}

func (r Region) String() string {
	return string(r)
}

type MarketKey string

const (
	MarketH2H       MarketKey = "h2h"
	MarketSpreads   MarketKey = "spreads"
	MarketTotals    MarketKey = "totals"
	MarketOutrights MarketKey = "outrights"
)

func (m MarketKey) Valid() bool {
	switch m {
	case MarketH2H, MarketSpreads, MarketTotals, MarketOutrights:
		return true
	}
	return false
}

func (m MarketKey) String() string {
	return string(m)
}
