package replicaset

import (
	"fmt"
	"strings"

	"github.com/bpineau/katafygio/config"
	"github.com/bpineau/katafygio/pkg/controllers"

	v1 "k8s.io/api/apps/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/ghodss/yaml"
)

const objectKind = "ReplicaSet"

// Controller monitors Kubernetes' deployments objects in the cluster
type Controller struct {
	controllers.CommonController
}

func init() {
	controllers.Register(strings.ToLower(objectKind), New)
}

// New initialize controller
func New(conf *config.KdnConfig, ch chan<- controllers.Event) controllers.Controller {
	c := Controller{}
	c.CommonController = controllers.CommonController{
		Conf: conf,
		Name: strings.ToLower(objectKind),
		Send: ch,
	}

	client := c.Conf.ClientSet
	c.ObjType = &v1.ReplicaSet{}
	c.MarshalF = c.Marshal
	selector := meta.ListOptions{LabelSelector: conf.Filter}
	c.ListWatch = &cache.ListWatch{
		ListFunc: func(options meta.ListOptions) (runtime.Object, error) {
			return client.AppsV1().ReplicaSets(meta.NamespaceAll).List(selector)
		},
		WatchFunc: func(options meta.ListOptions) (watch.Interface, error) {
			return client.AppsV1().ReplicaSets(meta.NamespaceAll).Watch(selector)
		},
	}

	return &c
}

// Marshal filter irrelevant fields from the object, and export it as yaml string
func (c *Controller) Marshal(obj interface{}) (string, error) {
	f := obj.(*v1.ReplicaSet).DeepCopy()

	f.Status.Reset()
	f.ResourceVersion = ""
	f.Generation = 0
	f.SelfLink = ""
	f.UID = ""

	y, err := yaml.Marshal(f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("apiVersion: %s\nkind: %s\n%s",
		v1.SchemeGroupVersion.String(), objectKind, string(y)), nil
}
