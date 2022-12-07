package core

// Field path constants that are specific to the internal API
// representation.
const (
	NodeUnschedulableField = "spec.unschedulable"
	ObjectNameField        = "metadata.name"
	PodHostField           = "spec.nodeName"
	PodStatusField         = "status.phase"
	SecretTypeField        = "type"

	EventReasonField                  = "action"
	EventSourceField                  = "reportingComponent"
	EventTypeField                    = "type"
	EventInvolvedKindField            = "involvedObject.kind"
	EventInvolvedNamespaceField       = "involvedObject.namespace"
	EventInvolvedNameField            = "involvedObject.name"
	EventInvolvedUIDField             = "involvedObject.uid"
	EventInvolvedAPIVersionField      = "involvedObject.apiVersion"
	EventInvolvedResourceVersionField = "involvedObject.resourceVersion"
	EventInvolvedFieldPathField       = "involvedObject.fieldPath"
)
