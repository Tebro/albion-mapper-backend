package albion

import (
	"log"
	"time"

	"github.com/tebro/albion-mapper-backend/db"
)

// Zone describes a map in albion
type Zone struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Tier  int    `json:"tier"`
}

var validColors = []string{"black", "red", "yellow", "blue", "road"}
var validTiers = []int{4, 5, 6, 7, 8}

// IsValidZone checks if the color and tier are valid options
func IsValidZone(zone Zone) bool {
	validColor := false
	for _, s := range validColors {
		if zone.Color == s {
			validColor = true
		}
	}
	validTier := false
	for _, i := range validTiers {
		if zone.Tier == i {
			validTier = true
		}
	}

	return validColor && validTier
}

// Portal describes a roads portal between two zones
type Portal struct {
	id      int
	Source  string    `json:"source"`
	Target  string    `json:"target"`
	Size    int       `json:"size"`
	Expires time.Time `json:"expires"`
}

func zoneNameInZones(name string, zones []Zone) bool {
	for _, z := range zones {
		if z.Name == name {
			return true
		}
	}
	return false
}

// IsValidPortal checks that the portal is OK
func IsValidPortal(portal Portal) (bool, error) {
	if !(portal.Size == 2 || portal.Size == 7 || portal.Size == 20) {
		return false, nil
	}

	zones, err := GetZones()
	if err != nil {
		return false, err
	}

	return zoneNameInZones(portal.Source, zones) && zoneNameInZones(portal.Target, zones), nil
}

// GetZones returns all zones in the DB
func GetZones() ([]Zone, error) {
	db, err := db.GetDb()
	if err != nil {
		return []Zone{}, err
	}
	rows, err := db.Query("SELECT name, color, tier from zones;")
	if err != nil {
		return []Zone{}, err
	}
	defer rows.Close()

	zones := []Zone{}
	for rows.Next() {
		var zone Zone
		err = rows.Scan(&zone.Name, &zone.Color, &zone.Tier)
		if err != nil {
			return []Zone{}, err
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

// SetZone adds or updates a zones information
func SetZone(zone Zone) error {
	db, err := db.GetDb()
	if err != nil {
		return err
	}
	q, err := db.Query("REPLACE INTO zones (name, color, tier) VALUES (?, ?, ?);", zone.Name, zone.Color, zone.Tier)
	if err != nil {
		return err
	}
	defer q.Close()

	return err
}

// GetPortals returns the portals in the DB
func GetPortals() ([]Portal, error) {
	db, err := db.GetDb()
	if err != nil {
		return []Portal{}, err
	}
	rows, err := db.Query("SELECT id, source, target, size, expires from portals;")
	if err != nil {
		return []Portal{}, err
	}
	defer rows.Close()

	portals := []Portal{}
	for rows.Next() {
		var portal Portal
		var expires []uint8
		err = rows.Scan(&portal.id, &portal.Source, &portal.Target, &portal.Size, &expires)
		if err != nil {
			return []Portal{}, err
		}

		portal.Expires, err = time.Parse("2006-01-02 15:04:05", string(expires))
		if err != nil {
			return []Portal{}, err
		}
		portals = append(portals, portal)
	}

	return portals, nil
}

// AddPortal adds a new portal to the DB
func AddPortal(portal Portal) error {
	db, err := db.GetDb()
	if err != nil {
		return err
	}
	q, err := db.Query("INSERT INTO portals (source, target, size, expires) VALUES (?, ?, ?, ?);", portal.Source, portal.Target, portal.Size, portal.Expires)
	defer q.Close()

	return err
}

func deletePortal(portal Portal) error {
	db, err := db.GetDb()
	if err != nil {
		return nil
	}

	q, err := db.Query("DELETE FROM portals WHERE id = ?;", portal.id)
	defer q.Close()

	return err
}

// CleanupExpiredPortals does what you think
func CleanupExpiredPortals() error {
	portals, err := GetPortals()
	if err != nil {
		return err
	}

	toDelete := []Portal{}
	now := time.Now()

	for _, p := range portals {
		if p.Expires.Before(now) {
			toDelete = append(toDelete, p)
		}
	}

	for _, p := range toDelete {
		err = deletePortal(p)
		if err != nil {
			log.Printf("Unable to delete portal: %v", err)
		}
	}
	return err
}
