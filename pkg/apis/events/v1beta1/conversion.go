package v1beta1

import (
	v1beta1 "k8s.io/api/events/v1beta1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	k8s_api "k8s.io/kubernetes/pkg/apis/core"
	k8s_api_v1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

func Convert_v1beta1_Event_To_core_Event(in *v1beta1.Event, out *k8s_api.Event, s conversion.Scope) error {
	if err := autoConvert_v1beta1_Event_To_core_Event(in, out, s); err != nil {
		return err
	}
	if err := k8s_api_v1.Convert_v1_ObjectReference_To_core_ObjectReference(&in.Regarding, &out.InvolvedObject, s); err != nil {
		return err
	}
	if err := k8s_api_v1.Convert_v1_EventSource_To_core_EventSource(&in.DeprecatedSource, &out.Source, s); err != nil {
		return err
	}
	out.Message = in.Note
	out.FirstTimestamp = in.DeprecatedFirstTimestamp
	out.LastTimestamp = in.DeprecatedLastTimestamp
	out.Count = in.DeprecatedCount
	return nil
}

func Convert_core_Event_To_v1beta1_Event(in *k8s_api.Event, out *v1beta1.Event, s conversion.Scope) error {
	if err := autoConvert_core_Event_To_v1beta1_Event(in, out, s); err != nil {
		return err
	}
	if err := k8s_api_v1.Convert_core_ObjectReference_To_v1_ObjectReference(&in.InvolvedObject, &out.Regarding, s); err != nil {
		return err
	}
	if err := k8s_api_v1.Convert_core_EventSource_To_v1_EventSource(&in.Source, &out.DeprecatedSource, s); err != nil {
		return err
	}
	out.Note = in.Message
	out.DeprecatedFirstTimestamp = in.FirstTimestamp
	out.DeprecatedLastTimestamp = in.LastTimestamp
	out.DeprecatedCount = in.Count
	return nil
}
