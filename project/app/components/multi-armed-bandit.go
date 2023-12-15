package components

import (
	"app/database"
	"fmt"
	"log"
	"math"
)

func GetNeedBanned(slotId, groupId int) int {
	resultBannerId := 0

	// находим баннеры для данного слота
	bannersForSlot, err := GetSlotBanners(slotId)
	if err != nil {
		log.Fatal("error while search banners.")
	}
	for _, bannerId := range bannersForSlot {
		allShowsBanner := float64(GetBannerEvents(bannerId, "show"))
		allClickBanner := float64(GetBannerEvents(bannerId, "click"))
		allShows := float64(GetShows())
		averageRating := allClickBanner / allShowsBanner

		rate := GetRating(averageRating, allShowsBanner, allShows)

		fmt.Println(rate)
	}

	return resultBannerId
}

func GetShows() int {
	query := "SELECT COUNT(*) from events WHERE type = 'show'"
	stmt, err := database.DB.Query(query)

	if err != nil {
		return 0
	}

	defer stmt.Close()

	count := 0

	for stmt.Next() {
		if err := stmt.Scan(&count); err != nil {
			return 0
		}
	}

	return count
}

func GetBannerEvents(bannerId int, eventType string) int {
	query := "SELECT COUNT(*) from events WHERE banner_id = ? AND type = ?"
	stmt, err := database.DB.Query(query, bannerId, eventType)

	if err != nil {
		return 0
	}

	defer stmt.Close()

	count := 0

	for stmt.Next() {
		if err := stmt.Scan(&count); err != nil {
			return 0
		}
	}

	return count
}

func GetSlotBanners(slotId int) ([]int, error) {
	query := "SELECT banner_id from relations_banner_slot WHERE slot_id = ?"
	rows, err := database.DB.Query(query, slotId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	banners := make([]int, 0)
	banner := 0
	for rows.Next() {
		if err = rows.Scan(&banner); err != nil {
			log.Println(err.Error())
			continue
		}
		banners = append(banners, banner)
	}

	return banners, nil
}

func GetRating(averageRating float64, currentCount float64, allCounts float64) float64 {
	return averageRating + (math.Sqrt((2 * math.Log(allCounts)) / currentCount))
}
