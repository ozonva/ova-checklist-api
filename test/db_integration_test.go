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

	cl "github.com/ozonva/ova-checklist-api/internal/client"
	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
)

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

			descResponse, err := client.DescribeChecklist(context.Background(), &pb.DescribeChecklistRequest{
				UserId:      1,
				ChecklistId: createResponse.ChecklistId,
			})
			Expect(err).To(BeNil())
			Expect(proto.Equal(descResponse.Checklist, checklist)).To(BeTrue())
		})
	})

	Describe("When we write a batch of checklists to the service", func() {
		It("should be able to return all of them", func() {
			checklists := []*pb.Checklist{
				makeChecklist(1, "First checklist"),
				makeChecklist(1, "Second checklist"),
			}
			createResponse, err := client.MultiCreateChecklist(context.Background(), &pb.MultiCreateChecklistRequest{
				Checklists: checklists,
			})
			Expect(err).To(BeNil())
			Expect(createResponse.TotalSaved).To(Equal(uint32(2)))
			waitForDatabaseUpdate()

			listResponse, err := client.ListChecklists(context.Background(), &pb.ListChecklistsRequest{
				UserId: 1,
				Limit:  10,
				Offset: 0,
			})
			Expect(err).To(BeNil())
			Expect(len(listResponse.Checklists)).To(Equal(2))
			Expect(listResponse.Checklists[0].Checklist.Title).To(Equal("First checklist"))
			Expect(listResponse.Checklists[1].Checklist.Title).To(Equal("Second checklist"))
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

	Describe("When there is a checklist saved in the service storage", func() {
		It("should be possible to modify it", func() {
			checklist := makeChecklist(1, "First checklist")
			createResponse, _ := client.CreateChecklist(context.Background(), &pb.CreateChecklistRequest{
				Checklist: checklist,
			})
			descResponse, err := client.DescribeChecklist(context.Background(), &pb.DescribeChecklistRequest{
				UserId:      1,
				ChecklistId: createResponse.ChecklistId,
			})
			Expect(len(descResponse.Checklist.Items)).To(Equal(1))
			Expect(descResponse.Checklist.Items[0].IsComplete).To(BeFalse())

			updated := descResponse.Checklist
			updated.Items[0].IsComplete = true
			_, err = client.UpdateChecklist(context.Background(), &pb.UpdateChecklistRequest{
				ChecklistId: createResponse.ChecklistId,
				Checklist:   updated,
			})
			Expect(err).To(BeNil())

			descResponse, err = client.DescribeChecklist(context.Background(), &pb.DescribeChecklistRequest{
				UserId:      1,
				ChecklistId: createResponse.ChecklistId,
			})
			Expect(len(descResponse.Checklist.Items)).To(Equal(1))
			Expect(descResponse.Checklist.Items[0].IsComplete).To(BeTrue())
		})
	})
})

func cleanUpDatabase(dbConnect *pgx.Conn) {
	dbConnect.Exec(context.Background(), `
		DELETE FROM checklists;
	`)
}

func waitForDatabaseUpdate() {
	const waitTime = 50 * time.Millisecond
	time.Sleep(waitTime) // TODO: remove sleeping from tests
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
