/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	userv1 "github.com/fghimpeteanu-ionos/user-op/api/v1"
)

var _ = Describe("User Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		user := &userv1.User{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind User")
			err := k8sClient.Get(ctx, typeNamespacedName, user)
			if err != nil && errors.IsNotFound(err) {
				resource := &userv1.User{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &userv1.User{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance User")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &UserReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})

func TestUserReconciler_Reconcile(t *testing.T) {
	type fields struct {
		Client          client.Client
		Scheme          *runtime.Scheme
		userPersistence UserPersistence
	}
	type args struct {
		ctx context.Context
		req controllerruntime.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    controllerruntime.Result
		wantErr bool
	}{
		{
			name: "Test User creation",
			fields: fields{
				Client:          stubK8SClient(),
				Scheme:          runtime.NewScheme(),
				userPersistence: stubUserPersistence(),
			},
			args: args{
				ctx: context.Background(),
				req: controllerruntime.Request{
					NamespacedName: types.NamespacedName{
						Name: "user1",
					},
				},
			},
			want:    controllerruntime.Result{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserReconciler{
				Client:          tt.fields.Client,
				Scheme:          tt.fields.Scheme,
				userPersistence: tt.fields.userPersistence,
			}
			got, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type k8sClientStub struct {
	client.Client
	mock.Mock
}

func (k *k8sClientStub) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	args := k.Called(ctx, key, obj, opts)
	return args.Error(0)
}

func (k *k8sClientStub) Status() client.SubResourceWriter {
	args := k.Called()
	return args.Get(0).(client.SubResourceWriter)
}

func stubK8SClient() client.Client {
	k8sClient := new(k8sClientStub)
	k8sClient.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return k8sClient
}

type UserPersistenceStub struct {
	UserPersistence
	mock.Mock
}

func (u *UserPersistenceStub) Persist(user *userv1.User) (string, error) {
	args := u.Called(user)
	return args.String(0), args.Error(1)
}

func (u *UserPersistenceStub) Read(uuid string) (*userv1.User, error) {
	args := u.Called(uuid)
	return args.Get(0).(*userv1.User), args.Error(1)
}

func stubUserPersistence() *UserPersistenceStub {
	userPersistence := new(UserPersistenceStub)
	userPersistence.On("Read", mock.Anything).Return(basicUser(), nil)
	userPersistence.On("Persist", mock.Anything).Return("uuid", nil)
	return userPersistence
}

func basicUser() *userv1.User {
	return &userv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "user1",
		},
		Spec: userv1.UserSpec{
			FirstName: "Mike",
			LastName:  "Davidson",
			Age:       34,
			Address:   "Here and there Street, 1234",
			Email:     "mike.davidson@mail.com",
		},
	}
}
