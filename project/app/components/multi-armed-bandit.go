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

		rate := getRating(averageRating, allShowsBanner, allShows)

		fmt.Println(rate)
	}

	return resultBannerId
}

func GetShows() int {
	stmt, err := database.Db.Prepare(fmt.Sprintf("SELECT COUNT(*) from %s WHERE type = %s", "events", "show"))

	if err != nil {
		return 0
	}

	defer stmt.Close()

	count := 0
	if err = stmt.QueryRow().Scan(&count); err != nil {
		return 0
	}

	return 0
}

func GetBannerEvents(bannerId int, eventType string) int {
	stmt, err := database.Db.Prepare(fmt.Sprintf("SELECT COUNT(*) from %s WHERE id = ? AND type = ?", "events"))

	if err != nil {
		return 0
	}

	defer stmt.Close()

	count := 0
	if err = stmt.QueryRow(bannerId, eventType).Scan(&count); err != nil {
		return 0
	}

	return 0
}

func GetSlotBanners(slotTd int) ([]int, error) {
	rows, err := database.Db.Query(fmt.Sprintf("SELECT banner_id from %s WHERE slot_id = ?", "relations_banner_slot"), slotTd)

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

func getRating(averageRating float64, currentCount float64, allCounts float64) float32 {
	return float32(averageRating + (math.Sqrt((2 * math.Log(allCounts)) / currentCount)))
}
