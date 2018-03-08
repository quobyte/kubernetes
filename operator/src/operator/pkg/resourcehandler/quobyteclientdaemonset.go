package resourcehandler

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	extentionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

const (
	podCreateMaxRetries int = 5
)

func updateVersionLabel(dep *extentionsv1beta1.DaemonSet, version string) {
	labels := dep.Spec.Template.ObjectMeta.GetLabels()
	labels["version"] = version
	dep.Spec.Template.ObjectMeta.SetLabels(labels)
}

//DeleteQuobyteDeployment deletes Quobyte daemonset with the given name.
func DeleteQuobyteDeployment(name string) error {
	err = KubernetesClient.ExtensionsV1beta1().DaemonSets(quobyteNameSpace).Delete(name, &metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}

func GetDaemonsetByName(name string) (*extentionsv1beta1.DaemonSet, error) {
	return KubernetesClient.ExtensionsV1beta1().DaemonSets(quobyteNameSpace).Get(name, metav1.GetOptions{})
}

// GetClientPodsByVersion returns client pods with specified version.
func GetClientPodsByVersion(version string) (podlist *v1.PodList) {
	podlist, err := KubernetesClient.CoreV1().Pods(quobyteNameSpace).List(metav1.ListOptions{LabelSelector: fmt.Sprintf("version=%s,role=client", version)})
	if err != nil {
		glog.Errorf("Failed to get pods with version %s due to %s", version, err)
		return nil
	}
	return
}

// DeleteClientPods deletes Quobyte client pods with given version.
func DeleteClientPods(version string) {
	list := GetClientPodsByVersion(version)
	if list != nil {
		for _, pod := range list.Items {
			err := KubernetesClient.CoreV1().Pods(quobyteNameSpace).Delete(pod.Name, &metav1.DeleteOptions{})
			if err != nil {
				glog.Errorf("Failed to delete pod %s due to %s", pod.Name, err)
			}
		}
	}
}

// ControlledClientPodUpdate Deletes given list of pods. This triggers redeployment of updated client.
// Currently, this method checks pod recreation on node 5 times with each try dealyed by 60 secs.
// If Pod not running within 5 tries past the Pending pod status, this reports Failure and exits client update process.
func ControlledClientPodUpdate(list *v1.PodList) {
	for _, pod := range list.Items {
		glog.Infof("Deleting client pod %s on node %s \n", pod.Name, pod.Spec.NodeName)

		// drain all the pods on the current node except client pod.
		// err := KubernetesClient.CoreV1().Pods(quobyteNameSpace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "role!=client", FieldSelector: fmt.Sprintf("spec.nodeName=%s", pod.Spec.NodeName)})
		// if err != nil {
		// 	fmt.Println("Failed to drain the node ", pod.Spec.NodeName)
		// 	break // don't proceed with update.
		// }

		err = KubernetesClient.CoreV1().Pods(quobyteNameSpace).Delete(pod.Name, &metav1.DeleteOptions{})
		if err != nil {
			glog.Errorf("Exiting client update. Failed to delete client pod %s due to %s", pod.Name, err)
			break // don't proceed with update.
		}
		podRunning := false
		retryCount := 0
		for !podRunning && retryCount < podCreateMaxRetries {
			podList, err := KubernetesClient.CoreV1().Pods(quobyteNameSpace).List(metav1.ListOptions{LabelSelector: "role=client", FieldSelector: fmt.Sprintf("spec.nodeName=%s", pod.Spec.NodeName)})
			if err != nil {
				glog.Errorf("Failed to list updated pod on node %s", pod.Spec.NodeName)
			}

			podRunning, err = isPodRunning(podList.Items[0])

			if !podRunning {
				if err == nil { // Count retry only if pod is past pod Pending status
					retryCount++
				}

				glog.Info("Updated client Pod is not running on node %s\nretry: %d of configured max retries: %d", pod.Spec.NodeName, retryCount, podCreateMaxRetries)
				time.Sleep(time.Minute)
			}
		}
		if !podRunning {
			glog.Error("Exiting client update: Failed to create updated client pod")
			break
		}
	}
}

func isPodRunning(pod v1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case v1.PodFailed, v1.PodSucceeded:
		return false, nil
	case v1.PodRunning:
		for _, cond := range pod.Status.Conditions {
			if cond.Type != v1.PodReady {
				continue
			}
			return cond.Status == v1.ConditionTrue, nil
		}
		return false, nil
	case v1.PodPending:
		return false, fmt.Errorf("Pending")
	}

	return false, nil
}

// UpdateDaemonSet updates daemonset, with ondelete rolling update strategy and given version.
func UpdateDaemonSet(daemonsetname string, version string) error {

	ds, err := GetDaemonsetByName(daemonsetname)
	if err != nil {
		fmt.Printf("Unable to read client daemonset: %s\n", err)
		return err
	}

	oldData, err := json.Marshal(ds)
	if err != nil {
		return err
	}

	updateVersionLabel(ds, version)
	ds.Spec.UpdateStrategy.Type = "OnDelete" // TODO: move it to the definition in yaml to make it default
	image := ds.Spec.Template.Spec.Containers[0].Image
	ds.Spec.Template.Spec.Containers[0].Image = image[0:strings.LastIndex(image, ":")+1] + version

	newJSON, err := json.Marshal(ds)
	patchbytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newJSON, extentionsv1beta1.DaemonSet{})

	if len(patchbytes) > 2 {
		updatedDS, err := KubernetesClient.ExtensionsV1beta1().DaemonSets(quobyteNameSpace).Patch(ds.Name, types.StrategicMergePatchType, patchbytes)
		if err != nil {
			glog.Errorf("update of client daemonset failed: %s", err)
			return err
		}
		glog.Infof("updated daemonset %s", updatedDS.Name)
	}
	return nil
}
