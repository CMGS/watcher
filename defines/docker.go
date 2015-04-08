package defines

import (
	"reflect"

	"../logs"
	"../utils"
	"github.com/fsouza/go-dockerclient"
)

type DockerWrapper struct {
	*docker.Client
	PushImage        func(docker.PushImageOptions, docker.AuthConfiguration) error
	PullImage        func(docker.PullImageOptions, docker.AuthConfiguration) error
	CreateContainer  func(docker.CreateContainerOptions) (*docker.Container, error)
	StartContainer   func(string, *docker.HostConfig) error
	BuildImage       func(docker.BuildImageOptions) error
	KillContainer    func(docker.KillContainerOptions) error
	StopContainer    func(string, uint) error
	InspectContainer func(string) (*docker.Container, error)
	ListContainers   func(docker.ListContainersOptions) ([]docker.APIContainers, error)
	ListImages       func(docker.ListImagesOptions) ([]docker.APIImages, error)
	RemoveContainer  func(docker.RemoveContainerOptions) error
	WaitContainer    func(string) (int, error)
	RemoveImage      func(string) error
	CreateExec       func(docker.CreateExecOptions) (*docker.Exec, error)
	StartExec        func(string, docker.StartExecOptions) error
}

func NewDocker(endpoint string) *DockerWrapper {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		logs.Assert(err, "Docker")
	}
	d := &DockerWrapper{Client: client}
	var makeDockerWrapper func(*DockerWrapper, *docker.Client) *DockerWrapper
	utils.MakeWrapper(&makeDockerWrapper)
	return makeDockerWrapper(d, client)
}

func MockDocker(d *DockerWrapper) {
	var makeMockedDockerWrapper func(*DockerWrapper, *docker.Client) *DockerWrapper
	MakeMockedWrapper(&makeMockedDockerWrapper)
	makeMockedDockerWrapper(d, d.Client)
}

func MakeMockedWrapper(fptr interface{}) {
	var maker = func(in []reflect.Value) []reflect.Value {
		wrapper := in[0].Elem()
		client := in[1]
		wrapperType := wrapper.Type()
		for i := 1; i < wrapperType.NumField(); i++ {
			field := wrapper.Field(i)
			fd, ok := client.Type().MethodByName(wrapperType.Field(i).Name)
			if !ok {
				logs.Info("Reflect Failed")
				continue
			}
			fdt := fd.Type
			f := reflect.MakeFunc(field.Type(), func(in []reflect.Value) []reflect.Value {
				ret := make([]reflect.Value, 0, fdt.NumOut())
				for i := 0; i < fdt.NumOut(); i++ {
					ret = append(ret, reflect.Zero(fdt.Out(i)))
				}
				return ret
			})
			field.Set(f)
		}
		return []reflect.Value{in[0]}
	}
	fn := reflect.ValueOf(fptr).Elem()
	v := reflect.MakeFunc(fn.Type(), maker)
	fn.Set(v)
}
