package sflow

import (
	"fmt"
	"io"
)

const (
	SFDataRawHeader      = 1
	SFDataEthernetHeader = 2
	SFDataIPV4Header     = 3
	SFDataIPV6Header     = 4

	SFDataExtSwitch     = 1001
	SFDataExtRouter     = 1002
	SFDataExtGateway    = 1003
	SFDataExtUser       = 1004
	SFDataExtURL        = 1005
	SFDataExtMPLS       = 1006
	SFDataExtNAT        = 1007
	SFDataExtMPLSTunnel = 1008
	SFDataExtMPLSVC     = 1009
	SFDataExtMPLSFEC    = 1010
	SFDataExtMPLSLVPFEC = 1011
	SFDataExtVLANTunnel = 1012
)

type FlowSample struct {
	SequenceNo   uint32 // Incremented with each flow sample
	SourceId     byte   // fsSourceId
	SamplingRate uint32 // fsPacketSamplingRate
	SamplePool   uint32 // Total number of packets that could have been sampled
	Drops        uint32 // Number of times a packet was dropped due to lack of resources
	Input        uint32 // SNMP ifIndex of input interface
	Output       uint32 // SNMP ifIndex of input interface
	RecordsNo    uint32 // Number of records to follow
}

type SampledHeader struct {
	Protocol     uint32 // (enum SFLHeader_protocol)
	FrameLength  uint32 // Original length of packet before sampling
	Stripped     uint32 // Header/trailer bytes stripped by sender
	HeaderLength uint32 // Length of sampled header bytes to follow
	Header       []byte // Header bytes
}

func decodeFlowSample(r io.ReadSeeker) error {
	var (
		fs          = new(FlowSample)
		rTypeFormat uint32
		rTypeLength uint32
		err         error
	)

	if err = read(r, &fs.SequenceNo); err != nil {
		return err
	}

	if err = read(r, &fs.SourceId); err != nil {
		return err
	}

	r.Seek(3, 1) // skip counter sample decoding

	if err = read(r, &fs.SamplingRate); err != nil {
		return err
	}

	if err = read(r, &fs.SamplePool); err != nil {
		return err
	}

	if err = read(r, &fs.Drops); err != nil {
		return err
	}

	if err = read(r, &fs.Input); err != nil {
		return err
	}

	if err = read(r, &fs.Output); err != nil {
		return err
	}

	if err = read(r, &fs.RecordsNo); err != nil {
		return err
	}

	fmt.Printf("%#v\n", fs) // just for test

	for i := uint32(0); i < fs.RecordsNo; i++ {
		if err = read(r, &rTypeFormat); err != nil {
			return err
		}
		if err = read(r, &rTypeLength); err != nil {
			return err
		}

		switch rTypeFormat {
		case SFDataRawHeader:
			decodeSampledHeader(r)
			r.Seek(int64(rTypeLength), 1)
		case SFDataExtSwitch:
			decodeSampledExtSwitch(r)
			r.Seek(int64(rTypeLength), 1)
		default:
			r.Seek(int64(rTypeLength), 1)
		}
	}

	return nil
}

func decodeSampledHeader(r io.Reader) error {
	// TODO
	var (
		h   = new(SampledHeader)
		err error
	)
	if err = read(r, &h.Protocol); err != nil {
		return err
	}

	if err = read(r, &h.FrameLength); err != nil {
		return err
	}

	if err = read(r, &h.Stripped); err != nil {
		return err
	}

	if err = read(r, &h.HeaderLength); err != nil {
		return err
	}

	// TODO: make sure the padding works!!
	// cut off a header length mod 4 == 0 number of bytes
	tmp := 4 - (h.HeaderLength % 4)

	h.Header = make([]byte, h.HeaderLength+tmp)
	r.Read(h.Header)
	h.Header = h.Header[:h.HeaderLength]

	fmt.Printf("%#v\n", h)

	return nil
}

func decodeSampledExtSwitch(r io.Reader) {
	// TODO
}

func decodeSampledIPv4() {
	// TODO
}

func decodeSampledIPv6() {
	// TODO
}
