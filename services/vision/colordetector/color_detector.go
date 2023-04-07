// Package colordetector uses a heuristic based on hue and connected components to create
// bounding boxes around objects of a specified color.
package colordetector

import (
	"context"

	"github.com/edaniels/golog"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/robot"
	"go.viam.com/rdk/services/vision"
	"go.viam.com/rdk/utils"
	objdet "go.viam.com/rdk/vision/objectdetection"
)

var model = resource.NewDefaultModel("color_detector")

func init() {
	resource.RegisterService(vision.Subtype, model, resource.Registration[vision.Service, *objdet.ColorDetectorConfig]{
		DeprecatedRobotConstructor: func(ctx context.Context, r any, c resource.Config, logger golog.Logger) (vision.Service, error) {
			attrs, err := resource.NativeConfig[*objdet.ColorDetectorConfig](c)
			if err != nil {
				return nil, err
			}
			actualR, err := utils.AssertType[robot.Robot](r)
			if err != nil {
				return nil, err
			}
			return registerColorDetector(ctx, c.ResourceName().Name, attrs, actualR)
		},
	})
}

// registerColorDetector creates a new Color Detector from the config.
func registerColorDetector(
	ctx context.Context,
	name string,
	conf *objdet.ColorDetectorConfig,
	r robot.Robot,
) (vision.Service, error) {
	_, span := trace.StartSpan(ctx, "service::vision::registerColorDetector")
	defer span.End()
	if conf == nil {
		return nil, errors.New("object detection config for color detector cannot be nil")
	}
	detector, err := objdet.NewColorDetector(conf)
	if err != nil {
		return nil, errors.Wrapf(err, "register color detector %s", name)
	}
	return vision.NewService(name, r, nil, nil, detector, nil)
}
