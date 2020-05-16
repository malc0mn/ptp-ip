package ptp

// The most significant nibble (4 bits) is used to indicate the category of the code and whether the code value is
// standard or vendor-extended: 0100 = standard, 1100 = vendor-extended.
type EventCode uint16

const (
	EC_Undefined EventCode = 0x4000

	// This formal event is used to cancel a transaction for transports that do not have a specified or standard way of
	// canceling transactions. The particular method used to cancel transactions may be ip-specific. When an
	// Initiator or Responder receives a CancelTransaction event, it should abort the transaction referred to by the
	// TransactionID in the event dataset. If that transaction is already complete, the event should be ignored.
	// After receiving a CancelTransfer event from the Initiator, the Responder shall send an IncompleteTransfer
	// response for the operation that was cancelled. Both devices will then be ready for the next transaction.
	EC_CancelTransaction EventCode = 0x4001

	// A new data object was added to the device. The new handle assigned by the device to the new object should be
	// passed in the Parameter1 field of the event. If more than one object was added, each new object should generate a
	// separate ObjectAdded event. The appearance of a new store on the device should not cause the creation of new
	// ObjectAdded events for the new objects present on the new store, but should instead cause the generation of a
	// StoreAdded event.
	EC_ObjectAdded EventCode = 0x4002

	// A data object was removed from the device unexpectedly due to something external to the current session. The
	// handle of the object that was removed should be passed in the Parameter1 field of the event. If more than one
	// image was removed, the separate ObjectRemoved events should be generated for each. If the data object that was
	// removed was removed because of a previous operation that is a part of this session, no event needs to be sent to
	// the opposing device. The removal of a store on the device should not cause the creation of ObjectRemoved events
	// for the objects present on the removed store, but should instead cause the generation of one StoreRemoved event
	// with the appropriate PhysicalStorageID.
	EC_ObjectRemoved EventCode = 0x4003

	// A new store was added to the device. If this is a new physical store that contains only one logical store, then
	// the complete StorageID of the new store should be indicated in the first parameter. If the new store contains
	// more than one logical store, then the first parameter should be set to 0x00000000. This indicates that the list
	// of StorageIDs should be re-obtained using the GetStorageIDs operation, and examined appropriately. Any new
	// StorageIDs discovered should result in the appropriate invocations of GetStorageInfo operations.
	EC_StoreAdded EventCode = 0x4004

	// The indicated stores are no longer available. The opposing device may assume that the StorageInfo datasets and
	// ObjectHandles associated with those stores are no longer valid. The first parameter is used to indicate the
	// StorageID of the store that is no longer available. If the store removed is only a single logical store within a
	// physical store, the entire StorageID should be sent, which indicates that any other logical stores on that
	// physical store are still available. If the physical store and all logical stores upon it are removed (e.g.
	// removal of an ejectable media with multiple partitions), the first parameter should contain the PhysicalStorageID
	// in the most significant sixteen bits, with the least significant sixteen bits set to 0xFFFF.
	EC_StoreRemoved EventCode = 0x4005

	// A property changed on the device due to something external to this session. The appropriate property dataset
	// should be requested from the opposing device.
	EC_DevicePropChanged EventCode = 0x4006

	// Indicates that the ObjectInfo dataset for a particular object has changed, and that it should be re-requested.
	EC_ObjectInfoChanged EventCode = 0x4007

	// Indicates that the capabilities of the device have changed, and that the DeviceInfo should be re-requested. This
	// may be caused by the device going into or out of a sleep state, or by the device losing or gaining some
	// functionality in some way.
	EC_DeviceInfoChanged EventCode = 0x4008

	// This event can be used by a Responder to ask the Initiator to initiate a GetObject operation on the handle
	// specified in the first parameter. This allows for push-mode to be enabled on devices/transports that are
	// intrinsically pull mode.
	EC_RequestObjectTransfer EventCode = 0x4009

	// This event shall be sent when a store becomes full. Any multi-object capture that may be occurring should retain
	// the objects that were written to a store before the store became full.
	EC_StoreFull EventCode = 0x400A

	// This event needs only to be supported for devices that support multiple sessions or in the case if the device is
	// capable of resetting itself automatically or manually through user intervention while connected. This event shall
	// be sent to all open sessions other than the session that initiated the operation. This event shall be interpreted
	// as indicating that the sessions are about to be closed.
	EC_DeviceReset EventCode = 0x400B

	// This event is used when information in the StorageInfo dataset for a store changes. This can occur due to device
	// properties changing, such as ImageSize, which can cause changes in fields such as FreeSpaceInImages. This event
	// is typically not needed if the change is caused by an in-session operation that affects whole objects in a
	// deterministic manner. This includes changes in FreeSpaceInImages or FreeSpaceInBytes caused by operations such as
	// InitiateCapture or CopyObject, where the Initiator can recognize the changes due to the successful response code
	// of the operation, and/or related required events.
	EC_StorageInfoChanged EventCode = 0x400C

	// This event is used to indicate that a capture session, previously initiated by the InitiateCapture operation, is
	// complete, and that no more ObjectAdded events will occur as the result of this asynchronous operation. This
	// operation is not used for InitiateOpenCapture operations.
	EC_CaptureComplete EventCode = 0x400D

	// This event may be implemented for certain transports where situations can arise where the Responder was unable to
	// report events to the Initiator regarding changes in its internal status. When an Initiator receives this event,
	// it is responsible for doing whatever is necessary to ensure that its knowledge of the Responder is up to date.
	// This may include re-obtaining individual datasets, ObjectHandle lists, etc., or may even result in the session
	// being closed and re-opened. This event is typically only needed in situations where the ip used by the
	// device supports a suspend/resume/remote-wakeup feature and the Responder has gone into a suspend state and has
	// been unable to report state changes during that time period. This prevents the need for queuing of these
	// unreportable events. The details of the use of this event are ip-specific and should be fully specified in
	// the specific ip implementation specification.
	EC_UnreportedStatus EventCode = 0x400E
)

type Event struct {
	// Indicates the event.
	EventCode EventCode

	// Indicates the SessionID of the session for which the event is relevant. If the event is relevant to all open
	// sessions, this field should be set to 0xFFFFFFFF.
	SessionID SessionID

	// If the event corresponds to a previously initiated transaction, this field shall hold the TransactionID of that
	// operation. If the event is not specific to a particular transaction, this field shall be set to 0xFFFFFFFF.
	TransactionID TransactionID

	// These fields hold the event-specific nth parameter. Events may have at most three parameters. The interpretation
	// of any parameter is dependent upon the EventCode. Any unused parameter fields should be set to 0x00000000. If a
	// parameter holds a value that is less than 32 bits, the lowest significant bits shall be used to store the value,
	// with the most significant bits being set to zeros.
	Parameter1 interface{}
	Parameter2 interface{}
	Parameter3 interface{}
}

func (e *Event) Session() SessionID {
	return e.SessionID
}
