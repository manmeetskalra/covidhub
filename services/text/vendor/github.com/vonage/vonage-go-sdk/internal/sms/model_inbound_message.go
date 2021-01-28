/*
 * SMS API
 *
 * With the Nexmo SMS API you can send SMS from your account and lookup messages both messages that you've sent as well as messages sent to your virtual numbers. Numbers are specified in E.164 format. More SMS API documentation is at <https://developer.nexmo.com/messaging/sms/overview>
 *
 * API version: 1.0.5
 * Contact: devrel@nexmo.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package sms
import (
	"os"
)
// InboundMessage struct for InboundMessage
type InboundMessage struct {
	// The phone number that this inbound message was sent from. Numbers are specified in E.164 format.
	Msisdn string `json:"msisdn"`
	// The phone number the message was sent to. **This is your virtual number**. Numbers are specified in E.164 format.
	To string `json:"to"`
	// The ID of the message
	MessageId string `json:"messageId"`
	// The message body for this inbound message.
	Text string `json:"text"`
	// Possible values are:    - `text` - standard text.   - `unicode` - URLencoded   unicode  . This is valid for standard GSM, Arabic, Chinese, double-encoded characters and so on.   - `binary` - a binary message. 
	Type string `json:"type"`
	// The first word in the message body. Converted to upper case.
	Keyword string `json:"keyword"`
	// The time when Nexmo started to push this Delivery Receipt to your webhook endpoint.
	MessageTimestamp string `json:"message-timestamp"`
	// A unix timestamp representation of message-timestamp.
	Timestamp string `json:"timestamp,omitempty"`
	// A random string that forms part of the signed set of parameters, it adds an extra element of unpredictability into the signature for the request. You use the nonce and timestamp parameters with your shared secret to calculate and validate the signature for inbound messages.
	Nonce string `json:"nonce,omitempty"`
	// True - if this is a concatenated message. This field does not exist if it is a single message
	Concat string `json:"concat,omitempty"`
	// The transaction reference. All parts of this message share this value.
	ConcatRef string `json:"concat-ref,omitempty"`
	// The number of parts in this concatenated message.
	ConcatTotal string `json:"concat-total,omitempty"`
	// The number of this part in the message. Counting starts at 1.
	ConcatPart string `json:"concat-part,omitempty"`
	// The content of this message, if type is binary.
	Data *os.File `json:"data,omitempty"`
	// The hex encoded User Data Header, if type is binary
	Udh string `json:"udh,omitempty"`
}
