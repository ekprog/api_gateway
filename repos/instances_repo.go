package repos

import (
	"context"
	"database/sql"
	"fmt"
	"microservice/app/core"
	"microservice/domain"
	"microservice/tools"
)

type InstancesRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewInstancesRepo(log core.Logger, db *sql.DB) *InstancesRepo {
	return &InstancesRepo{
		db:  db,
		log: log,
	}
}

func (r *InstancesRepo) All(ctx context.Context) ([]*domain.Instance, error) {
	var items []*domain.Instance

	query := `SELECT id, 
       			folder, 
       			endpoint, 
       			is_active
			FROM services 
			WHERE deleted_at is null
			ORDER BY created_at;`
	raws, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for raws.Next() {
		item := &domain.Instance{}
		err = raws.Scan(&item.Id,
			&item.Folder,
			&item.Endpoint,
			&item.IsActive)
		if err != nil {
			return nil, err
		}
		item.Name = item.Folder
		items = append(items, item)
	}
	return items, nil
}

func (r *InstancesRepo) GetByFolder(ctx context.Context, folder string) (*domain.Instance, error) {
	item := &domain.Instance{}

	query := `SELECT id, 
       			folder, 
       			endpoint, 
       			is_active
			FROM services 
			WHERE deleted_at is null and folder=$1
			ORDER BY created_at;`
	err := r.db.QueryRowContext(ctx, query, folder).Scan(&item.Id,
		&item.Folder,
		&item.Endpoint,
		&item.IsActive)
	switch err {
	case nil:
		return item, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *InstancesRepo) Insert(ctx context.Context, item *domain.Instance) error {
	var id int32
	query := "INSERT INTO services (folder, endpoint) VALUES ($1, $2) returning id"

	err := r.db.QueryRowContext(ctx, query, item.Folder, item.Endpoint).Scan(&id)
	if err != nil {
		return err
	}
	item.Id = id
	return nil
}

func (r *InstancesRepo) Delete(ctx context.Context, id int32) error {
	query := "UPDATE services SET deleted_at=now() WHERE id=$1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *InstancesRepo) Update(ctx context.Context, req *tools.UpdateReq) error {
	k, v := req.BuildFor("folder", "endpoint", "is_active")
	query := fmt.Sprintf("UPDATE services SET %s, updated_at=now() WHERE id=$1", k)
	_, err := r.db.ExecContext(ctx, query, v...)
	if err != nil {
		return err
	}
	return nil
}
