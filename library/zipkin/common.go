package zipkin

import (
	metadataMicro "github.com/micro/go-micro/v2/metadata"
	"google.golang.org/grpc"
	metadataGrpc "google.golang.org/grpc/metadata"
	"io"
	"net"
	"strings"
)

type MetadataMircoReaderWriter struct {
	metadataMicro.Metadata
}

func (w MetadataMircoReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	w.Metadata[key] = val
}

func (w MetadataMircoReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, v := range w.Metadata {
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

type MetadataGrpcReaderWriter struct {
	metadataGrpc.MD
}

func (w MetadataGrpcReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

func (w MetadataGrpcReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

type OpenTracingClientStream struct {
	grpc.ClientStream
	Desc       *grpc.StreamDesc
	FinishFunc func(error)
}

func (cs *OpenTracingClientStream) Header() (metadataGrpc.MD, error) {
	md, err := cs.ClientStream.Header()
	if err != nil {
		cs.FinishFunc(err)
	}
	return md, err
}

func (cs *OpenTracingClientStream) SendMsg(m interface{}) error {
	err := cs.ClientStream.SendMsg(m)
	if err != nil {
		cs.FinishFunc(err)
	}
	return err
}

func (cs *OpenTracingClientStream) RecvMsg(m interface{}) error {
	err := cs.ClientStream.RecvMsg(m)
	if err == io.EOF {
		cs.FinishFunc(nil)
		return err
	} else if err != nil {
		cs.FinishFunc(err)
		return err
	}
	if !cs.Desc.ServerStreams {
		cs.FinishFunc(nil)
	}
	return err
}

func (cs *OpenTracingClientStream) CloseSend() error {
	err := cs.ClientStream.CloseSend()
	if err != nil {
		cs.FinishFunc(err)
	}
	return err
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
