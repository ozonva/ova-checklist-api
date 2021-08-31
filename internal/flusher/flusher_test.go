package flusher

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mrepo "github.com/ozonva/ova-checklist-api/internal/generated/repo"
	"github.com/ozonva/ova-checklist-api/internal/types"
)

func TestFlusher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flusher Suite")
}

var _ = Describe("Flusher", func() {
	var (
		ctrl *gomock.Controller
		repo *mrepo.MockRepo
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repo = mrepo.NewMockRepo(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Flushing to a repository", func() {
		Context("When input slice is empty", func() {
			It("should not call repo.AddChecklists at all", func() {
				f := New(1, repo)
				repo.
					EXPECT().
					AddChecklists(gomock.Any()).
					Times(0)
				Expect(f.Flush([]types.Checklist{})).To(Equal([]types.Checklist{}))
			})
		})

		Context("When chunkSize is greater than the length of an input slice", func() {
			It("should call repo.AddChecklists only once", func() {
				repo.
					EXPECT().
					AddChecklists(gomock.Any()).
					DoAndReturn(addChecklistsSuccess).
					Times(1)
				f := New(10, repo)
				input := []types.Checklist{checklist(0), checklist(1)}
				Expect(f.Flush(input)).To(Equal([]types.Checklist{}))
			})
		})

		Context("When chunkSize is equal to the length of an input slice", func() {
			It("should call repo.AddChecklists only once", func() {
				repo.
					EXPECT().
					AddChecklists(gomock.Any()).
					DoAndReturn(addChecklistsSuccess).
					Times(1)
				f := New(3, repo)
				input := []types.Checklist{checklist(0), checklist(1), checklist(2)}
				Expect(f.Flush(input)).To(Equal([]types.Checklist{}))
			})
		})

		Context("When chunkSize is lesser than the length of an input slice and divides it", func() {
			It("should call repo.AddChecklists exactly len(input)/chunkSize times", func() {
				repo.
					EXPECT().
					AddChecklists(gomock.Any()).
					DoAndReturn(func (_ []types.Checklist) error {
						return nil
					}).
					Times(2)
				f := New(2, repo)
				input := []types.Checklist{checklist(0), checklist(1), checklist(2), checklist(3)}
				Expect(f.Flush(input)).To(Equal([]types.Checklist{}))
			})
		})

		Context("When chunkSize is lesser than the length of an input slice and does not divide it", func() {
			It("should call repo.AddChecklists exactly len(input)/chunkSize + 1 times", func() {
				repo.
					EXPECT().
					AddChecklists(gomock.Any()).
					DoAndReturn(addChecklistsSuccess).
					Times(4)
				f := New(3, repo)
				input := []types.Checklist{
					checklist(0), checklist(1), checklist(2), checklist(3), checklist(4),
					checklist(5), checklist(6), checklist(7), checklist(8), checklist(9),
				}
				Expect(f.Flush(input)).To(Equal([]types.Checklist{}))
			})
		})

		Context("When repo.AddChecklists fails", func() {
			It("should collect not pushed checklists", func() {
				repo.
					EXPECT().
					AddChecklists(gomock.Any()).
					DoAndReturn(func (chunk []types.Checklist) error {
						if chunk[0].UserID % 2 == 0 {
							return errors.New("let's fail when a user ID is an even number")
						}
						return nil
					}).
					Times(10)
				f := New(1, repo)
				input := []types.Checklist{
					checklist(0), checklist(1), checklist(2), checklist(3), checklist(4),
					checklist(5), checklist(6), checklist(7), checklist(8), checklist(9),
				}
				notPushed := []types.Checklist{
					checklist(0), checklist(2), checklist(4), checklist(6), checklist(8),
				}
				Expect(f.Flush(input)).To(Equal(notPushed))
			})
		})
	})
})

func addChecklistsSuccess(_ []types.Checklist) error {
	return nil
}

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
