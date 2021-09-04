package repo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ozonva/ova-checklist-api/internal/types"
)

// repoDB implements Repo
type repoDB struct {
	pool *pgxpool.Pool
}

type queryBuilderConsumer func(*squirrel.StatementBuilderType) (squirrel.Sqlizer, error)

func NewRepoOverDB(pool *pgxpool.Pool) Repo {
	return &repoDB{
		pool: pool,
	}
}

func (r *repoDB) AddChecklists(ctx context.Context, checklists []types.Checklist) error {
	if len(checklists) == 0 {
		return nil
	}

	return r.writeWithPool(ctx, func(builder *squirrel.StatementBuilderType) (squirrel.Sqlizer, error) {
		inserter := builder.Insert("checklists").Columns("user_id", "checklist_id", "data")
		for _, checklist := range checklists {
			serialized, err := checklist.ToJSON()
			if err != nil {
				return nil, err
			}
			inserter = inserter.Values(checklist.UserID, checklist.ID, serialized)
		}
		return inserter, nil
	})
}

func (r *repoDB) ListChecklists(ctx context.Context, userId, limit, offset uint64) ([]types.Checklist, error) {
	var serializedChecklists []string
	err := r.readWithPool(ctx, func(builder *squirrel.StatementBuilderType) (squirrel.Sqlizer, error) {
		selector := builder.
			Select("data").
			From("checklists").
			Where(squirrel.Eq{
				"user_id": userId,
			}).
			OrderBy("created_at").
			Limit(limit).
			Offset(offset)
		return selector, nil
	}, &serializedChecklists)

	if err != nil {
		return nil, err
	}

	return deserializeChecklists(serializedChecklists)
}

func (r *repoDB) DescribeChecklist(ctx context.Context, checklistId string) (*types.Checklist, error) {
	var serializedChecklists []string
	err := r.readWithPool(ctx, func(builder *squirrel.StatementBuilderType) (squirrel.Sqlizer, error) {
		selector := builder.
			Select("data").
			From("checklists").
			Where(squirrel.Eq{
				"checklist_id": checklistId,
			}).
			Limit(1)
		return selector, nil
	}, &serializedChecklists)

	if err != nil {
		return nil, err
	}

	checklists, err := deserializeChecklists(serializedChecklists)
	if err != nil || len(checklists) == 0 {
		return nil, err
	}
	return &checklists[0], nil
}

func (r *repoDB) RemoveChecklist(ctx context.Context, checklistId string) error {
	return r.writeWithPool(ctx, func(builder *squirrel.StatementBuilderType) (squirrel.Sqlizer, error) {
		remover := builder.
			Delete("checklists").
			Where(squirrel.Eq{
				"checklist_id": checklistId,
			})
		return remover, nil
	})
}

func (r *repoDB) writeWithPool(ctx context.Context, consumer queryBuilderConsumer) error {
	conn, query, args, err := r.prepareSqlRequest(ctx, consumer)
	defer closeConnection(conn)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (r *repoDB) readWithPool(ctx context.Context, consumer queryBuilderConsumer, result interface{}) error {
	conn, query, args, err := r.prepareSqlRequest(ctx, consumer)
	defer closeConnection(conn)
	if err != nil {
		return err
	}
	return pgxscan.Select(ctx, conn, result, query, args...)
}

func (r *repoDB) prepareSqlRequest(ctx context.Context, consumer queryBuilderConsumer) (*pgxpool.Conn, string, []interface{}, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	builder, err := consumer(newPgQuery())
	if err != nil {
		return nil, "", nil, err
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, "", nil, err
	}

	return conn, query, args, nil
}

func closeConnection(conn *pgxpool.Conn) {
	if conn != nil {
		conn.Release()
	}
}

func newPgQuery() *squirrel.StatementBuilderType {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &builder
}

func deserializeChecklists(serializedChecklists []string) ([]types.Checklist, error) {
	result := make([]types.Checklist, 0, len(serializedChecklists))
	for _, serialized := range serializedChecklists {
		checklist, err := types.ChecklistFromJSON(serialized)
		if err != nil {
			return nil, err
		}
		result = append(result, checklist)
	}
	return result, nil
}
