package mlmodel

import (
	"context"

	"github.com/edaniels/golog"
	pb "go.viam.com/api/service/mlmodel/v1"
	vprotoutils "go.viam.com/utils/protoutils"
	"go.viam.com/utils/rpc"
)

// client implements MLModelServiceClient.
type client struct {
	name   string
	conn   rpc.ClientConn
	client pb.MLModelServiceClient
	logger golog.Logger
}

// NewClientFromConn constructs a new Client from connection passed in.
func NewClientFromConn(ctx context.Context, conn rpc.ClientConn, name string, logger golog.Logger) Service {
	grpcClient := pb.NewMLModelServiceClient(conn)
	c := &client{
		name:   name,
		conn:   conn,
		client: grpcClient,
		logger: logger,
	}
	return c
}

func (c *client) Infer(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	inProto, err := vprotoutils.StructToStructPb(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Infer(ctx, &pb.InferRequest{
		Name:      c.name,
		InputData: inProto,
	})
	if err != nil {
		return nil, err
	}
	return resp.OutputData.AsMap(), nil
}

func (c *client) Metadata(ctx context.Context) (MLMetadata, error) {
	resp, err := c.client.Metadata(ctx, &pb.MetadataRequest{
		Name: c.name,
	})
	if err != nil {
		return MLMetadata{}, err
	}
	metadata, err := protoToMetadata(resp.Metadata)
	if err != nil {
		return MLMetadata{}, err
	}
	return metadata, nil
}

// protoToMetadata takes a pb.Metadata protobuf message and turns it into an MLMetadata struct.
func protoToMetadata(pbmd *pb.Metadata) (MLMetadata, error) {
	metadata := MLMetadata{
		ModelName:        pbmd.Name,
		ModelType:        pbmd.Type,
		ModelDescription: pbmd.Description,
	}
	inputData := make([]TensorInfo, 0, len(pbmd.InputInfo))
	for _, idproto := range pbmd.InputInfo {
		id, err := protoToTensorInfo(idproto)
		if err != nil {
			return MLMetadata{}, err
		}
		inputData = append(inputData, id)
	}
	metadata.Inputs = inputData
	outputData := make([]TensorInfo, 0, len(pbmd.OutputInfo))
	for _, odproto := range pbmd.OutputInfo {
		od, err := protoToTensorInfo(odproto)
		if err != nil {
			return MLMetadata{}, err
		}
		outputData = append(outputData, od)
	}
	metadata.Outputs = outputData
	return metadata, nil
}

// protoToTensorInfo takes a pb.TensorInfo protobuf message and turns it into an TensorInfo struct.
func protoToTensorInfo(pbti *pb.TensorInfo) (TensorInfo, error) {
	ti := TensorInfo{
		Name:        pbti.Name,
		Description: pbti.Description,
		DataType:    pbti.DataType,
		NDim:        int(pbti.NDim),
		Extra:       pbti.Extra.AsMap(),
	}
	associatedFiles := make([]File, 0, len(pbti.AssociatedFiles))
	for _, afproto := range pbti.AssociatedFiles {
		af, err := protoToFile(afproto)
		if err != nil {
			return TensorInfo{}, err
		}
		associatedFiles = append(associatedFiles, af)
	}
	ti.AssociatedFiles = associatedFiles
	return ti, nil
}

// protoToFile takes a pb.File protobuf message and turns it into an File struct.
func protoToFile(pbf *pb.File) (File, error) {
	f := File{
		Name:        pbf.Name,
		Description: pbf.Description,
	}
	switch pbf.LabelType {
	case pb.LabelType_LABEL_TYPE_UNSPECIFIED:
		f.LabelType = LabelTypeUnspecified
	case pb.LabelType_LABEL_TYPE_TENSOR_VALUE:
		f.LabelType = LabelTypeTensorValue
	case pb.LabelType_LABEL_TYPE_TENSOR_AXIS:
		f.LabelType = LabelTypeTensorAxis
	default:
		// this should never happen as long as all possible enums are included in the switch
		f.LabelType = LabelTypeUnspecified
	}
	return f, nil
}
