// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/openshift/hive/apis/hiveinternal/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FakeClusterInstallLister helps list FakeClusterInstalls.
// All objects returned here must be treated as read-only.
type FakeClusterInstallLister interface {
	// List lists all FakeClusterInstalls in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FakeClusterInstall, err error)
	// FakeClusterInstalls returns an object that can list and get FakeClusterInstalls.
	FakeClusterInstalls(namespace string) FakeClusterInstallNamespaceLister
	FakeClusterInstallListerExpansion
}

// fakeClusterInstallLister implements the FakeClusterInstallLister interface.
type fakeClusterInstallLister struct {
	indexer cache.Indexer
}

// NewFakeClusterInstallLister returns a new FakeClusterInstallLister.
func NewFakeClusterInstallLister(indexer cache.Indexer) FakeClusterInstallLister {
	return &fakeClusterInstallLister{indexer: indexer}
}

// List lists all FakeClusterInstalls in the indexer.
func (s *fakeClusterInstallLister) List(selector labels.Selector) (ret []*v1alpha1.FakeClusterInstall, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FakeClusterInstall))
	})
	return ret, err
}

// FakeClusterInstalls returns an object that can list and get FakeClusterInstalls.
func (s *fakeClusterInstallLister) FakeClusterInstalls(namespace string) FakeClusterInstallNamespaceLister {
	return fakeClusterInstallNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FakeClusterInstallNamespaceLister helps list and get FakeClusterInstalls.
// All objects returned here must be treated as read-only.
type FakeClusterInstallNamespaceLister interface {
	// List lists all FakeClusterInstalls in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FakeClusterInstall, err error)
	// Get retrieves the FakeClusterInstall from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FakeClusterInstall, error)
	FakeClusterInstallNamespaceListerExpansion
}

// fakeClusterInstallNamespaceLister implements the FakeClusterInstallNamespaceLister
// interface.
type fakeClusterInstallNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FakeClusterInstalls in the indexer for a given namespace.
func (s fakeClusterInstallNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FakeClusterInstall, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FakeClusterInstall))
	})
	return ret, err
}

// Get retrieves the FakeClusterInstall from the indexer for a given namespace and name.
func (s fakeClusterInstallNamespaceLister) Get(name string) (*v1alpha1.FakeClusterInstall, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fakeclusterinstall"), name)
	}
	return obj.(*v1alpha1.FakeClusterInstall), nil
}
