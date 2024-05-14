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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	userv1 "github.com/fghimpeteanu-ionos/user-op/api/v1"
)

const waitTime = 10 * time.Second

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	userPersistence UserPersistence
}

//+kubebuilder:rbac:groups=filip.org,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=filip.org,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=filip.org,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	user := &userv1.User{}
	err := r.Client.Get(ctx, req.NamespacedName, user)
	if err != nil {
		logger.Error(err, "unable to fetch User with name: "+req.Name)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if user.Status.State == userv1.READY {
		logger.Info("Nothing to reconcile. Will go to sleep for " + waitTime.String() + " ...")
		return ctrl.Result{RequeueAfter: waitTime}, nil
	}

	logger.Info("Reconciling User", "User Spec", user.Spec)
	user.Status.State = userv1.PENDING

	err = r.addUserInDB(user)
	if err != nil {
		logger.Error(err, "unable to add user in DB")
		user.Status.State = userv1.FAILED
		return ctrl.Result{}, err
	}

	err = r.Client.Status().Update(ctx, user)
	if err != nil {
		logger.Error(err, "unable to update User resource")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *UserReconciler) addUserInDB(user *userv1.User) error {
	uuid, err := r.userPersistence.Persist(user)
	if err != nil {
		return err
	}
	user.Status.UUID = uuid
	user.Status.State = userv1.READY
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.userPersistence = &UserPersistenceImpl{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&userv1.User{}).
		Complete(r)
}
