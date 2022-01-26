/*
Copyright 2022.

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

package controllers

import (
	"bufio"
	"context"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	networkv1 "github.com/burnyd/networkconfigoperator/api/v1"
	eos "github.com/burnyd/networkconfigoperator/pkg/eosconfig"
)

// NetDevsReconciler reconciles a NetDevs object
type NetDevsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=network.networkconfigoperator,resources=netdevs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=network.networkconfigoperator,resources=netdevs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=network.networkconfigoperator,resources=netdevs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NetDevs object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile

func TrimThings(config string) string {
	var a string
	x := bufio.NewScanner(strings.NewReader(config))
	for x.Scan() {
		line := strings.TrimSpace(x.Text())
		if strings.HasPrefix(line, "!") {
			continue
		}
		if strings.HasPrefix(line, "\n") {
			continue
		}
		a += line + "\n"
	}
	return a
}

func (r *NetDevsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	_ = log.FromContext(ctx, "netdevs", req.Namespace)
	log.Log.Info("Starting to check on the device")
	netdevs := &networkv1.NetDevs{}
	err := r.Client.Get(ctx, req.NamespacedName, netdevs)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("Resource not found")
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "Failed to get Netdev controller")
		return ctrl.Result{}, err
	}
	if netdevs.Spec.NetworkOs == "eos" {
		log.Log.Info("Device is EOS")
		log.Log.Info("Checking on Device " + netdevs.Spec.Host)
		EosDev := eos.Conn{
			Transport: netdevs.Spec.Transport,
			Username:  netdevs.Spec.Username,
			Password:  netdevs.Spec.Password,
			Host:      netdevs.Spec.Host,
			Port:      netdevs.Spec.Port,
			Config:    netdevs.Spec.RunningConfig,
		}
		switchcfg := TrimThings(EosDev.Compare())
		crdcfg := TrimThings(netdevs.Spec.RunningConfig)
		if switchcfg == crdcfg {
			log.Log.Info("Switch config matches CR ")
		} else {
			log.Log.Info("Switch config does not match CR ")
			log.Log.Info("Preparing to push config to " + netdevs.ObjectMeta.Name)
			EosDev.Configure(netdevs.Spec.RunningConfig)
		}
	}

	return ctrl.Result{RequeueAfter: time.Second * 5}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *NetDevsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.NetDevs{}).
		Complete(r)
}
