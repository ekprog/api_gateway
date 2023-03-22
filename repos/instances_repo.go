package repos

import (
	"database/sql"
	"microservice/app/core"
	"microservice/domain"
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

func (r *InstancesRepo) All() ([]*domain.Instance, error) {
	var items []*domain.Instance

	query := `SELECT id, 
       			folder, 
       			endpoint, 
       			is_active
			FROM services 
			WHERE deleted_at is null
			ORDER BY created_at;`
	raws, err := r.db.Query(query)
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

func (r *InstancesRepo) GetByFolder(folder string) (*domain.Instance, error) {
	item := &domain.Instance{}

	query := `SELECT id, 
       			folder, 
       			endpoint, 
       			is_active
			FROM services 
			WHERE deleted_at is null and folder=$1
			ORDER BY created_at;`
	err := r.db.QueryRow(query, folder).Scan(&item.Id,
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

func (r *InstancesRepo) Insert(item *domain.Instance) error {
	var id int64
	query := "INSERT INTO services (folder, endpoint) VALUES ($1, $2) returning id"
	err := r.db.QueryRow(query, item.Folder, item.Endpoint).Scan(&id)
	if err != nil {
		return err
	}
	item.Id = id
	return nil
}

func (r *InstancesRepo) Delete(id int64) error {
	query := "UPDATE services SET deleted_at=now() WHERE id=$1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *InstancesRepo) Update(instance *domain.Instance) error {
	query := "UPDATE services SET folder=$2, endpoint=$3, is_active=$4, updated_at=now() WHERE id=$1"
	_, err := r.db.Exec(query, instance.Id, instance.Folder, instance.IsActive)
	if err != nil {
		return err
	}
	return nil
}
