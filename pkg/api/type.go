package api

import 	"time"

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


