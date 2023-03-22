package repos

import (
	"database/sql"
	"microservice/app/core"
	"microservice/domain"
)

type RoutesRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewRoutesRepo(log core.Logger, db *sql.DB) *RoutesRepo {
	return &RoutesRepo{
		db:  db,
		log: log,
	}
}

func (r *RoutesRepo) All() ([]*domain.Route, error) {
	var items []*domain.Route

	query := `SELECT id, 
       			from_method, 
       			from_address,
       			instance,
       			proto_service, 
       			proto_method, 
       			access_role,
       			is_active
			FROM routes 
			WHERE deleted_at is null
			ORDER BY created_at;`
	raws, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for raws.Next() {
		item := &domain.Route{}
		err = raws.Scan(&item.Id,
			&item.HttpMethod,
			&item.HttpAddress,
			&item.Instance,
			&item.ProtoService,
			&item.ProtoMethod,
			&item.AccessRole,
			&item.IsActive)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *RoutesRepo) GetByAddress(addr string) (*domain.Route, error) {
	item := &domain.Route{}

	query := `SELECT id, 
       			from_method, 
       			from_address, 
       			instance,
       			proto_service, 
       			proto_method, 
       			access_role,
       			is_active
			FROM routes 
			WHERE deleted_at is null and from_address=$1
			ORDER BY created_at;`
	err := r.db.QueryRow(query, addr).Scan(&item.Id,
		&item.HttpMethod,
		&item.HttpAddress,
		&item.Instance,
		&item.ProtoService,
		&item.ProtoMethod,
		&item.AccessRole,
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

func (r *RoutesRepo) Insert(item *domain.Route) error {
	var id int64
	query := "INSERT INTO routes (from_method, from_address, instance, proto_service, proto_method, access_role) VALUES ($1, $2, $3, $4, $5, $6) returning id"
	err := r.db.QueryRow(query,
		item.HttpMethod,
		item.HttpAddress,
		item.Instance,
		item.ProtoService,
		item.ProtoMethod).Scan(&id)
	if err != nil {
		return err
	}
	item.Id = id
	return nil
}

func (r *RoutesRepo) Delete(id int64) error {
	query := "UPDATE services SET deleted_at=now() WHERE id=$1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
