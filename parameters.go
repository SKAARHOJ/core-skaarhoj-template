package main

import (
	ib "github.com/SKAARHOJ/ibeam-corelib-go"
	pb "github.com/SKAARHOJ/ibeam-corelib-go/ibeam-core"
	b "github.com/SKAARHOJ/ibeam-corelib-go/paramhelpers"
)

func configureParameters(r *ib.IBeamParameterRegistry) {

	/*
		// If not registered the connection parameter gets added automatically with ID 1
		r.RegisterParameter(&pb.ParameterDetail{
			Id:            &pb.ModelParameterID{Parameter: 1},
			Path:          "config",
			Name:          "connection",
			Label:         "Connected",
			Description:   "Connection status of device",
			GenericType:   pb.GenericType_ConnectionState,
			ControlStyle:  pb.ControlStyle_NoControl,
			FeedbackStyle: pb.FeedbackStyle_NormalFeedback,
			ValueType:     pb.ValueType_Binary,
		})
	*/

	r.RegisterParameter(&pb.ParameterDetail{
		Id:            &pb.ModelParameterID{Parameter: 2},
		Path:          "tests",
		Name:          "device_test",
		Label:         "Device Test",
		ShortLabel:    "Test",
		Description:   "Testfunction to identify your device, it will: ADD A DESCRIPTION HERE THAT EXPLAINS THE TEST ROUTINE",
		GenericType:   pb.GenericType_TestTrigger,
		ControlStyle:  pb.ControlStyle_Oneshot,
		FeedbackStyle: pb.FeedbackStyle_NoFeedback,
		ValueType:     pb.ValueType_NoValue,
	})

	r.RegisterParameterForModel(1, &pb.ParameterDetail{
		Id:             &pb.ModelParameterID{Parameter: 3},
		Path:           "tests/binarys",
		Name:           "normal_integer",
		Label:          "Normal Integer",
		ShortLabel:     "NormalBinary",
		Description:    "Testing Integers",
		ControlStyle:   pb.ControlStyle_Normal,
		FeedbackStyle:  pb.FeedbackStyle_NormalFeedback,
		ValueType:      pb.ValueType_Integer,
		RetryCount:     2,
		ControlDelayMs: 1,
		Minimum:        0,
		Maximum:        10,
		DefaultValue:   b.Int(5),
	})
}
