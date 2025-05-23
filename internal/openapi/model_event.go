/*
Orchestra API

Микросервис API для \"Клуба друзей оркестра\". **Все пользователи считаются равными**, а доступ из внешнего мира осуществляется через Telegram-бот.

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"time"
)

// checks if the Event type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Event{}

// Event struct for Event
type Event struct {
	Id          *int32     `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	EventType   *int32     `json:"event_type,omitempty"`
	EventDate   *time.Time `json:"event_date,omitempty"`
	Location    *int32     `json:"location,omitempty"`
	Capacity    *int32     `json:"capacity,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// NewEvent instantiates a new Event object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEvent() *Event {
	this := Event{}
	return &this
}

// NewEventWithDefaults instantiates a new Event object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEventWithDefaults() *Event {
	this := Event{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Event) GetId() int32 {
	if o == nil || IsNil(o.Id) {
		var ret int32
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetIdOk() (*int32, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Event) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int32 and assigns it to the Id field.
func (o *Event) SetId(v int32) {
	o.Id = &v
}

// GetTitle returns the Title field value if set, zero value otherwise.
func (o *Event) GetTitle() string {
	if o == nil || IsNil(o.Title) {
		var ret string
		return ret
	}
	return *o.Title
}

// GetTitleOk returns a tuple with the Title field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetTitleOk() (*string, bool) {
	if o == nil || IsNil(o.Title) {
		return nil, false
	}
	return o.Title, true
}

// HasTitle returns a boolean if a field has been set.
func (o *Event) HasTitle() bool {
	if o != nil && !IsNil(o.Title) {
		return true
	}

	return false
}

// SetTitle gets a reference to the given string and assigns it to the Title field.
func (o *Event) SetTitle(v string) {
	o.Title = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Event) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Event) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Event) SetDescription(v string) {
	o.Description = &v
}

// GetEventType returns the EventType field value if set, zero value otherwise.
func (o *Event) GetEventType() int32 {
	if o == nil || IsNil(o.EventType) {
		var ret int32
		return ret
	}
	return *o.EventType
}

// GetEventTypeOk returns a tuple with the EventType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetEventTypeOk() (*int32, bool) {
	if o == nil || IsNil(o.EventType) {
		return nil, false
	}
	return o.EventType, true
}

// HasEventType returns a boolean if a field has been set.
func (o *Event) HasEventType() bool {
	if o != nil && !IsNil(o.EventType) {
		return true
	}

	return false
}

// SetEventType gets a reference to the given int32 and assigns it to the EventType field.
func (o *Event) SetEventType(v int32) {
	o.EventType = &v
}

// GetEventDate returns the EventDate field value if set, zero value otherwise.
func (o *Event) GetEventDate() time.Time {
	if o == nil || IsNil(o.EventDate) {
		var ret time.Time
		return ret
	}
	return *o.EventDate
}

// GetEventDateOk returns a tuple with the EventDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetEventDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.EventDate) {
		return nil, false
	}
	return o.EventDate, true
}

// HasEventDate returns a boolean if a field has been set.
func (o *Event) HasEventDate() bool {
	if o != nil && !IsNil(o.EventDate) {
		return true
	}

	return false
}

// SetEventDate gets a reference to the given time.Time and assigns it to the EventDate field.
func (o *Event) SetEventDate(v time.Time) {
	o.EventDate = &v
}

// GetLocation returns the Location field value if set, zero value otherwise.
func (o *Event) GetLocation() int32 {
	if o == nil || IsNil(o.Location) {
		var ret int32
		return ret
	}
	return *o.Location
}

// GetLocationOk returns a tuple with the Location field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetLocationOk() (*int32, bool) {
	if o == nil || IsNil(o.Location) {
		return nil, false
	}
	return o.Location, true
}

// HasLocation returns a boolean if a field has been set.
func (o *Event) HasLocation() bool {
	if o != nil && !IsNil(o.Location) {
		return true
	}

	return false
}

// SetLocation gets a reference to the given int32 and assigns it to the Location field.
func (o *Event) SetLocation(v int32) {
	o.Location = &v
}

// GetCapacity returns the Capacity field value if set, zero value otherwise.
func (o *Event) GetCapacity() int32 {
	if o == nil || IsNil(o.Capacity) {
		var ret int32
		return ret
	}
	return *o.Capacity
}

// GetCapacityOk returns a tuple with the Capacity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetCapacityOk() (*int32, bool) {
	if o == nil || IsNil(o.Capacity) {
		return nil, false
	}
	return o.Capacity, true
}

// HasCapacity returns a boolean if a field has been set.
func (o *Event) HasCapacity() bool {
	if o != nil && !IsNil(o.Capacity) {
		return true
	}

	return false
}

// SetCapacity gets a reference to the given int32 and assigns it to the Capacity field.
func (o *Event) SetCapacity(v int32) {
	o.Capacity = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *Event) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *Event) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *Event) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *Event) GetUpdatedAt() time.Time {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret time.Time
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Event) GetUpdatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *Event) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given time.Time and assigns it to the UpdatedAt field.
func (o *Event) SetUpdatedAt(v time.Time) {
	o.UpdatedAt = &v
}

func (o Event) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Event) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Title) {
		toSerialize["title"] = o.Title
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.EventType) {
		toSerialize["event_type"] = o.EventType
	}
	if !IsNil(o.EventDate) {
		toSerialize["event_date"] = o.EventDate
	}
	if !IsNil(o.Location) {
		toSerialize["location"] = o.Location
	}
	if !IsNil(o.Capacity) {
		toSerialize["capacity"] = o.Capacity
	}
	if !IsNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	if !IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	return toSerialize, nil
}

type NullableEvent struct {
	value *Event
	isSet bool
}

func (v NullableEvent) Get() *Event {
	return v.value
}

func (v *NullableEvent) Set(val *Event) {
	v.value = val
	v.isSet = true
}

func (v NullableEvent) IsSet() bool {
	return v.isSet
}

func (v *NullableEvent) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEvent(val *Event) *NullableEvent {
	return &NullableEvent{value: val, isSet: true}
}

func (v NullableEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEvent) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
