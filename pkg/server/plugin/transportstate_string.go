// Code generated by "stringer -type=TransportState"; DO NOT EDIT

package plugin

import "fmt"

const _TransportState_name = "TransportStateUninitializedTransportStateInitializedTransportStateReadyTransportStateWorkerUnavailableTransportStateError"

var _TransportState_index = [...]uint8{0, 27, 52, 71, 102, 121}

func (i TransportState) String() string {
	if i < 0 || i >= TransportState(len(_TransportState_index)-1) {
		return fmt.Sprintf("TransportState(%d)", i)
	}
	return _TransportState_name[_TransportState_index[i]:_TransportState_index[i+1]]
}
