package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const DefaultFinalizer string = "io.apicurio.finalizer"

type ApicurioType string


const (
	ApicurioFullType         ApicurioType = "apicurio-full"
	ApicurioExternalAuthType ApicurioType = "apicurio-external-auth"
)

type MinMaxCfg struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type ApicurioResource struct {
	Route  string     `json:"route"`
	Url    string     `json:"url"`
	Memory *MinMaxCfg `json:"memory"`
	CPU    *MinMaxCfg `json:"cpu"`
	JVM    *MinMaxCfg `json:"jvm"`
}

type ApicurioAuthRresource struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Realm    string `json:"realm"`
	ApicurioResource
}

type ApicurioSpec struct {
	Version   string                 `json:"version"`
	Type      ApicurioType           `json:"type"`
	AppDomain string                 `json:"app_domain"`
	Auth      *ApicurioAuthRresource `json:"auth"`
	Api       *ApicurioResource      `json:"api"`
	WebSocket *ApicurioResource      `json:"websocket"`
	Studio    *ApicurioResource      `json:"studio"`
}

type ApicurioConditionType string

const (
	ApicurioNone      ApicurioConditionType = ""
	ApicurioNew       ApicurioConditionType = "New"
	ApicurioReconcile ApicurioConditionType = "Reconcile"
	ApicurioReady     ApicurioConditionType = "Ready"
	ApicurioDelete    ApicurioConditionType = "Delete"
)

// ApicurioStatus defines the observed state of Apicurio
type ApicurioStatus struct {
	Type    ApicurioConditionType `json:"type"`
	Reason  *string               `json:"reason"`
	Message *string               `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Apicurio is the Schema for the apicuriodeployments API
// +k8s:openapi-gen=true
type Apicurio struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApicurioSpec   `json:"spec,omitempty"`
	Status ApicurioStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApicurioList contains a list of Apicurio
type ApicurioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Apicurio `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Apicurio{}, &ApicurioList{})
}
