// Code generated by protoc-gen-go-drpc. DO NOT EDIT.
// protoc-gen-go-drpc version: v0.0.20
// source: certificate.proto

package certificatepb

import (
	bytes "bytes"
	context "context"
	errors "errors"

	jsonpb "github.com/gogo/protobuf/jsonpb"
	proto "github.com/gogo/protobuf/proto"

	drpc "storj.io/drpc"
	drpcerr "storj.io/drpc/drpcerr"
)

type drpcEncoding_File_certificate_proto struct{}

func (drpcEncoding_File_certificate_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_certificate_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_certificate_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	var buf bytes.Buffer
	err := new(jsonpb.Marshaler).Marshal(&buf, msg.(proto.Message))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (drpcEncoding_File_certificate_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return jsonpb.Unmarshal(bytes.NewReader(buf), msg.(proto.Message))
}

type DRPCCertificatesClient interface {
	DRPCConn() drpc.Conn

	Sign(ctx context.Context, in *SigningRequest) (*SigningResponse, error)
}

type drpcCertificatesClient struct {
	cc drpc.Conn
}

func NewDRPCCertificatesClient(cc drpc.Conn) DRPCCertificatesClient {
	return &drpcCertificatesClient{cc}
}

func (c *drpcCertificatesClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcCertificatesClient) Sign(ctx context.Context, in *SigningRequest) (*SigningResponse, error) {
	out := new(SigningResponse)
	err := c.cc.Invoke(ctx, "/node.Certificates/Sign", drpcEncoding_File_certificate_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCCertificatesServer interface {
	Sign(context.Context, *SigningRequest) (*SigningResponse, error)
}

type DRPCCertificatesUnimplementedServer struct{}

func (s *DRPCCertificatesUnimplementedServer) Sign(context.Context, *SigningRequest) (*SigningResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

type DRPCCertificatesDescription struct{}

func (DRPCCertificatesDescription) NumMethods() int { return 1 }

func (DRPCCertificatesDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/node.Certificates/Sign", drpcEncoding_File_certificate_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCCertificatesServer).
					Sign(
						ctx,
						in1.(*SigningRequest),
					)
			}, DRPCCertificatesServer.Sign, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterCertificates(mux drpc.Mux, impl DRPCCertificatesServer) error {
	return mux.Register(impl, DRPCCertificatesDescription{})
}

type DRPCCertificates_SignStream interface {
	drpc.Stream
	SendAndClose(*SigningResponse) error
}

type drpcCertificates_SignStream struct {
	drpc.Stream
}

func (x *drpcCertificates_SignStream) SendAndClose(m *SigningResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_certificate_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}
