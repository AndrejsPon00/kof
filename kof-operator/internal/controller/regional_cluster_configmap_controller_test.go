/*
Copyright 2025.

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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	kcmv1beta1 "github.com/K0rdent/kcm/api/v1beta1"
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	grafanav1beta1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	kofv1beta1 "github.com/k0rdent/kof/kof-operator/api/v1beta1"
	istio "github.com/k0rdent/kof/kof-operator/internal/controller/istio"
	"github.com/k0rdent/kof/kof-operator/internal/controller/istio/cert"
	remotesecret "github.com/k0rdent/kof/kof-operator/internal/controller/istio/remote-secret"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	sveltosv1beta1 "github.com/projectsveltos/addon-controller/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("RegionalConfigMap Controller", func() {
	Context("When reconciling a resource", func() {
		ctx := context.Background()
		var clusterDeploymentReconciler *ClusterDeploymentReconciler
		var regionalClusterConfigmapReconciler *RegionalClusterConfigMapReconciler

		const regionalClusterDeploymentName = "test-regional-cm"

		regionalClusterConfigmapNamespacedName := types.NamespacedName{
			Name:      GetRegionalClusterConfigMapName(regionalClusterDeploymentName),
			Namespace: defaultNamespace,
		}

		// child ClusterDeployment

		const childClusterDeploymentName = "test-child-cm"

		childClusterDeploymentNamespacedName := types.NamespacedName{
			Name:      childClusterDeploymentName,
			Namespace: defaultNamespace,
		}

		childClusterDeploymentLabels := map[string]string{
			IstioRoleLabel:              "child",
			KofClusterRoleLabel:         "child",
			KofRegionalClusterNameLabel: regionalClusterDeploymentName,
		}

		childClusterDeploymentAnnotations := map[string]string{}

		const childClusterDeploymentConfig = `{"region": "us-east-2"}`

		// child cluster ConfigMap

		childClusterConfigMapNamespacedName := types.NamespacedName{
			Name:      "kof-cluster-config-test-child-cm", // prefix + child cluster name
			Namespace: defaultNamespace,
		}

		// istio child

		const clusterCertificateName = "kof-istio-test-child-cm-ca"

		clusterCertificateNamespacedName := types.NamespacedName{
			Name:      clusterCertificateName,
			Namespace: istio.IstioSystemNamespace,
		}

		const secretName = "test-child-cm-kubeconfig"

		kubeconfigSecretNamespacedName := types.NamespacedName{
			Name:      secretName,
			Namespace: defaultNamespace,
		}

		remoteSecretNamespacedName := types.NamespacedName{
			Name:      remotesecret.GetRemoteSecretName(childClusterDeploymentName),
			Namespace: istio.IstioSystemNamespace,
		}

		profileDeploymentName := types.NamespacedName{
			Name:      remotesecret.CopyRemoteSecretProfileName(childClusterDeploymentName),
			Namespace: defaultNamespace,
		}

		// create regional cluster configmap

		createRegionalClusterConfigMap := func(
			name,
			namespace,
			istioRole,
			readMetricsEndpoint,
			writeMetricsEndpoint,
			readLogsEndpoint,
			writeLogsEndpoint,
			writeTracesEndpoint,
			cloud,
			awsRegion,
			azureLocation,
			openstackRegion,
			vshereDatacenter string,
		) {
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      GetRegionalClusterConfigMapName(name),
					Namespace: namespace,
					Labels: map[string]string{
						"k0rdent.mirantis.com/kof-regional-cluster": name,
					},
				},
				Data: map[string]string{
					RegionalClusterNameKey:      name,
					RegionalClusterNamespaceKey: defaultNamespace,
					RegionalIstioRoleKey:        istioRole,
					RegionalClusterCloudKey:     cloud,
					ReadMetricsKey:              readMetricsEndpoint,
					WriteMetricsKey:             writeMetricsEndpoint,
					ReadLogsKey:                 readLogsEndpoint,
					WriteLogsKey:                writeLogsEndpoint,
					WriteTracesKey:              writeTracesEndpoint,
					AwsRegionKey:                awsRegion,
					AzureLocationKey:            azureLocation,
					OpenstackRegionKey:          openstackRegion,
					VSphereDatacenterKey:        vshereDatacenter,
				},
			}
			Expect(k8sClient.Create(ctx, configMap)).To(Succeed())
		}
		// createClusterDeployment

		createClusterDeployment := func(
			name string,
			namespace string,
			labels map[string]string,
			annotations map[string]string,
			config string,
		) {
			clusterDeployment := &kcmv1beta1.ClusterDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        name,
					Namespace:   namespace,
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: kcmv1beta1.ClusterDeploymentSpec{
					Template: "aws-cluster-template",
					Config:   &apiextensionsv1.JSON{Raw: []byte(config)},
				},
			}
			Expect(k8sClient.Create(ctx, clusterDeployment)).To(Succeed())

			clusterDeployment.Status = kcmv1beta1.ClusterDeploymentStatus{
				Conditions: []metav1.Condition{
					{
						Type:               kcmv1beta1.ReadyCondition,
						Status:             metav1.ConditionTrue,
						LastTransitionTime: metav1.Time{Time: time.Now()},
						Reason:             "ClusterReady",
						Message:            "Cluster is ready",
					},
					{
						Type:               kcmv1beta1.CAPIClusterSummaryCondition,
						Status:             metav1.ConditionTrue,
						LastTransitionTime: metav1.Time{Time: time.Now()},
						Reason:             "InfrastructureReady",
						Message:            "Infrastructure is ready",
					},
				},
			}
			Expect(k8sClient.Status().Update(ctx, clusterDeployment)).To(Succeed())
		}

		// createSecret

		createSecret := func(secretName string) {
			kubeconfigSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: defaultNamespace,
					Labels:    map[string]string{},
				},
				Data: map[string][]byte{"value": []byte("")},
			}
			Expect(k8sClient.Create(ctx, kubeconfigSecret)).To(Succeed())
		}

		// before each test case

		BeforeEach(func() {

			clusterDeploymentReconciler = &ClusterDeploymentReconciler{
				Client:              k8sClient,
				Scheme:              k8sClient.Scheme(),
				RemoteSecretManager: remotesecret.NewFakeManager(k8sClient),
				IstioCertManager:    cert.New(k8sClient),
			}

			regionalClusterConfigmapReconciler = &RegionalClusterConfigMapReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			By(fmt.Sprintf("creating the %s namespace", istio.IstioSystemNamespace))
			certNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: istio.IstioSystemNamespace,
				},
			}
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      istio.IstioSystemNamespace,
				Namespace: istio.IstioSystemNamespace,
			}, certNamespace)
			if err != nil && errors.IsNotFound(err) {
				Expect(k8sClient.Create(ctx, certNamespace)).To(Succeed())
			}

			By("creating regional cluster configMap")
			createRegionalClusterConfigMap(
				regionalClusterDeploymentName,
				defaultNamespace,
				"",
				"https://vmauth.test-aws-ue2.kof.example.com/vm/select/0/prometheus",
				"https://vmauth.test-aws-ue2.kof.example.com/vm/insert/0/prometheus/api/v1/write",
				"https://vmauth.test-aws-ue2.kof.example.com/vls/select/opentelemetry/v1/logs",
				"https://vmauth.test-aws-ue2.kof.example.com/vli/insert/opentelemetry/v1/logs",
				"https://jaeger.test-aws-ue2.kof.example.com/collector",
				"aws",
				"us-east-2",
				"", "", "",
			)

			By("creating child ClusterDeployment")
			createClusterDeployment(
				childClusterDeploymentName,
				defaultNamespace,
				childClusterDeploymentLabels,
				childClusterDeploymentAnnotations,
				childClusterDeploymentConfig,
			)

			By("creating the fake Secret for the cluster deployment kubeconfig")
			createSecret(secretName)
		})

		// after each test case

		AfterEach(func() {
			cd := &kcmv1beta1.ClusterDeployment{}

			if err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, cd); err == nil {
				By("Cleanup child ClusterDeployment")
				Expect(k8sClient.Delete(ctx, cd)).To(Succeed())
			}

			configMap := &corev1.ConfigMap{}
			if err := k8sClient.Get(ctx, childClusterConfigMapNamespacedName, configMap); err == nil {
				By("Cleanup child cluster ConfigMap")
				Expect(k8sClient.Delete(ctx, configMap)).To(Succeed())
			}

			kubeconfigSecret := &corev1.Secret{}
			if err := k8sClient.Get(ctx, kubeconfigSecretNamespacedName, kubeconfigSecret); err == nil {
				By("Cleanup the Kubeconfig Secret")
				Expect(k8sClient.Delete(ctx, kubeconfigSecret)).To(Succeed())
			}

			remoteSecret := &corev1.Secret{}
			if err := k8sClient.Get(ctx, remoteSecretNamespacedName, remoteSecret); err == nil {
				By("Cleanup the Remote Secret")
				Expect(k8sClient.Delete(ctx, remoteSecret)).To(Succeed())
			}

			cert := &cmv1.Certificate{}
			if err := k8sClient.Get(ctx, clusterCertificateNamespacedName, cert); err == nil {
				By("Cleanup the Certificate")
				Expect(k8sClient.Delete(ctx, cert)).To(Succeed())
			}

			regionalConfigMap := &corev1.ConfigMap{}
			if err := k8sClient.Get(ctx, regionalClusterConfigmapNamespacedName, regionalConfigMap); err == nil {
				By("Cleanup regional Configmap")
				Expect(k8sClient.Delete(ctx, regionalConfigMap)).To(Succeed())
			}
		})

		// test cases

		It("should successfully reconcile the CA resource", func() {
			_, err := clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			cert := &cmv1.Certificate{}
			err = k8sClient.Get(ctx, clusterCertificateNamespacedName, cert)
			Expect(err).NotTo(HaveOccurred())
			Expect(cert.Spec.CommonName).To(Equal(fmt.Sprintf("%s CA", childClusterDeploymentName)))
		})

		It("should successfully reconcile the resource when deleted", func() {
			By("Reconciling the deleted resource")
			clusterDeployment := &kcmv1beta1.ClusterDeployment{}
			err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Delete(ctx, clusterDeployment)).To(Succeed())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			secret := &corev1.Secret{}
			err = k8sClient.Get(ctx, remoteSecretNamespacedName, secret)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should successfully reconcile the resource when not ready", func() {
			By("Reconciling the not ready resource")
			clusterDeployment := &kcmv1beta1.ClusterDeployment{}
			err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())

			for i := range clusterDeployment.Status.Conditions {
				if clusterDeployment.Status.Conditions[i].Type == kcmv1beta1.ReadyCondition {
					clusterDeployment.Status.Conditions[i].Status = metav1.ConditionFalse
					break
				}
			}

			err = k8sClient.Status().Update(ctx, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())

			_, err = regionalClusterConfigmapReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: regionalClusterConfigmapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			secret := &corev1.Secret{}
			err = k8sClient.Get(ctx, remoteSecretNamespacedName, secret)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should successfully reconcile the resource if special label not provided", func() {
			By("Reconciling the resource without labels")
			clusterDeployment := &kcmv1beta1.ClusterDeployment{}
			err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())

			clusterDeployment.ObjectMeta.Labels = nil

			err = k8sClient.Update(ctx, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			secret := &corev1.Secret{}
			err = k8sClient.Get(ctx, remoteSecretNamespacedName, secret)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should successfully reconcile when remote secret already exists", func() {
			By("Reconciling the resource with existed remote secret")
			clusterDeployment := &kcmv1beta1.ClusterDeployment{}
			err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())

			_, err = regionalClusterConfigmapReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: regionalClusterConfigmapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			secret := &corev1.Secret{}
			err = k8sClient.Get(ctx, remoteSecretNamespacedName, secret)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should successfully reconcile after creating and deleting resource", func() {
			By("Verifying resource reconciliation after creation and deletion")
			cd := &kcmv1beta1.ClusterDeployment{}
			err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, cd)
			Expect(err).NotTo(HaveOccurred())
			cdUID := cd.GetUID()

			_, err = regionalClusterConfigmapReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: regionalClusterConfigmapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(k8sClient.Delete(ctx, cd)).To(Succeed())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			secret := &corev1.Secret{}
			err = k8sClient.Get(ctx, remoteSecretNamespacedName, secret)
			Expect(errors.IsNotFound(err)).To(BeTrue())

			cert := &cmv1.Certificate{}
			err = k8sClient.Get(ctx, clusterCertificateNamespacedName, cert)
			Expect(errors.IsNotFound(err)).To(BeTrue())

			configMap := &corev1.ConfigMap{}
			err = k8sClient.Get(ctx, childClusterConfigMapNamespacedName, configMap)
			Expect(err).NotTo(HaveOccurred())
			// There is no garbage collector in the `envtest`,
			// so we should test that `OwnerReference` is set correctly,
			// and assume that Kubernetes garbage collection works:
			// https://github.com/kubernetes-sigs/controller-runtime/issues/626#issuecomment-538529534
			owner := configMap.OwnerReferences[0]
			Expect(owner.APIVersion).To(Equal("k0rdent.mirantis.com/v1beta1"))
			Expect(owner.Kind).To(Equal("ClusterDeployment"))
			Expect(owner.Name).To(Equal(childClusterDeploymentName))
			Expect(owner.UID).To(Equal(cdUID))
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			clusterDeployment := &kcmv1beta1.ClusterDeployment{}
			err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())

			_, err = regionalClusterConfigmapReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: regionalClusterConfigmapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			remoteSecret := &corev1.Secret{}
			err = k8sClient.Get(ctx, remoteSecretNamespacedName, remoteSecret)
			Expect(err).NotTo(HaveOccurred())
		})

		DescribeTable("should create PromxyServerGroup and GrafanaDatasource for regional cluster", func(
			istioRole,
			cloud,
			awsRegion,
			azureLocation,
			openstackRegion,
			vshereDatacenter,
			expectedMetricsScheme,
			expectedMetricsTarget,
			expectedMetricsPathPrefix string,
			expectedMetricsHttpConfig kofv1beta1.HTTPClientConfig,
			expectedGrafanaDatasourceURL string,
			expectedGrafanaDatasourceJsonData string,
		) {
			By("creating regional ClusterDeployment with labels and config from the table")
			const tableRegionalClusterDeploymentName = "test-regional-from-table-cm"

			regionalClusterConfigmapNamespacedName := types.NamespacedName{
				Name:      GetRegionalClusterConfigMapName(tableRegionalClusterDeploymentName),
				Namespace: defaultNamespace,
			}

			promxyServerGroupNamespacedName := types.NamespacedName{
				Name:      tableRegionalClusterDeploymentName + "-metrics",
				Namespace: ReleaseNamespace,
			}

			grafanaDatasourceNamespacedName := types.NamespacedName{
				Name:      tableRegionalClusterDeploymentName + "-logs",
				Namespace: ReleaseNamespace,
			}

			secretName := tableRegionalClusterDeploymentName + "-kubeconfig"
			createSecret(secretName)

			createRegionalClusterConfigMap(
				tableRegionalClusterDeploymentName,
				defaultNamespace,
				istioRole,
				"https://vmauth.test-aws-ue2.kof.example.com/vm/select/0/prometheus",
				"https://vmauth.test-aws-ue2.kof.example.com/vm/insert/0/prometheus/api/v1/write",
				expectedGrafanaDatasourceURL,
				"https://vmauth.test-aws-ue2.kof.example.com/vli/insert/opentelemetry/v1/logs",
				"https://jaeger.test-aws-ue2.kof.example.com/collector",
				cloud,
				awsRegion,
				azureLocation,
				openstackRegion,
				vshereDatacenter,
			)

			DeferCleanup(func() {
				regionalClusterConfigMap := &corev1.ConfigMap{}
				if err := k8sClient.Get(ctx, regionalClusterConfigmapNamespacedName, regionalClusterConfigMap); err == nil {
					By("cleanup regional cluster configMap")
					Expect(k8sClient.Delete(ctx, regionalClusterConfigMap)).To(Succeed())
				}

				kubeconfigSecretNamespacedName := types.NamespacedName{
					Name:      secretName,
					Namespace: defaultNamespace,
				}
				kubeconfigSecret := &corev1.Secret{}
				if err := k8sClient.Get(ctx, kubeconfigSecretNamespacedName, kubeconfigSecret); err == nil {
					By("cleanup kubeconfig Secret")
					Expect(k8sClient.Delete(ctx, kubeconfigSecret)).To(Succeed())
				}

				promxyServerGroup := &kofv1beta1.PromxyServerGroup{}
				if err := k8sClient.Get(ctx, promxyServerGroupNamespacedName, promxyServerGroup); err == nil {
					By("cleanup PromxyServerGroup")
					Expect(k8sClient.Delete(ctx, promxyServerGroup)).To(Succeed())
				}

				grafanaDatasource := &grafanav1beta1.GrafanaDatasource{}
				if err := k8sClient.Get(ctx, grafanaDatasourceNamespacedName, grafanaDatasource); err == nil {
					By("cleanup GrafanaDatasource")
					Expect(k8sClient.Delete(ctx, grafanaDatasource)).To(Succeed())
				}
			})

			By("reconciling regional cluster ConfigMap")
			_, err := regionalClusterConfigmapReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: regionalClusterConfigmapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("reading PromxyServerGroup")
			promxyServerGroup := &kofv1beta1.PromxyServerGroup{}
			err = k8sClient.Get(ctx, promxyServerGroupNamespacedName, promxyServerGroup)
			Expect(err).NotTo(HaveOccurred())
			Expect(promxyServerGroup.Spec.Scheme).To(Equal(expectedMetricsScheme))
			Expect(promxyServerGroup.Spec.Targets).To(Equal([]string{expectedMetricsTarget}))
			Expect(promxyServerGroup.Spec.PathPrefix).To(Equal(expectedMetricsPathPrefix))
			Expect(promxyServerGroup.Spec.HttpClient).To(Equal(expectedMetricsHttpConfig))

			By("reading GrafanaDatasource")
			grafanaDatasource := &grafanav1beta1.GrafanaDatasource{}
			err = k8sClient.Get(ctx, grafanaDatasourceNamespacedName, grafanaDatasource)
			Expect(err).NotTo(HaveOccurred())
			Expect(grafanaDatasource.Spec.Datasource.URL).To(Equal(expectedGrafanaDatasourceURL))
			if expectedGrafanaDatasourceJsonData != "" {
				Expect(grafanaDatasource.Spec.Datasource.JSONData).To(MatchJSON(json.RawMessage(expectedGrafanaDatasourceJsonData)))
			}
		},

			/*
				Entry(
					description,
					istioRole,
					cloud,
					awsRegion,
					azureLocation,
					openstackRegion,
					vshereDatacenter,
					expectedMetricsScheme,
					expectedMetricsTarget,
					expectedMetricsPathPrefix,
					expectedMetricsBasicAuth,
					expectedGrafanaDatasourceURL,
				),
			*/

			Entry(
				"Default endpoints",
				"",
				"aws",
				"us-east-2",
				"", "", "",
				"https",
				"vmauth.test-aws-ue2.kof.example.com:443",
				"/vm/select/0/prometheus",
				kofv1beta1.HTTPClientConfig{
					DialTimeout: defaultDialTimeout,
					TLSConfig: kofv1beta1.TLSConfig{
						InsecureSkipVerify: false,
					},
					BasicAuth: kofv1beta1.BasicAuth{
						CredentialsSecretName: "storage-vmuser-credentials",
						UsernameKey:           "username",
						PasswordKey:           "password"},
				},
				"https://vmauth.test-aws-ue2.kof.example.com/vls", "",
			),
			/*
				Entry(
					"Istio endpoints",
					"child",
					"aws",
					"us-east-2",
					"", "", "",
					"https",
					"vmauth.test-aws-ue2.kof.example.com:443",
					"/vm/select/0/prometheus",
					kofv1beta1.HTTPClientConfig{
						DialTimeout: defaultDialTimeout,
						TLSConfig: kofv1beta1.TLSConfig{
							InsecureSkipVerify: false,
						},
						BasicAuth: kofv1beta1.BasicAuth{},
					},
					"http://test-regional-from-table-logs-select:9471", "",
				),
			*/
		)

		It("should create ConfigMap for child cluster", func() {
			By("reconciling regional cluster ConfigMap")
			_, err := regionalClusterConfigmapReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: regionalClusterConfigmapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("reconciling child ClusterDeployment")
			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("reading created ConfigMap")
			configMap := &corev1.ConfigMap{}
			err = k8sClient.Get(ctx, childClusterConfigMapNamespacedName, configMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(configMap.Data[RegionalClusterNameKey]).To(Equal(regionalClusterDeploymentName))
			Expect(configMap.Data[ReadMetricsKey]).To(Equal(
				"https://vmauth.test-aws-ue2.kof.example.com/vm/select/0/prometheus",
			))
			Expect(configMap.Data[WriteMetricsKey]).To(Equal(
				"https://vmauth.test-aws-ue2.kof.example.com/vm/insert/0/prometheus/api/v1/write",
			))
			Expect(configMap.Data[WriteLogsKey]).To(Equal(
				"https://vmauth.test-aws-ue2.kof.example.com/vli/insert/opentelemetry/v1/logs",
			))
			Expect(configMap.Data[WriteTracesKey]).To(Equal(
				"https://jaeger.test-aws-ue2.kof.example.com/collector",
			))
		})

		DescribeTable("should discover regional cluster by AWS region or label", func(
			withLabel bool,
			crossNamespace bool,
			childNamespace string,
			shouldSucceed bool,
		) {

			By("setting CROSS_NAMESPACE environment variable")
			err := os.Setenv("CROSS_NAMESPACE", strconv.FormatBool(crossNamespace))
			Expect(err).NotTo(HaveOccurred())

			const childClusterDeploymentName = "test-child-aws"
			const regionalClusterName = "test-regional-aws"

			childClusterDeploymentNamespacedName := types.NamespacedName{
				Name:      childClusterDeploymentName,
				Namespace: childNamespace,
			}

			childClusterConfigMapNamespacedName := types.NamespacedName{
				Name:      "kof-cluster-config-" + childClusterDeploymentName,
				Namespace: childNamespace,
			}

			regionalClusterConfigmapNamespacedName := types.NamespacedName{
				Name:      GetRegionalClusterConfigMapName(regionalClusterName),
				Namespace: defaultNamespace,
			}

			DeferCleanup(func() {
				err := os.Unsetenv("CROSS_NAMESPACE")
				Expect(err).NotTo(HaveOccurred())

				childClusterDeployment := &kcmv1beta1.ClusterDeployment{}
				if err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, childClusterDeployment); err == nil {
					By("cleanup child ClusterDeployment")
					Expect(k8sClient.Delete(ctx, childClusterDeployment)).To(Succeed())
				}

				configMap := &corev1.ConfigMap{}
				if err := k8sClient.Get(ctx, childClusterConfigMapNamespacedName, configMap); err == nil {
					By("cleanup child cluster ConfigMap")
					Expect(k8sClient.Delete(ctx, configMap)).To(Succeed())
				}

				regionalClusterConfigMap := &corev1.ConfigMap{}
				if err := k8sClient.Get(ctx, regionalClusterConfigmapNamespacedName, regionalClusterConfigMap); err == nil {
					By("cleanup regional cluster configMap")
					Expect(k8sClient.Delete(ctx, regionalClusterConfigMap)).To(Succeed())
				}
			})

			By("creating child ClusterDeployment without kof-regional-cluster-name label")

			childClusterDeploymentLabels := map[string]string{
				KofClusterRoleLabel: "child",
				// Note no `KofRegionalClusterNameLabel` here, it will be auto-discovered!
			}
			if withLabel {
				childClusterDeploymentLabels[KofRegionalClusterNameLabel] = regionalClusterName
				if crossNamespace {
					childClusterDeploymentLabels[KofRegionalClusterNamespaceLabel] = "default"
				}
			}

			childClusterDeploymentAnnotations := map[string]string{}

			createRegionalClusterConfigMap(
				regionalClusterName,
				regionalClusterConfigmapNamespacedName.Namespace,
				"", "", "",
				"", "", "",
				"aws",
				"us-east-2",
				"", "", "",
			)

			createClusterDeployment(
				childClusterDeploymentName,
				childNamespace,
				childClusterDeploymentLabels,
				childClusterDeploymentAnnotations,
				childClusterDeploymentConfig,
			)

			By("reconciling child ClusterDeployment")
			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			if shouldSucceed {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			}

			By("reading created ConfigMap")
			configMap := &corev1.ConfigMap{}
			err = k8sClient.Get(ctx, childClusterConfigMapNamespacedName, configMap)
			if shouldSucceed {
				Expect(err).NotTo(HaveOccurred())
				Expect(configMap.Data[RegionalClusterNameKey]).To(Equal(regionalClusterName))
			} else {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			}
		},
			EntryDescription("withLabel=%t, crossNamespace=%t, childNamespace=%s, ok=%t"),
			Entry(nil, false, false, defaultNamespace, true),
			Entry(nil, false, false, ReleaseNamespace, false),
			Entry(nil, false, true, defaultNamespace, true),
			Entry(nil, false, true, ReleaseNamespace, true),
			Entry(nil, true, false, defaultNamespace, true),
			Entry(nil, true, false, ReleaseNamespace, false),
			Entry(nil, true, true, defaultNamespace, true),
			Entry(nil, true, true, ReleaseNamespace, true),
		)

		It("should create profile", func() {
			By("reading child ClusterDeployment")
			clusterDeployment := &kcmv1beta1.ClusterDeployment{}
			err := k8sClient.Get(ctx, childClusterDeploymentNamespacedName, clusterDeployment)
			Expect(err).NotTo(HaveOccurred())

			By("reconciling regional cluster ConfigMap")
			_, err = regionalClusterConfigmapReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: regionalClusterConfigmapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("reconciling child ClusterDeployment")
			_, err = clusterDeploymentReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: childClusterDeploymentNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("reading profile")
			profile := &sveltosv1beta1.Profile{}
			err = k8sClient.Get(ctx, profileDeploymentName, profile)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
