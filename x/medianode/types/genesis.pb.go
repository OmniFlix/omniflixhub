// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: OmniFlix/medianode/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the medianode module's genesis state
type GenesisState struct {
	// params defines all the parameters of the module
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// media_nodes is the list of registered media nodes
	Nodes []MediaNode `protobuf:"bytes,2,rep,name=nodes,proto3" json:"nodes"`
	// leases is the list of active leases
	Leases []Lease `protobuf:"bytes,3,rep,name=leases,proto3" json:"leases"`
	// next_medianode_id is the current usable media node ID
	NodeCounter uint64 `protobuf:"varint,4,opt,name=node_counter,json=nodeCounter,proto3" json:"node_counter,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ee63eaf92493f5e, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetNodes() []MediaNode {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *GenesisState) GetLeases() []Lease {
	if m != nil {
		return m.Leases
	}
	return nil
}

func (m *GenesisState) GetNodeCounter() uint64 {
	if m != nil {
		return m.NodeCounter
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "OmniFlix.medianode.v1beta1.GenesisState")
}

func init() {
	proto.RegisterFile("OmniFlix/medianode/v1beta1/genesis.proto", fileDescriptor_1ee63eaf92493f5e)
}

var fileDescriptor_1ee63eaf92493f5e = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x31, 0x4b, 0xc3, 0x40,
	0x18, 0x86, 0x73, 0xb6, 0x76, 0xb8, 0x76, 0x0a, 0x0e, 0xa1, 0xc3, 0x99, 0x16, 0xc4, 0xe0, 0x70,
	0x47, 0xab, 0xb8, 0xaa, 0x15, 0x74, 0xd1, 0x2a, 0x75, 0x73, 0x91, 0x4b, 0x73, 0xa6, 0x07, 0x49,
	0x2e, 0xe4, 0x2e, 0xa5, 0xfe, 0x0b, 0x7f, 0x56, 0xc7, 0x8e, 0x4e, 0x22, 0xc9, 0x8f, 0x70, 0x95,
	0xcb, 0xa5, 0xd5, 0xc5, 0x6c, 0x1f, 0x1f, 0xcf, 0xfb, 0xf0, 0xf2, 0x42, 0xef, 0x21, 0x4e, 0xf8,
	0x4d, 0xc4, 0x57, 0x24, 0x66, 0x01, 0xa7, 0x89, 0x08, 0x18, 0x59, 0x8e, 0x7c, 0xa6, 0xe8, 0x88,
	0x84, 0x2c, 0x61, 0x92, 0x4b, 0x9c, 0x66, 0x42, 0x09, 0xbb, 0xbf, 0x25, 0xf1, 0x8e, 0xc4, 0x35,
	0xd9, 0x3f, 0x08, 0x45, 0x28, 0x2a, 0x8c, 0xe8, 0xcb, 0x24, 0xfa, 0x27, 0x0d, 0xee, 0x5f, 0x87,
	0x61, 0x8f, 0x1b, 0xd8, 0x94, 0x66, 0x34, 0xae, 0x6b, 0x0c, 0xbf, 0x01, 0xec, 0xdd, 0x9a, 0x62,
	0x4f, 0x8a, 0x2a, 0x66, 0x5f, 0xc2, 0x8e, 0x01, 0x1c, 0xe0, 0x02, 0xaf, 0x3b, 0x1e, 0xe2, 0xff,
	0x8b, 0xe2, 0xc7, 0x8a, 0x9c, 0xb4, 0xd7, 0x9f, 0x87, 0xd6, 0xac, 0xce, 0xd9, 0x57, 0x70, 0x5f,
	0x43, 0xd2, 0xd9, 0x73, 0x5b, 0x5e, 0x77, 0x7c, 0xd4, 0x24, 0xb8, 0xd7, 0x9f, 0xa9, 0x08, 0x58,
	0xed, 0x30, 0x49, 0xfb, 0x02, 0x76, 0x22, 0x46, 0x25, 0x93, 0x4e, 0xab, 0x72, 0x0c, 0x9a, 0x1c,
	0x77, 0x9a, 0xdc, 0x76, 0x30, 0x31, 0x7b, 0x00, 0x7b, 0x9a, 0x79, 0x99, 0x8b, 0x3c, 0x51, 0x2c,
	0x73, 0xda, 0x2e, 0xf0, 0xda, 0xb3, 0xae, 0xfe, 0x5d, 0x9b, 0xd7, 0x64, 0xba, 0x2e, 0x10, 0xd8,
	0x14, 0x08, 0x7c, 0x15, 0x08, 0xbc, 0x97, 0xc8, 0xda, 0x94, 0xc8, 0xfa, 0x28, 0x91, 0xf5, 0x7c,
	0x16, 0x72, 0xb5, 0xc8, 0x7d, 0x3c, 0x17, 0x31, 0xd9, 0xed, 0x28, 0xe2, 0x84, 0xbf, 0x46, 0x7c,
	0xb5, 0xc8, 0x7d, 0xb2, 0x3c, 0x27, 0x7f, 0x87, 0x55, 0x6f, 0x29, 0x93, 0x7e, 0xa7, 0x1a, 0xf4,
	0xf4, 0x27, 0x00, 0x00, 0xff, 0xff, 0xd5, 0x82, 0x89, 0x1d, 0x03, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.NodeCounter != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.NodeCounter))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Leases) > 0 {
		for iNdEx := len(m.Leases) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Leases[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Nodes) > 0 {
		for iNdEx := len(m.Nodes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Nodes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Nodes) > 0 {
		for _, e := range m.Nodes {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Leases) > 0 {
		for _, e := range m.Leases {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.NodeCounter != 0 {
		n += 1 + sovGenesis(uint64(m.NodeCounter))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nodes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nodes = append(m.Nodes, MediaNode{})
			if err := m.Nodes[len(m.Nodes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leases", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Leases = append(m.Leases, Lease{})
			if err := m.Leases[len(m.Leases)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeCounter", wireType)
			}
			m.NodeCounter = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NodeCounter |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
