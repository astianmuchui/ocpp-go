package core

import (
	"encoding/json"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Data Transfer (CP -> CS / CS -> CP) --------------------

const DataTransferFeatureName = "DataTransfer"

// Status in DataTransferConfirmation messages.
type DataTransferStatus string

const (
	DataTransferStatusAccepted         DataTransferStatus = "Accepted"
	DataTransferStatusRejected         DataTransferStatus = "Rejected"
	DataTransferStatusUnknownMessageId DataTransferStatus = "UnknownMessageId"
	DataTransferStatusUnknownVendorId  DataTransferStatus = "UnknownVendorId"
)

func isValidDataTransferStatus(fl validator.FieldLevel) bool {
	status := DataTransferStatus(fl.Field().String())
	switch status {
	case DataTransferStatusAccepted, DataTransferStatusRejected, DataTransferStatusUnknownMessageId, DataTransferStatusUnknownVendorId:
		return true
	default:
		return false
	}
}

// The field definition of the DataTransfer request payload sent by an endpoint to ther other endpoint.
//
// Vendors routinely put their custom fields at the top level of the payload
// instead of nesting them under "data" (VMOTO's rateReq/deviceData/chargePointData
// are one example). Those fields are preserved in Custom rather than dropped,
// and are written back out on marshalling, so a caller can read and echo them
// without this package knowing anything about the vendor's schema.
type DataTransferRequest struct {
	VendorId  string       `json:"vendorId" validate:"required,max=255"`
	MessageId string       `json:"messageId,omitempty" validate:"max=50"`
	Data      interface{}  `json:"data,omitempty"`
	Custom    CustomFields `json:"-"`
}

func (r *DataTransferRequest) UnmarshalJSON(b []byte) error {
	type alias DataTransferRequest
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	custom, err := extractCustomFields(b, "vendorId", "messageId", "data")
	if err != nil {
		return err
	}
	*r = DataTransferRequest(a)
	r.Custom = custom
	return nil
}

func (r DataTransferRequest) MarshalJSON() ([]byte, error) {
	type alias DataTransferRequest
	return marshalWithCustomFields(alias(r), r.Custom)
}

// This field definition of the DataTransfer confirmation payload, sent by an endpoint in response to a DataTransferRequest, coming from the other endpoint.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
//
// Custom carries vendor fields placed at the top level of the payload, and is
// how a tariff response of the VMOTO rateReq shape is built; see DataTransferRequest.
type DataTransferConfirmation struct {
	Status DataTransferStatus `json:"status" validate:"required,dataTransferStatus16"`
	Data   interface{}        `json:"data,omitempty"`
	Custom CustomFields       `json:"-"`
}

func (c *DataTransferConfirmation) UnmarshalJSON(b []byte) error {
	type alias DataTransferConfirmation
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	custom, err := extractCustomFields(b, "status", "data")
	if err != nil {
		return err
	}
	*c = DataTransferConfirmation(a)
	c.Custom = custom
	return nil
}

func (c DataTransferConfirmation) MarshalJSON() ([]byte, error) {
	type alias DataTransferConfirmation
	return marshalWithCustomFields(alias(c), c.Custom)
}

// If a Charge Point needs to send information to the Central System for a function not supported by OCPP, it SHALL use a DataTransfer message.
// The same functionality may also be offered the other way around, allowing a Central System to send arbitrary custom commands to a Charge Point.
type DataTransferFeature struct{}

func (f DataTransferFeature) GetFeatureName() string {
	return DataTransferFeatureName
}

func (f DataTransferFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(DataTransferRequest{})
}

func (f DataTransferFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(DataTransferConfirmation{})
}

func (r DataTransferRequest) GetFeatureName() string {
	return DataTransferFeatureName
}

func (c DataTransferConfirmation) GetFeatureName() string {
	return DataTransferFeatureName
}

// Creates a new DataTransferRequest, containing all required fields. Optional fields may be set afterwards.
func NewDataTransferRequest(vendorId string) *DataTransferRequest {
	return &DataTransferRequest{VendorId: vendorId}
}

// Creates a new DataTransferConfirmation. Optional fields may be set afterwards.
func NewDataTransferConfirmation(status DataTransferStatus) *DataTransferConfirmation {
	return &DataTransferConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("dataTransferStatus16", isValidDataTransferStatus)
}
