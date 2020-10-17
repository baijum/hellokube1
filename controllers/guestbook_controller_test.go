package controllers

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gbv1 "github.com/baijum/hellokube1/api/v1"
)

var _ = Describe("Guestbook Controller:", func() {

	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 10
	)

	Context("When creating Guestbook", func() {

		var testNamespace string
		var ns *corev1.Namespace
		ctx := context.Background()

		BeforeEach(func() {
			//k8sClient = k8sManager.GetClient()
			testNamespace = "testns-" + uuid.Must(uuid.NewRandom()).String()
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: testNamespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, ns)).Should(Succeed())
		})

		It("should update the Guestbook", func() {

			gb := &gbv1.Guestbook{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "webapp.baiju.dev/v1",
					Kind:       "Guestbook",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gb",
					Namespace: testNamespace,
				},
				Spec: gbv1.GuestbookSpec{
					Foo: "hello",
				},
			}
			Expect(k8sClient.Create(ctx, gb)).Should(Succeed())

			guestbookLookupKey := client.ObjectKey{Name: "guestbook", Namespace: testNamespace}
			createdGuestbook := &gbv1.Guestbook{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, guestbookLookupKey, createdGuestbook)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

		})
	})
})
