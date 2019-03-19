package template

import (
	"encoding/base64"
	"fmt"
	"github.com/rancher/types/apis/management.cattle.io/v3"

	"github.com/rancher/rio/cli/pkg/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/norman/types/convert/schemaconvert"
	"github.com/rancher/rio/pkg/namespace"
	"github.com/rancher/rio/pkg/pretty"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1"
	riov1 "github.com/rancher/rio/types/apis/rio.cattle.io/v1"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FromClientStack(stack *riov1.Stack) (*Template, error) {
	stackSchema := schema.Schemas.Schema(&schema.Version, types.StackType)
	internalStack := &v1.Stack{}
	err := schemaconvert.ToInternal(stack, stackSchema, internalStack)
	if err != nil {
		return nil, err
	}

	return FromStack(internalStack)
}

func FromStack(stack *v1.Stack) (*Template, error) {
	result := &Template{
		Namespace:       namespace.StackToNamespace(stack),
		Content:         []byte(stack.Spec.Template),
		Answers:         map[string]string{},
		AdditionalFiles: map[string][]byte{},
	}

	for name, value := range stack.Spec.AdditionalFiles {
		content, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template [%s]: %v", name, err)
		}
		result.AdditionalFiles[name] = content
	}

	if stack.Spec.Answers != nil {
		result.Answers = stack.Spec.Answers
	}

	return result, nil
}

func (t *Template) ToStack(namespace, name string) *v1.Stack {
	s := &v1.Stack{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Stack",
			APIVersion: "rio.cattle.io/v1",
		},
		Spec: v1.StackSpec{
			Template:  string(t.Content),
			Answers:   t.Answers,
			Questions: t.Questions,
		},
	}

	for name, value := range t.AdditionalFiles {
		s.Spec.AdditionalFiles[name] = base64.StdEncoding.EncodeToString(value)
	}

	s.Name = name
	s.Namespace = namespace

	return s
}

func (t *Template) ToInternalStack() (*v1.InternalStack, error) {
	data, err := t.parseYAML()
	if err != nil {
		return nil, err
	}

	return pretty.ToInternalStack(data)
}

func (t *Template) ToClientStack() (*riov1.Stack, error) {
	result := &riov1.Stack{
		Spec: riov1.StackSpec{
			Answers:         t.Answers,
			Template:        string(t.Content),
			AdditionalFiles: map[string]string{},
		},
	}

	for name, content := range t.AdditionalFiles {
		encoded := base64.StdEncoding.EncodeToString(content)
		result.Spec.AdditionalFiles[name] = encoded
	}

	for _, q := range t.Questions {
		newQ := v3.Question{}
		if err := convert.ToObj(q, &newQ); err != nil {
			return nil, err
		}
		result.Spec.Questions = append(result.Spec.Questions, newQ)
	}

	return result, nil
}
