package api

import 	"time"


/*
//Sample stream:

 "streamName": "USA_SD_NAT_4195_0_7503892744946620000",
    "streamId": "7503892744946620001",
    "expiry": "2018-01-29T20:00:00Z",
    "Transports": [{
        "url": "http://ccr.linear-nat-dash.xcr.comcast.net/dash/USA_SD_NAT_4183_0_7503892744946620183/USA_SD_NAT_4183_0_7503892744946620163_DASH.mpd",
        "avgBitrate": 1.875,
        "maxBitrate": 1.875
    }],
    "archiveConfig": {
        "archiveTime": 72,
        "reArchiveTime": 24,
        "archivalStartTime": "02:00AM",
        "archivalDuration": 10,
        "archivalPause": 5
    }
*/

type StreamCfg struct {
	StreamName    string          `json:"streamName"`
	StreamID      string          `json:"streamId"`
	Expiry        string          `json:"expiry,omitempty"`
	State         uint8           `json:"state, omitempty"`
	Transports    [1]TransportCfg `json:"Transports"`
	ArchiveConfig ArchiveCfg      `json:"archiveConfig"`
	NoDeDupe      bool            `json:"noDeDupe,omitempty"`
}

type ArchiveCfg struct {
	ArchiveTime       uint32 `json:"archiveTime"`
	ReArchiveTime     uint32 `json:"reArchiveTime"`
	ArchivalStartTime string `json:"archivalStartTime"`
	ArchivalDuration  uint32 `json:"archivalDuration"`
	ArchivalPause     uint32 `json:"archivalPause"`
}

type TransportCfg struct {
	URL        string  `json:"url"`
	AvgBitrate float32 `json:"avgBitrate"`
	MaxBitrate float32 `json:"maxBitrate"`
}



type Endpoints struct {
	Subsets []EndpointSubset `json:"subsets"`
}

type EndpointSubset struct {
	// IP addresses which offer the related ports that are marked as ready. These endpoints
	// should be considered safe for load balancers and clients to utilize.
	Addresses []EndpointAddress `json:"addresses,omitempty"`
	// IP addresses which offer the related ports but are not currently marked as ready
	// because they have not yet finished starting, have recently failed a readiness check,
	// or have recently failed a liveness check.
	NotReadyAddresses []EndpointAddress `json:"notReadyAddresses,omitempty"`
	// Port numbers available on the related IP addresses.
	Ports []EndpointPort `json:"ports,omitempty"`
}

// EndpointAddress is a tuple that describes single IP address.
type EndpointAddress struct {
	// The IP of this endpoint.
	// May not be loopback (127.0.0.0/8), link-local (169.254.0.0/16),
	// or link-local multicast ((224.0.0.0/24).
	// TODO: This should allow hostname or IP, See #4447.
	IP string `json:"ip"`

	// Reference to object providing the endpoint.
	TargetRef *ObjectReference `json:"targetRef,omitempty"`
}

// EndpointPort is a tuple that describes a single port.
type EndpointPort struct {
	// The name of this port (corresponds to ServicePort.Name).
	// Must be a DNS_LABEL.
	// Optional only if one port is defined.
	Name string `json:"name,omitempty"`

	// The port number of the endpoint.
	Port int32 `json:"port"`

	// The IP protocol for this port.
	// Must be UDP or TCP.
	// Default is TCP.
	//Protocol Protocol `json:"protocol,omitempty"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// Kind of the referent.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#types-kinds
	Kind string `json:"kind,omitempty"`
	// Namespace of the referent.
	// More info: http://releases.k8s.io/release-1.2/docs/user-guide/namespaces.md
	Namespace string `json:"namespace,omitempty"`
	// Name of the referent.
	// More info: http://releases.k8s.io/release-1.2/docs/user-guide/identifiers.md#names
	Name string `json:"name,omitempty"`
	// UID of the referent.
	// More info: http://releases.k8s.io/release-1.2/docs/user-guide/identifiers.md#uids
	//UID types.UID `json:"uid,omitempty"`
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
	// Specific resourceVersion to which this reference is made, if any.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#concurrency-control-and-consistency
	ResourceVersion string `json:"resourceVersion,omitempty"`

	// If referring to a piece of an object instead of an entire object, this string
	// should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].
	// For example, if the object reference is to a container within a pod, this would take on a value like:
	// "spec.containers{name}" (where "name" refers to the name of the container that triggered
	// the event) or if no container name is specified "spec.containers[2]" (container with
	// index 2 in this pod). This syntax is chosen only to have some well-defined way of
	// referencing a part of an object.
	// TODO: this design is not final and this field is subject to change in the future.
	FieldPath string `json:"fieldPath,omitempty"`
}



//Manifest Agent

// MAStatus of manifest fetch operations.
type MAStatus struct {
	StartTime time.Time // time when application starts running

	MpdURL               string        `json:"MpdURL,omitempty"`// URL for fetching manifests
	StreamID             string        `json:"StreamID,omitempty"`// The stream we are processing
	ISID                 uint64        `json:"ISID,omitempty"` // The stream ID we are processing
	MpdLastRequestedTime time.Time     `json:"MpdLastRequestedTime,omitempty"`// For tracking mpd polling
	MpdLastReceivedTime  time.Time     `json:"MpdLastReceivedTime,omitempty"`// For tracking mpd polling
	MpdPollIntervalMax   time.Duration `json:"MpdPollIntervalMax,omitempty"`// link to config parameter
	MpdPollIntervalMin   time.Duration `json:"MpdPollIntervalMin,omitempty"`// link to config parameter
	MpdLivePoint         time.Time     `json:"MpdLivePoint,omitempty"`// time of most recent segment in mpd
	MpdLastPublishTime   time.Time     `json:"MpdLastPublishTime,omitempty"`// publish time from most recent mpd
	MpdPollSuccessCount  int           `json:"MpdPollSuccessCount,omitempty"`// number of good polls since start
	MpdPollFailureCount  int           `json:"MpdPollFailureCount,omitempty"`// number of bad polls since start
	MpdPollLastError     string        `json:"MpdPollLastError,omitempty"`// last error from fetching MPD
	MpdLivePointDrift    time.Duration `json:"MpdLivePointDrift,omitempty"`// Time difference between livepoint and system time (positive means LP's livepoint is ahead)
	MpdDriftCorrect      time.Duration `json:"MpdDriftCorrect,omitempty"`// The drift correction being applied to the manifests (drift in wall clock time seen within the manifests)

	ActiveRecordings        int `json:"ActiveRecordings,omitempty"`// number of active recordings
	ActiveRecordingsWithErr int `json:"ActiveRecordingsWithErr,omitempty"`// number of active recordings with error set
	ActiveBatches           int `json:"ActiveBatches,omitempty"`// number of active batches

	MpdDropped     int `json:"MpdDropped,omitempty"`// number of manifest dropped because it contained errors
	ContentDropped int `json:"ContentDropped,omitempty"`// number of times content has been dropped from a manifest because it contained errors
}


// ServiceList holds a list of services.
type ServiceList struct {
	// List of services
	Items []Service `json:"items" protobuf:"bytes,2,rep,name=items"`
}


// Service is a named abstraction of software service (for example, mysql) consisting of local port
// (for example 3306) that the proxy listens on, and the selector that determines which pods
// will answer requests sent through the proxy.
type Service struct {

	ObjectMeta
	// Spec defines the behavior of a service.
	// https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Spec ServiceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

}


// ServiceSpec describes the attributes that a user creates on a service.
type ServiceSpec struct {
	// Route service traffic to pods with label keys and values matching this
	// selector. If empty or not present, the service is assumed to have an
	// external process managing its endpoints, which Kubernetes will not
	// modify. Only applies to types ClusterIP, NodePort, and LoadBalancer.
	// Ignored if type is ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/
	// +optional
	Selector map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`

	// clusterIP is the IP address of the service and is usually assigned
	// randomly by the master. If an address is specified manually and is not in
	// use by others, it will be allocated to the service; otherwise, creation
	// of the service will fail. This field can not be changed through updates.
	// Valid values are "None", empty string (""), or a valid IP address. "None"
	// can be specified for headless services when proxying is not required.
	// Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if
	// type is ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	ClusterIP string `json:"clusterIP,omitempty" protobuf:"bytes,3,opt,name=clusterIP"`

}

// TypeMeta describes an individual object in an API response or request
// with strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// users must create.
type ObjectMeta struct {
	// Name must be unique within a namespace. Is required when creating resources, although
	// some resources may allow a client to request the generation of an appropriate name
	// automatically. Name is primarily intended for creation idempotence and configuration
	// definition.
	// Cannot be updated.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// GenerateName is an optional prefix, used by the server, to generate a unique
	// name ONLY IF the Name field has not been provided.
	// If this field is used, the name returned to the client will be different
	// than the name passed. This value will also be combined with a unique suffix.
	// The provided value has the same validation rules as the Name field,
	// and may be truncated by the length of the suffix required to make the value
	// unique on the server.
	//
	// If this field is specified and the generated name exists, the server will
	// NOT return a 409 - instead, it will either return 201 Created or 500 with Reason
	// ServerTimeout indicating a unique name could not be found in the time allotted, and the client
	// should retry (optionally after the time indicated in the Retry-After header).
	//
	// Applied only if Name is not specified.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#idempotency
	// +optional
	GenerateName string `json:"generateName,omitempty" protobuf:"bytes,2,opt,name=generateName"`

	// Namespace defines the space within each name must be unique. An empty namespace is
	// equivalent to the "default" namespace, but "default" is the canonical representation.
	// Not all objects are required to be scoped to a namespace - the value of this field for
	// those objects will be empty.
	//
	// Must be a DNS_LABEL.
	// Cannot be updated.
	// More info: http://kubernetes.io/docs/user-guide/namespaces
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`

	// SelfLink is a URL representing this object.
	// Populated by the system.
	// Read-only.
	// +optional
	SelfLink string `json:"selfLink,omitempty" protobuf:"bytes,4,opt,name=selfLink"`

	// An opaque value that represents the internal version of this object that can
	// be used by clients to determine when objects have changed. May be used for optimistic
	// concurrency, change detection, and the watch operation on a resource or set of resources.
	// Clients must treat these values as opaque and passed unmodified back to the server.
	// They may only be valid for a particular resource or set of resources.
	//
	// Populated by the system.
	// Read-only.
	// Value must be treated as opaque by clients and .
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#concurrency-control-and-consistency
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,6,opt,name=resourceVersion"`

	// A sequence number representing a specific generation of the desired state.
	// Populated by the system. Read-only.
	// +optional
	Generation int64 `json:"generation,omitempty" protobuf:"varint,7,opt,name=generation"`

	// Number of seconds allowed for this object to gracefully terminate before
	// it will be removed from the system. Only set when deletionTimestamp is also set.
	// May only be shortened.
	// Read-only.
	// +optional
	DeletionGracePeriodSeconds *int64 `json:"deletionGracePeriodSeconds,omitempty" protobuf:"varint,10,opt,name=deletionGracePeriodSeconds"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`


	// Must be empty before the object is deleted from the registry. Each entry
	// is an identifier for the responsible component that will remove the entry
	// from the list. If the deletionTimestamp of the object is non-nil, entries
	// in this list can only be removed.
	// +optional
	// +patchStrategy=merge
	Finalizers []string `json:"finalizers,omitempty" patchStrategy:"merge" protobuf:"bytes,14,rep,name=finalizers"`

	// The name of the cluster which the object belongs to.
	// This is used to distinguish resources with same name and namespace in different clusters.
	// This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
	// +optional
	ClusterName string `json:"clusterName,omitempty" protobuf:"bytes,15,opt,name=clusterName"`
}

