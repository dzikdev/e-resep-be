package model

import "time"

type (
	Province struct {
		ID        int        `db:"id" json:"id"`
		AreaID    string     `db:"area_id" json:"area_id"`
		AreaName  string     `db:"area_name" json:"area_name"`
		Level     string     `db:"level" json:"level"`
		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		Delete    int        `db:"delete" json:"delete"`
	}

	City struct {
		ID         int        `db:"id" json:"id"`
		AreaID     string     `db:"area_id" json:"area_id"`
		AreaName   string     `db:"area_name" json:"area_name"`
		Level      string     `db:"level" json:"level"`
		ProvinceID int        `db:"province_id" json:"province_id"`
		CreatedAt  time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
		Delete     int        `db:"delete" json:"delete"`
	}

	District struct {
		ID         int        `db:"id" json:"id"`
		AreaID     string     `db:"area_id" json:"area_id"`
		AreaName   string     `db:"area_name" json:"area_name"`
		Level      string     `db:"level" json:"level"`
		ProvinceID int        `db:"province_id" json:"province_id"`
		CityID     int        `db:"city_id" json:"city_id"`
		CreatedAt  time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
		Delete     int        `db:"delete" json:"delete"`
	}

	SubDistrict struct {
		ID         int        `db:"id" json:"id"`
		AreaID     string     `db:"area_id" json:"area_id"`
		AreaName   string     `db:"area_name" json:"area_name"`
		Level      string     `db:"level" json:"level"`
		ProvinceID int        `db:"province_id" json:"province_id"`
		CityID     int        `db:"city_id" json:"city_id"`
		DistrictID int        `db:"district_id" json:"district_id"`
		CreatedAt  time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
		Delete     int        `db:"delete" json:"delete"`
	}
)
