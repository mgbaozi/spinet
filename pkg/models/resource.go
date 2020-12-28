package models

import (
	"errors"
	"fmt"
)

type Namespace struct {
	Name  string
	Tasks map[string]*Task
}

func NewNamespace(name string) Namespace {
	return Namespace{
		name,
		make(map[string]*Task),
	}
}

func (ns Namespace) GetTask(name string) (*Task, error) {
	if task, ok := ns.Tasks[name]; ok {
		return task, nil
	} else {
		return nil, errors.New(fmt.Sprintf("task %s in namespace %s not found", name, ns.Name))
	}
}

func (ns Namespace) CreateTask(task *Task) error {
	if _, ok := ns.Tasks[task.Name]; !ok {
		ns.Tasks[task.Name] = task
		return nil
	} else {
		return errors.New(fmt.Sprintf("task %s in namespace %s is already exists", task.Name, ns.Name))
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
		return errors.New(fmt.Sprintf("namespace %s already exists", name))
	}
}

func (res Resource) GetNamespace(name string) (Namespace, error) {
	if ns, ok := res[name]; ok {
		return ns, nil
	} else {
		return ns, errors.New(fmt.Sprintf("namespace %s not found", name))
	}
}

func (res Resource) ListNamespaces() (result []Namespace) {
	for _, value := range res {
		result = append(result, value)
	}
	return
}

func (res Resource) GetTask(name, namespace string) (*Task, error) {
	if len(namespace) == 0 {
		namespace = "default"
	}
	if ns, err := res.GetNamespace(namespace); err != nil {
		return nil, err
	} else {
		return ns.GetTask(name)
	}
}

func (res Resource) CreateTask(task *Task) error {
	if len(task.Namespace) == 0 {
		task.Namespace = "default"
	}
	if ns, err := res.GetNamespace(task.Namespace); err != nil {
		return err
	} else {
		return ns.CreateTask(task)
	}
}
