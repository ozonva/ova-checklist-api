// +build integration

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"

	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
	cl "github.com/ozonva/ova-checklist-api/pkg/client"
)

const waitUpdateFor = 50 * time.Millisecond // TODO: remove sleeping from tests

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = Describe("Database integration", func() {
	var (
		client, _    = cl.NewClient("localhost", 8080)
		dbConnect, _ = pgx.Connect(context.Background(), "postgres://gopher:gopher@localhost:5432/general")
	)

	BeforeEach(func() {
		cleanUpDatabase(dbConnect)
	})

	AfterEach(func() {
		cleanUpDatabase(dbConnect)
	})

	Describe("When we write a value to the service", func() {
		It("should be readable", func() {
			checklist := makeChecklist(1, "First checklist")
			createResponse, err := client.CreateChecklist(context.Background(), &pb.CreateChecklistRequest{
				Checklist: checklist,
			})
			Expect(err).To(BeNil())
			Expect(len(createResponse.ChecklistId)).To(BeNumerically(">", 0))
			time.Sleep(waitUpdateFor)

			descResponse, err := client.DescribeChecklist(context.Background(), &pb.DescribeChecklistRequest{
				UserId:      1,
				ChecklistId: createResponse.ChecklistId,
			})
			Expect(err).To(BeNil())
			Expect(proto.Equal(descResponse.Checklist, checklist)).To(BeTrue())
		})
	})

	Describe("When there are many checklists", func() {
		const userId = 1

		var (
			checklists []*pb.UserChecklist
		)

		BeforeEach(func() {
			ids := []uint64{1, 2, 3}
			checklists = make([]*pb.UserChecklist, 0, len(ids))
			for _, id := range ids {
				checklist := makeChecklist(userId, fmt.Sprintf("List #%d", id))
				response, _ := client.CreateChecklist(context.Background(), &pb.CreateChecklistRequest{
					Checklist: checklist,
				})
				checklists = append(checklists, &pb.UserChecklist{
					Checklist:   checklist,
					ChecklistId: response.ChecklistId,
				})
			}
			time.Sleep(waitUpdateFor)
		})

		It("should be possible to list all of them", func() {
			response, err := client.ListChecklists(context.Background(), &pb.ListChecklistsRequest{
				UserId: userId,
				Limit:  3,
				Offset: 0,
			})
			Expect(err).To(BeNil())
			Expect(len(response.Checklists)).To(Equal(3))
			for i, expectedChecklist := range checklists {
				actualChecklist := response.Checklists[i]
				Expect(proto.Equal(actualChecklist, expectedChecklist)).To(BeTrue())
			}
		})

		It("should be possible to list some part of them", func() {
			response, err := client.ListChecklists(context.Background(), &pb.ListChecklistsRequest{
				UserId: userId,
				Limit:  2,
				Offset: 1,
			})
			Expect(err).To(BeNil())
			Expect(len(response.Checklists)).To(Equal(2))
			Expect(proto.Equal(response.Checklists[0], checklists[1])).To(BeTrue())
			Expect(proto.Equal(response.Checklists[1], checklists[2])).To(BeTrue())
		})

		It("should be possible to remove any of them", func() {
			_, err := client.RemoveChecklist(context.Background(), &pb.RemoveChecklistRequest{
				UserId:      userId,
				ChecklistId: checklists[1].ChecklistId,
			})
			Expect(err).To(BeNil())

			// Ensure that the checklist was removed
			response, err := client.ListChecklists(context.Background(), &pb.ListChecklistsRequest{
				UserId: userId,
				Limit:  3,
				Offset: 0,
			})
			Expect(err).To(BeNil())
			Expect(len(response.Checklists)).To(Equal(2))
			Expect(proto.Equal(response.Checklists[0], checklists[0])).To(BeTrue())
			Expect(proto.Equal(response.Checklists[1], checklists[2])).To(BeTrue())
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
		UserId:      userId,
		Title:       title,
		Description: "Default description",
		Items: []*pb.ChecklistItem{
			&pb.ChecklistItem{
				Title:      "Item 1",
				IsComplete: false,
			},
		},
	}
}
