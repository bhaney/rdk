package mlvision

import (
	"context"
	"image"
	"math"
	"strconv"

	"github.com/montanaflynn/stats"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"

	"go.viam.com/rdk/rimage"
	"go.viam.com/rdk/services/mlmodel"
	"go.viam.com/rdk/vision/classification"
)

func attemptToBuildClassifier(mlm mlmodel.Service) (classification.Classifier, error) {
	md, err := mlm.Metadata(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "could not find metadata")
	}

	// Set up input type, height, width, and labels
	var inHeight, inWidth uint
	inType := md.Inputs[0].DataType
	labels := getLabelsFromMetadata(md)
	if shape := md.Inputs[0].Shape; getIndex(shape, 3) == 1 {
		inHeight, inWidth = uint(shape[2]), uint(shape[3])
	} else {
		inHeight, inWidth = uint(shape[1]), uint(shape[2])
	}

	return func(ctx context.Context, img image.Image) (classification.Classifications, error) {
		resized := resize.Resize(inWidth, inHeight, img, resize.Bilinear)
		inMap := make(map[string]interface{})
		switch inType {
		case UInt8:
			inMap["image"] = rimage.ImageToUInt8Buffer(resized)
		case Float32:
			inMap["image"] = rimage.ImageToFloatBuffer(resized)
		default:
			return nil, errors.New("invalid input type. try uint8 or float32")
		}
		outMap, err := mlm.Infer(ctx, inMap)
		if err != nil {
			return nil, err
		}

		probs := unpack(outMap, "probability", md)

		confs := checkClassificationScores(probs)
		classifications := make(classification.Classifications, 0, len(confs))
		for i := 0; i < len(confs); i++ {
			if labels != nil {
				classifications = append(classifications, classification.NewClassification(confs[i], labels[i]))
			} else {
				classifications = append(classifications, classification.NewClassification(confs[i], strconv.Itoa(i)))
			}
		}
		return classifications, nil
	}, nil
}

// softmax takes the input slice and applies the softmax function.
func softmax(in []float64) []float64 {
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

// checkClassification scores ensures that the input scores (output of classifier)
// will represent confidence values (from 0-1).
func checkClassificationScores(in []float64) []float64 {
	if len(in) > 1 {
		for _, p := range in {
			if p < 0 || p > 1 { // is logit, needs softmax
				confs := softmax(in)
				return confs
			}
		}
		return in // no need to softmax
	}
	// otherwise, this is a binary classifier
	if in[0] < -1 || in[0] > 1 { // needs sigmoid
		out, err := stats.Sigmoid(in)
		if err != nil {
			return in
		}
		return out
	}
	return in // no need to sigmoid
}
