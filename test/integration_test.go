// +build integration

package test

import (
	"context"
	"github.com/jackc/pgx/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"testing"

	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
	cl "github.com/ozonva/ova-checklist-api/pkg/client"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = Describe("Saver", func() {
	var (
		client, _ = cl.NewClient("localhost", 8080)
		dbConnect, _ = pgx.Connect(context.Background(), "postgres://gopher:gopher@localhost:5432/general")
	)

	AfterSuite(func() {
		cleanUpDatabase(dbConnect)
	})

	BeforeEach(func() {
		cleanUpDatabase(dbConnect)
	})

	Describe("Normal service usage", func() {
		Context("When a user adds a checklist", func() {
			checklist := makeChecklist(1, "First checklist")
			var checklistId string
			It("should be created successfully and return checklist id", func() {
				response, err := client.CreateChecklist(context.Background(), &pb.CreateChecklistRequest{
					Checklist: checklist,
				})
				Expect(err).To(BeNil())
				Expect(len(response.ChecklistId)).To(BeNumerically(">", 0))
				checklistId = response.ChecklistId
			})

			It("should be acquirable from the service", func() {
				response, err := client.DescribeChecklist(context.Background(), &pb.DescribeChecklistRequest{
					ChecklistId: checklistId,
				})
				Expect(err).To(BeNil())
				Expect(proto.Equal(response.Checklist, checklist)).To(BeTrue())
			})
		})

		Context("When there are multiple checklists", func() {
			checklists := []*pb.Checklist{
				makeChecklist(1, "first"),
				makeChecklist(1, "second"),
				makeChecklist(1, "third"),
			}

			pushChecklists := func() {
				for _, checklist := range checklists {
					client.CreateChecklist(context.Background(), &pb.CreateChecklistRequest{
						Checklist: checklist,
					})
				}
				time.Sleep(500 * time.Millisecond) // TODO: fix this place
			}

			It("should be possible to list all of them", func() {
				pushChecklists()
				response, err := client.ListChecklists(context.Background(), &pb.ListChecklistsRequest{
					UserId: 1,
					Limit: 3,
					Offset: 0,
				})
				Expect(err).To(BeNil())
				Expect(len(response.Checklists)).To(Equal(3))
				for i, expectedChecklist := range checklists {
					actualChecklist := response.Checklists[i]
					Expect(proto.Equal(actualChecklist, expectedChecklist)).To(BeTrue())
				}

				response, err = client.ListChecklists(context.Background(), &pb.ListChecklistsRequest{
					UserId: 1,
					Limit: 2,
					Offset: 1,
				})
				Expect(err).To(BeNil())
				Expect(len(response.Checklists)).To(Equal(2))
				Expect(proto.Equal(response.Checklists[0], checklists[1])).To(BeTrue())
				Expect(proto.Equal(response.Checklists[1], checklists[2])).To(BeTrue())
			})
		})
	})
})

func cleanUpDatabase(dbConnect *pgx.Conn) {
	dbConnect.Exec(context.Background(), `
		DELETE FROM checklists;
	`)
}

func makeChecklist(userId uint64, title string) *pb.Checklist {
	return &pb.Checklist{
		UserId: userId,
		Title: title,
		Description: "Default description",
		Items: []*pb.ChecklistItem {
			&pb.ChecklistItem {
				Title: "Item 1",
				IsComplete: false,
			},
		},
	}
}
