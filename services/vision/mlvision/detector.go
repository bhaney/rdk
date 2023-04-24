package mlvision

import (
	"context"
	"image"
	"strconv"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"go.viam.com/utils/artifact"

	"go.viam.com/rdk/rimage"
	"go.viam.com/rdk/services/mlmodel"
	"go.viam.com/rdk/utils"
	"go.viam.com/rdk/vision/objectdetection"
)

func attemptToBuildDetector(mlm mlmodel.Service) (objectdetection.Detector, error) {
	md, err := mlm.Metadata(context.Background())
	if err != nil {
		return nil, errors.New("could not get any metadata")
	}

	// Set up input type, height, width, and labels
	var inHeight, inWidth uint
	inType := md.Inputs[0].DataType
	labels := getLabelsFromMetadata(md)
	boxOrder, err := getBoxOrderFromMetadata(md)
	if err != nil || len(boxOrder) < 4 {
		boxOrder = []int{1, 0, 3, 2}
	}
	if shape := md.Inputs[0].Shape; getIndex(shape, 3) == 1 {
		inHeight, inWidth = uint(shape[2]), uint(shape[3])
	} else {
		inHeight, inWidth = uint(shape[1]), uint(shape[2])
	}

	return func(ctx context.Context, img image.Image) ([]objectdetection.Detection, error) {
		origW, origH := img.Bounds().Dx(), img.Bounds().Dy()
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

		locations, err := unpack(outMap, "location", md)
		if err != nil {
			return nil, err
		}
		if len(locations) == 0 {
			locations, err = unpack(outMap, "output0", md)
			if err != nil {
				return nil, err
			}
		}
		categories, err := unpack(outMap, "category", md)
		if err != nil {
			return nil, err
		}
		if len(categories) == 0 {
			categories, err = unpack(outMap, "output1", md)
			if err != nil {
				return nil, err
			}
		}
		scores, err := unpack(outMap, "score", md)
		if err != nil {
			return nil, err
		}
		if len(scores) == 0 {
			scores, err = unpack(outMap, "output2", md)
			if err != nil {
				return nil, err
			}
		}

		// Now reshape outMap into Detections
		detections := make([]objectdetection.Detection, 0, len(categories))
		for i := 0; i < len(scores); i++ {
			xmin, ymin, xmax, ymax := utils.Clamp(locations[4*i+getIndex(boxOrder, 0)], 0, 1)*float64(origW),
				utils.Clamp(locations[4*i+getIndex(boxOrder, 1)], 0, 1)*float64(origH),
				utils.Clamp(locations[4*i+getIndex(boxOrder, 2)], 0, 1)*float64(origW),
				utils.Clamp(locations[4*i+getIndex(boxOrder, 3)], 0, 1)*float64(origH)
			rect := image.Rect(int(xmin), int(ymin), int(xmax), int(ymax))
			labelNum := int(categories[i])

			if labels != nil {
				detections = append(detections, objectdetection.NewDetection(rect, scores[i], labels[labelNum]))
			} else {
				detections = append(detections, objectdetection.NewDetection(rect, scores[i], strconv.Itoa(labelNum)))
			}
		}
		return detections, nil
	}, nil
}

func checkIfDetectorWorks(ctx context.Context, df objectdetection.Detector) (objectdetection.Detector, error) {
	if df == nil {
		return nil, errors.New("Nil detector function")
	}

	img, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/dogscute.jpeg"))
	if err != nil {
		return nil, err
	}

	_, err = df(ctx, img)
	if err != nil {
		return nil, errors.New("Cannot use model as a detector")
	}
	return df, nil
}
