package resourcehandler

import (
	"encoding/json"
	"fmt"
	"os"

	quobytev1 "operator/pkg/kubernetes-actors/clientset/versioned"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/util/strategicpatch"

	"github.com/golang/glog"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

// TODO: no need of yaml to json, use yaml to interface unmarshal. See Namespace creation createQuobyteNameSpace.
var (
	config           *rest.Config
	QclientConfig    *quobytev1.Clientset
	err              error
	KubernetesClient *kubernetes.Clientset
	APIServerClient  *apiextensionsclient.Clientset
	quobyteNameSpace = "quobyte"
)

// GetQuobyteNodes get quobyte nodes
// TODO: only get quobyte nodes, currently gets all the nodes
func GetQuobyteNodes() (*v1.NodeList, error) {
	nodeList, err := KubernetesClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing nodes: %v", err)
		return nodeList, err
	}
	return nodeList, err
}

// GetPods Returns all the pods in cluster
func GetPods(client *kubernetes.Clientset) (*v1.PodList, error) {
	podList, err := client.CoreV1().Pods(quobyteNameSpace).List(metav1.ListOptions{})
	return podList, err
}

// LabelNodes updates labels on nodes.
// op add,remove
// TODO: handle GetQuobyteNodes error, and change the way node retrieval
func LabelNodes(labelNodes []string, op, label string) {

	for _, nodeName := range labelNodes {
		if nodeName == "" {
			break
		}
		node, err := KubernetesClient.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("Failed to get the node %s due to %s", nodeName, err)
		}
		oldData, _ := json.Marshal(node)
		labels := node.GetLabels()
		switch op {
		case "add":
			labels[label] = "true"
		case "remove":
			delete(labels, label)
		}
		node.SetLabels(labels)
		newJSON, _ := json.Marshal(node)
		patchbytes, _ := strategicpatch.CreateTwoWayMergePatch(oldData, newJSON, v1.Node{})
		if len(patchbytes) > 2 {
			_, err = KubernetesClient.CoreV1().Nodes().Patch(node.Name, types.StrategicMergePatchType, patchbytes)
			if err != nil {
				glog.Errorf("Failed labeling node %s", node.Name)
			}
		}
	}
}
