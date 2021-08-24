package saver

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync/atomic"

	mflusher "ova-checklist-api/internal/generated/flusher"
	"ova-checklist-api/internal/types"

	"sync"
	"testing"
	"time"
)

func TestSaver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Saver Suite")
}

var _ = Describe("Saver", func() {
	var (
		ctrl    *gomock.Controller
		flusher *mflusher.MockFlusher
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		flusher = mflusher.NewMockFlusher(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Checklist saver", func() {
		Context("When the internal buffer is out of space", func() {
			It("should flush data without waiting for a timer", func() {
				const bufferSize = 10
				var repo []types.Checklist
				var expectedRepo []types.Checklist
				var wg sync.WaitGroup
				wg.Add(1)
				flusher.
					EXPECT().
					Flush(gomock.Any()).
					Times(1).
					DoAndReturn(func(values []types.Checklist) []types.Checklist {
						defer wg.Done()
						repo = append(repo, values...)
						return nil
					})
				s := NewSaver(flusher, bufferSize, 100500*time.Hour)
				for i := 0; i < bufferSize; i++ {
					value := checklist(uint64(i))
					expectedRepo = append(expectedRepo, value)
					Expect(s.TrySave(value)).To(Equal(true))
				}
				wg.Wait() // Ensure that the flush happens because of the internal buffer overflow
				s.Close()

				Expect(repo).To(Equal(expectedRepo))
			})
		})

		Context("When there are pending values but a saver is being closed", func() {
			It("should flush all pending values", func() {
				const bufferSize = 10
				var repo []types.Checklist
				var expectedRepo []types.Checklist
				flusher.
					EXPECT().
					Flush(gomock.Any()).
					Times(1).
					DoAndReturn(func(values []types.Checklist) []types.Checklist {
						repo = append(repo, values...)
						return nil
					})
				s := NewSaver(flusher, bufferSize, 100500*time.Hour)
				for i := 0; i < bufferSize/2; i++ {
					value := checklist(uint64(i))
					expectedRepo = append(expectedRepo, value)
					Expect(s.TrySave(value)).To(Equal(true))
				}
				s.Close()

				Expect(repo).To(Equal(expectedRepo))
			})
		})

		Context("When a timer ticks", func() {
			It("should flush all pending values", func() {
				const bufferSize = 5000
				var repo []types.Checklist
				var wg sync.WaitGroup
				wg.Add(1)
				flusher.
					EXPECT().
					Flush(gomock.Any()).
					Times(1).
					DoAndReturn(func(values []types.Checklist) []types.Checklist {
						defer wg.Done()
						repo = append(repo, values...)
						return nil
					})
				s := NewSaver(flusher, bufferSize, 50*time.Millisecond)
				Expect(s.TrySave(checklist(0))).To(Equal(true))
				wg.Wait() // Ensure that the flush happens because of a timer tick
				s.Close()

				Expect(repo).To(Equal([]types.Checklist{
					checklist(0),
				}))
			})
		})

		Context("When a flusher fails", func() {
			It("should keep not flushed values and try to flush them in future", func() {
				valuesToSend := []types.Checklist{checklist(0), checklist(1), checklist(2), checklist(3)}
				firstFlushed := []types.Checklist{checklist(1), checklist(3)}
				expectedRepo := []types.Checklist{checklist(1), checklist(3), checklist(0), checklist(2)}

				const bufferSize = 20
				var repo []types.Checklist
				var wg sync.WaitGroup
				var phase int32
				flusher.
					EXPECT().
					Flush(gomock.Any()).
					AnyTimes().
					DoAndReturn(func(values []types.Checklist) []types.Checklist {
						var failed []types.Checklist
						var flushed []types.Checklist

						oldRepoSize := len(repo)
						currentPhase := atomic.LoadInt32(&phase)
						if currentPhase == 0 {
							for _, checklist := range values {
								if checklist.UserID%2 == 0 {
									failed = append(failed, checklist)
								} else {
									flushed = append(flushed, checklist)
								}
							}
						} else {
							flushed = values
						}

						repo = append(repo, flushed...)
						if oldRepoSize != len(repo) {
							wg.Done()
						}
						return failed
					})

				s := NewSaver(flusher, bufferSize, 50*time.Millisecond)

				// First phase: trying to save all values. Only 1 and 3 will be saved
				wg.Add(1)
				for _, value := range valuesToSend {
					Expect(s.TrySave(value)).To(Equal(true))
				}
				wg.Wait()
				Expect(repo).To(Equal(firstFlushed))

				// Second phase: wait until 0 and 2 will be saved
				wg.Add(1)
				atomic.StoreInt32(&phase, 1)
				wg.Wait()
				Expect(repo).To(Equal(expectedRepo))

				s.Close()
			})
		})
	})
})

func checklist(userId uint64) types.Checklist {
	return types.Checklist{
		UserID:      userId,
		Title:       "Default checklist",
		Description: "Testing checklist utils",
		Items: []types.ChecklistItem{
			{"Step 1", false},
		},
	}
}
