// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: img_downloader/v1/img_downloader.proto

package img_downloaderv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "img_downloader/gen/img_downloader/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// ImageServiceName is the fully-qualified name of the ImageService service.
	ImageServiceName = "img_downloader.v1.ImageService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ImageServiceDownloadImagesProcedure is the fully-qualified name of the ImageService's
	// DownloadImages RPC.
	ImageServiceDownloadImagesProcedure = "/img_downloader.v1.ImageService/DownloadImages"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	imageServiceServiceDescriptor              = v1.File_img_downloader_v1_img_downloader_proto.Services().ByName("ImageService")
	imageServiceDownloadImagesMethodDescriptor = imageServiceServiceDescriptor.Methods().ByName("DownloadImages")
)

// ImageServiceClient is a client for the img_downloader.v1.ImageService service.
type ImageServiceClient interface {
	DownloadImages(context.Context, *connect.Request[v1.DownloadImagesRequest]) (*connect.Response[v1.DownloadImagesResponse], error)
}

// NewImageServiceClient constructs a client for the img_downloader.v1.ImageService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewImageServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ImageServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &imageServiceClient{
		downloadImages: connect.NewClient[v1.DownloadImagesRequest, v1.DownloadImagesResponse](
			httpClient,
			baseURL+ImageServiceDownloadImagesProcedure,
			connect.WithSchema(imageServiceDownloadImagesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// imageServiceClient implements ImageServiceClient.
type imageServiceClient struct {
	downloadImages *connect.Client[v1.DownloadImagesRequest, v1.DownloadImagesResponse]
}

// DownloadImages calls img_downloader.v1.ImageService.DownloadImages.
func (c *imageServiceClient) DownloadImages(ctx context.Context, req *connect.Request[v1.DownloadImagesRequest]) (*connect.Response[v1.DownloadImagesResponse], error) {
	return c.downloadImages.CallUnary(ctx, req)
}

// ImageServiceHandler is an implementation of the img_downloader.v1.ImageService service.
type ImageServiceHandler interface {
	DownloadImages(context.Context, *connect.Request[v1.DownloadImagesRequest]) (*connect.Response[v1.DownloadImagesResponse], error)
}

// NewImageServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewImageServiceHandler(svc ImageServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	imageServiceDownloadImagesHandler := connect.NewUnaryHandler(
		ImageServiceDownloadImagesProcedure,
		svc.DownloadImages,
		connect.WithSchema(imageServiceDownloadImagesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/img_downloader.v1.ImageService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ImageServiceDownloadImagesProcedure:
			imageServiceDownloadImagesHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedImageServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedImageServiceHandler struct{}

func (UnimplementedImageServiceHandler) DownloadImages(context.Context, *connect.Request[v1.DownloadImagesRequest]) (*connect.Response[v1.DownloadImagesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("img_downloader.v1.ImageService.DownloadImages is not implemented"))
}
