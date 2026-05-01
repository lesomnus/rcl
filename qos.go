// Copied from
//
//	https://github.com/tiiuae/rclgo/blob/508dd42245daa2f122ce9d35e16d704d84927fa7/pkg/rclgo/qos.go
/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rcl

import (
	"math"
	"time"
)

type Qos struct {
	History                      HistoryPolicy
	Depth                        int
	Reliability                  ReliabilityPolicy
	Durability                   DurabilityPolicy
	Deadline                     time.Duration
	Lifespan                     time.Duration
	Liveliness                   LivelinessPolicy
	LivelinessLeaseDuration      time.Duration
	AvoidRosNamespaceConventions bool
}

func NewQos() *Qos {
	return &Qos{
		History:                      HistoryKeepLast,
		Depth:                        10,
		Reliability:                  ReliabilityReliable,
		Durability:                   DurabilityVolatile,
		Deadline:                     DeadlineDefault,
		Lifespan:                     LifespanDefault,
		Liveliness:                   LivelinessSystemDefault,
		LivelinessLeaseDuration:      LivelinessLeaseDurationDefault,
		AvoidRosNamespaceConventions: false,
	}
}

const (
	DurationInfinite    = time.Duration(math.MaxInt64)
	DurationUnspecified = time.Duration(0)

	DeadlineDefault                = DurationUnspecified
	LifespanDefault                = DurationUnspecified
	LivelinessLeaseDurationDefault = DurationUnspecified
)

type HistoryPolicy int

const (
	HistorySystemDefault HistoryPolicy = iota
	HistoryKeepLast
	HistoryKeepAll
	HistoryUnknown
)

type ReliabilityPolicy int

const (
	ReliabilitySystemDefault ReliabilityPolicy = iota
	ReliabilityReliable
	ReliabilityBestEffort
	ReliabilityUnknown
)

type DurabilityPolicy int

const (
	DurabilitySystemDefault DurabilityPolicy = iota
	DurabilityTransientLocal
	DurabilityVolatile
	DurabilityUnknown
)

type LivelinessPolicy int

const (
	LivelinessSystemDefault LivelinessPolicy = iota
	LivelinessAutomatic
	_
	LivelinessManualByTopic
	LivelinessUnknown
)
