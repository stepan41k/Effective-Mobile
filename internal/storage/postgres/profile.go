package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"github.com/stepan41k/Effective-Mobile/internal/domain/models"
	"github.com/stepan41k/Effective-Mobile/internal/storage"
)


func (s *PStorage) TakeProfiles(ctx context.Context, person models.GetPerson) ([]models.Person, error) {
	const op = "storage.postgres.profile.GetProfiles"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()
	
	arguments, values, ind := []string{}, []any{}, 1
	query := `SELECT name, surname, patronymic, age, gender, nationalize FROM profiles`

	if person.Name != "" || person.Surname != "" || person.Patronymic != "" || person.Gender != "" || person.Nationalize != "" || person.Age != 0 {
		query += ` WHERE `
		if person.Name != "" {
			arguments = append(arguments, fmt.Sprintf(`name LIKE $%d`, ind))
			values = append(values, person.Name)
			ind++
		}
			
		if person.Surname != "" {
			arguments = append(arguments, fmt.Sprintf(`surname LIKE $%d`, ind))
			values = append(values, person.Surname)
			ind++
		}
			
		if person.Patronymic != "" {
			arguments = append(arguments, fmt.Sprintf(`patronymic LIKE $%d`, ind))
			values = append(values, person.Patronymic)
			ind++
		}
			
		if person.Gender != "" {
			arguments = append(arguments, fmt.Sprintf(`gender LIKE $%d`, ind))
			values = append(values, person.Gender)
			ind++
		}
			
		if person.Nationalize != "" {
			arguments = append(arguments, fmt.Sprintf(`nationalize LIKE $%d`, ind))
			values = append(values, person.Nationalize)
			ind++
		}
			
		if person.Age != 0 {
			if person.Greater {
				arguments = append(arguments, fmt.Sprintf(`age > $%d`, ind))
			} else {
				arguments = append(arguments, fmt.Sprintf(`age < $%d`, ind))
			}
			values = append(values, person.Age)
			ind++
		}
	
		query += strings.Join(arguments, " AND ")
	}

	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d;`, ind, ind+1)
	values = append(values, person.PageSize, (person.Page-1)*person.PageSize)


	rows, err := tx.Query(ctx, query, values...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrProfilesNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var item models.Person
		err = rows.Scan(&item.Name, &item.Surname, &item.Patronymic, &item.Age, &item.Gender, &item.Nationalize)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		persons = append(persons, item)
	}
	
	return persons, err
}


func (s *PStorage) RemoveProfile(ctx context.Context, person models.DeletePerson) (guid []byte, err error) {
	const op = "storage.postgres.profile.DeleteProfile"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	cTag, err := tx.Exec(ctx, `
		DELETE FROM profiles
		WHERE guid = $1;
	`, person.GUID)

	if cTag.RowsAffected() == 0 {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrProfileNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return []byte(person.GUID), nil
}


func (s *PStorage) UpdateProfile(ctx context.Context, person models.UpdatedPerson) (guid []byte, err error) {
	const op = "storage.postgres.profile.UpdateProfile"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func ()  {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	arguments, values, ind := []string{}, []any{}, 1
	query := `UPDATE profiles SET `

	if person.Name != "" || person.Surname != "" || person.Patronymic != "" || person.Age != 0 || person.Gender != "" || person.Nationalize != ""  {
		if person.Name != "" {
			arguments = append(arguments, fmt.Sprintf(`name = $%d`, ind))
			values = append(values, person.Name)
			ind++
		}	
		if person.Surname != "" {
			arguments = append(arguments, fmt.Sprintf(`surname = $%d`, ind))
			values = append(values, person.Surname)
			ind++
		}	
		if person.Patronymic != "" {
			arguments = append(arguments, fmt.Sprintf(`patronymic = $%d`, ind))
			values = append(values, person.Patronymic)
			ind++
		}	
		if person.Age != 0 {
			arguments = append(arguments, fmt.Sprintf(`age = $%d`, ind))
			values = append(values, person.Age)
			ind++
		}	
		if person.Gender != "" {
			arguments = append(arguments, fmt.Sprintf(`gender = $%d`, ind))
			values = append(values, person.Gender)
			ind++
		}
		if person.Nationalize != "" {
			arguments = append(arguments, fmt.Sprintf(`nationalize = $%d`, ind))
			values = append(values, person.Nationalize)
			ind++
		}
	} else {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrNoChanges)
	}

	query += strings.Join(arguments, ",")
	query += fmt.Sprintf(` WHERE guid = $%d RETURNING guid;`, ind)
	values = append(values, []byte(person.GUID))


	row := tx.QueryRow(ctx, query, values...)
	err = row.Scan(&guid)

	if err != nil {
		if guid == nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrProfileNotFound)
		}

		if errors.Is(err, pgx.ErrNoRows){
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNoChanges)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return guid, nil
}


func (s *PStorage) NewProfile(ctx context.Context, person models.EnrichedPerson) (guid []byte, err error) {
	const op = "storage.postgres.profile.NewProfile"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	row := tx.QueryRow(ctx, `
		INSERT INTO profiles (guid, name, surname, patronymic, age, gender, nationalize)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING guid;
	`, []byte(person.GUID), person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationalize)

	err = row.Scan(&guid)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return guid, nil
}