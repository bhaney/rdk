package mlvision

import (
	"context"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"go.viam.com/rdk/rimage"
	"go.viam.com/rdk/services/mlmodel"
	"go.viam.com/rdk/vision/classification"
	"image"
	"math"
)

func attemptToBuildClassifier(mlm mlmodel.Service) (classification.Classifier, error) {
	md, err := mlm.Metadata(context.Background())
	if err != nil {
		// If the metadata isn't there
		return nil, err
	}

	// Set up input type, height, width, and labels
	var inHeight, inWidth uint
	inType := md.Inputs[0].DataType
	labels, err := getLabelsFromMetadata(md)
	if err != nil {
		// Not true, still do something if we can't get labels
		return nil, err
	}

	if shape := md.Inputs[0].Shape; getIndex(shape, 3) == 1 {
		inHeight, inWidth = uint(shape[2]), uint(shape[3])
	} else {
		inHeight, inWidth = uint(shape[1]), uint(shape[2])
	}

	return func(ctx context.Context, img image.Image) (classification.Classifications, error) {
		resized := resize.Resize(inWidth, inHeight, img, resize.Bilinear)
		inMap := make(map[string]interface{}, 1)
		outMap := make(map[string]interface{}, 5)
		switch inType {
		case "uint8":
			inMap["image"] = rimage.ImageToUInt8Buffer(resized)
			outMap, err = mlm.Infer(ctx, inMap)
		case "float32":
			inMap["image"] = rimage.ImageToFloatBuffer(resized)
			outMap, err = mlm.Infer(ctx, inMap)
		default:
			return nil, errors.New("invalid input type. try uint8 or float32")
		}
		if err != nil {
			return nil, err
		}

		probs := unpackMe(outMap, "probability", md)

		// TODO: Khari, somewhere around here, softmax(probs) --> confs and then use confs
		confs := softmaxMe(probs)

		classifications := make(classification.Classifications, 0, len(confs))
		for i := 0; i < len(confs); i++ {
			classifications = append(classifications, classification.NewClassification(confs[i], labels[i]))
		}
		return classifications, nil

	}, nil

}

func softmaxMe(in []float64) []float64 {
	out := make([]float64, 0, len(in))
	bigSum := 0.0
	for _, x := range in {
		bigSum += math.Exp(x)
	}
	for _, x := range in {
		out = append(out, math.Exp(x)/bigSum)
	}
	return out
}
