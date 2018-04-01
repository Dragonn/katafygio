package namespace

import (
	"fmt"
	"strings"

	"github.com/bpineau/katafygio/config"
	"github.com/bpineau/katafygio/pkg/controllers"

	"k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/ghodss/yaml"
)

const objectKind = "Namespace"

// Controller monitors Kubernetes' namespaces objects in the cluster
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
	c.ObjType = &v1.Namespace{}
	c.MarshalF = c.Marshal
	selector := meta.ListOptions{LabelSelector: conf.Filter}
	c.ListWatch = &cache.ListWatch{
		ListFunc: func(options meta.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Namespaces().List(selector)
		},
		WatchFunc: func(options meta.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Namespaces().Watch(selector)
		},
	}

	return &c
}

// Marshal filter irrelevant fields from the object, and export it as yaml string
func (c *Controller) Marshal(obj interface{}) (string, error) {
	f := obj.(*v1.Namespace).DeepCopy()

	// Should we also remove selfLink and resourceVersion ?
	f.Status.Reset()
	f.ResourceVersion = ""
	f.SelfLink = ""
	f.UID = ""

	y, err := yaml.Marshal(f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("apiVersion: %s\nkind: %s\n%s",
		v1.SchemeGroupVersion.String(), objectKind, string(y)), nil
}
