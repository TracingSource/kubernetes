package testing

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type FakeImageService struct {
	sync.Mutex

	FakeImageSize uint64
	Called        []string
	Images        map[string]*runtimeapi.Image

	pulledImages []*pulledImage

	FakeFilesystemUsage []*runtimeapi.FilesystemUsage
}

func (r *FakeImageService) SetFakeImages(images []string) {
	r.Lock()
	defer r.Unlock()

	r.Images = make(map[string]*runtimeapi.Image)
	for _, image := range images {
		r.Images[image] = r.makeFakeImage(image)
	}
}

func (r *FakeImageService) SetFakeImageSize(size uint64) {
	r.Lock()
	defer r.Unlock()

	r.FakeImageSize = size
}

func (r *FakeImageService) SetFakeFilesystemUsage(usage []*runtimeapi.FilesystemUsage) {
	r.Lock()
	defer r.Unlock()

	r.FakeFilesystemUsage = usage
}

func NewFakeImageService() *FakeImageService {
	return &FakeImageService{
		Called: make([]string, 0),
		Images: make(map[string]*runtimeapi.Image),
	}
}

func (r *FakeImageService) makeFakeImage(image string) *runtimeapi.Image {
	return &runtimeapi.Image{
		Id:       image,
		Size_:    r.FakeImageSize,
		RepoTags: []string{image},
	}
}

// stringInSlice returns true if s is in list
func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}

	return false
}

func (r *FakeImageService) ListImages(filter *runtimeapi.ImageFilter) ([]*runtimeapi.Image, error) {
	r.Lock()
	defer r.Unlock()

	r.Called = append(r.Called, "ListImages")

	images := make([]*runtimeapi.Image, 0)
	for _, img := range r.Images {
		if filter != nil && filter.Image != nil {
			if !stringInSlice(filter.Image.Image, img.RepoTags) {
				continue
			}
		}

		images = append(images, img)
	}
	return images, nil
}

func (r *FakeImageService) ImageStatus(image *runtimeapi.ImageSpec) (*runtimeapi.Image, error) {
	r.Lock()
	defer r.Unlock()

	r.Called = append(r.Called, "ImageStatus")

	return r.Images[image.Image], nil
}

func (r *FakeImageService) PullImage(image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	r.Lock()
	defer r.Unlock()

	r.Called = append(r.Called, "PullImage")

	r.pulledImages = append(r.pulledImages, &pulledImage{imageSpec: image, authConfig: auth})
	// ImageID should be randomized for real container runtime, but here just use
	// image's name for easily making fake images.
	imageID := image.Image
	if _, ok := r.Images[imageID]; !ok {
		r.Images[imageID] = r.makeFakeImage(image.Image)
	}

	return imageID, nil
}

func (r *FakeImageService) RemoveImage(image *runtimeapi.ImageSpec) error {
	r.Lock()
	defer r.Unlock()

	r.Called = append(r.Called, "RemoveImage")

	// Remove the image
	delete(r.Images, image.Image)

	return nil
}

// ImageFsInfo returns information of the filesystem that is used to store images.
func (r *FakeImageService) ImageFsInfo() ([]*runtimeapi.FilesystemUsage, error) {
	r.Lock()
	defer r.Unlock()

	r.Called = append(r.Called, "ImageFsInfo")

	return r.FakeFilesystemUsage, nil
}

func (r *FakeImageService) AssertImagePulledWithAuth(t *testing.T, image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, failMsg string) {
	r.Lock()
	defer r.Unlock()
	expected := &pulledImage{imageSpec: image, authConfig: auth}
	assert.Contains(t, r.pulledImages, expected, failMsg)
}

type pulledImage struct {
	imageSpec  *runtimeapi.ImageSpec
	authConfig *runtimeapi.AuthConfig
}
