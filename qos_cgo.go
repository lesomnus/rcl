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

// #include "rmw/rmw.h"
import "C"
import (
	"time"
)

func (p *Qos) toCStruct(dst *C.rmw_qos_profile_t) {
	dst.history = uint32(p.History)
	dst.depth = C.size_t(p.Depth)
	dst.reliability = uint32(p.Reliability)
	dst.durability = uint32(p.Durability)
	dst.deadline = C.rmw_time_t{nsec: C.uint64_t(p.Deadline)}
	dst.lifespan = C.rmw_time_t{nsec: C.uint64_t(p.Lifespan)}
	dst.liveliness = uint32(p.Liveliness)
	dst.liveliness_lease_duration = C.rmw_time_t{nsec: C.uint64_t(p.LivelinessLeaseDuration)}
	dst.avoid_ros_namespace_conventions = C.bool(p.AvoidRosNamespaceConventions)
}

func (p *Qos) fromCStruct(src *C.rmw_qos_profile_t) {
	p.History = HistoryPolicy(src.history)
	p.Depth = int(src.depth)
	p.Reliability = ReliabilityPolicy(src.reliability)
	p.Durability = DurabilityPolicy(src.durability)
	p.Deadline = time.Duration(src.deadline.sec)*time.Second + time.Duration(src.deadline.nsec)
	p.Lifespan = time.Duration(src.lifespan.sec)*time.Second + time.Duration(src.lifespan.nsec)
	p.Liveliness = LivelinessPolicy(src.liveliness)
	p.LivelinessLeaseDuration = time.Duration(src.liveliness_lease_duration.sec)*time.Second + time.Duration(src.liveliness_lease_duration.nsec)
	p.AvoidRosNamespaceConventions = bool(src.avoid_ros_namespace_conventions)
}
