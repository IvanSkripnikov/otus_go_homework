package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"app/components"
	"app/database"
	"app/kafka"
	"app/models"

	"github.com/gin-gonic/gin"
)

func GetAllBanners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := database.DB.Query(fmt.Sprintf("SELECT * from %s", "banners"))

	if checkError(w, err) {
		return
	}

	defer rows.Close()

	var banners []models.Banner
	for rows.Next() {
		banner := models.Banner{}
		if err = rows.Scan(&banner.Id, &banner.Title, &banner.Body, &banner.CreatedAt, &banner.Active); err != nil {
			log.Println(err.Error())
			continue
		}
		banners = append(banners, banner)
	}

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)

	err = je.Encode(&banners)
	if checkError(w, err) {
		return
	}

	writeSuccess(w, buf.String())
}

func GetBanner(w http.ResponseWriter, r *http.Request) {
	var banner models.Banner
	banner.Id, _ = getIdFromRequestString(r.URL.Path)

	if banner.Id == 0 {
		wrongParamsResponse(w)
		return
	}

	stmt, err := database.DB.Prepare(fmt.Sprintf("SELECT * from %s WHERE id = ?", "banners"))

	if checkError(w, err) {
		return
	}

	defer stmt.Close()

	if err = stmt.QueryRow(banner.Id).Scan(&banner.Id, &banner.Title, &banner.Body, &banner.CreatedAt, &banner.Active); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "{ \"message\": \"Not Found\"}")
		return
	}

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)

	err = je.Encode(&banner)
	if checkError(w, err) {
		return
	}

	writeSuccess(w, buf.String())
}

func AddBannerToSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)

	slotId, okSlot := params["slot"]
	bannerId, okBanner := params["banner"]

	if !okSlot || !okBanner || resultString != "" {
		wrongParamsResponse(w)
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (banner_id, slot_id) VALUES (?, ?)", "relations_banner_slot")
	stmt, err := database.DB.Prepare(query)

	if checkError(w, err) {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(bannerId, slotId)

	if checkError(w, err) {
		return
	}

	writeSuccess(w, "{\"message\": \"Successfully added!\"}")
}

func RemoveBannerFromSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)

	slotId, okSlot := params["slot"]
	bannerId, okBanner := params["banner"]

	if !okSlot || !okBanner || resultString != "" {
		wrongParamsResponse(w)
		return
	}

	stmt, err := database.DB.Prepare(fmt.Sprintf("DELETE FROM %s WHERE banner_id=? AND slot_id=?", "relations_banner_slot"))

	if checkError(w, err) {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(bannerId, slotId)

	if checkError(w, err) {
		return
	}

	writeSuccess(w, "{\"message\": \"Successfully removed!\"}")
}

func GetBannerForShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)

	slotId, okSlot := params["slot"]
	groupId, okGroup := params["group"]

	if !okSlot || !okGroup || resultString != "" {
		wrongParamsResponse(w)
		return
	}

	bannerId := components.GetNeedBanned(slotId, groupId)

	fmt.Fprint(w, strconv.Itoa(bannerId))
}

func EventClick(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)

	slotId, okSlot := params["slot"]
	groupId, okGroup := params["group"]
	bannerId, okBanner := params["banner"]

	if !okSlot || !okGroup || !okBanner || resultString != "" {
		wrongParamsResponse(w)
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (`type`, `banner_id`, `slot_id`, `group_id`) VALUES (?, ?, ?, ?)", "events")
	stmt, err := database.DB.Prepare(query)

	// отправляем событие в кафку
	sendEventToKafka("click", bannerId, slotId, groupId)

	if checkError(w, err) {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec("click", bannerId, slotId, groupId)

	if checkError(w, err) {
		return
	}

	fmt.Fprint(w, r)
}

// -------------PRIVATE----------------------

func checkError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return true
	}

	return false
}

func writeSuccess(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, message)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func wrongParamsResponse(w http.ResponseWriter) {
	resultString := "{\"message\": \"Invalid request GetHandler\"}"
	log.Println(resultString)
	fmt.Fprint(w, resultString)
	w.WriteHeader(http.StatusBadRequest)
}

func getIdFromRequestString(url string) (int, error) {
	vars := strings.Split(url, "/")

	return strconv.Atoi(vars[len(vars)-1])
}

func getParamsFromQueryString(url string) (map[string]int, string) {
	resultMap := map[string]int{}

	outMessage := ""
	queryParams := strings.Split(url, "/")

	params := strings.Split(queryParams[len(queryParams)-1], "&")
	if len(params) == 1 {
		outMessage = "not all params is set"
		return resultMap, outMessage
	}

	for _, v := range params {
		pair := strings.Split(v, "=")
		if len(pair) == 1 {
			outMessage = "incorrect params value: " + v
			return resultMap, outMessage
		} else {
			resultMap[pair[0]], _ = strconv.Atoi(pair[1])
		}
	}

	return resultMap, outMessage
}

func sendEventToKafka(eventName string, bannerId, slotId, groupId int) {
	message := kafka.Message{Type: eventName, BannerId: bannerId, SlotId: slotId, GroupId: groupId}

	producer, err := kafka.SetupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/send", kafka.SendMessageHandler(producer, message))

	fmt.Printf("Kafka PRODUCER 📨 started at http://localhost%s\n", kafka.ProducerPort)

	if err := router.Run(kafka.ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

/*
func executeQuery(w http.ResponseWriter, query string, ids ...int) bool {
	stmt, err := database.Db.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec("click", pq.Array(ids))

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return false
	}

	return true
}
*/
