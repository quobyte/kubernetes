package controller

import (
	"encoding/json"
	quobytev1 "operator/pkg/kubernetes-actors/clientset/versioned"
	quobyte_crd_informer_factory "operator/pkg/kubernetes-actors/informers/externalversions"
	"operator/pkg/resourcehandler"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	glog "github.com/golang/glog"
	"github.com/mattbaird/jsonpatch"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"operator/pkg/api/quobyte.com/v1"
)

type controller struct {
	clientinformer cache.SharedIndexInformer
}

func (controller *controller) Run(stopCh chan struct{}) {
	defer utilruntime.HandleCrash()
	go controller.clientinformer.Run(stopCh)
}

func newQuobyteClientController(quobyteclient *quobytev1.Clientset) (*controller, error) {
	informerFactory := quobyte_crd_informer_factory.NewSharedInformerFactory(quobyteclient, 5*time.Minute)
	quobyteclientInformer := informerFactory.Quobyte().V1().QuobyteClients().Informer()

	quobyteclientInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			glog.Info("Client CR create event triggered")
			handleClientCRAdd(obj)
			// err := resourcehandler.CreateClientDaemonSet(obj.(*v1.QuobyteClient).Spec.Version)
			// if err != nil {
			// 	glog.Errorf("Failed to create client daemonset due to :\n %s", err)
			// }
		},
		UpdateFunc: func(old, new interface{}) {
			glog.Info("Client CR update event triggered")
			handleClientCRUpdate(old, new)
		},
		DeleteFunc: func(obj interface{}) {
			glog.Info("Client CR delete event triggered")
			handleClientCRDelete(obj)
		},
	})
	controllerDef := &controller{
		clientinformer: quobyteclientInformer,
	}
	return controllerDef, nil
}

func handleClientCRAdd(obj interface{}) {
	deployedClientDS, err := resourcehandler.GetDaemonsetByName("client")
	if err != nil {
		glog.Errorf("Failed to get client daemonset %v", err)
	}
	deployedClientLabels := deployedClientDS.GetLabels()
	version, ok := deployedClientLabels["version"]
	if ok {
		if obj.(*v1.QuobyteClient).Spec.Version != version {
			glog.Info("version change in CR, updating client daemonset to %v", version)
			err := resourcehandler.UpdateDaemonSet("client", obj.(*v1.QuobyteClient).Spec.Version) // TODO: get dynamic name
			if err != nil {
				glog.Errorf("Failed to update client daemonset with updated version %s\n", obj.(*v1.QuobyteClient).Spec.Version)
			}
		}
	}
	resourcehandler.LabelNodes(obj.(*v1.QuobyteClient).Spec.Nodes, "add", "quobyte_client")
}

func handleClientCRDelete(definition interface{}) {
	resourcehandler.LabelNodes(definition.(*v1.QuobyteClient).Spec.Nodes, "remove", "quobyte_client")

	// err := resourcehandler.DeleteQuobyteDeployment("client")

	// if err != nil {
	// 	glog.Errorf("Failed to delete client: %s\n", err)
	// }

	resourcehandler.DeleteClientPods(definition.(*v1.QuobyteClient).Spec.Version)
	glog.Info("Removed Quobyte client definition (Pod termination signalled)")
}

func handleClientCRUpdate(old, cur interface{}) {
	oldCr := old.(*v1.QuobyteClient)
	curCr := cur.(*v1.QuobyteClient)

	glog.Info("QuobyteClient CR updated: Updating Quobyte to the updated template.")

	if !reflect.DeepEqual(oldCr.Spec.Nodes, curCr.Spec.Nodes) {
		handleNodeChanges(oldCr, curCr, oldCr.Spec.Nodes, curCr.Spec.Nodes, "quobyte_client")
	}

	if oldCr.Spec.Version != curCr.Spec.Version {
		glog.Infof("Client version changed from %s to %s", oldCr.Spec.Version, curCr.Spec.Version)
		err := resourcehandler.UpdateDaemonSet("client", curCr.Spec.Version) // TODO: get dynamic name

		if err != nil {
			glog.Errorf("Failed to update client daemonset with updated version %s\n", curCr.Spec.Version)
		}
		resourcehandler.ControlledClientPodUpdate(resourcehandler.GetClientPodsByVersion(oldCr.Spec.Version))
	}
}

func handleNodeChanges(oldCr, curCr interface{}, oldNodes, curNodes []string, label string) {
	oldCrData, err := json.Marshal(oldCr)
	if err != nil {
		glog.Errorf("Failed Json conversion of old client CR")
	}
	curCrData, err := json.Marshal(curCr)

	if err != nil {
		glog.Errorf("Failed Json conversion of updated client CR")
	}

	patch, err := jsonpatch.CreatePatch([]byte(oldCrData), []byte(curCrData))
	if err != nil {
		glog.Errorf("Error creating JSON patch:%s", err)
		return
	}

	nodesPath := getNodesPath(label)
	removedNodes := make([]string, 0, len(oldNodes)) // at most, only nodes in old definition can be removed.
	addedNodes := make([]string, 0, len(curNodes))   // at most, only nodes in current definition must have been added.
	for _, operation := range patch {
		if strings.Contains(strings.ToLower(operation.Path), nodesPath) {
			if operation.Operation == "remove" {
				//get the removed node index, so that we don't remove/modify the unexpected node.
				index, _ := strconv.Atoi(operation.Path[strings.Index(strings.ToLower(operation.Path), nodesPath)+len(nodesPath) : len(operation.Path)])
				removedNode := oldNodes[index]
				removedNodes = append(removedNodes, removedNode)
			} else if operation.Operation == "add" {
				addedNodes = append(addedNodes, operation.Value.(string))
			} else if operation.Operation == "replace" {
				// remove old node and add new node
				index, _ := strconv.Atoi(operation.Path[strings.Index(strings.ToLower(operation.Path), nodesPath)+len(nodesPath) : len(operation.Path)])
				removedNode := oldNodes[index]
				removedNodes = append(removedNodes, removedNode)
				addedNodes = append(addedNodes, operation.Value.(string))
			}
		}
	}
	if len(removedNodes) > 0 {
		glog.Infof("Removed nodes:%s", removedNodes)
		resourcehandler.LabelNodes(removedNodes, "remove", label)
	}
	if len(addedNodes) > 0 {
		glog.Infof("Added nodes:%s", addedNodes)
		resourcehandler.LabelNodes(addedNodes, "add", label)
	}
}

func getNodesPath(label string) string {
	switch label {
	case "quobyte_client":
		return "/spec/nodes/"
	case "quobyte_registry":
		return "/spec/registry/nodes/"
	case "quobyte_data":
		return "/spec/data/nodes/"
	}
	return "none/node"
}

// Start starts quobyte crd monitoring
func Start(quobyteclient *quobytev1.Clientset) {
	controller, err := newQuobyteClientController(quobyteclient)
	if err != nil {
		glog.Errorf("Terminating operator due to:\n %v ", err)
	}
	stopCh := make(chan struct{})
	defer close(stopCh)

	go controller.Run(stopCh)

	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGTERM)
	signal.Notify(termSig, syscall.SIGINT)
	<-termSig
}
