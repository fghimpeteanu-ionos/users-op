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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	userv1 "github.com/fghimpeteanu-ionos/user-op/api/v1"
)

const waitTime = 15 * time.Second

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	userPersistence UserPersistence
	log             *logr.Logger
}

//+kubebuilder:rbac:groups=filip.org,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=filip.org,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=filip.org,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.setLogger(ctx)

	user := &userv1.User{}
	err := r.Client.Get(ctx, req.NamespacedName, user)
	if err != nil {
		r.logInfoFmt("Unable to fetch User with name: %s (might have been deleted)", req.Name)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	isUserPersisted := isUserPersisted(user, r)
	if !isUserPersisted && isUserStatusNotSet(user) || (!isUserPersisted && isUserReady(user)) {
		user.Status.State = userv1.CREATING
		r.logInfoFmt("User not persisted. Will reconcile ...")
		err = r.updateUserCR(ctx, err, user)
		return ctrl.Result{Requeue: true}, err
	}

	if isUserReady(user) {
		r.logInfoFmt("Nothing to reconcile. Will go to sleep for %s ...", waitTime.String())
		return ctrl.Result{RequeueAfter: waitTime}, nil
	}

	r.log.Info("Reconciling User", "User Spec", user.Spec)
	err = r.addUserInDB(user)
	if err != nil {
		r.log.Error(err, "Unable to add user in DB")
		user.Status.State = userv1.FAILED
		err = r.updateUserCR(ctx, err, user)
		return ctrl.Result{}, err
	}

	user.Status.State = userv1.READY
	err = r.updateUserCR(ctx, err, user)
	return ctrl.Result{}, err
}

func (r *UserReconciler) logInfoFmt(fmtMsg string, args ...any) {
	r.log.Info(fmt.Sprintf(fmtMsg, args...))
}

func isUserStatusNotSet(u *userv1.User) bool {
	return u.Status.State == ""
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.userPersistence = &UserPersistenceImpl{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&userv1.User{}).
		Complete(r)
}

func (r *UserReconciler) setLogger(ctx context.Context) {
	if r.log == nil {
		loggerFromContext := log.FromContext(ctx)
		r.log = &loggerFromContext
	}
}
func (r *UserReconciler) addUserInDB(user *userv1.User) error {
	createUserE := toCreateUserE(user)
	uuid, err := r.userPersistence.Persist(createUserE)
	if err != nil {
		return err
	}
	user.Status.UUID = uuid
	return nil
}

func (r *UserReconciler) updateUserCR(ctx context.Context, err error, user *userv1.User) error {
	err = r.Client.Status().Update(ctx, user)
	if err != nil {
		r.log.Error(err, "Unable to update User resource")
		return err
	}
	return nil
}
func isUserPersisted(user *userv1.User, r *UserReconciler) bool {
	userUUID := user.Status.UUID
	if userUUID == "" {
		return false
	}

	readUserE, err := r.userPersistence.Read(userUUID)
	if err != nil {
		r.log.Error(err, "Unable to read user from DB")
		return false
	}
	return readUserE != nil
}

func isUserReady(user *userv1.User) bool {
	return user.Status.State == userv1.READY
}

func toCreateUserE(user *userv1.User) *CreateUserE {
	createUserE := &CreateUserE{
		firstName: user.Spec.FirstName,
		lastName:  user.Spec.LastName,
		age:       user.Spec.Age,
		address:   user.Spec.Address,
		email:     user.Spec.Email,
	}
	return createUserE
}
