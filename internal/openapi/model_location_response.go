/*
Orchestra API

Микросервис API для \"Клуба друзей оркестра\". **Все пользователи считаются равными**, а доступ из внешнего мира осуществляется через Telegram-бот.

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the LocationResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LocationResponse{}

// LocationResponse struct for LocationResponse
type LocationResponse struct {
	Id       *int32  `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Route    *string `json:"route,omitempty"`
	Features *string `json:"features,omitempty"`
}

// NewLocationResponse instantiates a new LocationResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLocationResponse() *LocationResponse {
	this := LocationResponse{}
	return &this
}

// NewLocationResponseWithDefaults instantiates a new LocationResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLocationResponseWithDefaults() *LocationResponse {
	this := LocationResponse{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *LocationResponse) GetId() int32 {
	if o == nil || IsNil(o.Id) {
		var ret int32
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LocationResponse) GetIdOk() (*int32, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *LocationResponse) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int32 and assigns it to the Id field.
func (o *LocationResponse) SetId(v int32) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *LocationResponse) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LocationResponse) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *LocationResponse) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *LocationResponse) SetName(v string) {
	o.Name = &v
}

// GetRoute returns the Route field value if set, zero value otherwise.
func (o *LocationResponse) GetRoute() string {
	if o == nil || IsNil(o.Route) {
		var ret string
		return ret
	}
	return *o.Route
}

// GetRouteOk returns a tuple with the Route field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LocationResponse) GetRouteOk() (*string, bool) {
	if o == nil || IsNil(o.Route) {
		return nil, false
	}
	return o.Route, true
}

// HasRoute returns a boolean if a field has been set.
func (o *LocationResponse) HasRoute() bool {
	if o != nil && !IsNil(o.Route) {
		return true
	}

	return false
}

// SetRoute gets a reference to the given string and assigns it to the Route field.
func (o *LocationResponse) SetRoute(v string) {
	o.Route = &v
}

// GetFeatures returns the Features field value if set, zero value otherwise.
func (o *LocationResponse) GetFeatures() string {
	if o == nil || IsNil(o.Features) {
		var ret string
		return ret
	}
	return *o.Features
}

// GetFeaturesOk returns a tuple with the Features field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LocationResponse) GetFeaturesOk() (*string, bool) {
	if o == nil || IsNil(o.Features) {
		return nil, false
	}
	return o.Features, true
}

// HasFeatures returns a boolean if a field has been set.
func (o *LocationResponse) HasFeatures() bool {
	if o != nil && !IsNil(o.Features) {
		return true
	}

	return false
}

// SetFeatures gets a reference to the given string and assigns it to the Features field.
func (o *LocationResponse) SetFeatures(v string) {
	o.Features = &v
}

func (o LocationResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LocationResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Route) {
		toSerialize["route"] = o.Route
	}
	if !IsNil(o.Features) {
		toSerialize["features"] = o.Features
	}
	return toSerialize, nil
}

type NullableLocationResponse struct {
	value *LocationResponse
	isSet bool
}

func (v NullableLocationResponse) Get() *LocationResponse {
	return v.value
}

func (v *NullableLocationResponse) Set(val *LocationResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableLocationResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableLocationResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLocationResponse(val *LocationResponse) *NullableLocationResponse {
	return &NullableLocationResponse{value: val, isSet: true}
}

func (v NullableLocationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLocationResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
