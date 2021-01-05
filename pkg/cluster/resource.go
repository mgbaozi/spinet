package cluster

import (
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
)

type Namespace struct {
	Name  string
	Tasks map[string]*models.Task
}

func NewNamespace(name string) Namespace {
	return Namespace{
		name,
		make(map[string]*models.Task),
	}
}

func (ns Namespace) GetTask(name string) (*models.Task, error) {
	if task, ok := ns.Tasks[name]; ok {
		return task, nil
	} else {
		return nil, errors.New(fmt.Sprintf("task %s in namespace %s not found", name, ns.Name))
	}
}

func (ns Namespace) ListTasks() (res []*models.Task) {
	for _, task := range ns.Tasks {
		res = append(res, task)
	}
	return res
}

func (ns Namespace) CreateTask(task *models.Task) error {
	if _, ok := ns.Tasks[task.Name]; !ok {
		klog.V(2).Infof("Create new task %s in namespace %s", task.Name, ns.Name)
		ns.Tasks[task.Name] = task
		go task.Start()
		return nil
	} else {
		return errors.New(fmt.Sprintf("task `%s` in namespace `%s` is already exists", task.Name, ns.Name))
	}
}

type Resource map[string]Namespace

func NewResource() Resource {
	return make(Resource)
}

func (res Resource) CreateNamespace(name string) error {
	if _, ok := res[name]; !ok {
		res[name] = NewNamespace(name)
		return nil
	} else {
		return errors.New(fmt.Sprintf("namespace `%s` already exists", name))
	}
}

func (res Resource) GetNamespace(name string) (Namespace, error) {
	if name == "" {
		name = models.DefaultNamespace
	}
	if ns, ok := res[name]; ok {
		return ns, nil
	} else {
		return ns, errors.New(fmt.Sprintf("namespace `%s` not found", name))
	}
}

func (res Resource) ListNamespaces() (result []Namespace) {
	for _, value := range res {
		result = append(result, value)
	}
	return
}

func (res Resource) ListTasks(namespace string) ([]*models.Task, error) {
	if ns, err := res.GetNamespace(namespace); err != nil {
		return nil, err
	} else {
		return ns.ListTasks(), nil
	}
}

func (res Resource) GetTask(name, namespace string) (*models.Task, error) {
	if ns, err := res.GetNamespace(namespace); err != nil {
		return nil, err
	} else {
		return ns.GetTask(name)
	}
}

func (res Resource) CreateTask(task *models.Task) error {
	if ns, err := res.GetNamespace(task.Namespace); err != nil {
		return err
	} else {
		return ns.CreateTask(task)
	}
}
