package mlvision

import (
	"context"
	"sync"
	"testing"

	"github.com/edaniels/golog"
	"go.viam.com/test"
	"go.viam.com/utils/artifact"

	"go.viam.com/rdk/rimage"
	"go.viam.com/rdk/services/mlmodel/tflitecpu"
	"go.viam.com/rdk/testutils/inject"
	"go.viam.com/rdk/vision/classification"
)

func BenchmarkAddMLVisionModel(b *testing.B) {
	modelLoc := artifact.MustPath("vision/tflite/effdet0.tflite")

	name := "myMLModel"
	cfg := tflitecpu.TFLiteConfig{
		ModelPath:  modelLoc,
		NumThreads: 2,
	}
	ctx := context.Background()
	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, name)
	test.That(b, err, test.ShouldBeNil)
	test.That(b, out, test.ShouldNotBeNil)
	modelCfg := MLModelConfig{ModelName: name}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service, err := registerMLModelVisionService(ctx, name, &modelCfg, &inject.Robot{}, golog.NewLogger("benchmark"))
		test.That(b, err, test.ShouldBeNil)
		test.That(b, service, test.ShouldNotBeNil)
	}
}

func BenchmarkUseMLVisionModel(b *testing.B) {
	modelLoc := artifact.MustPath("vision/tflite/effdet0.tflite")
	pic, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/dogscute.jpeg"))
	test.That(b, err, test.ShouldBeNil)
	test.That(b, pic, test.ShouldNotBeNil)
	name := "myMLModel"
	cfg := tflitecpu.TFLiteConfig{
		ModelPath:  modelLoc,
		NumThreads: 2,
	}
	ctx := context.Background()
	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, name)
	test.That(b, err, test.ShouldBeNil)
	test.That(b, out, test.ShouldNotBeNil)
	modelCfg := MLModelConfig{ModelName: name}

	service, err := registerMLModelVisionService(ctx, name, &modelCfg, &inject.Robot{}, golog.NewLogger("benchmark"))
	test.That(b, err, test.ShouldBeNil)
	test.That(b, service, test.ShouldNotBeNil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Detections should be worst case (more to unpack)
		detections, err := service.Detections(ctx, pic, nil)
		test.That(b, err, test.ShouldBeNil)
		test.That(b, detections, test.ShouldNotBeNil)
	}
}

func TestNewMLDetector(t *testing.T) {
	// Test that a detector would give an expected output on the dog image
	// Set it up as a ML Model

	ctx := context.Background()
	logger := golog.NewLogger("testLogger")
	modelLoc := artifact.MustPath("vision/tflite/effdet0.tflite")
	labelLoc := artifact.MustPath("vision/tflite/effdetlabels.txt")
	cfg := tflitecpu.TFLiteConfig{ // detector config
		ModelPath:  modelLoc,
		NumThreads: 2,
		LabelPath:  &labelLoc,
	}
	noLabelCfg := tflitecpu.TFLiteConfig{ // detector config
		ModelPath:  modelLoc,
		NumThreads: 2,
	}

	pic, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/dogscute.jpeg"))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, pic, test.ShouldNotBeNil)

	// Test that a detector would give the expected output on the dog image
	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, "myMLDet")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, out, test.ShouldNotBeNil)
	check, err := out.Metadata(ctx)
	test.That(t, check, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, check.Inputs[0].Name, test.ShouldResemble, "image")
	test.That(t, check.Outputs[0].Name, test.ShouldResemble, "location")
	test.That(t, check.Outputs[1].Name, test.ShouldResemble, "category")
	test.That(t, check.Outputs[1].Extra["labels"], test.ShouldNotBeNil)

	gotDetector, err := attemptToBuildDetector(out, logger)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotDetector, test.ShouldNotBeNil)

	gotDetections, err := gotDetector(ctx, pic)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotDetections[0].Score(), test.ShouldBeGreaterThan, 0.789)
	test.That(t, gotDetections[1].Score(), test.ShouldBeGreaterThan, 0.7)
	test.That(t, gotDetections[0].BoundingBox().Min.X, test.ShouldEqual, 126)
	test.That(t, gotDetections[0].BoundingBox().Min.Y, test.ShouldEqual, 42)
	test.That(t, gotDetections[0].BoundingBox().Max.X, test.ShouldEqual, 199)
	test.That(t, gotDetections[0].BoundingBox().Max.Y, test.ShouldEqual, 162)
	test.That(t, gotDetections[0].Label(), test.ShouldResemble, "Dog")
	test.That(t, gotDetections[1].Label(), test.ShouldResemble, "Dog")

	// Ensure that the same model without labelpath responds similarly
	outNL, err := tflitecpu.NewTFLiteCPUModel(ctx, &noLabelCfg, "myOtherMLDet")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outNL, test.ShouldNotBeNil)
	gotDetectorNL, err := attemptToBuildDetector(outNL, logger)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotDetectorNL, test.ShouldNotBeNil)
	gotDetectionsNL, err := gotDetectorNL(ctx, pic)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotDetectionsNL[0].Score(), test.ShouldBeGreaterThan, 0.789)
	test.That(t, gotDetectionsNL[1].Score(), test.ShouldBeGreaterThan, 0.7)
	test.That(t, gotDetectionsNL[0].BoundingBox().Min.X, test.ShouldEqual, 126)
	test.That(t, gotDetectionsNL[0].BoundingBox().Min.Y, test.ShouldEqual, 42)
	test.That(t, gotDetectionsNL[0].BoundingBox().Max.X, test.ShouldEqual, 199)
	test.That(t, gotDetectionsNL[0].BoundingBox().Max.Y, test.ShouldEqual, 162)
	test.That(t, gotDetectionsNL[0].Label(), test.ShouldResemble, "17")
	test.That(t, gotDetectionsNL[1].Label(), test.ShouldResemble, "17")
}

func TestNewMLClassifier(t *testing.T) {
	ctx := context.Background()
	logger := golog.NewLogger("testLogger")
	modelLoc := artifact.MustPath("vision/tflite/effnet0.tflite")
	labelLoc := artifact.MustPath("vision/tflite/imagenetlabels.txt")

	cfg := tflitecpu.TFLiteConfig{ // detector config
		ModelPath:  modelLoc,
		NumThreads: 2,
		LabelPath:  &labelLoc,
	}
	noLabelCfg := tflitecpu.TFLiteConfig{ // detector config
		ModelPath:  modelLoc,
		NumThreads: 2,
	}
	pic, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/lion.jpeg"))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, pic, test.ShouldNotBeNil)

	// Test that a classifier would give the expected result on the lion image
	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, "myMLClassif")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, out, test.ShouldNotBeNil)
	check, err := out.Metadata(ctx)
	test.That(t, check, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, check.Inputs[0].Name, test.ShouldResemble, "image")
	test.That(t, check.Outputs[0].Name, test.ShouldResemble, "probability")
	test.That(t, check.Outputs[0].Extra["labels"], test.ShouldNotBeNil)

	gotClassifier, err := attemptToBuildClassifier(out, logger)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotClassifier, test.ShouldNotBeNil)

	gotClassifications, err := gotClassifier(ctx, pic)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotClassifications, test.ShouldNotBeNil)
	gotTop, err := gotClassifications.TopN(5)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotTop, test.ShouldNotBeNil)
	test.That(t, gotTop[0].Label(), test.ShouldContainSubstring, "lion")
	test.That(t, gotTop[0].Score(), test.ShouldBeGreaterThan, 0.99)
	test.That(t, gotTop[1].Score(), test.ShouldBeLessThan, 0.01)

	// Ensure that the same model without labelpath responds similarly
	outNL, err := tflitecpu.NewTFLiteCPUModel(ctx, &noLabelCfg, "myOtherMLClassif")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outNL, test.ShouldNotBeNil)
	gotClassifierNL, err := attemptToBuildClassifier(outNL, logger)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotClassifierNL, test.ShouldNotBeNil)
	gotClassificationsNL, err := gotClassifierNL(ctx, pic)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotClassificationsNL, test.ShouldNotBeNil)
	topNL, err := gotClassificationsNL.TopN(5)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, topNL, test.ShouldNotBeNil)
	test.That(t, topNL[0].Label(), test.ShouldContainSubstring, "291")
	test.That(t, topNL[0].Score(), test.ShouldBeGreaterThan, 0.99)
	test.That(t, topNL[1].Score(), test.ShouldBeLessThan, 0.01)
}

func TestMoreMLDetectors(t *testing.T) {
	// Test that a detector would give an expected output on the dog image
	pic, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/dogscute.jpeg"))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, pic, test.ShouldNotBeNil)

	name := "ssd"
	ctx := context.Background()
	modelLoc := artifact.MustPath("vision/tflite/ssdmobilenet.tflite")
	labelLoc := artifact.MustPath("vision/tflite/effdetlabels.txt")
	cfg := tflitecpu.TFLiteConfig{
		ModelPath:  modelLoc,
		NumThreads: 2,
		LabelPath:  &labelLoc,
	}

	// Test that a detector would give the expected output on the dog image
	outModel, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, name)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outModel, test.ShouldNotBeNil)
	check, err := outModel.Metadata(ctx)
	test.That(t, check, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeNil)
	// Even without metadata we should find
	test.That(t, check.Inputs[0].Shape, test.ShouldResemble, []int{1, 320, 320, 3})
	test.That(t, check.Inputs[0].DataType, test.ShouldResemble, "float32")

	gotDetector, err := attemptToBuildDetector(outModel, golog.NewLogger("testLogger"))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotDetector, test.ShouldNotBeNil)

	gotDetections, err := gotDetector(ctx, pic)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, gotDetections[0].Score(), test.ShouldBeGreaterThan, 0.82)
	test.That(t, gotDetections[1].Score(), test.ShouldBeGreaterThan, 0.8)
	test.That(t, gotDetections[0].Label(), test.ShouldResemble, "17")
	test.That(t, gotDetections[1].Label(), test.ShouldResemble, "17")

}

func TestMoreMLClassifiers(t *testing.T) {

}

func TestLabelReader(t *testing.T) {
	ctx := context.Background()
	modelLoc := artifact.MustPath("vision/tflite/effdet0.tflite")
	labelLoc := artifact.MustPath("vision/tflite/fakelabels.txt")
	cfg := tflitecpu.TFLiteConfig{ // detector config
		ModelPath:  modelLoc,
		NumThreads: 2,
		LabelPath:  &labelLoc,
	}
	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, "fakeLabels")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, out, test.ShouldNotBeNil)
	outMD, err := out.Metadata(ctx)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outMD, test.ShouldNotBeNil)
	outLabels := getLabelsFromMetadata(outMD)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outLabels[0], test.ShouldResemble, "this")
	test.That(t, outLabels[1], test.ShouldResemble, "could")
	test.That(t, outLabels[2], test.ShouldResemble, "be")
	test.That(t, len(outLabels), test.ShouldEqual, 12)
}

func TestSpaceDelineatedLabels(t *testing.T) {
	ctx := context.Background()
	modelLoc := artifact.MustPath("vision/tflite/effdet0.tflite")
	labelLoc := artifact.MustPath("vision/classification/lorem.txt")
	cfg := tflitecpu.TFLiteConfig{ // detector config
		ModelPath:  modelLoc,
		NumThreads: 2,
		LabelPath:  &labelLoc,
	}
	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, "spacedLabels")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, out, test.ShouldNotBeNil)
	outMD, err := out.Metadata(ctx)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outMD, test.ShouldNotBeNil)
	outLabels := getLabelsFromMetadata(outMD)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, len(outLabels), test.ShouldEqual, 10)
}

func TestOneClassifierOnManyCameras(t *testing.T) {
	ctx := context.Background()

	// Test that one classifier can be used in two goroutines
	picPanda, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/redpanda.jpeg"))
	test.That(t, err, test.ShouldBeNil)
	picLion, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/lion.jpeg"))
	test.That(t, err, test.ShouldBeNil)

	modelLoc := artifact.MustPath("vision/tflite/mobilenetv2_class.tflite")
	cfg := tflitecpu.TFLiteConfig{
		ModelPath:  modelLoc,
		NumThreads: 2,
	}

	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, "testClassifier")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, out, test.ShouldNotBeNil)
	outClassifier, err := attemptToBuildClassifier(out, golog.NewLogger("testLogger"))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outClassifier, test.ShouldNotBeNil)
	valuePanda, valueLion := classifyTwoImages(picPanda, picLion, outClassifier)
	test.That(t, valuePanda, test.ShouldNotBeNil)
	test.That(t, valueLion, test.ShouldNotBeNil)
}

func TestMultipleClassifiersOneModel(t *testing.T) {
	modelLoc := artifact.MustPath("vision/tflite/mobilenetv2_class.tflite")
	cfg := tflitecpu.TFLiteConfig{
		ModelPath:  modelLoc,
		NumThreads: 2,
	}
	ctx := context.Background()
	out, err := tflitecpu.NewTFLiteCPUModel(ctx, &cfg, "testClassifier")
	test.That(t, err, test.ShouldBeNil)

	Classifier1, err := attemptToBuildClassifier(out)
	test.That(t, err, test.ShouldBeNil)

	Classifier2, err := attemptToBuildClassifier(out)
	test.That(t, err, test.ShouldBeNil)

	picPanda, err := rimage.NewImageFromFile(artifact.MustPath("vision/tflite/redpanda.jpeg"))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		getNClassifications(t, ctx, picPanda, 10, Classifier1)
	}()

	go func() {
		defer wg.Done()
		getNClassifications(t, ctx, picPanda, 10, Classifier2)
	}()

	wg.Wait()
}

func classifyTwoImages(picPanda, picLion *rimage.Image,
	got classification.Classifier,
) (classification.Classifications, classification.Classifications) {
	resultPanda := make(chan classification.Classifications)
	resultLion := make(chan classification.Classifications)

	go gotWithCallback(picPanda, resultPanda, got)
	go gotWithCallback(picLion, resultLion, got)

	valuePanda := <-resultPanda
	valueLion := <-resultLion

	close(resultPanda)
	close(resultLion)

	return valuePanda, valueLion
}

func gotWithCallback(img *rimage.Image, result chan classification.Classifications, got classification.Classifier) {
	classifications, _ := got(context.Background(), img)
	result <- classifications
}

func getNClassifications(t *testing.T, ctx context.Context, img *rimage.Image, N int, c classification.Classifier) []classification.Classifications {
	t.Helper()
	results := make([]classification.Classifications, N)
	var err error

	for i := 0; i < N; i++ {
		results[i], err = c(ctx, img)
		test.That(t, err, test.ShouldBeNil)
		res, err := results[i].TopN(1)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, res[0].Score(), test.ShouldNotBeNil)
	}
	return results
}
